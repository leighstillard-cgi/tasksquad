package parser

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ParseCompletions(dir string) ([]CompletionReport, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []CompletionReport{}, nil
		}
		return nil, err
	}

	var reports []CompletionReport
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), "-completion.md") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		report, err := parseCompletionFile(path)
		if err != nil {
			continue
		}
		reports = append(reports, report)
	}

	return reports, nil
}

func parseCompletionFile(path string) (CompletionReport, error) {
	file, err := os.Open(path)
	if err != nil {
		return CompletionReport{}, err
	}
	defer file.Close()

	report := CompletionReport{FilePath: path}
	scanner := bufio.NewScanner(file)
	inFrontmatter := false
	frontmatterCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if trimmed == "---" {
			frontmatterCount++
			if frontmatterCount == 1 {
				inFrontmatter = true
				continue
			} else if frontmatterCount == 2 {
				inFrontmatter = false
				break
			}
		}

		if inFrontmatter {
			if strings.HasPrefix(trimmed, "title:") {
				report.Title = parseYAMLString(trimmed[6:])
			} else if strings.HasPrefix(trimmed, "id:") {
				id := parseYAMLString(trimmed[3:])
				if strings.HasPrefix(id, "COMPLETION-") {
					report.StoryID = strings.TrimPrefix(id, "COMPLETION-")
				}
			} else if strings.HasPrefix(trimmed, "status:") {
				report.Status = parseYAMLString(trimmed[7:])
			} else if strings.HasPrefix(trimmed, "created:") {
				if t, err := time.Parse(time.RFC3339, parseYAMLString(trimmed[8:])); err == nil {
					report.Created = t
				}
			} else if strings.HasPrefix(trimmed, "last_updated:") {
				if t, err := time.Parse(time.RFC3339, parseYAMLString(trimmed[13:])); err == nil {
					report.LastUpdated = t
				}
			} else if strings.HasPrefix(trimmed, "parent_epic:") {
				report.ParentEpic = parseYAMLString(trimmed[12:])
			} else if strings.HasPrefix(trimmed, "phase:") {
				report.Phase = parseYAMLString(trimmed[6:])
			} else if strings.HasPrefix(trimmed, "repos:") {
				report.Repos = parseYAMLArray(trimmed[6:])
			}
		}
	}

	return report, scanner.Err()
}

func parseYAMLString(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "\"'")
	return s
}

func parseYAMLArray(s string) []string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "[]")
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		p = strings.Trim(p, "\"'")
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
