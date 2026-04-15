package handlers

import (
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	hub         *Hub
	watcher     *Watcher
}

func New(worklogPath string, tmpl *template.Template, logger *slog.Logger) *Handler {
	h := &Handler{
		worklogPath: worklogPath,
		tmpl:        tmpl,
		data:        &DashboardData{},
		logger:      logger,
		hub:         NewHub(logger),
	}
	h.refresh()
	return h
}

// Hub returns the WebSocket hub.
func (h *Handler) Hub() *Hub {
	return h.hub
}

// StartLiveUpdates starts the WebSocket hub and file watcher.
func (h *Handler) StartLiveUpdates() error {
	// Start the hub
	go h.hub.Run()

	// Create and start the watcher
	watcher, err := NewWatcher(h.worklogPath, h.logger, h.broadcastUpdate)
	if err != nil {
		return err
	}
	h.watcher = watcher

	if err := h.watcher.Start(); err != nil {
		return err
	}

	h.logger.Info("live updates enabled")
	return nil
}

// broadcastUpdate refreshes data and broadcasts to all WebSocket clients.
func (h *Handler) broadcastUpdate() {
	h.refresh()

	h.mu.RLock()
	data := h.data
	h.mu.RUnlock()

	jsonData, err := json.Marshal(data)
	if err != nil {
		h.logger.Error("failed to marshal data for broadcast", "error", err)
		return
	}

	h.hub.Broadcast(jsonData)
	h.logger.Debug("broadcast update sent", "clients", h.hub.ClientCount())
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

// isValidStoryID validates story IDs match expected format (e.g., STORY-01.2)
func isValidStoryID(id string) bool {
	if id == "" || len(id) > 50 {
		return false
	}
	// Only allow alphanumeric, dash, underscore, and dot
	for _, c := range id {
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') ||
			(c >= '0' && c <= '9') || c == '-' || c == '_' || c == '.') {
			return false
		}
	}
	return true
}

// isValidRepoName validates repository names
func isValidRepoName(name string) bool {
	if name == "" || len(name) > 100 {
		return false
	}
	// Only allow alphanumeric, dash, underscore
	for _, c := range name {
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') ||
			(c >= '0' && c <= '9') || c == '-' || c == '_') {
			return false
		}
	}
	return true
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

	// Validate inputs to prevent path traversal and injection
	if !isValidStoryID(req.StoryID) {
		http.Error(w, "Invalid story_id format", http.StatusBadRequest)
		return
	}
	if !isValidRepoName(req.Repo) {
		http.Error(w, "Invalid repo format", http.StatusBadRequest)
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
		strings.Contains(output, "nothing to commit") ||
		strings.Contains(output, "no changes added"))
}


func (h *Handler) ServeStatic(w http.ResponseWriter, r *http.Request) {
	staticDir := filepath.Join(filepath.Dir(os.Args[0]), "static")
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		staticDir = "static"
	}
	http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))).ServeHTTP(w, r)
}
