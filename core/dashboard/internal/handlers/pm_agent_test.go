package handlers

import (
	"strings"
	"testing"
)

func TestPMAgentShellScriptUsesSupportedCodexCommand(t *testing.T) {
	script := pmAgentShellScript()

	invalid := `"$CODEX_BIN" exec -C "$WORKLOG_PATH" --sandbox workspace-write --ask-for-approval never - < "$prompt_path"`
	if strings.Contains(script, invalid) {
		t.Fatalf("pm agent script uses unsupported codex exec flag order: %s", invalid)
	}

	expected := `"$CODEX_BIN" --ask-for-approval never exec -C "$WORKLOG_PATH" --sandbox workspace-write - < "$prompt_path"`
	if !strings.Contains(script, expected) {
		t.Fatalf("pm agent script missing supported codex launch command %q in:\n%s", expected, script)
	}

	if !strings.Contains(script, "while :; do") {
		t.Fatalf("pm agent script should run as a persistent poll loop:\n%s", script)
	}
	if !strings.Contains(script, "sleep \"$POLL_SECONDS\"") {
		t.Fatalf("pm agent script should wait before each PM cycle:\n%s", script)
	}
}

func TestPMAgentPromptUsesPMInstructionsAndProtectsDirtyWork(t *testing.T) {
	prompt := pmAgentPrompt()
	for _, required := range []string{
		"PM_INSTRUCTIONS.md is your source of truth",
		"Run exactly one PM poll cycle",
		"Do not stage or commit pre-existing unrelated changes",
		"Never use git add .",
	} {
		if !strings.Contains(prompt, required) {
			t.Fatalf("pm agent prompt missing %q in:\n%s", required, prompt)
		}
	}
}
