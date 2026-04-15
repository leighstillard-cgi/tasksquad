package parser

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

var storyHeaderRegex = regexp.MustCompile(`^###\s+(STORY-[\d.]+)\s+[·•-]\s*(.*)$`)

func ParseBacklog(path string) (*BacklogOverview, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &BacklogOverview{}, nil
		}
		return nil, err
	}
	defer file.Close()

	overview := &BacklogOverview{}
	scanner := bufio.NewScanner(file)

	var currentStory *BacklogStory
	var collectingDescription bool

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if matches := storyHeaderRegex.FindStringSubmatch(trimmed); len(matches) >= 3 {
			if currentStory != nil {
				addStoryToOverview(overview, currentStory)
			}
			currentStory = &BacklogStory{
				StoryID: matches[1],
				Title:   strings.TrimSpace(matches[2]),
			}
			collectingDescription = false
			continue
		}

		if currentStory != nil {
			if strings.HasPrefix(trimmed, "**Status:**") {
				currentStory.Status = parseFieldValue(trimmed, "**Status:**")
			} else if strings.HasPrefix(trimmed, "**Repo:**") {
				currentStory.Repo = parseFieldValue(trimmed, "**Repo:**")
			} else if strings.HasPrefix(trimmed, "**Depends on:**") {
				currentStory.DependsOn = parseFieldValue(trimmed, "**Depends on:**")
			} else if strings.HasPrefix(trimmed, "**Priority:**") {
				currentStory.Priority = parseFieldValue(trimmed, "**Priority:**")
			} else if strings.HasPrefix(trimmed, "**Description:**") {
				collectingDescription = true
				desc := parseFieldValue(trimmed, "**Description:**")
				if desc != "" {
					currentStory.Description = desc
				}
			} else if collectingDescription && trimmed != "" && !strings.HasPrefix(trimmed, "**") && !strings.HasPrefix(trimmed, "---") && !strings.HasPrefix(trimmed, "##") {
				if currentStory.Description != "" {
					currentStory.Description += " "
				}
				currentStory.Description += trimmed
			} else if strings.HasPrefix(trimmed, "**") || strings.HasPrefix(trimmed, "---") || strings.HasPrefix(trimmed, "##") {
				collectingDescription = false
			}
		}

		if strings.HasPrefix(trimmed, "## ") && currentStory != nil {
			addStoryToOverview(overview, currentStory)
			currentStory = nil
		}
	}

	if currentStory != nil {
		addStoryToOverview(overview, currentStory)
	}

	return overview, scanner.Err()
}

func parseFieldValue(line, prefix string) string {
	value := strings.TrimPrefix(line, prefix)
	return strings.TrimSpace(value)
}

func addStoryToOverview(overview *BacklogOverview, story *BacklogStory) {
	status := strings.ToLower(story.Status)
	switch status {
	case "done":
		overview.Done = append(overview.Done, *story)
	case "in-progress":
		overview.InProgress = append(overview.InProgress, *story)
	case "ready":
		overview.Ready = append(overview.Ready, *story)
	case "blocked":
		overview.Blocked = append(overview.Blocked, *story)
	case "cancelled":
		overview.Cancelled = append(overview.Cancelled, *story)
	default:
		if story.Status != "" {
			overview.Ready = append(overview.Ready, *story)
		}
	}
}
