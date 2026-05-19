package config

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	ListenAddr          string
	PollInterval        time.Duration
	WorklogPath         string
	PMAgentEnabled      bool
	PMAgentPollInterval time.Duration
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

	pmAgentPollStr := os.Getenv("TASKSQUAD_PM_AGENT_INTERVAL")
	pmAgentPollInterval := 5 * time.Minute
	if pmAgentPollStr != "" {
		if d, err := time.ParseDuration(pmAgentPollStr); err == nil {
			pmAgentPollInterval = d
		}
	}

	return &Config{
		ListenAddr:          listenAddr,
		PollInterval:        pollInterval,
		WorklogPath:         worklogPath,
		PMAgentEnabled:      boolEnv("TASKSQUAD_PM_AGENT_ENABLED", true),
		PMAgentPollInterval: pmAgentPollInterval,
	}
}

func boolEnv(name string, defaultValue bool) bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(name)))
	switch value {
	case "":
		return defaultValue
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return defaultValue
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
