package parser

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParseCompletions(t *testing.T) {
	content := `---
title: "Completion: STORY-03.4"
id: "COMPLETION-STORY-03.4"
status: draft
page_class: agent-write
page_type: story
tags: [completion]
created: "2026-04-15T00:00:00Z"
last_updated: "2026-04-15T00:00:00Z"
parent_epic: "EPIC-03"
phase: "03"
repos: ["tasksquad"]
---

# Completion: STORY-03.4

Summary here.
`

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "STORY-03.4-completion.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	reports, err := ParseCompletions(dir)
	if err != nil {
		t.Fatalf("ParseCompletions failed: %v", err)
	}

	if len(reports) != 1 {
		t.Fatalf("expected 1 report, got %d", len(reports))
	}

	r := reports[0]
	if r.StoryID != "STORY-03.4" {
		t.Errorf("expected STORY-03.4, got %s", r.StoryID)
	}
	if r.Status != "draft" {
		t.Errorf("expected draft, got %s", r.Status)
	}
	if r.ParentEpic != "EPIC-03" {
		t.Errorf("expected EPIC-03, got %s", r.ParentEpic)
	}
	if len(r.Repos) != 1 || r.Repos[0] != "tasksquad" {
		t.Errorf("expected [tasksquad], got %v", r.Repos)
	}
	if r.CompletedAt.IsZero() {
		t.Errorf("expected completion timestamp fallback, got zero time")
	}
}

func TestParseCompletions_ParsesDateOnlyFrontmatterAndBodyTimestamp(t *testing.T) {
	content := `---
title: "Completion: STORY-10.3"
id: "COMPLETION-STORY-10.3"
status: complete
created: "2026-05-20"
last_updated: "2026-05-21"
parent_epic: "EPIC-10"
phase: "10"
repos: ["tasksquad"]
---

# Completion: STORY-10.3

**Timestamp:** 2026-05-19T23:51:39Z
`

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "STORY-10.3-completion.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	reports, err := ParseCompletions(dir)
	if err != nil {
		t.Fatalf("ParseCompletions failed: %v", err)
	}
	if len(reports) != 1 {
		t.Fatalf("expected 1 report, got %d", len(reports))
	}

	r := reports[0]
	if !r.Created.Equal(time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("expected date-only created timestamp, got %s", r.Created.Format(time.RFC3339))
	}
	if !r.LastUpdated.Equal(time.Date(2026, 5, 21, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("expected date-only last_updated timestamp, got %s", r.LastUpdated.Format(time.RFC3339))
	}
	if !r.CompletedAt.Equal(time.Date(2026, 5, 19, 23, 51, 39, 0, time.UTC)) {
		t.Errorf("expected body completion timestamp, got %s", r.CompletedAt.Format(time.RFC3339))
	}
}

func TestParseCompletions_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	reports, err := ParseCompletions(dir)
	if err != nil {
		t.Fatalf("should not error on empty dir: %v", err)
	}
	if len(reports) != 0 {
		t.Errorf("expected 0 reports, got %d", len(reports))
	}
}

func TestParseCompletions_NotFound(t *testing.T) {
	reports, err := ParseCompletions("/nonexistent/completions")
	if err != nil {
		t.Fatalf("should not error on missing dir: %v", err)
	}
	if len(reports) != 0 {
		t.Errorf("expected 0 reports, got %d", len(reports))
	}
}
