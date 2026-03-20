# STORY-00.1 · TaskSquad Monitoring Dashboard

**Status:** ready
**Repo:** tasksquad/dashboard
**Depends on:** Nothing
**Priority:** High — first story dispatched through the workflow

## Context

This is the first story built using the TaskSquad workflow. It serves a dual purpose: the dashboard is required infrastructure for monitoring the agent fleet, and building it exercises the full agent lifecycle (dispatch → implementation → completion report → PM validation) before any business-domain work begins.

If the workflow has problems, we discover them while building our own tooling — not while building something for a client.

## Required Reading

- `CLAUDE.md` — project standards
- `PM_INSTRUCTIONS.md` — understand the PM's poll cycle and dispatch protocol (the dashboard observes this)
- `templates/story-completion.md` — the format you'll write your completion report in

## Design

### Technology

Go binary, single executable. HTML templates rendered server-side. WebSocket for live log streaming. No JavaScript framework — vanilla JS for the WebSocket client and form interactions. Minimal dependencies.

### Architecture

```
Dashboard (Go binary)
├── Watches: worklog repo (git pull on configurable interval, default 30s)
├── Reads: dispatch-log.md, completions/, escalations/, backlog.md
├── Queries: K8s API for pod status + log streaming
├── Writes: dispatches/ (manual dispatch), observer-log.md (prompt audit)
└── Serves: HTML + WebSocket on configurable port (default 8080)
```

### Panels

**Active work panel:**
- Parse `dispatch-log.md` for rows where status = `dispatched`
- For each, query K8s API for pod status (Running/Completed/Failed/Pending)
- Show: story ID, repo, dispatched timestamp, elapsed time, pod status
- Auto-refresh on git pull interval

**Completion feed:**
- List files in `completions/` (excluding `archive/`)
- Parse YAML frontmatter for story ID, status, timestamp
- Show: story ID, status (complete/partial), timestamp
- Link to full completion report content

**Escalation panel:**
- List files in `escalations/`
- Show count badge (red if > 0)
- Each escalation shows: story ID, reason, timestamp
- "Resolve" button writes a resolution note and moves to archive

**Backlog overview:**
- Parse `backlog.md` for all stories
- Group by status: done / in-progress / ready / blocked
- Show counts per group

**Manual dispatch form:**
- Dropdown: select a story from backlog where status = `ready`
- Input: target repo
- Submit: writes a dispatch file to `dispatches/` using the template format
- Updates `dispatch-log.md` and `backlog.md`

### Live Job Monitor

**Observation (read):**
- Select an active worker pod from the active work panel
- Stream pod logs via K8s API: `GET /api/v1/namespaces/{ns}/pods/{pod}/log?follow=true`
- Display in a scrollable terminal-style pane (monospace font, dark background)
- WebSocket connection from browser to dashboard backend for real-time streaming

**Intervention (write):**
- Text input field below the terminal pane
- Toggle: `/btw` mode (non-blocking) or `Direct` mode (blocking)
- On submit: dashboard writes the prompt to `/workspace/observer-input.txt` on the worker pod (via K8s exec API or shared volume)
- Append to `observer-log.md`: timestamp, author (from session), target pod, mode, prompt text

### Git Integration

- On startup: clone the worklog repo to a local working directory
- On interval (configurable, default 30s): `git pull` to refresh state
- On write operations (dispatch, escalation resolution): `git add`, `git commit`, `git push`
- Handle merge conflicts gracefully: if pull fails due to conflict, log warning and retry on next interval

### Configuration

All via environment variables:

| Variable | Default | Description |
|---|---|---|
| `WORKLOG_REPO` | required | Git URL of the worklog repository |
| `WORKLOG_BRANCH` | `main` | Branch to track |
| `POLL_INTERVAL` | `30s` | How often to git pull |
| `LISTEN_ADDR` | `:8080` | HTTP server listen address |
| `K8S_NAMESPACE` | `default` | Kubernetes namespace for worker pods |
| `K8S_IN_CLUSTER` | `true` | Use in-cluster K8s config (false = use kubeconfig) |

## Acceptance Criteria

| Criteria | Verification |
|---|---|
| Go binary compiles and runs | `go build ./... && ./dashboard --help` succeeds |
| Active work panel shows dispatched stories | Manually add a row to dispatch-log.md, verify it appears |
| Completion feed shows reports | Drop a test completion file in completions/, verify it appears |
| Escalation panel shows unresolved items | Drop a test escalation, verify badge count and content |
| Backlog overview groups by status | Verify counts match backlog.md content |
| Manual dispatch creates dispatch file | Submit form, verify file in dispatches/ and dispatch-log.md updated |
| Live monitor streams logs | Start a test pod, verify output appears in terminal pane |
| Observer /btw prompt works | Send a /btw prompt, verify it appears in observer-log.md and doesn't block the agent |
| Observer direct prompt works | Send a direct prompt, verify it appears in observer-log.md |
| Git integration pulls changes | Modify worklog externally, verify dashboard reflects changes within poll interval |
| K8s pod status displayed | Verify Running/Completed/Failed status shown for active dispatches |
| Dockerfile builds | `docker build .` succeeds |

## Sub-tasks

- [ ] Set up Go module: `go mod init`, HTTP server skeleton, configuration from env vars
- [ ] Implement git integration: clone on startup, pull on interval, commit+push on writes
- [ ] Implement markdown parser: read dispatch-log.md, backlog.md, completion reports, escalations
- [ ] Build active work panel: parse dispatch log, query K8s pod status, render HTML
- [ ] Build completion feed: list completions/, parse frontmatter, render HTML
- [ ] Build escalation panel: list escalations/, render with badge count, resolve action
- [ ] Build backlog overview: parse backlog.md, group by status, render counts
- [ ] Build manual dispatch form: story selector, repo input, writes dispatch file + updates logs
- [ ] Implement WebSocket log streaming: K8s log API → WebSocket → browser terminal pane
- [ ] Build observer prompt input: text field, mode toggle, writes to observer-input.txt + observer-log.md
- [ ] Add K8s client: pod listing, status queries, log streaming, exec for observer input
- [ ] Write Dockerfile with multi-stage build
- [ ] Write tests: markdown parsing, dispatch file creation, observer log formatting
- [ ] Manual smoke test: run dashboard against this worklog repo, verify all panels render
- [ ] Write completion report with evidence for all acceptance criteria
