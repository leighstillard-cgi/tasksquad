package parser

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func ParseDispatches(dir string) ([]DispatchFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []DispatchFile{}, nil
		}
		return nil, err
	}

	var dispatches []DispatchFile
	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == ".gitkeep" {
			continue
		}
		if !strings.HasSuffix(entry.Name(), "-dispatch.md") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		dispatch, err := parseDispatchFile(path)
		if err != nil {
			continue
		}
		dispatches = append(dispatches, dispatch)
	}

	return dispatches, nil
}

func parseDispatchFile(path string) (DispatchFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return DispatchFile{}, err
	}
	defer file.Close()

	dispatch := DispatchFile{FilePath: path}
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
			if strings.HasPrefix(trimmed, "story_id:") {
				dispatch.StoryID = parseYAMLString(trimmed[9:])
			} else if strings.HasPrefix(trimmed, "dispatched_at:") {
				if t, err := time.Parse(time.RFC3339, parseYAMLString(trimmed[14:])); err == nil {
					dispatch.DispatchedAt = t
				}
			} else if strings.HasPrefix(trimmed, "dispatched_by:") {
				dispatch.DispatchedBy = parseYAMLString(trimmed[14:])
			} else if strings.HasPrefix(trimmed, "attempt:") {
				if n, err := strconv.Atoi(parseYAMLString(trimmed[8:])); err == nil {
					dispatch.Attempt = n
				}
			} else if strings.HasPrefix(trimmed, "max_retries:") {
				if n, err := strconv.Atoi(parseYAMLString(trimmed[12:])); err == nil {
					dispatch.MaxRetries = n
				}
			}
		}
	}

	return dispatch, scanner.Err()
}
