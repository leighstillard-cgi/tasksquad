package parser

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ParseEscalations(dir string) ([]EscalationReport, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []EscalationReport{}, nil
		}
		return nil, err
	}

	var reports []EscalationReport
	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == ".gitkeep" {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		report, err := parseEscalationFile(path)
		if err != nil {
			continue
		}
		reports = append(reports, report)
	}

	return reports, nil
}

func parseEscalationFile(path string) (EscalationReport, error) {
	file, err := os.Open(path)
	if err != nil {
		return EscalationReport{}, err
	}
	defer file.Close()

	report := EscalationReport{FilePath: path}
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
			if strings.HasPrefix(trimmed, "story_id:") {
				report.StoryID = parseYAMLString(trimmed[9:])
			} else if strings.HasPrefix(trimmed, "reason:") {
				report.Reason = parseYAMLString(trimmed[7:])
			} else if strings.HasPrefix(trimmed, "timestamp:") || strings.HasPrefix(trimmed, "created:") {
				idx := strings.Index(trimmed, ":")
				if t, err := time.Parse(time.RFC3339, parseYAMLString(trimmed[idx+1:])); err == nil {
					report.Timestamp = t
				}
			}
		}

		if !inFrontmatter && frontmatterCount >= 2 {
			if strings.HasPrefix(trimmed, "# ") && report.StoryID == "" {
				report.StoryID = extractStoryIDFromTitle(trimmed)
			}
			if strings.HasPrefix(trimmed, "**Reason:**") || strings.HasPrefix(trimmed, "## Reason") {
				report.Reason = strings.TrimPrefix(trimmed, "**Reason:**")
				report.Reason = strings.TrimSpace(report.Reason)
			}
		}
	}

	if report.StoryID == "" {
		base := filepath.Base(path)
		if strings.Contains(base, "STORY-") {
			idx := strings.Index(base, "STORY-")
			end := strings.Index(base[idx:], "-escalation")
			if end == -1 {
				end = strings.Index(base[idx:], ".")
			}
			if end > 0 {
				report.StoryID = base[idx : idx+end]
			}
		}
	}

	return report, scanner.Err()
}

func extractStoryIDFromTitle(title string) string {
	title = strings.TrimPrefix(title, "# ")
	if strings.Contains(title, "STORY-") {
		idx := strings.Index(title, "STORY-")
		end := idx
		for end < len(title) && (title[end] != ' ' && title[end] != ':') {
			end++
		}
		return title[idx:end]
	}
	return ""
}
