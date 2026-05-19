package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type DispatchRequest struct {
	StoryID      string    `json:"story_id"`
	Repo         string    `json:"repo"`
	Description  string    `json:"description"`
	MaxRetries   int       `json:"max_retries"`
	DispatchedBy string    `json:"dispatched_by"`
	DispatchedAt time.Time `json:"-"`
}

func WriteDispatchFile(dir string, req DispatchRequest) (string, error) {
	if req.MaxRetries == 0 {
		req.MaxRetries = 5
	}
	if req.DispatchedBy == "" {
		req.DispatchedBy = "dashboard-manual"
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	now := time.Now().UTC()
	if !req.DispatchedAt.IsZero() {
		now = req.DispatchedAt.UTC()
	}

	filename := fmt.Sprintf("%s-dispatch.md", req.StoryID)
	path := filepath.Join(dir, filename)

	content := fmt.Sprintf(`---
story_id: %s
dispatched_at: %s
dispatched_by: %s
attempt: 1
max_retries: %d
---

# %s: Manual Dispatch

## Story Spec

**Status:** ready -> in-progress
**Repo:** %s

**Description:** %s

## Completion Output

Write completion report to: data/completions/%s-completion.md
Use template at: core/templates/story-completion.md
`,
		req.StoryID,
		now.Format(time.RFC3339),
		req.DispatchedBy,
		req.MaxRetries,
		req.StoryID,
		req.Repo,
		req.Description,
		req.StoryID,
	)

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", err
	}

	return path, nil
}
