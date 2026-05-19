package parser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const dispatchLogHeader = `# Dispatch Log

PM agent maintains this file. Updated on every dispatch and completion.

| Story | Repo | Dispatched At | Status | Completion Report |
|---|---|---|---|---|
`

func ParseDispatchLog(path string) ([]DispatchEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []DispatchEntry{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var entries []DispatchEntry
	scanner := bufio.NewScanner(file)
	inTable := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "|") && strings.Contains(line, "Story") && strings.Contains(line, "Repo") {
			inTable = true
			continue
		}

		if strings.HasPrefix(line, "|---") || strings.HasPrefix(line, "| ---") {
			continue
		}

		if inTable && strings.HasPrefix(line, "|") {
			entry := parseDispatchLine(line)
			if entry.StoryID != "" {
				entries = append(entries, entry)
			}
		}
	}

	return entries, scanner.Err()
}

func AppendDispatchLogEntry(path string, req DispatchRequest) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	dispatchedAt := time.Now().UTC()
	if !req.DispatchedAt.IsZero() {
		dispatchedAt = req.DispatchedAt.UTC()
	}

	contentBytes, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	content := string(contentBytes)
	if strings.TrimSpace(content) == "" {
		content = dispatchLogHeader
	} else if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	row := fmt.Sprintf("| %s | %s | %s | dispatched | |\n", req.StoryID, req.Repo, dispatchedAt.Format(time.RFC3339))
	return os.WriteFile(path, []byte(content+row), 0644)
}

func parseDispatchLine(line string) DispatchEntry {
	parts := strings.Split(line, "|")
	if len(parts) < 6 {
		return DispatchEntry{}
	}

	var entry DispatchEntry
	idx := 0
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		switch idx {
		case 0:
			entry.StoryID = part
		case 1:
			entry.Repo = part
		case 2:
			if t, err := time.Parse(time.RFC3339, part); err == nil {
				entry.DispatchedAt = t
				entry.ElapsedTime = time.Since(t)
			}
		case 3:
			entry.Status = part
		case 4:
			entry.CompletionReport = part
		}
		idx++
	}
	return entry
}
