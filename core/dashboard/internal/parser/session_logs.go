package parser

import (
	"os"
	"path/filepath"
	"strings"
)

func ParseSessionLogs(dir string) ([]SessionLog, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []SessionLog{}, nil
		}
		return nil, err
	}

	var logs []SessionLog
	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == ".gitkeep" {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		status := inferSessionStatus(entry.Name(), path)

		logs = append(logs, SessionLog{
			FileName:  entry.Name(),
			FilePath:  path,
			Status:    status,
			Timestamp: info.ModTime(),
		})
	}

	return logs, nil
}

func inferSessionStatus(filename, path string) string {
	nameLower := strings.ToLower(filename)

	if strings.Contains(nameLower, "error") || strings.Contains(nameLower, "fail") {
		return "error"
	}
	if strings.Contains(nameLower, "pass") || strings.Contains(nameLower, "success") || strings.Contains(nameLower, "complete") {
		return "pass"
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return "unknown"
	}

	contentLower := strings.ToLower(string(content))
	if strings.Contains(contentLower, "error") || strings.Contains(contentLower, "failed") {
		return "error"
	}
	if strings.Contains(contentLower, "completed successfully") || strings.Contains(contentLower, "passed") {
		return "pass"
	}

	return "unknown"
}
