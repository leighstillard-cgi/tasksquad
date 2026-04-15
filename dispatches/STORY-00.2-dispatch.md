---
story_id: STORY-00.2
dispatched_at: 2026-04-15T03:50:00Z
dispatched_by: pm-agent
attempt: 1
max_retries: 5
---

# STORY-00.2 · Dashboard Live Updates and Detail Views

## Context

STORY-00.1 delivered a functional dashboard but users must manually refresh to see updates, and clicking on reports shows a placeholder. This story adds WebSocket live updates and proper detail views.

## Implementation Location

Continue work in `dashboard/` directory.

## Key Deliverables

1. **WebSocket live updates** — fsnotify watches files, broadcasts changes to connected browsers
2. **Completion report detail view** — Click to see full rendered markdown in modal
3. **Session log detail view** — Click to see full rendered markdown in modal
4. **Transcript path field** — Add to session-log metadata schema

## Dependencies to Add

```
go get github.com/gorilla/websocket
go get github.com/fsnotify/fsnotify
go get github.com/yuin/goldmark
```

## Architecture

```
fsnotify watcher → debounce (100ms) → re-parse → WebSocket hub → broadcast to clients
```

Server-side:
- `/ws` endpoint for WebSocket connections
- Hub pattern: register/unregister clients, broadcast JSON updates
- Watch: dispatch-log.md, backlog.md, completions/, escalations/, dispatches/, session-logs/

Client-side:
- Connect to /ws on page load
- Receive JSON, update DOM without reload
- Reconnect with backoff on disconnect

## Required Reading

- `dashboard/README.md`
- `dashboard/internal/handlers/handlers.go` — existing refresh loop to extend
- `story-specs/STORY-00.2-live-updates.md` — full spec

## Completion Output

Write to: `completions/STORY-00.2-completion.md`
