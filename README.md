# TaskSquad Worklog

Source of truth for architecture decisions, story specifications, and work tracking for the TaskSquad multi-agent development workflow.

This repo is the **shared brain** for all agents. It does not contain application code — that lives in the individual service repos. If you're trying to understand what needs to be built, what decisions have been made, or what's currently in progress, you're in the right place.

---

## Repo Contents

| File / Directory | Purpose |
|---|---|
| `CLAUDE.md` | Shared agent standards — symlinked into every code repo |
| `PM_INSTRUCTIONS.md` | Program Manager agent behaviour and rules |
| `backlog.md` | Full product backlog with story status and priorities |
| `dispatch-log.md` | PM-maintained: what's assigned, to whom, when |
| `observer-log.md` | Audit trail of live monitor interventions |
| `adrs/` | Architecture Decision Records |
| `guides/` | Living reference guides (e.g., conversion patterns) |
| `story-specs/` | Expanded story specifications for agent handoff |
| `completions/` | Repo agents write completion reports here |
| `completions/archive/` | PM moves processed reports here |
| `dispatches/` | PM writes dispatch files here to assign work |
| `escalations/` | Flagged items needing human review |
| `state-of-play/` | PM generates periodic status summaries |
| `schema.md` | Canonical data model documentation (stub) |
| `architecture.md` | System architecture and design rationale |
| `domain-knowledge.md` | Business context for agents |
| `templates/` | Standard formats for completions, dispatches, and stories |

---

## How Agents Use This Repo

**Program Manager agent** — working directory is this repo. Polls on a 5-minute cycle: checks for new completions, validates evidence, dispatches next stories, updates backlog and dispatch log.

**Repo agents (workers)** — receive this repo's docs via symlinks in their code repo's `docs/` directory. Read story specs and ADRs before starting work. Write completion reports to `completions/` when done.

**Architecture advisor** — receives state-of-play uploads for periodic sync. Generates briefs that are dropped into this repo for the PM to process.

---

## Quick Start: What To Do Right Now

1. **Phases B–D** — Build PM agent container, worker infrastructure, and create the `tasksquad/dashboard` repo (see implementation plan)
2. **STORY-00.1** — Once infrastructure is ready, dispatch the monitoring dashboard as the first automated story
