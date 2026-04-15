package parser

import (
	"os"
	"path/filepath"
	"testing"
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
