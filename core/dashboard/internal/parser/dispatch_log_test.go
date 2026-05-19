package parser

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParseDispatchLog(t *testing.T) {
	content := `# Dispatch Log

PM agent maintains this file. Updated on every dispatch and completion.

| Story | Repo | Dispatched At | Status | Completion Report |
|---|---|---|---|---|
| STORY-03.4 | tasksquad | 2026-04-15T02:13:00Z | done | data/completions/STORY-03.4-completion.md |
| STORY-03.5 | tasksquad | 2026-04-15T02:26:00Z | in-progress | |
`

	dir := t.TempDir()
	path := filepath.Join(dir, "dispatch-log.md")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	entries, err := ParseDispatchLog(path)
	if err != nil {
		t.Fatalf("ParseDispatchLog failed: %v", err)
	}

	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}

	if entries[0].StoryID != "STORY-03.4" {
		t.Errorf("expected STORY-03.4, got %s", entries[0].StoryID)
	}

	if entries[0].Status != "done" {
		t.Errorf("expected done, got %s", entries[0].Status)
	}

	if entries[1].Status != "in-progress" {
		t.Errorf("expected in-progress, got %s", entries[1].Status)
	}
}

func TestAppendDispatchLogEntry(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data", "dispatch-log.md")
	dispatchedAt := time.Date(2026, 5, 20, 1, 2, 3, 0, time.UTC)

	err := AppendDispatchLogEntry(path, DispatchRequest{
		StoryID:      "STORY-99.1",
		Repo:         "tasksquad",
		DispatchedAt: dispatchedAt,
	})
	if err != nil {
		t.Fatalf("AppendDispatchLogEntry failed: %v", err)
	}

	entries, err := ParseDispatchLog(path)
	if err != nil {
		t.Fatalf("ParseDispatchLog failed: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].StoryID != "STORY-99.1" {
		t.Fatalf("expected STORY-99.1, got %s", entries[0].StoryID)
	}
	if entries[0].Status != "dispatched" {
		t.Fatalf("expected dispatched status, got %s", entries[0].Status)
	}
	if !entries[0].DispatchedAt.Equal(dispatchedAt) {
		t.Fatalf("expected dispatched time %s, got %s", dispatchedAt, entries[0].DispatchedAt)
	}
}

func TestParseDispatchLog_NotFound(t *testing.T) {
	entries, err := ParseDispatchLog("/nonexistent/path/dispatch-log.md")
	if err != nil {
		t.Fatalf("should not error on missing file: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries for missing file, got %d", len(entries))
	}
}
