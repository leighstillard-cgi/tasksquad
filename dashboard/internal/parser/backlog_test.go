package parser

import (
	"os"
	"path/filepath"
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

func TestParseBacklog_NotFound(t *testing.T) {
	overview, err := ParseBacklog("/nonexistent/backlog.md")
	if err != nil {
		t.Fatalf("should not error on missing file: %v", err)
	}
	if overview == nil {
		t.Error("expected empty overview, got nil")
	}
}
