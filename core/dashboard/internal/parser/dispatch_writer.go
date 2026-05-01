package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type DispatchRequest struct {
	StoryID      string
	Repo         string
	Description  string
	MaxRetries   int
	DispatchedBy string
}

func WriteDispatchFile(dir string, req DispatchRequest) (string, error) {
	if req.MaxRetries == 0 {
		req.MaxRetries = 5
	}
	if req.DispatchedBy == "" {
		req.DispatchedBy = "dashboard-manual"
	}

	filename := fmt.Sprintf("%s-dispatch.md", req.StoryID)
	path := filepath.Join(dir, filename)

	now := time.Now().UTC()
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
