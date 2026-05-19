package handlers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/tasksquad/dashboard/internal/parser"
)

type workerLaunchRequest struct {
	Story        parser.BacklogStory
	Dispatch     parser.DispatchRequest
	DispatchPath string
}

type workerProcess struct {
	PID            int
	SessionLogPath string
	PromptPath     string
	ScriptPath     string
}

func workerShellScript() string {
	return `#!/bin/sh
set -u
printf '\n## Worker Output\n\n'
printf 'Started at: %s\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
printf 'Command: %s --ask-for-approval never exec -C %s --sandbox workspace-write -\n' "$CODEX_BIN" "$WORKLOG_PATH"
"$CODEX_BIN" --ask-for-approval never exec -C "$WORKLOG_PATH" --sandbox workspace-write - < "$PROMPT_PATH"
status=$?
printf '\nFinished at: %s\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
printf 'Exit status: %s\n' "$status"
if [ "$status" -eq 0 ]; then
  printf 'Completed successfully\n'
else
  printf 'Worker failed\n'
fi
exit "$status"
`
}

func (h *Handler) launchLocalWorker(req workerLaunchRequest) (*workerProcess, error) {
	codexBin := os.Getenv("TASKSQUAD_CODEX_BIN")
	if codexBin == "" {
		codexBin = "codex"
	}
	resolvedCodex, err := exec.LookPath(codexBin)
	if err != nil {
		return nil, fmt.Errorf("worker command %q not found", codexBin)
	}

	sessionDir := h.writableDataDir("session-logs")
	promptDir := filepath.Join(sessionDir, ".worker-prompts")
	scriptDir := filepath.Join(sessionDir, ".worker-scripts")
	if err := os.MkdirAll(promptDir, 0755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(scriptDir, 0755); err != nil {
		return nil, err
	}

	dispatchedAt := req.Dispatch.DispatchedAt.UTC()
	if dispatchedAt.IsZero() {
		dispatchedAt = time.Now().UTC()
	}
	baseName := fmt.Sprintf("%s-%s", req.Dispatch.StoryID, dispatchedAt.Format("20060102T150405Z"))
	promptPath := filepath.Join(promptDir, baseName+"-prompt.md")
	scriptPath := filepath.Join(scriptDir, baseName+"-worker.sh")
	sessionLogPath := filepath.Join(sessionDir, baseName+"-running.md")

	prompt, err := h.buildWorkerPrompt(req)
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(promptPath, []byte(prompt), 0644); err != nil {
		return nil, err
	}

	script := workerShellScript()
	if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
		return nil, err
	}

	logFile, err := os.OpenFile(sessionLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	if _, err := fmt.Fprintf(logFile, `# %s Worker Session

Status: running
Story: %s
Title: %s
Repository: %s
Dispatch: %s
Prompt: %s
Started: %s

`, req.Dispatch.StoryID, req.Dispatch.StoryID, req.Story.Title, req.Dispatch.Repo, relPathForLog(h.worklogPath, req.DispatchPath), relPathForLog(h.worklogPath, promptPath), time.Now().UTC().Format(time.RFC3339)); err != nil {
		logFile.Close()
		return nil, err
	}

	cmd := exec.Command("/bin/sh", scriptPath)
	cmd.Dir = h.worklogPath
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.Env = append(os.Environ(),
		"CODEX_BIN="+resolvedCodex,
		"WORKLOG_PATH="+h.worklogPath,
		"PROMPT_PATH="+promptPath,
	)

	if err := cmd.Start(); err != nil {
		logFile.Close()
		return nil, err
	}

	process := &workerProcess{
		PID:            cmd.Process.Pid,
		SessionLogPath: sessionLogPath,
		PromptPath:     promptPath,
		ScriptPath:     scriptPath,
	}

	go func() {
		err := cmd.Wait()
		if err != nil {
			h.logger.Warn("worker exited with error", "story_id", req.Dispatch.StoryID, "pid", process.PID, "error", err)
		} else {
			h.logger.Info("worker exited", "story_id", req.Dispatch.StoryID, "pid", process.PID)
		}
		logFile.Close()
		h.refresh()
		h.broadcastCurrentData()
	}()

	h.logger.Info("started worker", "story_id", req.Dispatch.StoryID, "pid", process.PID, "session_log", sessionLogPath)
	return process, nil
}

func (h *Handler) writableDataDir(name string) string {
	dataRoot := filepath.Join(h.worklogPath, "data")
	if info, err := os.Stat(dataRoot); err == nil && info.IsDir() {
		return filepath.Join(dataRoot, name)
	}
	return h.dataDir(name)
}

func (h *Handler) buildWorkerPrompt(req workerLaunchRequest) (string, error) {
	dispatchContent, err := os.ReadFile(req.DispatchPath)
	if err != nil {
		return "", err
	}

	storySpecPath := h.findStorySpecPath(req.Dispatch.StoryID)
	storySpecSection := "No separate story spec file was found; use the dispatch file and backlog entry."
	if storySpecPath != "" {
		content, err := os.ReadFile(storySpecPath)
		if err != nil {
			return "", err
		}
		storySpecSection = fmt.Sprintf("Path: %s\n\n```markdown\n%s\n```", relPathForLog(h.worklogPath, storySpecPath), strings.TrimSpace(string(content)))
	}

	return fmt.Sprintf(`You are a TaskSquad worker for %s: %s.

## Your Task
Execute the dispatched story locally in this repository.

## Dispatch File
Path: %s

`+"```markdown"+`
%s
`+"```"+`

## Story Spec
%s

## Required Reading
- PM_INSTRUCTIONS.md explains TaskSquad dispatch and completion flow.
- AGENTS.md contains Codex-facing repo instructions when present.
- CLAUDE.md contains shared TaskSquad standards; treat Claude-only hooks, plugins, and slash commands as unavailable unless this environment explicitly provides them.
- Use core/templates/story-completion.md for the completion report.

## Completion
When done, write the completion report to:
  data/completions/%s-completion.md

Include evidence for every acceptance criterion.

## Rules
- Do not modify backlog.md or data/dispatch-log.md.
- Do not dispatch other stories.
- If you hit an architectural decision or blocker, record it in the completion report or create an escalation under data/escalations/.
- Commit and push work changes before writing the completion report when that is appropriate for the story.
`, req.Dispatch.StoryID, req.Story.Title, relPathForLog(h.worklogPath, req.DispatchPath), strings.TrimSpace(string(dispatchContent)), storySpecSection, req.Dispatch.StoryID), nil
}

func (h *Handler) findStorySpecPath(storyID string) string {
	matches, err := filepath.Glob(filepath.Join(h.worklogPath, "data", "story-specs", storyID+"-*.md"))
	if err != nil || len(matches) == 0 {
		return ""
	}
	sort.Strings(matches)
	return matches[0]
}

func relPathForLog(base, target string) string {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return target
	}
	return rel
}
