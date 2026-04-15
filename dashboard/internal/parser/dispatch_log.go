package parser

import (
	"bufio"
	"os"
	"strings"
	"time"
)

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
