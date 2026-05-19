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
)

func newTestHandler(t *testing.T, worklogPath string) *Handler {
	t.Helper()

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	tmpl := template.Must(template.New("index.html").Parse(""))

	h := New(worklogPath, tmpl, logger)
	h.gitSync = func(filePath, message string) error { return nil }
	return h
}

func TestAPIDispatchAcceptsSnakeCaseJSON(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmpDir, "dispatches"), 0755); err != nil {
		t.Fatal(err)
	}

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
	if !strings.HasSuffix(result["path"], filepath.Join("dispatches", "STORY-99.1-dispatch.md")) {
		t.Fatalf("unexpected dispatch path %q", result["path"])
	}
}

func TestAPIDispatchDoesNotWaitForGitSync(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmpDir, "dispatches"), 0755); err != nil {
		t.Fatal(err)
	}

	h := newTestHandler(t, tmpDir)
	releaseGitSync := make(chan struct{})
	h.gitSync = func(filePath, message string) error {
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
