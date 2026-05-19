package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseSessionLogsTreatsTerminalSuccessAsPass(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "STORY-10.1-running.md")
	content := `# Worker Session

Status: running

error: transient read failed

Finished at: 2026-05-19T22:13:51Z
Exit status: 0
Completed successfully
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	logs, err := ParseSessionLogs(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(logs) != 1 {
		t.Fatalf("expected one log, got %d", len(logs))
	}
	if logs[0].Status != "pass" {
		t.Fatalf("expected pass status, got %q", logs[0].Status)
	}
}

func TestParseSessionLogsKeepsFailedWorkerAsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "STORY-10.1-running.md")
	content := `# Worker Session

Status: running

error: unexpected argument

Finished at: 2026-05-19T22:13:51Z
Exit status: 2
Worker failed
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	logs, err := ParseSessionLogs(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(logs) != 1 {
		t.Fatalf("expected one log, got %d", len(logs))
	}
	if logs[0].Status != "error" {
		t.Fatalf("expected error status, got %q", logs[0].Status)
	}
}
