package parser

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
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
				continue
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
				if t, ok := parseYAMLTime(trimmed[8:]); ok {
					report.Created = t
				}
			} else if strings.HasPrefix(trimmed, "last_updated:") {
				if t, ok := parseYAMLTime(trimmed[13:]); ok {
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

		if !inFrontmatter && frontmatterCount >= 2 && report.CompletedAt.IsZero() {
			if value, ok := parseMarkdownField(trimmed, "Timestamp"); ok {
				if t, ok := parseYAMLTime(value); ok {
					report.CompletedAt = t
				}
			}
		}
	}

	if report.CompletedAt.IsZero() {
		report.CompletedAt = report.LastUpdated
	}
	if report.CompletedAt.IsZero() {
		report.CompletedAt = report.Created
	}

	return report, scanner.Err()
}

func parseMarkdownField(line, name string) (string, bool) {
	boldPrefix := "**" + name + ":**"
	if strings.HasPrefix(line, boldPrefix) {
		return strings.TrimSpace(strings.TrimPrefix(line, boldPrefix)), true
	}

	plainPrefix := name + ":"
	if strings.HasPrefix(line, plainPrefix) {
		return strings.TrimSpace(strings.TrimPrefix(line, plainPrefix)), true
	}

	return "", false
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
