---
story_id: STORY-00.1
dispatched_at: 2026-04-15T03:28:00Z
dispatched_by: pm-agent
attempt: 1
max_retries: 5
---

# STORY-00.1 · TaskSquad Monitoring Dashboard

## Context

First infrastructure story. Go binary with HTTP server that monitors the TaskSquad agent fleet by reading tracking files (dispatch-log.md, data/completions/, data/escalations/, backlog.md, data/session-logs/).

## Implementation Location

Build in `core/dashboard/` directory within this repo. Initialize as Go module.

## Key Design Decisions

- **Standalone first**: K8s integration is optional. Dashboard MUST work fully without K8s access.
- **No JavaScript frameworks**: Vanilla JS only for WebSocket client and forms.
- **Server-side HTML**: Go templates, no SPA.
- **Git integration**: Clone/pull worklog repo, commit/push on writes.

## Scope Clarification

Per the story spec, the K8s live monitor and observer prompt delivery are nice-to-haves. The core deliverables are:
1. Go binary compiles and runs
2. Active work panel (reads dispatch-log.md + data/session-logs/)
3. Completion feed (reads data/completions/)
4. Escalation panel (reads data/escalations/)
5. Backlog overview (reads backlog.md)
6. Manual dispatch form (writes to data/dispatches/)
7. Session log viewer with status filtering
8. Git integration (pull on interval, push on writes)

The spec mentions K8s pod status queries, but standalone mode must work fully without K8s. Start with standalone mode.

## Dependencies Completed

- STORY-03.3: Dispatch mechanism delivered (skills, agents, worktree isolation)
- STORY-03.1: Wiki and lint tooling
- STORY-03.2: Core framework separation
- All EPIC-03 stories complete

## Required Reading

- `CLAUDE.md` — project standards
- `PM_INSTRUCTIONS.md` — understand dispatch protocol
- `core/templates/story-completion.md` — completion report format

## Completion Output

Write to: `data/completions/STORY-00.1-completion.md`
