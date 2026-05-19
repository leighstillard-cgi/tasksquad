package handlers

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type PMAgentProcess struct {
	PID            int
	SessionLogPath string
	PromptDir      string
	ScriptPath     string
	AlreadyRunning bool
}

func pmAgentPrompt() string {
	return `You are the TaskSquad PM agent for this worklog repository.

Run exactly one PM poll cycle, then stop. The shell wrapper waits the configured interval before each cycle and then starts the next cycle.

## Required Reading
- PM_INSTRUCTIONS.md is your source of truth for PM workflow.
- backlog.md is the product backlog.
- data/dispatch-log.md records active and completed dispatches.
- data/completions/ contains unprocessed completion reports.
- data/escalations/ contains open escalation records.

## Cycle Scope
1. Read PM_INSTRUCTIONS.md and inspect current backlog, dispatch log, completions, escalations, dispatch files, and session logs.
2. Process completion reports according to PM_INSTRUCTIONS.md when evidence satisfies acceptance criteria.
3. Detect active dispatches that have matching completion reports and close them only when validation passes.
4. Detect stalled or failed dispatches and write escalations instead of guessing.
5. Dispatch the next eligible ready story only when dependencies and same-repo concurrency rules allow it.

## Safety Rules
- Do not write application implementation code.
- Do not stage or commit pre-existing unrelated changes.
- Never use git add .; use explicit git add paths for files you changed during this PM cycle.
- If the worktree has unrelated dirty files, leave them untouched and continue only with PM tracking changes that do not conflict.
- Never force-push.
- If human approval or architectural judgment is needed, create or update an escalation and stop that action.

## Output
Summarize the cycle result, files changed, and any escalation or dispatch decision.`
}

func pmAgentShellScript() string {
	return `#!/bin/sh
set -u
printf '\n## PM Agent Loop\n\n'
printf 'Started at: %s\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
printf 'Poll interval: %s seconds\n' "$POLL_SECONDS"
while :; do
  printf '\nNext PM poll cycle after %s seconds.\n' "$POLL_SECONDS"
  sleep "$POLL_SECONDS"
  cycle="$(date -u +%Y%m%dT%H%M%SZ)"
  prompt_path="$PROMPT_DIR/pm-agent-$cycle-prompt.md"
  cat > "$prompt_path" <<'EOF_PROMPT'
` + pmAgentPrompt() + `
EOF_PROMPT
  printf '\n### PM Poll Cycle %s\n' "$cycle"
  printf 'Prompt: %s\n' "$prompt_path"
  printf 'Command: %s --ask-for-approval never exec -C %s --sandbox workspace-write -\n' "$CODEX_BIN" "$WORKLOG_PATH"
  "$CODEX_BIN" --ask-for-approval never exec -C "$WORKLOG_PATH" --sandbox workspace-write - < "$prompt_path"
  status=$?
  printf 'Cycle finished at: %s\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  printf 'Cycle exit status: %s\n' "$status"
done
`
}

func StartPMAgent(worklogPath string, pollInterval time.Duration, logger *slog.Logger) (*PMAgentProcess, error) {
	codexBin := os.Getenv("TASKSQUAD_CODEX_BIN")
	if codexBin == "" {
		codexBin = "codex"
	}
	resolvedCodex, err := exec.LookPath(codexBin)
	if err != nil {
		return nil, fmt.Errorf("pm agent command %q not found", codexBin)
	}
	if pollInterval <= 0 {
		pollInterval = 5 * time.Minute
	}

	sessionDir := filepath.Join(worklogPath, "data", "session-logs")
	pmDir := filepath.Join(sessionDir, ".pm-agent")
	promptDir := filepath.Join(pmDir, "prompts")
	scriptPath := filepath.Join(pmDir, "pm-agent.sh")
	if err := os.MkdirAll(promptDir, 0755); err != nil {
		return nil, err
	}

	if pid, ok := findProcessByArg(scriptPath); ok {
		return &PMAgentProcess{PID: pid, PromptDir: promptDir, ScriptPath: scriptPath, AlreadyRunning: true}, nil
	}

	if err := os.WriteFile(scriptPath, []byte(pmAgentShellScript()), 0755); err != nil {
		return nil, err
	}

	startedAt := time.Now().UTC()
	sessionLogPath := filepath.Join(sessionDir, fmt.Sprintf("pm-agent-%s-running.md", startedAt.Format("20060102T150405Z")))
	logFile, err := os.OpenFile(sessionLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	if _, err := fmt.Fprintf(logFile, `# TaskSquad PM Agent Session

Status: running
Role: pm-agent
Instructions: PM_INSTRUCTIONS.md
Script: %s
Prompt Directory: %s
Started: %s
`, relPathForPMAgentLog(worklogPath, scriptPath), relPathForPMAgentLog(worklogPath, promptDir), startedAt.Format(time.RFC3339)); err != nil {
		logFile.Close()
		return nil, err
	}

	pollSeconds := int(pollInterval.Seconds())
	if pollSeconds < 1 {
		pollSeconds = 1
	}

	cmd := exec.Command("/bin/sh", scriptPath)
	cmd.Dir = worklogPath
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.Env = append(os.Environ(),
		"CODEX_BIN="+resolvedCodex,
		"WORKLOG_PATH="+worklogPath,
		"PROMPT_DIR="+promptDir,
		fmt.Sprintf("POLL_SECONDS=%d", pollSeconds),
	)

	if err := cmd.Start(); err != nil {
		logFile.Close()
		return nil, err
	}

	process := &PMAgentProcess{
		PID:            cmd.Process.Pid,
		SessionLogPath: sessionLogPath,
		PromptDir:      promptDir,
		ScriptPath:     scriptPath,
	}

	go func() {
		err := cmd.Wait()
		if err != nil {
			logger.Warn("pm agent exited with error", "pid", process.PID, "error", err)
		} else {
			logger.Info("pm agent exited", "pid", process.PID)
		}
		logFile.Close()
	}()

	logger.Info("started pm agent", "pid", process.PID, "session_log", sessionLogPath, "poll_interval", pollInterval.String())
	return process, nil
}

func findProcessByArg(needle string) (int, bool) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return 0, false
	}
	self := os.Getpid()
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(entry.Name())
		if err != nil || pid == self {
			continue
		}
		cmdline, err := os.ReadFile(filepath.Join("/proc", entry.Name(), "cmdline"))
		if err != nil || len(cmdline) == 0 {
			continue
		}
		args := strings.ReplaceAll(string(cmdline), "\x00", " ")
		if strings.Contains(args, needle) {
			return pid, true
		}
	}
	return 0, false
}

func relPathForPMAgentLog(base, target string) string {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return target
	}
	return rel
}
