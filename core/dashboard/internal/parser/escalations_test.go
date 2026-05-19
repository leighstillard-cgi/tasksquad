package parser

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParseCompletionEscalations(t *testing.T) {
	content := `---
title: "Completion: STORY-10.1"
id: "COMPLETION-STORY-10.1"
status: draft
created: "2026-05-20T01:00:00Z"
last_updated: "2026-05-20T01:02:03Z"
parent_epic: "EPIC-10"
phase: "test"
repos: ["tasksquad"]
---

# Completion: STORY-10.1

## Summary

Done.

## Architectural Escalations

- First architecture decision needs review.
- Second escalation starts here
  and continues on the next line.

## Files Changed

- ignored.md
`

	dir := t.TempDir()
	path := filepath.Join(dir, "STORY-10.1-completion.md")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	reports, err := ParseCompletionEscalations(dir)
	if err != nil {
		t.Fatalf("ParseCompletionEscalations failed: %v", err)
	}
	if len(reports) != 2 {
		t.Fatalf("expected 2 escalation reports, got %d: %#v", len(reports), reports)
	}

	wantTimestamp := time.Date(2026, 5, 20, 1, 2, 3, 0, time.UTC)
	if reports[0].StoryID != "STORY-10.1" || reports[0].Reason != "First architecture decision needs review." {
		t.Fatalf("unexpected first escalation: %#v", reports[0])
	}
	if reports[1].Reason != "Second escalation starts here and continues on the next line." {
		t.Fatalf("unexpected continued escalation reason: %#v", reports[1])
	}
	if !reports[0].Timestamp.Equal(wantTimestamp) || !reports[1].Timestamp.Equal(wantTimestamp) {
		t.Fatalf("expected last_updated timestamp %s, got %#v", wantTimestamp, reports)
	}
	if reports[0].FilePath != path || reports[1].FilePath != path {
		t.Fatalf("expected completion file paths, got %#v", reports)
	}
}
