package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteDispatchFile(t *testing.T) {
	dir := t.TempDir()

	req := DispatchRequest{
		StoryID:     "STORY-99.1",
		Repo:        "test-repo",
		Description: "Test description",
		MaxRetries:  3,
	}

	path, err := WriteDispatchFile(dir, req)
	if err != nil {
		t.Fatalf("WriteDispatchFile failed: %v", err)
	}

	expectedPath := filepath.Join(dir, "STORY-99.1-dispatch.md")
	if path != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, path)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read dispatch file: %v", err)
	}

	contentStr := string(content)

	if !strings.Contains(contentStr, "story_id: STORY-99.1") {
		t.Error("missing story_id in dispatch file")
	}
	if !strings.Contains(contentStr, "max_retries: 3") {
		t.Error("missing max_retries in dispatch file")
	}
	if !strings.Contains(contentStr, "dispatched_by: dashboard-manual") {
		t.Error("missing dispatched_by in dispatch file")
	}
	if !strings.Contains(contentStr, "**Repo:** test-repo") {
		t.Error("missing repo in dispatch file")
	}
	if !strings.Contains(contentStr, "Test description") {
		t.Error("missing description in dispatch file")
	}
}

func TestWriteDispatchFile_DefaultValues(t *testing.T) {
	dir := t.TempDir()

	req := DispatchRequest{
		StoryID: "STORY-99.2",
		Repo:    "test-repo",
	}

	path, err := WriteDispatchFile(dir, req)
	if err != nil {
		t.Fatalf("WriteDispatchFile failed: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read dispatch file: %v", err)
	}

	contentStr := string(content)

	if !strings.Contains(contentStr, "max_retries: 5") {
		t.Error("default max_retries should be 5")
	}
}

func TestWriteDispatchFile_Roundtrip(t *testing.T) {
	dir := t.TempDir()

	req := DispatchRequest{
		StoryID:     "STORY-99.3",
		Repo:        "roundtrip-repo",
		Description: "Roundtrip test",
	}

	path, err := WriteDispatchFile(dir, req)
	if err != nil {
		t.Fatalf("WriteDispatchFile failed: %v", err)
	}

	dispatches, err := ParseDispatches(dir)
	if err != nil {
		t.Fatalf("ParseDispatches failed: %v", err)
	}

	if len(dispatches) != 1 {
		t.Fatalf("expected 1 dispatch, got %d", len(dispatches))
	}

	d := dispatches[0]
	if d.StoryID != "STORY-99.3" {
		t.Errorf("expected STORY-99.3, got %s", d.StoryID)
	}
	if d.FilePath != path {
		t.Errorf("expected path %s, got %s", path, d.FilePath)
	}
	if d.MaxRetries != 5 {
		t.Errorf("expected 5 retries, got %d", d.MaxRetries)
	}
}
