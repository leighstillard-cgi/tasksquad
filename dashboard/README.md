# TaskSquad Dashboard

A Go-based web dashboard for monitoring TaskSquad agent operations. Reads tracking files from the worklog repo and displays them in a browser UI with a dark theme.

## Features

- *Active Work Panel* — Shows dispatched stories from `dispatch-log.md` with status badges and elapsed time
- *Completion Feed* — Lists completion reports from `completions/` with YAML frontmatter parsing
- *Escalation Panel* — Displays escalations from `escalations/` with a red badge count when items exist
- *Backlog Overview* — Groups stories by status (Done/In Progress/Ready/Blocked) with counts
- *Manual Dispatch Form* — Create dispatch files directly from the UI, auto-commits and pushes to git
- *Session Log Viewer* — Browse `session-logs/` with filtering by status (pass/error/unknown)

## Quick Start

```bash
cd dashboard

# Run against the parent worklog directory
WORKLOG_PATH=.. ./dashboard

# Open http://localhost:8080
```

## Configuration

All configuration is via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `LISTEN_ADDR` | `:8080` | HTTP server address |
| `POLL_INTERVAL` | `30s` | How often the server re-reads tracking files |
| `WORKLOG_PATH` | `.` | Path to the worklog repository root |

## How Refresh Works

The dashboard uses a server-side polling model:

1. *Server polls files*: Every `POLL_INTERVAL` (default 30s), the server re-reads all tracking files and updates its in-memory data
2. *Browser shows current state*: When you load or refresh the page, you see the latest server-side data

*To see updates, refresh your browser.* The page does not auto-reload.

This is simpler than WebSocket streaming and sufficient for a monitoring dashboard. The trade-off:
- You won't see changes appear instantly without a page refresh
- Updates are visible within `POLL_INTERVAL` + time until you refresh

For most workflows, this is fine — you check the dashboard periodically rather than watching it continuously. If real-time streaming becomes necessary (e.g., watching live agent logs), WebSocket support can be added later.

## Panels

### Active Work

Shows all rows from `dispatch-log.md` where status = `dispatched`. Displays:
- Story ID with link to dispatch file
- Target repo
- Dispatched timestamp
- Elapsed time since dispatch
- Status badge

### Completions

Parses YAML frontmatter from files in `completions/` (excludes `archive/`). Shows:
- Story ID
- Status badge (complete/partial/draft)
- Completion timestamp
- Link to full report

### Escalations

Lists files in `escalations/`. Shows:
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
4. Submit — creates file in `dispatches/`, commits and pushes to git

### Session Logs

Lists files in `session-logs/` with status filtering:
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
dashboard/
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

## Standalone Mode

The dashboard works fully without any external dependencies:
- Reads local files only
- No database required
- No K8s/container runtime required
- Git operations are optional (fails gracefully if git not available)

This makes it suitable for local development and simple deployments.
