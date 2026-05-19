package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseBacklog(t *testing.T) {
	content := `# TaskSquad — Product Backlog

## Infrastructure

### STORY-00.1 · Dashboard

**Status:** ready
**Repo:** tasksquad
**Priority:** High

**Description:** Build monitoring dashboard.

---

### STORY-03.1 · Wiki Setup

**Status:** done
**Repo:** tasksquad

---

### STORY-03.2 · Core Framework

**Status:** in-progress
**Repo:** tasksquad

---

### STORY-03.3 · Skills

**Status:** blocked
**Repo:** tasksquad

---
`

	dir := t.TempDir()
	path := filepath.Join(dir, "backlog.md")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	overview, err := ParseBacklog(path)
	if err != nil {
		t.Fatalf("ParseBacklog failed: %v", err)
	}

	if len(overview.Ready) != 1 {
		t.Errorf("expected 1 ready story, got %d", len(overview.Ready))
	}
	if len(overview.Done) != 1 {
		t.Errorf("expected 1 done story, got %d", len(overview.Done))
	}
	if len(overview.InProgress) != 1 {
		t.Errorf("expected 1 in-progress story, got %d", len(overview.InProgress))
	}
	if len(overview.Blocked) != 1 {
		t.Errorf("expected 1 blocked story, got %d", len(overview.Blocked))
	}

	if overview.Ready[0].StoryID != "STORY-00.1" {
		t.Errorf("expected STORY-00.1, got %s", overview.Ready[0].StoryID)
	}
	if overview.Ready[0].Title != "Dashboard" {
		t.Errorf("expected Dashboard, got %s", overview.Ready[0].Title)
	}
}

func TestUpdateBacklogStoryStatus(t *testing.T) {
	content := `# TaskSquad Product Backlog

## Infrastructure

### STORY-00.1 - Dashboard
**Status:** ready
**Repo:** tasksquad

### STORY-00.2 - Live Updates
**Status:** ready
**Repo:** tasksquad
`

	dir := t.TempDir()
	path := filepath.Join(dir, "backlog.md")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	if err := UpdateBacklogStoryStatus(path, "STORY-00.1", "in-progress", "ready"); err != nil {
		t.Fatalf("UpdateBacklogStoryStatus failed: %v", err)
	}

	updated, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(updated), "### STORY-00.1 - Dashboard\n**Status:** in-progress") {
		t.Fatalf("target story status not updated: %s", string(updated))
	}
	if !strings.Contains(string(updated), "### STORY-00.2 - Live Updates\n**Status:** ready") {
		t.Fatalf("non-target story status changed: %s", string(updated))
	}
}

func TestParseBacklog_NotFound(t *testing.T) {
	overview, err := ParseBacklog("/nonexistent/backlog.md")
	if err != nil {
		t.Fatalf("should not error on missing file: %v", err)
	}
	if overview == nil {
		t.Error("expected empty overview, got nil")
	}
}
