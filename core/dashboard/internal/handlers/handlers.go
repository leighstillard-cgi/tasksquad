package handlers

import (
	"context"
	"encoding/json"
	"fmt"
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
	worklogPath    string
	tmpl           *template.Template
	data           *DashboardData
	mu             sync.RWMutex
	logger         *slog.Logger
	hub            *Hub
	watcher        *Watcher
	gitMu          sync.Mutex
	gitSync        func(filePaths []string, message string) error
	workerLauncher func(workerLaunchRequest) (*workerProcess, error)
}

func New(worklogPath string, tmpl *template.Template, logger *slog.Logger) *Handler {
	h := &Handler{
		worklogPath: worklogPath,
		tmpl:        tmpl,
		data:        &DashboardData{},
		logger:      logger,
		hub:         NewHub(logger),
	}
	h.gitSync = h.gitCommitAndPush
	h.workerLauncher = h.launchLocalWorker
	h.refresh()
	return h
}

func (h *Handler) dataDir(name string) string {
	consolidated := filepath.Join(h.worklogPath, "data", name)
	if _, err := os.Stat(consolidated); err == nil {
		return consolidated
	}
	return filepath.Join(h.worklogPath, name)
}

func (h *Handler) dataFile(name string) string {
	consolidated := filepath.Join(h.worklogPath, "data", name)
	if _, err := os.Stat(consolidated); err == nil {
		return consolidated
	}
	return filepath.Join(h.worklogPath, name)
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
	h.broadcastCurrentData()
}

func (h *Handler) broadcastCurrentData() {
	h.mu.RLock()
	data := *h.data
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

	dispatches, err := parser.ParseDispatchLog(h.dataFile("dispatch-log.md"))
	if err != nil {
		h.logger.Error("failed to parse dispatch log", "error", err)
	}
	h.data.Dispatches = dispatches

	completions, err := parser.ParseCompletions(h.dataDir("completions"))
	if err != nil {
		h.logger.Error("failed to parse completions", "error", err)
	}
	h.data.Completions = completions

	escalations, err := parser.ParseEscalations(h.dataDir("escalations"))
	if err != nil {
		h.logger.Error("failed to parse escalations", "error", err)
	}
	completionEscalations, err := parser.ParseCompletionEscalations(h.dataDir("completions"))
	if err != nil {
		h.logger.Error("failed to parse completion escalations", "error", err)
	}
	h.data.Escalations = append(escalations, completionEscalations...)

	backlog, err := parser.ParseBacklog(filepath.Join(h.worklogPath, "backlog.md"))
	if err != nil {
		h.logger.Error("failed to parse backlog", "error", err)
	}
	if backlog == nil {
		backlog = &parser.BacklogOverview{}
	}
	h.data.Backlog = backlog
	h.data.ReadyStories = backlog.Ready

	sessionLogs, err := parser.ParseSessionLogs(h.dataDir("session-logs"))
	if err != nil {
		h.logger.Error("failed to parse session logs", "error", err)
	}
	h.data.SessionLogs = sessionLogs

	dispatchFiles, err := parser.ParseDispatches(h.dataDir("dispatches"))
	if err != nil {
		h.logger.Error("failed to parse dispatch files", "error", err)
	}
	h.data.DispatchFiles = dispatchFiles
	h.data.Dispatches = projectDispatchFiles(dispatches, dispatchFiles, completions, backlog)
	h.data.Backlog = backlog
	h.data.ReadyStories = backlog.Ready

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

func writeJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
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

func findReadyStory(backlog *parser.BacklogOverview, storyID string) (parser.BacklogStory, bool) {
	if backlog == nil {
		return parser.BacklogStory{}, false
	}
	for _, story := range backlog.Ready {
		if story.StoryID == storyID {
			return story, true
		}
	}
	return parser.BacklogStory{}, false
}

func hasActiveDispatch(entries []parser.DispatchEntry, storyID string) bool {
	for _, entry := range entries {
		if entry.StoryID == storyID && isActiveDispatchStatus(entry.Status) {
			return true
		}
	}
	return false
}

func isActiveDispatchStatus(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "dispatched", "in-progress":
		return true
	default:
		return false
	}
}

func isTerminalDispatchStatus(status string) bool {
	if isDoneDispatchStatus(status) {
		return true
	}
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "cancelled", "canceled":
		return true
	default:
		return false
	}
}

func isDoneDispatchStatus(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "complete", "completed", "done":
		return true
	default:
		return false
	}
}

func projectDispatchFiles(dispatches []parser.DispatchEntry, dispatchFiles []parser.DispatchFile, completions []parser.CompletionReport, backlog *parser.BacklogOverview) []parser.DispatchEntry {
	completed := completionStoryIDs(completions)
	terminal := make(map[string]bool, len(dispatches)+len(completed))
	for storyID := range completed {
		terminal[storyID] = true
	}

	for _, entry := range dispatches {
		if entry.StoryID == "" {
			continue
		}
		if isDoneDispatchStatus(entry.Status) {
			completed[entry.StoryID] = true
			terminal[entry.StoryID] = true
			continue
		}
		if isTerminalDispatchStatus(entry.Status) {
			terminal[entry.StoryID] = true
		}
	}

	projectCompletedStoriesToDone(backlog, completed)
	removeStoriesFromInProgress(backlog, terminal)

	active := make(map[string]bool, len(dispatches))
	projected := make([]parser.DispatchEntry, 0, len(dispatches)+len(dispatchFiles))
	for _, entry := range dispatches {
		if entry.StoryID == "" || terminal[entry.StoryID] || !isActiveDispatchStatus(entry.Status) {
			continue
		}

		projected = append(projected, entry)
		active[entry.StoryID] = true
	}

	for _, dispatchFile := range dispatchFiles {
		if dispatchFile.StoryID == "" || active[dispatchFile.StoryID] || terminal[dispatchFile.StoryID] {
			continue
		}

		story, ok := moveBacklogStoryToInProgress(backlog, dispatchFile.StoryID)
		if !ok {
			continue
		}

		entry := parser.DispatchEntry{
			StoryID:      dispatchFile.StoryID,
			Repo:         story.Repo,
			DispatchedAt: dispatchFile.DispatchedAt,
			Status:       "dispatched",
		}
		if !entry.DispatchedAt.IsZero() {
			entry.ElapsedTime = time.Since(entry.DispatchedAt)
		}
		projected = append(projected, entry)
		active[dispatchFile.StoryID] = true
	}

	return projected
}

func completionStoryIDs(completions []parser.CompletionReport) map[string]bool {
	storyIDs := make(map[string]bool, len(completions))
	for _, completion := range completions {
		if completion.StoryID != "" {
			storyIDs[completion.StoryID] = true
		}
	}
	return storyIDs
}

func projectCompletedStoriesToDone(backlog *parser.BacklogOverview, storyIDs map[string]bool) {
	if backlog == nil || len(storyIDs) == 0 {
		return
	}

	done := make(map[string]bool, len(backlog.Done)+len(storyIDs))
	for _, story := range backlog.Done {
		done[story.StoryID] = true
	}

	for storyID := range storyIDs {
		if done[storyID] {
			continue
		}

		story, found := removeStoryFromBacklogBuckets(backlog, storyID)
		if !found {
			story = parser.BacklogStory{StoryID: storyID}
		}
		story.Status = "done"
		backlog.Done = append(backlog.Done, story)
		done[storyID] = true
	}
}

func removeStoryFromBacklogBuckets(backlog *parser.BacklogOverview, storyID string) (parser.BacklogStory, bool) {
	var story parser.BacklogStory
	var found bool
	backlog.InProgress, story, found = removeStoryFromBacklogStories(backlog.InProgress, storyID)
	if found {
		return story, true
	}
	backlog.Ready, story, found = removeStoryFromBacklogStories(backlog.Ready, storyID)
	if found {
		return story, true
	}
	backlog.Blocked, story, found = removeStoryFromBacklogStories(backlog.Blocked, storyID)
	if found {
		return story, true
	}
	backlog.Cancelled, story, found = removeStoryFromBacklogStories(backlog.Cancelled, storyID)
	return story, found
}

func removeStoryFromBacklogStories(stories []parser.BacklogStory, storyID string) ([]parser.BacklogStory, parser.BacklogStory, bool) {
	for i, story := range stories {
		if story.StoryID != storyID {
			continue
		}
		return append(stories[:i], stories[i+1:]...), story, true
	}
	return stories, parser.BacklogStory{}, false
}

func removeStoriesFromInProgress(backlog *parser.BacklogOverview, storyIDs map[string]bool) {
	if backlog == nil || len(storyIDs) == 0 || len(backlog.InProgress) == 0 {
		return
	}

	kept := backlog.InProgress[:0]
	for _, story := range backlog.InProgress {
		if !storyIDs[story.StoryID] {
			kept = append(kept, story)
		}
	}
	backlog.InProgress = kept
}

func moveBacklogStoryToInProgress(backlog *parser.BacklogOverview, storyID string) (parser.BacklogStory, bool) {
	if backlog == nil {
		return parser.BacklogStory{}, false
	}

	for _, story := range backlog.InProgress {
		if story.StoryID == storyID {
			return story, true
		}
	}

	for i, story := range backlog.Ready {
		if story.StoryID != storyID {
			continue
		}
		story.Status = "in-progress"
		backlog.Ready = append(backlog.Ready[:i], backlog.Ready[i+1:]...)
		backlog.InProgress = append(backlog.InProgress, story)
		return story, true
	}

	return parser.BacklogStory{}, false
}

func (h *Handler) APIDispatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req parser.DispatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.StoryID == "" || req.Repo == "" {
		writeJSONError(w, "story_id and repo are required", http.StatusBadRequest)
		return
	}

	// Validate inputs to prevent path traversal and injection
	if !isValidStoryID(req.StoryID) {
		writeJSONError(w, "Invalid story_id format", http.StatusBadRequest)
		return
	}
	if !isValidRepoName(req.Repo) {
		writeJSONError(w, "Invalid repo format", http.StatusBadRequest)
		return
	}

	req.DispatchedAt = time.Now().UTC()

	backlogPath := filepath.Join(h.worklogPath, "backlog.md")
	backlog, err := parser.ParseBacklog(backlogPath)
	if err != nil {
		h.logger.Error("failed to parse backlog before dispatch", "error", err)
		writeJSONError(w, "Failed to read backlog", http.StatusInternalServerError)
		return
	}
	story, ok := findReadyStory(backlog, req.StoryID)
	if !ok {
		writeJSONError(w, "story is not ready for dispatch", http.StatusConflict)
		return
	}

	dispatchLogPath := h.dataFile("dispatch-log.md")
	dispatches, err := parser.ParseDispatchLog(dispatchLogPath)
	if err != nil {
		h.logger.Error("failed to parse dispatch log before dispatch", "error", err)
		writeJSONError(w, "Failed to read dispatch log", http.StatusInternalServerError)
		return
	}
	if hasActiveDispatch(dispatches, req.StoryID) {
		writeJSONError(w, "story is already dispatched", http.StatusConflict)
		return
	}

	dispatchDir := h.dataDir("dispatches")
	path, err := parser.WriteDispatchFile(dispatchDir, req)
	if err != nil {
		h.logger.Error("failed to write dispatch file", "error", err)
		writeJSONError(w, "Failed to create dispatch file", http.StatusInternalServerError)
		return
	}
	if err := parser.AppendDispatchLogEntry(dispatchLogPath, req); err != nil {
		h.logger.Error("failed to update dispatch log", "error", err)
		writeJSONError(w, "Failed to update dispatch log", http.StatusInternalServerError)
		return
	}
	if err := parser.UpdateBacklogStoryStatus(backlogPath, req.StoryID, "in-progress", "ready"); err != nil {
		h.logger.Error("failed to update backlog", "error", err)
		writeJSONError(w, "Failed to update backlog", http.StatusInternalServerError)
		return
	}

	workerStatus := "disabled"
	workerPID := ""
	workerLog := ""
	workerError := ""
	if h.workerLauncher != nil {
		worker, err := h.workerLauncher(workerLaunchRequest{
			Story:        story,
			Dispatch:     req,
			DispatchPath: path,
		})
		if err != nil {
			workerStatus = "failed"
			workerError = err.Error()
			h.logger.Warn("failed to start worker", "story_id", req.StoryID, "error", err)
		} else {
			workerStatus = "started"
			workerPID = fmt.Sprintf("%d", worker.PID)
			workerLog = worker.SessionLogPath
		}
	}

	h.refresh()
	h.broadcastCurrentData()
	h.queueGitSync([]string{path, dispatchLogPath, backlogPath}, "Dispatch "+req.StoryID)

	response := map[string]string{
		"status":        "created",
		"path":          path,
		"git_status":    "queued",
		"worker_status": workerStatus,
	}
	if workerPID != "" {
		response["worker_pid"] = workerPID
	}
	if workerLog != "" {
		response["worker_log"] = workerLog
	}
	if workerError != "" {
		response["worker_error"] = workerError
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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

const gitSyncTimeout = 10 * time.Second

func (h *Handler) queueGitSync(filePaths []string, message string) {
	if h.gitSync == nil {
		return
	}
	paths := append([]string(nil), filePaths...)

	go func() {
		h.gitMu.Lock()
		defer h.gitMu.Unlock()

		if err := h.gitSync(paths, message); err != nil {
			h.logger.Warn("git commit/push failed", "error", err)
		}
	}()
}

func (h *Handler) gitCommitAndPush(filePaths []string, message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), gitSyncTimeout)
	defer cancel()

	for _, filePath := range filePaths {
		relPath, err := filepath.Rel(h.worklogPath, filePath)
		if err != nil {
			relPath = filePath
		}

		addCmd := gitCommand(ctx, h.worklogPath, "add", relPath)
		if output, err := addCmd.CombinedOutput(); err != nil {
			if ctx.Err() != nil {
				return fmt.Errorf("git commit/push timed out after %s", gitSyncTimeout)
			}
			h.logger.Warn("git add failed", "output", string(output), "error", err)
			return err
		}
	}

	commitCmd := gitCommand(ctx, h.worklogPath, "commit", "-m", message)
	if output, err := commitCmd.CombinedOutput(); err != nil {
		if ctx.Err() != nil {
			return fmt.Errorf("git commit/push timed out after %s", gitSyncTimeout)
		}
		if !isNothingToCommit(string(output)) {
			h.logger.Warn("git commit failed", "output", string(output), "error", err)
			return err
		}
	}

	pushCmd := gitCommand(ctx, h.worklogPath, "push")
	if output, err := pushCmd.CombinedOutput(); err != nil {
		if ctx.Err() != nil {
			return fmt.Errorf("git commit/push timed out after %s", gitSyncTimeout)
		}
		h.logger.Warn("git push failed", "output", string(output), "error", err)
		return err
	}

	return nil
}

func gitCommand(ctx context.Context, worklogPath string, args ...string) *exec.Cmd {
	gitArgs := append([]string{"-C", worklogPath}, args...)
	cmd := exec.CommandContext(ctx, "git", gitArgs...)
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	return cmd
}

func isNothingToCommit(output string) bool {
	return len(output) > 0 && (strings.Contains(output, "nothing to commit") ||
		strings.Contains(output, "no changes added"))
}

func (h *Handler) ServeStatic(w http.ResponseWriter, r *http.Request) {
	staticDir := filepath.Join(filepath.Dir(os.Args[0]), "static")
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		staticDir = "static"
	}
	http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))).ServeHTTP(w, r)
}
