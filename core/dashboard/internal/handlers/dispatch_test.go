package handlers

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/tasksquad/dashboard/internal/parser"
)

func newTestHandler(t *testing.T, worklogPath string) *Handler {
	t.Helper()

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	tmpl := template.Must(template.New("index.html").Parse(""))

	h := New(worklogPath, tmpl, logger)
	h.gitSync = func(filePaths []string, message string) error { return nil }
	h.workerLauncher = func(req workerLaunchRequest) (*workerProcess, error) {
		return &workerProcess{
			PID:            1234,
			SessionLogPath: filepath.Join(worklogPath, "data", "session-logs", req.Dispatch.StoryID+"-running.md"),
		}, nil
	}
	return h
}

func setupDispatchWorklog(t *testing.T, worklogPath, storyID string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Join(worklogPath, "data", "dispatches"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(worklogPath, "data", "dispatch-log.md"), []byte(testDispatchLogHeader), 0644); err != nil {
		t.Fatal(err)
	}

	backlog := `# TaskSquad Product Backlog

## Test Work

### ` + storyID + ` - Test Dispatch
**Status:** ready
**Repo:** tasksquad
**Priority:** High

**Description:** Create dashboard
`
	if err := os.WriteFile(filepath.Join(worklogPath, "backlog.md"), []byte(backlog), 0644); err != nil {
		t.Fatal(err)
	}
}

const testDispatchLogHeader = `# Dispatch Log

PM agent maintains this file. Updated on every dispatch and completion.

| Story | Repo | Dispatched At | Status | Completion Report |
|---|---|---|---|---|
`

func TestAPIDispatchAcceptsSnakeCaseJSON(t *testing.T) {
	tmpDir := t.TempDir()
	setupDispatchWorklog(t, tmpDir, "STORY-99.1")

	h := newTestHandler(t, tmpDir)
	body := bytes.NewBufferString(`{"story_id":"STORY-99.1","repo":"tasksquad","description":"Create dashboard"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/dispatch", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.APIDispatch(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("expected application/json content type, got %q", contentType)
	}

	var result map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("expected JSON response: %v", err)
	}
	if result["status"] != "created" {
		t.Fatalf("expected created status, got %q", result["status"])
	}
	if result["git_status"] != "queued" {
		t.Fatalf("expected queued git status, got %q", result["git_status"])
	}
	if result["worker_status"] != "started" {
		t.Fatalf("expected started worker status, got %q", result["worker_status"])
	}
	if result["worker_pid"] != "1234" {
		t.Fatalf("expected worker pid 1234, got %q", result["worker_pid"])
	}
	if !strings.HasSuffix(result["path"], filepath.Join("data", "dispatches", "STORY-99.1-dispatch.md")) {
		t.Fatalf("unexpected dispatch path %q", result["path"])
	}

	dispatchLog, err := os.ReadFile(filepath.Join(tmpDir, "data", "dispatch-log.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(dispatchLog), "| STORY-99.1 | tasksquad |") || !strings.Contains(string(dispatchLog), "| dispatched |") {
		t.Fatalf("dispatch log was not updated: %s", string(dispatchLog))
	}

	backlog, err := os.ReadFile(filepath.Join(tmpDir, "backlog.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(backlog), "**Status:** in-progress") {
		t.Fatalf("backlog status was not updated: %s", string(backlog))
	}

	h.mu.RLock()
	defer h.mu.RUnlock()
	if len(h.data.Dispatches) != 1 || h.data.Dispatches[0].StoryID != "STORY-99.1" {
		t.Fatalf("expected dispatched story in active data, got %#v", h.data.Dispatches)
	}
	if h.data.Backlog == nil || len(h.data.Backlog.InProgress) != 1 || h.data.Backlog.InProgress[0].StoryID != "STORY-99.1" {
		t.Fatalf("expected story in in-progress backlog data, got %#v", h.data.Backlog)
	}
}

func TestRefreshProjectsDispatchFileIntoActiveData(t *testing.T) {
	tmpDir := t.TempDir()
	setupDispatchWorklog(t, tmpDir, "STORY-99.3")

	_, err := parser.WriteDispatchFile(filepath.Join(tmpDir, "data", "dispatches"), parser.DispatchRequest{
		StoryID:      "STORY-99.3",
		Repo:         "tasksquad",
		DispatchedAt: time.Date(2026, 5, 20, 1, 2, 3, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("WriteDispatchFile failed: %v", err)
	}

	h := newTestHandler(t, tmpDir)

	h.mu.RLock()
	defer h.mu.RUnlock()
	if len(h.data.Dispatches) != 1 || h.data.Dispatches[0].StoryID != "STORY-99.3" {
		t.Fatalf("expected dispatch file to project into active data, got %#v", h.data.Dispatches)
	}
	if h.data.Backlog == nil || len(h.data.Backlog.InProgress) != 1 || h.data.Backlog.InProgress[0].StoryID != "STORY-99.3" {
		t.Fatalf("expected dispatch file to project into in-progress backlog data, got %#v", h.data.Backlog)
	}
	if len(h.data.ReadyStories) != 0 {
		t.Fatalf("expected projected story to be removed from ready stories, got %#v", h.data.ReadyStories)
	}
}

func TestRefreshRemovesCompletedDispatchFromActiveData(t *testing.T) {
	tmpDir := t.TempDir()
	setupDispatchWorklog(t, tmpDir, "STORY-99.5")

	backlogPath := filepath.Join(tmpDir, "backlog.md")
	backlog, err := os.ReadFile(backlogPath)
	if err != nil {
		t.Fatal(err)
	}
	backlog = []byte(strings.Replace(string(backlog), "**Status:** ready", "**Status:** in-progress", 1))
	if err := os.WriteFile(backlogPath, backlog, 0644); err != nil {
		t.Fatal(err)
	}

	dispatchedAt := time.Date(2026, 5, 20, 1, 2, 3, 0, time.UTC)
	_, err = parser.WriteDispatchFile(filepath.Join(tmpDir, "data", "dispatches"), parser.DispatchRequest{
		StoryID:      "STORY-99.5",
		Repo:         "tasksquad",
		DispatchedAt: dispatchedAt,
	})
	if err != nil {
		t.Fatalf("WriteDispatchFile failed: %v", err)
	}
	dispatchRow := "| STORY-99.5 | tasksquad | " + dispatchedAt.Format(time.RFC3339) + " | dispatched | |\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "data", "dispatch-log.md"), []byte(testDispatchLogHeader+dispatchRow), 0644); err != nil {
		t.Fatal(err)
	}

	completionDir := filepath.Join(tmpDir, "data", "completions")
	if err := os.MkdirAll(completionDir, 0755); err != nil {
		t.Fatal(err)
	}
	completion := `---
title: "Completion: STORY-99.5"
id: "COMPLETION-STORY-99.5"
status: draft
created: "2026-05-20T01:03:00Z"
last_updated: "2026-05-20T01:03:00Z"
parent_epic: "EPIC-99"
phase: "test"
repos: ["tasksquad"]
---

# Completion: STORY-99.5
`
	if err := os.WriteFile(filepath.Join(completionDir, "STORY-99.5-completion.md"), []byte(completion), 0644); err != nil {
		t.Fatal(err)
	}

	h := newTestHandler(t, tmpDir)

	h.mu.RLock()
	defer h.mu.RUnlock()
	if len(h.data.Dispatches) != 0 {
		t.Fatalf("expected completed story to be removed from active dispatches, got %#v", h.data.Dispatches)
	}
	if h.data.Backlog == nil || len(h.data.Backlog.InProgress) != 0 {
		t.Fatalf("expected completed story to be removed from in-progress backlog projection, got %#v", h.data.Backlog)
	}
	if h.data.Backlog == nil || len(h.data.Backlog.Done) != 1 || h.data.Backlog.Done[0].StoryID != "STORY-99.5" || h.data.Backlog.Done[0].Status != "done" {
		t.Fatalf("expected completed story to be projected into done backlog data, got %#v", h.data.Backlog)
	}
	if len(h.data.Completions) != 1 || h.data.Completions[0].StoryID != "STORY-99.5" {
		t.Fatalf("expected completion report to remain visible, got %#v", h.data.Completions)
	}
}

func TestRefreshIncludesCompletionEscalations(t *testing.T) {
	tmpDir := t.TempDir()
	setupDispatchWorklog(t, tmpDir, "STORY-99.6")

	completionDir := filepath.Join(tmpDir, "data", "completions")
	if err := os.MkdirAll(completionDir, 0755); err != nil {
		t.Fatal(err)
	}
	completion := `---
title: "Completion: STORY-99.6"
id: "COMPLETION-STORY-99.6"
status: draft
created: "2026-05-20T01:00:00Z"
last_updated: "2026-05-20T01:02:03Z"
parent_epic: "EPIC-99"
phase: "test"
repos: ["tasksquad"]
---

# Completion: STORY-99.6

## Architectural Escalations

- Missing Codex runtime files need architectural review.
- Dashboard should surface completion report escalations.
`
	if err := os.WriteFile(filepath.Join(completionDir, "STORY-99.6-completion.md"), []byte(completion), 0644); err != nil {
		t.Fatal(err)
	}

	h := newTestHandler(t, tmpDir)

	h.mu.RLock()
	defer h.mu.RUnlock()
	if len(h.data.Escalations) != 2 {
		t.Fatalf("expected completion escalations in dashboard data, got %#v", h.data.Escalations)
	}
	if h.data.Escalations[0].StoryID != "STORY-99.6" || h.data.Escalations[0].Reason != "Missing Codex runtime files need architectural review." {
		t.Fatalf("unexpected first escalation: %#v", h.data.Escalations[0])
	}
}

func TestAPIDispatchPassesStoryToWorkerLauncher(t *testing.T) {
	tmpDir := t.TempDir()
	setupDispatchWorklog(t, tmpDir, "STORY-99.4")

	h := newTestHandler(t, tmpDir)
	var launched workerLaunchRequest
	h.workerLauncher = func(req workerLaunchRequest) (*workerProcess, error) {
		launched = req
		return &workerProcess{PID: 4321, SessionLogPath: filepath.Join(tmpDir, "data", "session-logs", "worker.md")}, nil
	}

	body := bytes.NewBufferString(`{"story_id":"STORY-99.4","repo":"tasksquad","description":"Create dashboard"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/dispatch", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.APIDispatch(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if launched.Dispatch.StoryID != "STORY-99.4" {
		t.Fatalf("worker launcher got story %q", launched.Dispatch.StoryID)
	}
	if launched.Story.Title != "Test Dispatch" {
		t.Fatalf("worker launcher got title %q", launched.Story.Title)
	}
	if !strings.HasSuffix(launched.DispatchPath, filepath.Join("data", "dispatches", "STORY-99.4-dispatch.md")) {
		t.Fatalf("worker launcher got dispatch path %q", launched.DispatchPath)
	}
}

func TestAPIDispatchDoesNotWaitForGitSync(t *testing.T) {
	tmpDir := t.TempDir()
	setupDispatchWorklog(t, tmpDir, "STORY-99.2")

	h := newTestHandler(t, tmpDir)
	releaseGitSync := make(chan struct{})
	h.gitSync = func(filePaths []string, message string) error {
		<-releaseGitSync
		return nil
	}
	defer close(releaseGitSync)

	body := bytes.NewBufferString(`{"story_id":"STORY-99.2","repo":"tasksquad","description":"Create dashboard"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/dispatch", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	done := make(chan struct{})
	go func() {
		h.APIDispatch(rec, req)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("dispatch response waited for git sync")
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestAPIDispatchReturnsJSONError(t *testing.T) {
	tmpDir := t.TempDir()
	h := newTestHandler(t, tmpDir)
	req := httptest.NewRequest(http.MethodPost, "/api/dispatch", bytes.NewBufferString(`{"repo":"tasksquad"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.APIDispatch(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
	if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("expected application/json content type, got %q", contentType)
	}

	var result map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("expected JSON error response: %v", err)
	}
	if result["error"] != "story_id and repo are required" {
		t.Fatalf("unexpected error %q", result["error"])
	}
}

func TestWorkerShellScriptUsesGlobalCodexApprovalFlag(t *testing.T) {
	script := workerShellScript()

	invalid := `"$CODEX_BIN" exec -C "$WORKLOG_PATH" --sandbox workspace-write --ask-for-approval never - < "$PROMPT_PATH"`
	if strings.Contains(script, invalid) {
		t.Fatalf("worker script uses unsupported codex exec flag order: %s", invalid)
	}

	expected := `"$CODEX_BIN" --ask-for-approval never exec -C "$WORKLOG_PATH" --sandbox workspace-write - < "$PROMPT_PATH"`
	if !strings.Contains(script, expected) {
		t.Fatalf("worker script missing supported codex launch command %q in:\n%s", expected, script)
	}

	logged := `Command: %s --ask-for-approval never exec -C %s --sandbox workspace-write -`
	if !strings.Contains(script, logged) {
		t.Fatalf("worker script logs an unexpected command shape:\n%s", script)
	}
}
