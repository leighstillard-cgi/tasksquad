package main

import (
	"embed"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/tasksquad/dashboard/internal/config"
	"github.com/tasksquad/dashboard/internal/handlers"
)

//go:embed templates/*.html
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg := config.Load()

	worklogPath, err := filepath.Abs(cfg.WorklogPath)
	if err != nil {
		logger.Error("failed to resolve worklog path", "error", err)
		os.Exit(1)
	}

	funcMap := template.FuncMap{
		"base": filepath.Base,
	}
	tmpl, err := template.New("").Funcs(funcMap).ParseFS(templateFS, "templates/*.html")
	if err != nil {
		logger.Error("failed to parse templates", "error", err)
		os.Exit(1)
	}

	h := handlers.New(worklogPath, tmpl, logger)
	h.StartRefreshLoop(cfg.PollInterval)

	// Start live updates (WebSocket hub and file watcher)
	if err := h.StartLiveUpdates(); err != nil {
		logger.Warn("failed to start live updates", "error", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", h.Index)
	mux.HandleFunc("/api/data", h.APIData)
	mux.HandleFunc("/api/refresh", h.APIRefresh)
	mux.HandleFunc("/api/dispatch", h.APIDispatch)
	mux.HandleFunc("/api/session-logs", h.SessionLogsFiltered)
	mux.HandleFunc("/api/completion/", h.APICompletion)
	mux.HandleFunc("/api/session-log/", h.APISessionLog)
	mux.HandleFunc("/ws", h.Hub().ServeWs)
	mux.Handle("/static/", http.FileServerFS(staticFS))

	logger.Info("starting server",
		"address", cfg.ListenAddr,
		"worklog_path", worklogPath,
		"poll_interval", cfg.PollInterval.String(),
	)

	if err := http.ListenAndServe(cfg.ListenAddr, mux); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}
