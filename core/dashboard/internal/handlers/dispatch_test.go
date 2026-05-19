package handlers

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAPIDispatchReturnsBeforeGitSyncCompletes(t *testing.T) {
	worklogPath := t.TempDir()
	if err := os.MkdirAll(filepath.Join(worklogPath, "data", "dispatches"), 0755); err != nil {
		t.Fatalf("create dispatch dir: %v", err)
	}

	h := New(worklogPath, template.New("test"), slog.New(slog.NewTextHandler(io.Discard, nil)))
	gitStarted := make(chan struct{})
	gitRelease := make(chan struct{})
	h.gitSync = func(filePath, message string) error {
		close(gitStarted)
		<-gitRelease
		return nil
	}
	t.Cleanup(func() { close(gitRelease) })

	body := bytes.NewBufferString(`{"story_id":"STORY-10.1","repo":"tasksquad","description":"Create a dispatch"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/dispatch", body)
	rec := httptest.NewRecorder()

	h.APIDispatch(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var payload map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload["status"] != "created" {
		t.Fatalf("expected created status, got %q", payload["status"])
	}
	if payload["git_status"] != "queued" {
		t.Fatalf("expected queued git_status, got %q", payload["git_status"])
	}
	if payload["path"] == "" {
		t.Fatal("expected dispatch path in response")
	}

	select {
	case <-gitStarted:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("expected git sync to be queued")
	}
}
