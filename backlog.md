# TaskSquad — Product Backlog

**Last updated:** 2026-04-15
**Format:** Each story maps to a dispatch unit. PM agent reads this file to determine what to dispatch next.

> **Note:** Client-specific stories and epics are maintained separately in `.client/backlog-client.md` (gitignored). This file contains only framework and infrastructure stories safe for the public repo.

---

## Status Key

| Status | Meaning |
|---|---|
| `ready` | All dependencies met, available for dispatch |
| `in-progress` | Dispatched to a worker agent |
| `done` | Completion report validated, story closed |
| `blocked` | Has unresolved dependencies or requires human approval |
| `cancelled` | Story removed from active tracking |

---

## Infrastructure — TaskSquad Bootstrap

### STORY-00.1 · Subagent Monitor and Progress Viewer

**Status:** blocked (waiting for STORY-04.4 — dispatch mechanism must exist first)
**Repo:** tasksquad/dashboard (repo to be created)
**Depends on:** STORY-04.4
**Priority:** High — first story dispatched through the workflow once dispatch mechanism is ready

**Description:** Build a terminal-native monitoring tool for TaskSquad subagent operations. Monitors subagent dispatches, completion reports, escalations, and batch progress. Serves dual purpose: required infrastructure + validates the agent lifecycle.

**Acceptance criteria:**
- [ ] Go binary with HTTP server compiles and runs (single container, no K8s dependency)
- [ ] Active work panel: reads `dispatch-log.md` and `session-logs/`, shows dispatched stories with status
- [ ] Completion feed: reads `completions/`, shows recent reports with pass/fail
- [ ] Escalation panel: reads `escalations/`, shows unresolved items with badge count
- [ ] Backlog overview: reads `backlog.md`, shows stories by status with progress counters
- [ ] Batch progress: shows N/total validated, M failures, K conversions needed
- [ ] Manual dispatch form: writes dispatch files to `dispatches/`
- [ ] Session log viewer: browse session-logs/ with filtering by status (pass/fail/error)
- [ ] Refreshes on worklog repo changes (git pull on interval or fsnotify)
- [ ] Standalone mode: works without any external APIs, reads only the filesystem

**Sub-tasks:**
- [ ] Set up Go module with HTTP server skeleton
- [ ] Implement git integration: clone/pull worklog repo, parse markdown files
- [ ] Build active work panel (reads dispatch-log.md + session-logs/)
- [ ] Build completion feed (reads completions/)
- [ ] Build escalation panel (reads escalations/)
- [ ] Build backlog overview with batch progress counters
- [ ] Build manual dispatch form (writes to dispatches/)
- [ ] Build session log viewer with status filtering
- [ ] Write tests for markdown parsing, panel rendering, dispatch file creation
- [ ] Write completion report with evidence for all acceptance criteria

**Spec:** `story-specs/STORY-00.1-dashboard.md` (needs updating to match revised scope)

---

## EPIC-03 · Framework Port and Harness Setup

*Port proven processes, tooling, and ways of working from the reference harness into TaskSquad. Establish the wiki, knowledge graph, cross-session memory, hooks, skills, and agent infrastructure needed to operate at scale.*

### STORY-03.1 · Wiki Structure and Lint Tooling

**Status:** done
**Repo:** tasksquad (this repo)
**Depends on:** none
**Priority:** High — prerequisite for graphify and structured documentation

**Description:** Port the wiki structure definition, lint tooling, and generated manuals pattern from the reference harness. Create the wiki directory hierarchy, frontmatter schema (STRUCTURE.md), lint scripts (lint-wiki.sh + lint-wiki-helper.py), and manual generation scripts. Adapt page types for the project context (drop domain-specific types, add migration/validation types as needed).

**Acceptance criteria:**
- [ ] `core/docs/wiki/STRUCTURE.md` created with page types adapted for project context
- [ ] `core/scripts/lint-wiki.sh` and `lint-wiki-helper.py` ported and working
- [ ] `core/scripts/generate-manuals.sh` ported and working
- [ ] Wiki directory created with initial index page
- [ ] Lint passes clean on initial wiki pages
- [ ] At least one generated manual produced from wiki content

---

### STORY-03.2 · Core Framework Separation

**Status:** done
**Repo:** tasksquad (this repo)
**Depends on:** none
**Priority:** High — establishes the framework + project overlay pattern

**Description:** Restructure TaskSquad into the core/ (reusable framework) + project/ (engagement-specific) overlay pattern. Move generic standards, templates, and scripts into core/. Create Code-SOP.md for worker agents. Port the bootstrap and cascade scripts for multi-repo consistency.

**Acceptance criteria:**
- [ ] `core/` directory created with CLAUDE.md (PM template), Code-SOP.md (worker template)
- [ ] Existing docs/standards/ moved under core/docs/standards/
- [ ] Existing templates/ moved under core/templates/
- [ ] Additional templates ported: ADR, epic, concept, component, runbook, draft, lint-report
- [ ] `core/scripts/bootstrap-repo.sh` ported and adapted for TaskSquad context
- [ ] `core/scripts/cascade.sh` ported
- [ ] `core/scripts/check-repo-health.sh` ported
- [ ] Project overlay directory created (project/ or data/) for engagement-specific content
- [ ] CLAUDE.md updated to load via @-imports from core/ and project/

---

### STORY-03.3 · Port Skills and Agents

**Status:** done
**Repo:** tasksquad (this repo)
**Depends on:** STORY-03.2
**Priority:** High

**Description:** Port the PM agent skills from the reference harness that apply to TaskSquad. Adapt dispatch skill from Slack-based to terminal/subagent-based. Port end-session handoff skill. Create .claude/ directory structure with skills and agents.

**Acceptance criteria:**
- [ ] `.claude/skills/process-completion/` ported and adapted (file-based tracking)
- [ ] `.claude/skills/state-of-play/` ported and adapted
- [ ] `.claude/skills/audit/` ported (backlog vs completion reports sync check)
- [ ] `.claude/skills/end-session/` ported (session handoff to session-state.md)
- [ ] `.claude/skills/dispatch/` created — terminal-native dispatch (writes dispatch file, optionally launches subagent)
- [ ] `.claude/agents/backlog-auditor` ported (Haiku background auditor)
- [ ] All skills documented in CLAUDE.md

---

### STORY-03.4 · Hooks and Safety Controls

**Status:** done
**Repo:** tasksquad (this repo)
**Depends on:** STORY-03.2
**Priority:** Medium

**Description:** Port applicable hooks from the reference harness. Create canonical-infra-inject for the target environment (DB connection details, environment facts). Adapt gh-guard if GitHub is used. Create project-level settings.json with hook configuration.

**Acceptance criteria:**
- [ ] `.claude/settings.json` created with project-level hook configuration
- [ ] Canonical infrastructure injection hook created (template — populated when environment details are known)
- [ ] RTK rewrite hook documented (user-level, not project-level — users will need RTK installed)
- [ ] gh-guard adapted if GitHub issues are used for tracking
- [ ] Hook documentation in CLAUDE.md explaining what each hook does and why

---

### STORY-03.5 · Graphify Knowledge Graph Setup

**Status:** done
**Repo:** tasksquad (this repo)
**Depends on:** STORY-03.1
**Priority:** Medium

**Description:** Set up graphify to index the TaskSquad wiki and generate a knowledge graph. Configure CLAUDE.md to reference graphify output. Create the maintenance workflow (rebuild after edits).

**Acceptance criteria:**
- [ ] graphify installed and runnable in the environment
- [ ] Initial graph generated from wiki + standards + guides
- [ ] `graphify-out/GRAPH_REPORT.md` generated with god nodes and community structure
- [ ] `graphify-out/graph.html` interactive visualization available
- [ ] CLAUDE.md updated with graphify instructions (read GRAPH_REPORT before architecture questions)
- [ ] Rebuild command documented for post-edit maintenance

---

### STORY-03.6 · Claude-Mem Cross-Session Memory

**Status:** ready
**Repo:** tasksquad (this repo)
**Depends on:** none
**Priority:** Medium

**Description:** Configure claude-mem plugin for TaskSquad. Set up the PostToolUse/Stop/SessionStart hook chain for observation recording and session summaries. Document the three-layer query workflow (search → timeline → get_observations).

**Acceptance criteria:**
- [ ] claude-mem plugin enabled in project or user settings
- [ ] Worker service running and accessible
- [ ] PostToolUse hook recording observations
- [ ] Stop hook generating session summaries
- [ ] SessionStart hook injecting recent context
- [ ] Query workflow documented in CLAUDE.md or tooling.md
- [ ] Verified: observations persist across sessions

---

### STORY-03.8 · Environment Bootstrap Install Script

**Status:** ready
**Repo:** tasksquad (this repo)
**Depends on:** none
**Priority:** High — needed before any target environment can use the harness

**Description:** Create an install script that bootstraps a fresh Claude Code environment with the plugins and tools TaskSquad depends on. Must handle: claude-mem plugin, graphify, RTK, and any required Claude Code plugins/settings. Script should be idempotent (safe to re-run) and report what was installed vs already present.

**Acceptance criteria:**
- [ ] `scripts/install.sh` created and executable
- [ ] Installs claude-mem plugin (or verifies already installed)
- [ ] Installs graphify (or verifies already installed)
- [ ] Installs RTK (or verifies already installed)
- [ ] Configures Claude Code settings.json with required hooks and plugin enablement
- [ ] Configures .NET LSP plugin if .NET SDK is detected
- [ ] Idempotent: safe to run multiple times, skips already-installed components
- [ ] Reports summary: what was installed, what was already present, what failed
- [ ] Works in Docker container environments (no sudo assumed)
- [ ] Documents prerequisites (Node.js, Python, etc.) and fails early with clear message if missing

---

### STORY-03.7 · Full Gap Analysis and Port Recommendations

**Status:** ready
**Repo:** tasksquad (this repo)
**Depends on:** none
**Priority:** Low — produces a checklist for review

**Description:** Systematic comparison of reference harness processes vs TaskSquad. Document every process, pattern, and way of working present in the reference but absent in TaskSquad. Tag each as: must-port, nice-to-have, or not-applicable. Present for decision.

**Acceptance criteria:**
- [ ] Gap analysis document created covering: skills, hooks, templates, standards, wiki features, scripts, agents, dispatch patterns, completion workflows, escalation handling
- [ ] Each gap tagged with recommendation and rationale
- [ ] Document reviewed
- [ ] Decisions recorded (which gaps to close, which to defer, which to skip)

---

## Dependency Map

```
┌──────────────────────────────────────────────────────┐
│              EPIC-03: Framework Port                  │
│                                                      │
│  STORY-03.1 (Wiki) ──→ STORY-03.5 (Graphify)        │
│  STORY-03.2 (Core) ──→ STORY-03.3 (Skills)          │
│                    ──→ STORY-03.4 (Hooks)            │
│  STORY-03.6 (Claude-mem) — independent               │
│  STORY-03.7 (Gap analysis) — independent             │
│                                                      │
│  STORY-03.3 ──→ STORY-04.4 (Dispatch) [client]      │
│  STORY-04.4 ──→ STORY-00.1 (Monitor)                │
└──────────────────────────────────────────────────────┘

Critical path: 03.2 → 03.3 → 04.4 → 00.1
See .client/backlog-client.md for full dependency map including client stories.
```

---

## Quick Start: What To Do Right Now

1. **STORY-03.1 + 03.2 + 03.6 + 03.7** — Can start immediately in parallel (wiki, core structure, claude-mem, gap analysis)
2. **Critical path** — Unblock dispatch mechanism: 03.2 → 03.3 → 04.4 (client)
3. See `.client/backlog-client.md` for client-specific priorities and environment setup stories
