# TaskSquad Dashboard

A Go-based web dashboard for monitoring TaskSquad agent operations. Reads tracking files from the worklog repo and displays them in a browser UI with a dark theme.

## Features

- *Live Updates* — WebSocket pushes changes to browser instantly when files change (no refresh needed)
- *Detail Views* — Click completions or session logs to view rendered markdown in modal
- *Active Work Panel* — Shows dispatched stories from `data/dispatch-log.md` with status badges and elapsed time
- *Completion Feed* — Lists completion reports from `data/completions/` with YAML frontmatter parsing
- *Escalation Panel* — Displays escalations from `data/escalations/` with a red badge count when items exist
- *Backlog Overview* — Groups stories by status (Done/In Progress/Ready/Blocked) with counts
- *Manual Dispatch Form* — Create dispatch files directly from the UI, auto-commits and pushes to git
- *Session Log Viewer* — Browse `data/session-logs/` with filtering by status (pass/error/unknown)

## Quick Start

```bash
cd core/dashboard

# Run against the parent worklog directory
WORKLOG_PATH=../.. ./dashboard

# Open http://localhost:8080
```

## Configuration

All configuration is via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `LISTEN_ADDR` | `:8080` | HTTP server address |
| `POLL_INTERVAL` | `30s` | How often the server re-reads tracking files |
| `WORKLOG_PATH` | auto-detected | Path to the worklog repository root |
| `TASKSQUAD_CODEX_BIN` | `codex` | Codex executable used to start the PM agent |
| `TASKSQUAD_PM_AGENT_ENABLED` | `true` | Start a persistent local PM-agent loop with the dashboard |
| `TASKSQUAD_PM_AGENT_INTERVAL` | `5m` | How often the PM agent runs one PM poll cycle |

## Live Updates

The dashboard uses WebSocket for real-time updates:

1. *File watcher*: fsnotify monitors worklog directories for changes
2. *WebSocket broadcast*: When files change, the server pushes updates to all connected browsers
3. *Auto-reconnect*: If the connection drops, the browser reconnects with exponential backoff (1s to 30s)

*Panels update automatically* — no manual refresh needed. A connection status indicator shows WebSocket state.

## Detail Views

Click on a completion report or session log to view the full content in a modal:

- Rendered markdown (GitHub Flavored Markdown via goldmark)
- Close with X button, Escape key, or clicking outside
- Dark theme styling consistent with dashboard

## Panels

### Active Work

Shows all rows from `data/dispatch-log.md` where status = `dispatched`. Displays:
- Story ID with link to dispatch file
- Target repo
- Dispatched timestamp
- Elapsed time since dispatch
- Status badge

### Completions

Parses YAML frontmatter from files in `data/completions/` (excludes `archive/`). Shows:
- Story ID
- Status badge (complete/partial/draft)
- Completion timestamp
- Link to full report

### Escalations

Lists files in `data/escalations/`. Shows:
- Red badge with count (visible in header)
- Story ID
- Escalation reason
- Timestamp

### Backlog Overview

Parses `backlog.md` and groups stories by status. Shows:
- Status cards with counts
- List of ready stories (available for dispatch)

### Manual Dispatch

Form to create new dispatch files:
1. Select a ready story from dropdown
2. Confirm target repo
3. Add optional description
4. Submit — creates file in `data/dispatches/`, commits and pushes to git

### Session Logs

Lists files in `data/session-logs/` with status filtering:
- Filter dropdown: all / pass / error / unknown
- Status badge per log
- Timestamp

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | Main dashboard HTML |
| `/api/data` | GET | JSON dump of all parsed data |
| `/api/refresh` | POST | Force immediate re-read of files |
| `/api/dispatch` | POST | Create a new dispatch file |
| `/api/session-logs?status=X` | GET | Filtered session logs |

## Development

```bash
# Build
go build ./...

# Run tests
go test ./... -v

# Run with live reload (if air is installed)
air
```

## File Structure

```
core/dashboard/
├── main.go                 # Entry point, HTTP server
├── go.mod
├── internal/
│   ├── config/            # Environment variable parsing
│   ├── handlers/          # HTTP handlers + refresh loop
│   └── parser/            # Markdown/YAML parsers
│       ├── dispatch_log.go
│       ├── backlog.go
│       ├── completions.go
│       ├── escalations.go
│       ├── dispatches.go
│       ├── session_logs.go
│       ├── dispatch_writer.go
│       └── *_test.go
├── templates/
│   └── index.html         # Dashboard template
└── static/
    ├── style.css          # Dark theme styles
    └── app.js             # Client-side JS
```

## PM Agent Launch

The dashboard starts a local PM-agent loop by default. The process is identifiable as `data/session-logs/.pm-agent/pm-agent.sh`, writes output to `data/session-logs/pm-agent-*-running.md`, and runs one Codex PM poll cycle every `TASKSQUAD_PM_AGENT_INTERVAL`. Set `TASKSQUAD_PM_AGENT_ENABLED=false` to run the dashboard without the PM agent.

The PM-agent shell wrapper waits for the configured interval before each cycle, so startup makes the PM process visible without immediately mutating backlog or dispatch state.

## Standalone Mode

The dashboard works fully without any external dependencies:
- Reads local files only
- No database required
- No K8s/container runtime required
- Git operations are optional (fails gracefully if git not available)
- PM agent launch requires the local Codex CLI

This makes it suitable for local development and simple deployments.
