# TaskSquad Product Backlog

**Last updated:** 2026-05-01
**Format:** Each story maps to a dispatch unit. The PM agent reads this file to determine what to dispatch next.

Client-specific stories and epics should live in `.client/backlog-client.md`, which is ignored by Git.

## Status Key

| Status | Meaning |
|---|---|
| `ready` | All dependencies met, available for dispatch |
| `in-progress` | Dispatched to a worker agent |
| `done` | Completion report validated, story closed |
| `blocked` | Has unresolved dependencies or requires human approval |
| `cancelled` | Story removed from active tracking |

## Active Framework Work

No active framework stories are queued.

Add new framework stories here using [core/templates/story.md](core/templates/story.md). Move completed dispatches and completion reports into `data/archive/` after the PM validates them.

## Archived History

Pre-release framework bootstrap work has been archived under `data/archive/`.
