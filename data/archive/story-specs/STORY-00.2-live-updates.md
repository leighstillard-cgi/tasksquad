# STORY-00.2 · Dashboard Live Updates and Detail Views

**Status:** ready
**Repo:** tasksquad/dashboard
**Depends on:** STORY-00.1
**Priority:** High — improves usability of the monitoring dashboard

## Context

STORY-00.1 delivered a functional dashboard but with two limitations:
1. Users must manually refresh the browser to see updates
2. Clicking on completion reports or session logs shows a placeholder alert, not the actual content

This story adds WebSocket-based live updates and proper detail views for agent output.

## Required Reading

- `core/dashboard/README.md` — current architecture
- `core/dashboard/internal/handlers/handlers.go` — existing refresh loop
- `data/completions/STORY-00.1-completion.md` — example completion report format
- `data/session-logs/` — example session log format

## Scope

### In Scope

1. **WebSocket live updates** — Panels update automatically when files change
2. **Completion report detail view** — Click to see full rendered markdown
3. **Session log detail view** — Click to see full rendered markdown
4. **Transcript path capture** — Record Claude Code session path in session-log metadata for future linkability

### Out of Scope

- K8s pod log streaming (deferred)
- Live agent output during execution (deferred)
- Raw JSONL transcript viewer (future story if needed)

## Design

### WebSocket Architecture

```
┌─────────────────────────────────────────────────────┐
│                    Dashboard                         │
│                                                     │
│  fsnotify ──► file change detected                  │
│       │                                             │
│       ▼                                             │
│  parse files ──► broadcast JSON to WebSocket hub    │
│                         │                           │
│                         ▼                           │
│              ┌─────────────────────┐                │
│              │   WebSocket Hub     │                │
│              │  (manages clients)  │                │
│              └─────────────────────┘                │
│                    │    │    │                      │
│                    ▼    ▼    ▼                      │
│              Browser connections                    │
└─────────────────────────────────────────────────────┘
```

**Server-side:**
- Add `fsnotify` watcher on worklog directories: `data/dispatch-log.md`, `backlog.md`, `data/completions/`, `data/escalations/`, `data/dispatches/`, `data/session-logs/`
- On file change: debounce (100ms), re-parse affected files, broadcast update
- WebSocket endpoint at `/ws` using gorilla/websocket
- Hub pattern: register/unregister clients, broadcast to all

**Client-side:**
- Connect to `/ws` on page load
- Receive JSON updates with panel data
- Update DOM without full page reload
- Reconnect with exponential backoff on disconnect

### Detail View

**Modal approach:**
- Click completion report or session log → fetch markdown via API
- Render in a modal overlay with close button
- Parse markdown to HTML server-side (use goldmark or similar)

**API endpoints:**
- `GET /api/completion/:filename` → returns rendered HTML
- `GET /api/session-log/:filename` → returns rendered HTML

### Transcript Path Capture

Update dispatch skill to record the session JSONL path in session-log metadata:

```yaml
---
story_id: STORY-XX.X
transcript_path: ~/.claude/projects/.../session-id.jsonl
---
```

This makes it possible to link to or retrieve the raw transcript in future if needed.

## Acceptance Criteria

| Criteria | Verification |
|----------|--------------|
| WebSocket connection established | Browser console shows `WebSocket connected` on page load |
| Panel updates without refresh | Create a dispatch file externally, verify dashboard shows it within 2s without F5 |
| Completion detail view works | Click a completion report, verify modal shows full rendered markdown |
| Session log detail view works | Click a session log, verify modal shows full rendered markdown |
| WebSocket reconnects on disconnect | Kill and restart server, verify browser reconnects automatically |
| Transcript path in session logs | Check session-log files include `transcript_path` field |

## Sub-tasks

- [ ] Add fsnotify dependency and file watcher for worklog directories
- [ ] Implement WebSocket hub (register, unregister, broadcast)
- [ ] Add `/ws` endpoint with gorilla/websocket
- [ ] Implement debounced file change handler that triggers re-parse and broadcast
- [ ] Add client-side WebSocket connection with reconnect logic
- [ ] Add client-side DOM update handlers for each panel
- [ ] Add goldmark dependency for markdown rendering
- [ ] Implement `/api/completion/:filename` endpoint
- [ ] Implement `/api/session-log/:filename` endpoint
- [ ] Add modal component to index.html with close button
- [ ] Add client-side modal open/close logic
- [ ] Update dispatch skill to capture transcript path
- [ ] Write tests for WebSocket hub, file watcher, markdown rendering
- [ ] Manual smoke test: verify live updates and detail views work
- [ ] Write completion report

## Technical Notes

**Dependencies to add:**
- `github.com/gorilla/websocket` — WebSocket support
- `github.com/fsnotify/fsnotify` — File system notifications
- `github.com/yuin/goldmark` — Markdown to HTML rendering

**Debouncing:**
File systems can emit multiple events for a single change (write + chmod). Debounce with 100ms window to avoid redundant broadcasts.

**Security:**
- Validate filenames in detail view endpoints (no path traversal)
- Only serve files from expected directories (`data/completions/`, `data/session-logs/`)

## Estimated Effort

~300-400 lines of Go (WebSocket hub, watcher, endpoints)
~100 lines of JS (WebSocket client, DOM updates, modal)
~50 lines of HTML/CSS (modal styling)
