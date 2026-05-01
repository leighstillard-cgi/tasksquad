package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAPICompletion_Success(t *testing.T) {
	// Create temp directory with a test completion file
	tmpDir := t.TempDir()
	completionsDir := filepath.Join(tmpDir, "completions")
	if err := os.MkdirAll(completionsDir, 0755); err != nil {
		t.Fatal(err)
	}

	testContent := `# Test Completion

This is a **test** completion report.

- Item 1
- Item 2
`
	testFile := filepath.Join(completionsDir, "STORY-01-completion.md")
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatal(err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	tmpl := template.Must(template.New("index.html").Parse(""))

	h := New(tmpDir, tmpl, logger)

	req := httptest.NewRequest("GET", "/api/completion/STORY-01-completion.md", nil)
	rec := httptest.NewRecorder()

	h.APICompletion(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "<h1") {
		t.Error("expected HTML heading in response")
	}
	if !strings.Contains(body, "<strong>test</strong>") {
		t.Error("expected bold text rendered")
	}
}

func TestAPICompletion_FileNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	completionsDir := filepath.Join(tmpDir, "completions")
	if err := os.MkdirAll(completionsDir, 0755); err != nil {
		t.Fatal(err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	tmpl := template.Must(template.New("index.html").Parse(""))

	h := New(tmpDir, tmpl, logger)

	req := httptest.NewRequest("GET", "/api/completion/nonexistent.md", nil)
	rec := httptest.NewRecorder()

	h.APICompletion(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}
}

func TestAPICompletion_PathTraversal(t *testing.T) {
	tmpDir := t.TempDir()
	completionsDir := filepath.Join(tmpDir, "completions")
	if err := os.MkdirAll(completionsDir, 0755); err != nil {
		t.Fatal(err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	tmpl := template.Must(template.New("index.html").Parse(""))

	h := New(tmpDir, tmpl, logger)

	testCases := []struct {
		name     string
		filename string
	}{
		{"parent directory", "../secret.md"},
		{"absolute path", "/etc/passwd"},
		{"double dots", "..%2F..%2Fetc%2Fpasswd"},
		{"backslash", "..\\secret.md"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/completion/"+tc.filename, nil)
			rec := httptest.NewRecorder()

			h.APICompletion(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("expected status 400 for %q, got %d", tc.filename, rec.Code)
			}
		})
	}
}

func TestAPISessionLog_Success(t *testing.T) {
	tmpDir := t.TempDir()
	logsDir := filepath.Join(tmpDir, "session-logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		t.Fatal(err)
	}

	testContent := `# Session Log

Session started at 2026-04-15.

## Summary

Everything worked.
`
	testFile := filepath.Join(logsDir, "session-2026-04-15.md")
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatal(err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	tmpl := template.Must(template.New("index.html").Parse(""))

	h := New(tmpDir, tmpl, logger)

	req := httptest.NewRequest("GET", "/api/session-log/session-2026-04-15.md", nil)
	rec := httptest.NewRecorder()

	h.APISessionLog(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Session Log") {
		t.Error("expected content in response")
	}
}

func TestIsValidFilename(t *testing.T) {
	testCases := []struct {
		filename string
		valid    bool
	}{
		{"STORY-01-completion.md", true},
		{"session-log.md", true},
		{"test_file-123.md", true},
		{"../secret.md", false},
		{"foo/bar.md", false},
		{"foo\\bar.md", false},
		{"..secret.md", false},
		{"", false},
		{"noextension", false},
		{"file.txt", false},
		{"file.md.exe", false},
	}

	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			result := isValidFilename(tc.filename)
			if result != tc.valid {
				t.Errorf("isValidFilename(%q) = %v, want %v", tc.filename, result, tc.valid)
			}
		})
	}
}
