package handlers

import (
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/tasksquad/dashboard/internal/parser"
)

type DashboardData struct {
	Dispatches    []parser.DispatchEntry
	Completions   []parser.CompletionReport
	Escalations   []parser.EscalationReport
	Backlog       *parser.BacklogOverview
	SessionLogs   []parser.SessionLog
	ReadyStories  []parser.BacklogStory
	DispatchFiles []parser.DispatchFile
	LastRefresh   time.Time
}

type Handler struct {
	worklogPath string
	tmpl        *template.Template
	data        *DashboardData
	mu          sync.RWMutex
	logger      *slog.Logger
}

func New(worklogPath string, tmpl *template.Template, logger *slog.Logger) *Handler {
	h := &Handler{
		worklogPath: worklogPath,
		tmpl:        tmpl,
		data:        &DashboardData{},
		logger:      logger,
	}
	h.refresh()
	return h
}

func (h *Handler) refresh() {
	h.mu.Lock()
	defer h.mu.Unlock()

	dispatches, err := parser.ParseDispatchLog(filepath.Join(h.worklogPath, "dispatch-log.md"))
	if err != nil {
		h.logger.Error("failed to parse dispatch log", "error", err)
	}
	h.data.Dispatches = dispatches

	completions, err := parser.ParseCompletions(filepath.Join(h.worklogPath, "completions"))
	if err != nil {
		h.logger.Error("failed to parse completions", "error", err)
	}
	h.data.Completions = completions

	escalations, err := parser.ParseEscalations(filepath.Join(h.worklogPath, "escalations"))
	if err != nil {
		h.logger.Error("failed to parse escalations", "error", err)
	}
	h.data.Escalations = escalations

	backlog, err := parser.ParseBacklog(filepath.Join(h.worklogPath, "backlog.md"))
	if err != nil {
		h.logger.Error("failed to parse backlog", "error", err)
	}
	h.data.Backlog = backlog
	h.data.ReadyStories = backlog.Ready

	sessionLogs, err := parser.ParseSessionLogs(filepath.Join(h.worklogPath, "session-logs"))
	if err != nil {
		h.logger.Error("failed to parse session logs", "error", err)
	}
	h.data.SessionLogs = sessionLogs

	dispatchFiles, err := parser.ParseDispatches(filepath.Join(h.worklogPath, "dispatches"))
	if err != nil {
		h.logger.Error("failed to parse dispatch files", "error", err)
	}
	h.data.DispatchFiles = dispatchFiles

	h.data.LastRefresh = time.Now()
}

func (h *Handler) StartRefreshLoop(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			h.refresh()
			h.logger.Info("refreshed dashboard data")
		}
	}()
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	data := h.data
	h.mu.RUnlock()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		h.logger.Error("failed to render template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *Handler) APIData(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	data := h.data
	h.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("failed to encode JSON", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *Handler) APIRefresh(w http.ResponseWriter, r *http.Request) {
	h.refresh()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "refreshed"})
}

func (h *Handler) APIDispatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req parser.DispatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.StoryID == "" || req.Repo == "" {
		http.Error(w, "story_id and repo are required", http.StatusBadRequest)
		return
	}

	dispatchDir := filepath.Join(h.worklogPath, "dispatches")
	path, err := parser.WriteDispatchFile(dispatchDir, req)
	if err != nil {
		h.logger.Error("failed to write dispatch file", "error", err)
		http.Error(w, "Failed to create dispatch file", http.StatusInternalServerError)
		return
	}

	if err := h.gitCommitAndPush(path, "Add dispatch: "+req.StoryID); err != nil {
		h.logger.Warn("git commit/push failed", "error", err)
	}

	h.refresh()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "created",
		"path":   path,
	})
}

func (h *Handler) SessionLogsFiltered(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	h.mu.RLock()
	logs := h.data.SessionLogs
	h.mu.RUnlock()

	var filtered []parser.SessionLog
	for _, log := range logs {
		if status == "" || log.Status == status {
			filtered = append(filtered, log)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filtered)
}

func (h *Handler) gitCommitAndPush(filePath, message string) error {
	relPath, err := filepath.Rel(h.worklogPath, filePath)
	if err != nil {
		relPath = filePath
	}

	addCmd := exec.Command("git", "-C", h.worklogPath, "add", relPath)
	if output, err := addCmd.CombinedOutput(); err != nil {
		h.logger.Warn("git add failed", "output", string(output), "error", err)
		return err
	}

	commitCmd := exec.Command("git", "-C", h.worklogPath, "commit", "-m", message)
	if output, err := commitCmd.CombinedOutput(); err != nil {
		if !isNothingToCommit(string(output)) {
			h.logger.Warn("git commit failed", "output", string(output), "error", err)
			return err
		}
	}

	pushCmd := exec.Command("git", "-C", h.worklogPath, "push")
	if output, err := pushCmd.CombinedOutput(); err != nil {
		h.logger.Warn("git push failed", "output", string(output), "error", err)
		return err
	}

	return nil
}

func isNothingToCommit(output string) bool {
	return len(output) > 0 && (
		contains(output, "nothing to commit") ||
		contains(output, "no changes added"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func (h *Handler) ServeStatic(w http.ResponseWriter, r *http.Request) {
	staticDir := filepath.Join(filepath.Dir(os.Args[0]), "static")
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		staticDir = "static"
	}
	http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))).ServeHTTP(w, r)
}
