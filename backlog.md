# TaskSquad — Product Backlog

**Last updated:** 2026-03-20
**Format:** Each story maps to a dispatch unit. PM agent reads this file to determine what to dispatch next.

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

### STORY-00.1 · Monitoring Dashboard

**Status:** blocked (infrastructure prerequisite: Phases B–D must be completed before first dispatch)
**Repo:** tasksquad/dashboard (repo must be created during Phase D)
**Depends on:** PM agent + worker infrastructure + dashboard repo creation (Phases B–D of implementation plan)
**Priority:** High — first story dispatched through the workflow once infrastructure is ready

**Description:** Build the TaskSquad monitoring dashboard. This serves a dual purpose: the dashboard is required infrastructure, and building it validates the entire agent lifecycle before any business-domain work begins.

**Acceptance criteria:**
- [ ] Go binary with HTTP server compiles and runs
- [ ] Active work panel: reads `dispatch-log.md`, shows dispatched stories with status
- [ ] Completion feed: reads `completions/`, shows recent reports with pass/fail
- [ ] Escalation panel: reads `escalations/`, shows unresolved items with badge count
- [ ] Backlog overview: reads `backlog.md`, shows stories by status
- [ ] Manual dispatch form: writes dispatch files to `dispatches/`
- [ ] Live job monitor: streams worker pod stdout via K8s log API over WebSocket
- [ ] Observer prompt input: supports `/btw` (non-blocking) and direct (blocking) modes
- [ ] Observer prompts logged to `observer-log.md` with timestamp, author, target pod
- [ ] K8s API integration for pod status
- [ ] Refreshes on worklog repo changes (git pull on interval or fsnotify)
- [ ] Deploys as a K8s Deployment or runs standalone

**Sub-tasks:**
- [ ] Set up Go module with HTTP server skeleton
- [ ] Implement git integration: clone/pull worklog repo, parse markdown files
- [ ] Build active work panel (reads dispatch-log.md)
- [ ] Build completion feed (reads completions/)
- [ ] Build escalation panel (reads escalations/)
- [ ] Build backlog overview (reads backlog.md)
- [ ] Build manual dispatch form (writes to dispatches/)
- [ ] Implement WebSocket log streaming from K8s pod log API
- [ ] Build observer prompt input with /btw and direct mode toggle
- [ ] Implement observer audit logging to observer-log.md
- [ ] Add K8s API client for pod status queries
- [ ] Write Dockerfile for K8s deployment
- [ ] Write tests for markdown parsing, panel rendering, dispatch file creation
- [ ] Write completion report with evidence for all acceptance criteria

**Spec:** `story-specs/STORY-00.1-dashboard.md`

---

## SAP ASE to MS SQL Conversion

### STORY-01.1 · SSMA Bulk Conversion and Problem Job Identification

**Status:** blocked (waiting for STORY-00.1 to validate the workflow)
**Repo:** tasksquad/ase-mssql-conversion
**Depends on:** STORY-00.1

**Description:** Run SSMA v10.5 bulk conversion against the ASE source database. Identify all queries and stored procedures that fail conversion or convert incorrectly. Classify each problem job by severity (low/medium/high) and create individual conversion stories in the backlog.

**Acceptance criteria:**
- [ ] SSMA bulk conversion completed
- [ ] All problem jobs identified and listed with failure reason
- [ ] Each problem job classified by severity
- [ ] Individual STORY-02.X stories created in backlog for each problem job
- [ ] Baseline outputs captured for each problem job (result sets, row counts, execution times)

---

### STORY-01.2 · Seed Conversion Patterns Guide

**Status:** blocked (waiting for STORY-01.1)
**Repo:** tasksquad/worklog (this repo)
**Depends on:** STORY-01.1

**Description:** Create the initial conversion patterns guide in the worklog based on the ASE-to-MS-SQL divergence patterns documented in the approach document. This guide is read by every worker agent before starting a conversion.

**Acceptance criteria:**
- [ ] `guides/conversion-patterns-guide.md` formalised with patterns from SSMA triage output
- [ ] Covers: @@error→TRY/CATCH, RAISERROR syntax, SET ROWCOUNT→TOP, CONVERT style codes, cursor syntax, string functions, date functions, temp table scope, cross-database references, identity handling
- [ ] Each pattern has ASE before and MS SQL after examples
- [ ] Guide is referenced in CLAUDE.md as required reading for conversion stories

**Note:** A seed version of the patterns guide already exists at `guides/conversion-patterns-guide.md`, created from the approach document during repo scaffolding. This story formalises and extends it based on actual SSMA triage output from STORY-01.1.

---

## Conversion Stories (created dynamically after STORY-01.1)

*Individual conversion stories (STORY-02.X) will be created here by the PM agent after SSMA triage identifies the specific problem jobs. Each story follows the conversion story template defined in the approach document.*

---

## Dependency Map

```
STORY-00.1 (Dashboard — validates workflow)
    │
    ▼
STORY-01.1 (SSMA bulk conversion + triage)
    │
    ▼
STORY-01.2 (Seed patterns guide)
    │
    ▼
STORY-02.X (Individual conversion stories — parallel, no dependencies between them)
```

---

## Quick Start: What To Do Right Now

1. **Phases B–D** — Build PM agent, worker infrastructure, and set up the dashboard repo (manual infrastructure work)
2. **STORY-00.1** — Once infrastructure is ready, dispatch the monitoring dashboard as the first automated story
