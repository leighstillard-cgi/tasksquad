package handlers

import (
	"bytes"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var (
	// safeFilenameRegex allows only alphanumeric, dash, underscore, and dot
	safeFilenameRegex = regexp.MustCompile(`^[a-zA-Z0-9._-]+\.md$`)

	md = goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithXHTML()),
	)
)

// APICompletion serves a rendered completion report as HTML.
func (h *Handler) APICompletion(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/api/completion/")
	h.serveMarkdown(w, r, filename, h.dataDir("completions"))
}

// APISessionLog serves a rendered session log as HTML.
func (h *Handler) APISessionLog(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/api/session-log/")
	h.serveMarkdown(w, r, filename, h.dataDir("session-logs"))
}

func (h *Handler) serveMarkdown(w http.ResponseWriter, r *http.Request, filename, baseDir string) {
	// Security: validate filename
	if !isValidFilename(filename) {
		h.logger.Warn("invalid filename requested", "filename", filename)
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	// Build the full path and verify it stays within baseDir
	fullPath := filepath.Join(baseDir, filename)
	cleanPath := filepath.Clean(fullPath)

	if !strings.HasPrefix(cleanPath, filepath.Clean(baseDir)) {
		h.logger.Warn("path traversal attempt", "filename", filename, "resolved", cleanPath)
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	// Read the file
	content, err := os.ReadFile(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		h.logger.Error("failed to read file", "path", cleanPath, "error", err)
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Convert markdown to HTML
	var buf bytes.Buffer
	if err := md.Convert(content, &buf); err != nil {
		h.logger.Error("failed to convert markdown", "path", cleanPath, "error", err)
		http.Error(w, "Failed to render content", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}

func isValidFilename(filename string) bool {
	if filename == "" {
		return false
	}

	// No path separators allowed
	if strings.ContainsAny(filename, "/\\") {
		return false
	}

	// No parent directory references
	if strings.Contains(filename, "..") {
		return false
	}

	// Must match our safe pattern
	return safeFilenameRegex.MatchString(filename)
}
