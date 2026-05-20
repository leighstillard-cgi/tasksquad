package parser

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParseEscalations_ParsesDateOnlyTimestamp(t *testing.T) {
	content := `---
story_id: "STORY-10.1"
reason: "Needs follow-up"
created: "2026-05-20"
---
`

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "STORY-10.1-escalation.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	reports, err := ParseEscalations(dir)
	if err != nil {
		t.Fatalf("ParseEscalations failed: %v", err)
	}
	if len(reports) != 1 {
		t.Fatalf("expected 1 report, got %d", len(reports))
	}

	r := reports[0]
	if r.StoryID != "STORY-10.1" {
		t.Errorf("expected STORY-10.1, got %s", r.StoryID)
	}
	if r.Reason != "Needs follow-up" {
		t.Errorf("expected reason, got %q", r.Reason)
	}
	if !r.Timestamp.Equal(time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("expected date-only timestamp, got %s", r.Timestamp.Format(time.RFC3339))
	}
}
