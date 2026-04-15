package config

import (
	"os"
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
		worklogPath = "."
	}

	return &Config{
		ListenAddr:   listenAddr,
		PollInterval: pollInterval,
		WorklogPath:  worklogPath,
	}
}
