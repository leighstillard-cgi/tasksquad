package config

import (
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	ListenAddr   string
	PollInterval time.Duration
	WorklogPath  string
}

func Load() *Config {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":8080"
	}

	pollStr := os.Getenv("POLL_INTERVAL")
	pollInterval := 30 * time.Second
	if pollStr != "" {
		if d, err := time.ParseDuration(pollStr); err == nil {
			pollInterval = d
		}
	}

	worklogPath := os.Getenv("WORKLOG_PATH")
	if worklogPath == "" {
		worklogPath = defaultWorklogPath()
	}

	return &Config{
		ListenAddr:   listenAddr,
		PollInterval: pollInterval,
		WorklogPath:  worklogPath,
	}
}

func defaultWorklogPath() string {
	for _, candidate := range []string{".", "../..", ".."} {
		if _, err := os.Stat(filepath.Join(candidate, "data")); err == nil {
			if _, err := os.Stat(filepath.Join(candidate, "data", "dispatch-log.md")); err == nil {
				return candidate
			}
			if _, err := os.Stat(filepath.Join(candidate, "dispatch-log.md")); err == nil {
				return candidate
			}
		}
	}
	return "."
}
