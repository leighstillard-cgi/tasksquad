# TaskSquad — Implementation Plan

**Date:** 2026-03-20
**Organisation:** CGI
**Status:** Ready for implementation

> **Authoritative files:** The operational files in this repository (`CLAUDE.md`, `PM_INSTRUCTIONS.md`, `backlog.md`, templates/) are the source of truth. Where inline examples in this plan differ from the actual files, the actual files take precedence. This plan is the design reference; the repo files are the living configuration.

---

## Changes from Reference Architecture

This plan adapts the TaskSquad workflow (originally built around Slack + a single development VM) for enterprise deployment. The following changes are called out explicitly:

| Original Design | This Plan | Rationale |
|---|---|---|
| Slack as agent communication layer | Filesystem-only protocol + lightweight monitoring dashboard | Slack not approved for agent-to-agent communication in this environment |
| cc-connect (Slack ↔ Claude Code bridge) | Replaced by filesystem polling + hooks | No chat-based bridge needed when agents communicate via files |
| Single development VM with tmux sessions | Kubernetes pods as workers (one per repo/story) | Scalable, ephemeral, matches existing infrastructure |
| Boss manually dispatches stories to repo agents | PM agent polls backlog and dispatches autonomously via filesystem protocol | Removes human from the dispatch loop; boss monitors and intervenes only on escalations |
| Repo agent notifies PM via Slack channel | Repo agent writes completion report; PM detects via filesystem watch or hook trigger | No chat layer needed — the completion report IS the notification |
| Persistent tmux sessions per repo | Ephemeral K8s pods spun up per story, terminated on completion | Better resource utilisation, cleaner isolation, no stale sessions |

**Key architectural principle preserved:** The filesystem is the communication layer. Agents never talk to each other directly. All state lives in markdown files in a shared git repo. The dashboard is a read-only observer of this state, not a message broker.

---

## System Overview

```
┌─────────────────────────────────────────────────────────┐
│                     Dashboard (web UI)                    │
│  Monitors: dispatch log, completions, escalations        │
│  Controls: manual dispatch, story approval, live monitor │
│  Observes agent activity + provides intervention channel │
└────────────────────────────┬────────────────────────────┘
                             │ reads + writes (dispatch, observer input)
                             ▼
┌─────────────────────────────────────────────────────────┐
│                  Shared Git Repository                    │
│  "worklog" — the single source of truth                  │
│                                                          │
│  ├── backlog.md            (stories, status, priorities)  │
│  ├── adrs/                 (architecture decisions)       │
│  ├── story-specs/          (expanded story specs)         │
│  ├── dispatch-log.md       (what's assigned where)        │
│  ├── completions/          (repo agent completion reports)│
│  ├── state-of-play/        (periodic summaries)           │
│  ├── templates/            (completion report template)   │
│  └── CLAUDE.md             (shared agent standards)       │
└─────┬──────────────┬──────────────┬─────────────────────┘
      │              │              │
      ▼              ▼              ▼
┌──────────┐  ┌──────────┐  ┌──────────┐
│ PM Agent │  │ Repo     │  │ Repo     │
│ (K8s pod │  │ Agent A  │  │ Agent B  │
│  or VM)  │  │ (K8s pod)│  │ (K8s pod)│
└──────────┘  └──────────┘  └──────────┘
```

---

## Part 1: Core Components

### 1.1 The Worklog Repository

A git repository that serves as the shared brain for all agents. Contains no application code — only documentation, backlog, story specs, completion reports, and agent coordination files.

**Structure:**

```
worklog/
├── CLAUDE.md                          # Shared agent standards (symlinked into every code repo)
├── backlog.md                         # Full product backlog with story status
├── dispatch-log.md                    # PM-maintained: what's assigned, to whom, when
├── schema.md                          # Canonical schema / data model documentation
├── architecture.md                    # System architecture and design rationale
├── domain-knowledge.md                # Business context for agents
│
├── adrs/                              # Architecture decision records
│   ├── ADR-001-event-schema.md
│   └── ...
│
├── guides/                            # Living reference guides
│   └── conversion-patterns-guide.md
│
├── story-specs/                       # Expanded story specifications
│   ├── STORY-01.1-short-name.md
│   └── ...
│
├── completions/                       # Repo agents write here when done
│   ├── archive/                       # PM moves processed reports here
│   └── STORY-01.1-completion.md
│
├── dispatches/                        # PM writes here to assign work to repo agents
│   └── STORY-01.1-dispatch.md
│
├── state-of-play/                     # PM generates periodic summaries
│   └── 2026-03-20.md
│
├── escalations/                       # Flagged items needing human review
│   ├── archive/                       # PM moves resolved escalations here
│   └── STORY-01.1-escalation.md
│
├── observer-log.md                    # Audit trail of live monitor interventions
│
└── templates/
    ├── story-completion.md            # Template for completion reports
    ├── state-of-play.md              # Template for periodic summaries
    ├── dispatch.md                    # Template for dispatch files
    └── conversion-story.md           # Template for conversion story specs
```

**Story Completion Template** (`templates/story-completion.md`):

```markdown
---
story: STORY-XX.X
status: complete
repo: <org>/<repo-name>
branch: <branch-name>
agent: <agent-identifier>
timestamp: <ISO-8601>
---

## Summary

<One paragraph: what was built, key design choices made during implementation>

## Sub-task Evidence

- [x] <Sub-task from story spec> — <evidence: query result, test output, or file reference>
- [x] <Sub-task from story spec> — <evidence>
- [ ] <Any incomplete sub-task> — <reason>

## Verification

<Include test output, query results, or other evidence that acceptance criteria are met.
Format depends on the story — SQL results for data stories, test suite output for code stories,
screenshot or log snippet for infrastructure stories.>

## Deviations from Spec

<Any changes made during implementation that differ from the story spec.
If none, write "None — implemented as specified.">

## Architectural Escalations

<Any design decisions that should be reviewed by the architecture advisor.
If none, write "None.">

## Files Changed

- `path/to/file.go` — <what changed>
- `path/to/file_test.go` — <what changed>

## Git Reference

- Branch: `<branch-name>`
- Final commit: `<short sha>`
```

**State of Play Template** (`templates/state-of-play.md`):

```markdown
# State of Play

**Generated:** <timestamp>
**Period:** Since <last sync date>
**Generated by:** Program Manager Agent

## Stories Completed Since Last Sync

| Story | Summary | Completion Date |
|---|---|---|
| STORY-XX.X | <one-line summary> | <date> |

## Stories In Progress

| Story | Status | Blocker (if any) |
|---|---|---|
| STORY-XX.X | <brief status> | <blocker or "none"> |

## New ADRs

| ADR | Summary |
|---|---|
| ADR-XXX | <one-line summary> |

## Decisions Made (not yet in ADRs)

- <Decision description — needs ADR if significant>

## Schema Changes

- <Table/column added/modified/dropped>

## Open Escalations

- <Questions flagged by repo agents that need architecture review>

## Current Blockers

- <Anything blocking progress across the project>

## Metrics

- Total open stories: X
- Stories completed this period: X
- Issues: X open, X closed this period
```

### 1.2 Program Manager Agent

A long-running agent (persistent K8s pod or VM process) that manages the project lifecycle. It does NOT write code. It manages documentation, tracks progress, dispatches work, and validates completions.

**Runs as:** A persistent process with Claude Code, working directory set to the worklog repo. Can be a K8s pod with a persistent volume for the git checkout, or a VM process.

**Lifecycle:** Runs a poll loop on a configurable interval (default: 5 minutes).

**Each poll cycle:**

1. `git pull` to get latest state
2. Check `completions/` for new unarchived reports → validate each
3. Check `dispatch-log.md` for stories marked `dispatched` that now have completions → process
4. Check `backlog.md` for the next unblocked, undispatched story → dispatch if available
5. `git add . && git commit && git push` any changes made

### 1.3 Repo Agents (Workers)

Ephemeral agents that receive a story, implement it, and write a completion report. Each runs in its own K8s pod with access to the relevant code repo and a symlinked copy of the worklog docs.

**Runs as:** K8s Job or Pod, spun up by the PM agent's dispatch mechanism, terminated after completion report is written and committed.

**Lifecycle:**

1. Pod starts with: code repo checked out, worklog docs symlinked, story ID passed as environment variable
2. Claude Code session starts, reads `CLAUDE.md` (which points to symlinked worklog docs)
3. Agent reads its assigned story spec
4. Agent implements the story, writes tests, verifies acceptance criteria
5. Agent writes completion report to `worklog/completions/STORY-XX.X-completion.md`
6. Agent commits and pushes both the code repo and the completion report
7. Pod terminates

### 1.4 Architecture Advisor

A separate Claude session (web app or dedicated instance) used for design validation and strategic decisions. Not part of the automated loop. Receives periodic state-of-play uploads and generates briefs for the PM to process.

**Interface:** Manual — a human uploads the state-of-play document and any specific items for review, discusses with the advisor, downloads briefs, and drops them in the worklog repo.

---

## Part 2: Agent Instructions

### 2.1 Shared Standards (CLAUDE.md)

A single `CLAUDE.md` file lives in the worklog repo and is symlinked into every code repo. Contains project-wide standards: language defaults, behavioural rules, critical rules (security, testing, error handling), and the completion report protocol.

The file should NOT contain repo-specific context. Each code repo has its own minimal `CLAUDE.md` that says:

```markdown
# Read docs/shared_standards.md first — it contains project-wide standards.

## Repo-Specific Context
- This repo is <description>
- Language: <Go/Python/etc>
- Key dependencies: <list>
- See docs/plans/ for active design docs
```

The symlinked shared standards file is the single source of truth. One edit propagates to all repos.

### 2.2 PM Agent Instructions

> **Authoritative source:** `PM_INSTRUCTIONS.md` in the worklog repo root. The excerpt below is a summary; see the actual file for the complete instructions including git conflict handling, stalled story detection, escalation resolution, and story generation.

Stored in `PM_INSTRUCTIONS.md`:

```markdown
# Program Manager Agent

You manage the development workflow for this project. You maintain
documentation, dispatch stories to worker agents, and validate completion
reports.

## Poll Cycle (every 5 minutes)

1. git pull
2. Process new completion reports in completions/
3. Check for dispatchable stories in backlog.md
4. Commit and push any changes

## Dispatching Stories

Check backlog.md for the next story where:
- Status is "ready" (not "in-progress", "blocked", or "done")
- All dependencies are marked "done"
- The story is not already in dispatch-log.md as "dispatched"

To dispatch:
1. Update dispatch-log.md: add row with story ID, timestamp, status "dispatched"
2. Create a dispatch file: dispatches/STORY-XX.X-dispatch.md containing:
   - Story ID and title
   - Repo to work in
   - Link to the story spec
   - Any additional context from recent completions
3. Update backlog.md: set story status to "in-progress"
4. Commit and push

## Processing Completion Reports

For each new file in completions/ (not in archive/):
1. Read the completion report
2. Find the matching story spec in story-specs/ or backlog.md
3. For each acceptance criterion:
   - Verify the report provides evidence (query results, test outputs)
4. If ALL criteria met:
   - Update backlog.md: mark story "done", check all sub-task boxes
   - Update dispatch-log.md: set status to "complete"
   - Move completion report to completions/archive/
5. If ANY criterion lacks evidence:
   - Write an escalation to escalations/STORY-XX.X-escalation.md
   - Do NOT mark the story as done
6. If the completion report has architectural escalations:
   - Copy them to escalations/ for human review

## Dispatch Rules

- Never dispatch two stories to the same repo simultaneously
- Never dispatch stories out of dependency order
- Stories touching schema, event contracts, or adapter interfaces
  require human approval — write to escalations/ and wait
- Maximum 3 concurrent dispatches (configurable)

## What You Do NOT Do

- Write application code
- Make architecture decisions
- Approve schema changes
- Override acceptance criteria
- Dispatch stories marked as "blocked"
```

### 2.3 Repo Agent Instructions

> **Authoritative source:** `CLAUDE.md` in the worklog repo root (symlinked into code repos as `docs/shared_standards.md`). The excerpt below is a summary; see the actual file for the complete instructions.

Delivered via the shared `CLAUDE.md` symlinked into every repo. The key sections:

```markdown
## Workflow: Receiving a Story

On session start:
1. Read CLAUDE.md (this file) and all symlinked docs
2. Check for a dispatch file: ../worklog/dispatches/STORY-XX.X-dispatch.md
   (or receive story ID as environment variable)
3. Read the full story spec from docs/story-specs/ or docs/backlog.md
4. Read relevant ADRs, schema docs, and domain knowledge
5. Implement the story as specified

## Workflow: Completion Reports

When done, create: ../worklog/completions/STORY-XX.X-completion.md

The report MUST include:
- Evidence for every acceptance criterion
- Any deviations from spec with rationale
- Any architectural questions that need escalation
- Git branch name and commit hash of the final state

Commit and push both the code repo and the completion report.

## Workflow: Boundaries

Do not:
- Make architecture decisions — escalate in the completion report
- Skip acceptance criteria — document why if you can't meet one
- Modify worklog files other than completions/ and dispatches/
- Continue working after writing the completion report — terminate
```

---

## Part 3: Dispatch and Notification Mechanism

### 3.1 PM Dispatches Work (replacing manual dispatch)

The PM agent autonomously dispatches stories by writing dispatch files:

```
worklog/dispatches/STORY-XX.X-dispatch.md
```

**Dispatch file format:**

```markdown
---
story: STORY-XX.X
repo: org/repo-name
dispatched_at: 2026-03-20T10:00:00Z
dispatched_by: pm-agent
---

## Story
<Title from backlog>

## Spec Location
story-specs/STORY-XX.X-spec.md

## Context
<Any relevant context from recent completions or decisions>

## Dependencies Completed
- STORY-XX.Y (completed 2026-03-19)
```

### 3.2 Worker Pod Lifecycle (K8s)

A controller process (cron job or custom controller) watches `worklog/dispatches/` for new dispatch files and spins up K8s Jobs:

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: story-xx-x-worker
  labels:
    agent: repo-worker
    story: STORY-XX.X
spec:
  backoffLimit: 0
  template:
    spec:
      containers:
      - name: claude-code
        image: agent-worker:latest
        env:
        - name: STORY_ID
          value: "STORY-XX.X"
        - name: CODE_REPO
          value: "org/repo-name"
        - name: WORKLOG_REPO
          value: "org/worklog"
        volumeMounts:
        - name: workdir
          mountPath: /workspace
        - name: ssh-keys
          mountPath: /root/.ssh
          readOnly: true
      restartPolicy: Never
      volumes:
      - name: workdir
        emptyDir: {}
      - name: ssh-keys
        secret:
          secretName: agent-ssh-keys
```

**The worker container's entrypoint script:**

```bash
#!/bin/bash
set -euo pipefail

# Clone repos
git clone "$CODE_REPO" /workspace/code
git clone "$WORKLOG_REPO" /workspace/worklog

# Set up symlinks (worklog docs → code repo's docs/ directory)
mkdir -p /workspace/code/docs
ln -sf /workspace/worklog/CLAUDE.md /workspace/code/docs/shared_standards.md
ln -sf /workspace/worklog/adrs /workspace/code/docs/adrs
ln -sf /workspace/worklog/guides /workspace/code/docs/guides
ln -sf /workspace/worklog/schema.md /workspace/code/docs/schema.md
ln -sf /workspace/worklog/domain-knowledge.md /workspace/code/docs/domain-knowledge.md
ln -sf /workspace/worklog/story-specs /workspace/code/docs/story-specs

# Create feature branch (per CLAUDE.md branching convention)
cd /workspace/code
git checkout -b "feat/${STORY_ID}"

# Start Claude Code with the story assignment
claude --story "$STORY_ID" \
       --instruction "Read docs/shared_standards.md first. Your assigned story is $STORY_ID. Find the spec in docs/story-specs/. Write your completion report to /workspace/worklog/completions/${STORY_ID}-completion.md when done."

# After Claude Code exits, push code repo (feature branch — no conflict risk)
cd /workspace/code
git push -u origin "feat/${STORY_ID}"

# Push completion report to worklog (main branch — retry on conflict)
cd /workspace/worklog
git add completions/ && git commit -m "completion: $STORY_ID"
push_attempts=0
until git push || [ $push_attempts -ge 3 ]; do
    push_attempts=$((push_attempts + 1))
    echo "Push failed, rebasing and retrying (attempt $push_attempts/3)..."
    git pull --rebase
done
if [ $push_attempts -ge 3 ]; then
    echo "ERROR: Failed to push completion report after 3 attempts" >&2
    exit 1
fi
```

### 3.3 Completion Notification (replacing Slack)

When a worker pod finishes and pushes a completion report, the PM needs to know. Three options, from simplest to most robust:

**Option A — Git polling (simplest, no infrastructure):**
The PM agent already runs a poll cycle every 5 minutes. It does `git pull` and checks `completions/` for new files. Latency: up to 5 minutes. No additional infrastructure needed.

**Option B — Git webhook (low latency, minimal infrastructure):**
Configure a webhook on the worklog repo that fires on push events. A small HTTP listener (running alongside the PM or as a sidecar) receives the webhook and triggers the PM's processing cycle immediately. Latency: seconds.

```
worklog repo → push event → webhook → PM trigger endpoint → PM processes
```

The webhook handler is a few lines of Go:

```go
http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
    // Verify webhook signature
    // Trigger PM processing cycle
    pmTrigger <- struct{}{}
    w.WriteHeader(200)
})
```

**Option C — Filesystem watcher (if repos are on shared storage):**
If the worklog repo is on a shared filesystem (NFS, EFS, or a K8s PersistentVolume accessible to both PM and workers), use `inotify` or `fsnotify` to watch the `completions/` directory. The PM gets notified instantly when a new file appears. No git pull needed for detection (though still needed for the actual file content).

**Recommendation:** Start with Option A (git polling). It's zero additional infrastructure and 5-minute latency is fine for most development workflows. If the team finds the delay frustrating, add Option B — a single webhook endpoint is trivial to deploy.

---

## Part 4: Monitoring Dashboard

### 4.1 What It Shows

The dashboard is a read-only web UI that watches the worklog repo and presents:

**Active work panel:**
- Stories currently dispatched (from dispatch-log.md)
- Which repo/pod each is assigned to
- Time since dispatch
- Pod status (running/completed/failed) from K8s API

**Completion feed:**
- Recent completion reports with pass/fail status
- Evidence summary for each acceptance criterion
- Link to the full completion report

**Escalation panel:**
- Items in escalations/ awaiting human review
- Red badge count for unresolved escalations

**Backlog overview:**
- Stories by status: done / in-progress / ready / blocked
- Dependency graph showing what unblocks what

**Dispatch controls (the only write capability):**
- Manual dispatch override: select a story + repo, write a dispatch file
- Story approval: for stories requiring human sign-off before dispatch
- Escalation resolution: mark an escalation as resolved with notes

**Live job monitor:**
- Select any active worker pod to observe
- Streams Claude Code's stdout/stderr in real time (via K8s pod log API or a log sidecar)
- Output displayed in a scrollable terminal-style pane within the dashboard
- Observer can send a prompt to the running agent via an input field — the prompt is written to a file the agent watches (e.g., `/workspace/observer-input.txt`), and the agent's Claude Code session picks it up as a `/btw` interjection or standard prompt depending on a toggle
- Use cases: unsticking a blocked agent, providing a clarification mid-story, redirecting an agent that's going off-track, or simply watching the work happen
- All observer prompts are logged with timestamp and author for audit

### 4.2 Technology

Go backend (single binary), HTML frontend, no JavaScript framework needed. Reads from git repo on disk. Optionally queries K8s API for pod status and log streaming. Serves on an internal port.

```
Dashboard (Go binary)
├── Watches: worklog repo (git pull on interval or fsnotify)
├── Reads: dispatch-log.md, completions/, escalations/, backlog.md
├── Queries: K8s API for pod status + log streaming
├── Writes: dispatches/ (manual override), observer-input (live prompts)
└── Serves: HTML + WebSocket (for live log streaming) on internal port
```

### 4.3 Live Monitor Implementation

The live monitor has two channels: observation (read) and intervention (write).

**Observation channel:**
Worker pods write Claude Code output to stdout. The dashboard streams this via the K8s log API (`/api/v1/namespaces/{ns}/pods/{pod}/log?follow=true`). A WebSocket connection from the browser to the dashboard backend keeps the terminal pane updating in real time. If K8s log streaming is not available, the worker can write to a shared volume that the dashboard tails via `fsnotify`.

**Intervention channel:**
The observer types a prompt in the dashboard input field and selects the delivery mode:

| Mode | Behaviour | When to use |
|---|---|---|
| `/btw` interjection | Written to the agent's observer input file as a background note — the agent sees it but doesn't interrupt its current task | Providing context, flagging something for later, non-urgent guidance |
| Direct prompt | Written as a blocking prompt that the agent must respond to before continuing | Unsticking a blocked agent, correcting course, asking for status |

The exact delivery mechanism for getting observer prompts into a running Claude Code session is TBD — it depends on Claude Code's CLI capabilities at deployment time. Options include piping to stdin, using the `/btw` slash command, writing to a session instruction file that Claude Code polls, or K8s exec to write a file the agent watches. **For the initial build (STORY-00.1), the dashboard implements the UI and audit logging only; prompt delivery to the agent is a follow-up integration.**

**Audit:** Every observer prompt is appended to `worklog/observer-log.md` with timestamp, author, target pod, delivery mode, and prompt text. This ensures full traceability of human interventions.

---

## Part 5: Security and Access Control

### 5.1 Agent Permissions

| Agent | Git (worklog) | Git (code repos) | Database | K8s API | Dashboard |
|---|---|---|---|---|---|
| PM Agent | Read + write (docs, dispatch, completions) | No access | No access | Read (pod status) | No access |
| Repo Agents | Write to completions/ only (see note) | Read + write (assigned repo only) | Read-only (if needed for story) | No access | No access |
| Dashboard | Read + write (dispatches/, observer-log) | No access | No access | Read (pod status + log streaming) | Serves UI |
| Human (boss) | Full access | Full access | Full access | Full access | Full access + live monitor |

> **Note on repo agent worklog access:** Git does not support directory-level write permissions within a single repository. The "completions/ only" constraint is enforced behaviourally via `CLAUDE.md` instructions. To enforce technically, add a server-side pre-receive hook that rejects pushes from worker deploy keys if they modify files outside `completions/`.

### 5.2 Agent Isolation

- Each repo agent pod has access ONLY to its assigned code repo and the worklog repo
- Agents cannot access other agents' repos or running pods
- No agent has database write access — all DB writes happen through the application code the agent produces, reviewed via the completion report before merge
- SSH keys for each agent are scoped: PM gets worklog-only deploy key, repo agents get per-repo deploy keys
- Observer prompt delivery mechanism is TBD pending Claude Code CLI integration (see STORY-00.1 spec for details). The dashboard logs all observer prompts to `observer-log.md` regardless of delivery status

### 5.3 Secrets Management

- No secrets in the worklog repo or any agent-accessible file
- Database credentials for application code are in the organisation's secrets manager (Vault, AWS Secrets Manager, K8s Secrets)
- Agent API keys (for Claude) are injected via K8s Secrets into pod environment variables
- SSH deploy keys are K8s Secrets mounted read-only

---

## Part 6: Implementation Steps

> **Bootstrap strategy:** The first project task delivered using this workflow will be the monitoring dashboard itself (Phase E). This serves a dual purpose: the dashboard is required infrastructure for the workflow, and building it exercises the full agent lifecycle (story spec → dispatch → implementation → completion report → validation) before any business-domain work begins. If the workflow has problems, we discover them while building our own tooling — not while building something for a client.

### Phase A — Worklog Repository

1. Create the worklog git repository
2. Create the directory structure (adrs/, story-specs/, completions/, completions/archive/, dispatches/, escalations/, state-of-play/, templates/)
3. Write the shared `CLAUDE.md` with project standards
4. Write PM agent instructions (`PM_INSTRUCTIONS.md`)
5. Create the completion report template (`templates/story-completion.md`)
6. Create the state-of-play template (`templates/state-of-play.md`)
7. Create `dispatch-log.md` with empty table
8. Create initial `backlog.md` with current stories
9. Create `observer-log.md` with header row
10. Commit and push

### Phase B — PM Agent

1. Build the PM agent container image:
   - Base: Ubuntu + Claude Code CLI
   - Entrypoint: poll loop script
   - Config: poll interval, max concurrent dispatches, worklog repo URL
2. Write the poll loop script:
   ```bash
   while true; do
     git pull
     # Check for new completions
     # Check for dispatchable stories
     # Commit and push changes
     sleep 300  # 5 minute interval
   done
   ```
3. Deploy as a K8s Deployment (replicas: 1) or run on a dedicated VM
4. Test: manually create a completion report, verify PM processes it

### Phase C — Worker Agent Infrastructure

1. Build the worker agent container image:
   - Base: Ubuntu + Claude Code CLI + git + language toolchains (Go, Python)
   - Entrypoint: the worker bootstrap script (clone repos, set up symlinks, run Claude Code)
2. Build the dispatch controller:
   - Watches `worklog/dispatches/` for new dispatch files
   - Creates K8s Jobs for each new dispatch
   - Cleans up completed Jobs after a retention period
3. Configure K8s RBAC:
   - PM ServiceAccount: can read pod status
   - Worker ServiceAccount: minimal permissions (no K8s API access needed)
4. Set up deploy keys:
   - PM: read/write to worklog repo
   - Workers: read/write to assigned code repo + write to worklog completions/
5. Test: manually create a dispatch file, verify a worker pod spins up, runs, and terminates

### Phase D — Symlinks and Code Repo Setup

For each code repository:

1. Create a minimal `CLAUDE.md`:
   ```markdown
   # Read docs/shared_standards.md first.
   ## Repo-Specific Context
   - This repo is <description>
   - Language: <language>
   ```
2. Create `docs/` directory
3. Add symlinks (or document them in the worker bootstrap script since symlinks won't persist in git — the worker script creates them at pod startup)
4. Write any repo-specific ADRs or story specs to the worklog repo

### Phase E — Dashboard

**This is the first story dispatched through the workflow.** Write a full story spec for the dashboard in `story-specs/`, add it to the backlog, and let the PM agent dispatch it to a worker pod. The dashboard implementation validates the entire agent lifecycle end-to-end.

1. Write `story-specs/STORY-00.1-dashboard.md` with acceptance criteria
2. Add to backlog as the first "ready" story
3. PM agent dispatches to a worker pod
4. Worker implements: Go binary with HTTP server
5. Git integration: pull worklog repo on interval, parse markdown files
6. HTML templates for:
   - Active work panel (reads dispatch-log.md)
   - Completion feed (reads completions/)
   - Escalation panel (reads escalations/)
   - Backlog overview (reads backlog.md)
   - Manual dispatch form (writes to dispatches/)
7. Live job monitor: WebSocket-based log streaming from active worker pods, observer prompt input with `/btw` and direct modes
8. K8s API integration for pod status and log streaming
9. Worker writes completion report with evidence
10. PM validates completion report and closes story
11. Deploy dashboard as K8s Deployment or run on the PM's VM

### Phase F — End-to-End Verification

If the dashboard was successfully built and deployed via the agent workflow (Phase E), the end-to-end test is already complete. Verify:

1. The dashboard story went through the full lifecycle: dispatch → implement → completion report → PM validation → closed
2. The dashboard itself shows this lifecycle in its own panels (the dashboard's first entry is its own completion)
3. Dispatch a second small story to confirm the cycle is repeatable
4. Intentionally submit a completion with a missing criterion — verify the PM creates an escalation instead of closing

---

## Part 7: Day-to-Day Operations

### Normal Flow (no human intervention needed)

```
PM polls backlog → finds unblocked story → writes dispatch file
    → controller creates worker pod → worker implements story
    → worker writes completion report → worker pod terminates
    → PM detects completion → validates evidence → marks done
    → PM dispatches next story
```

### Human Intervention Points

| Trigger | What happens | Where you see it |
|---|---|---|
| Escalation created | PM couldn't validate a completion or found architectural questions | Dashboard escalation panel |
| Schema-touching story ready | PM writes to escalations/ instead of dispatching | Dashboard escalation panel |
| Worker pod fails | K8s Job shows failed status | Dashboard active work panel |
| Architecture question | Completion report has escalation section | Dashboard escalation feed |
| Periodic review | PM generates state-of-play on schedule or request | state-of-play/ in worklog |

### Scaling

- Add more worker pods by increasing the PM's max concurrent dispatches
- Add more repos by adding deploy keys and updating the worker bootstrap script
- Dashboard handles any number of stories — it just reads markdown files
- PM agent is single-instance by design (avoids dispatch races)

---

## Part 8: Comparison with Reference Architecture

| Aspect | Reference (Slack-based) | This Plan (filesystem + K8s) |
|---|---|---|
| Communication | Slack channels | Git repo + markdown files |
| Agent hosting | tmux on single VM | K8s pods (ephemeral workers, persistent PM) |
| Dispatch | Human posts in Slack | PM agent writes dispatch files, controller creates pods |
| Notification | Slack message | Git polling or webhook |
| Monitoring | Slack channel scrollback | Purpose-built dashboard |
| Scalability | Limited by VM resources | K8s scales horizontally |
| Cost per story | Always-on tmux sessions | Pay only for active pods |
| Isolation | tmux sessions share filesystem | Each pod has its own filesystem |
| State persistence | Chat history + files | Files only (git is the audit trail) |

---

## Verification Checklist

### Worklog Repo
- [ ] Directory structure created and committed
- [ ] CLAUDE.md contains project standards
- [ ] PM instructions are complete
- [ ] Templates exist for completion reports and state-of-play

### PM Agent
- [ ] Polls on interval and processes completions
- [ ] Dispatches stories in correct dependency order
- [ ] Creates escalations for stories requiring human approval
- [ ] Does not dispatch to repos with active stories

### Worker Agents
- [ ] Pod starts, clones repos, sets up symlinks
- [ ] Reads assigned story spec
- [ ] Writes completion report on finish
- [ ] Pushes both code and completion report
- [ ] Pod terminates after completion

### Dashboard
- [ ] Shows active dispatches with pod status
- [ ] Shows completion feed with validation results
- [ ] Shows escalation panel with unresolved items
- [ ] Manual dispatch override works
- [ ] Refreshes on worklog repo changes
- [ ] Live monitor streams worker pod output in real time
- [ ] Observer can send `/btw` interjection to a running agent
- [ ] Observer can send a direct prompt to a running agent
- [ ] All observer prompts logged to observer-log.md

### End-to-End
- [ ] Story dispatched → implemented → completed → validated → closed
- [ ] Failed validation creates escalation, not closure
- [ ] Dependency ordering respected
- [ ] No two stories dispatched to same repo simultaneously
- [ ] State-of-play generation works
- [ ] Live monitor shows real output during an active story

---

## Part 9: RACI Matrix

**R** = Responsible (does the work), **A** = Accountable (approves/owns the outcome), **C** = Consulted (provides input before), **I** = Informed (told after)

### Story Lifecycle

| Activity | Boss | Arch. Advisor | PM Agent | Repo Agent |
|---|---|---|---|---|
| Identify requirement / idea | R | C | I | — |
| Design validation | C | R | I | — |
| Generate brief / story spec | I | R | I | — |
| Create GitHub issues | I | — | R | — |
| Update backlog and docs | A | — | R | — |
| Decide dispatch priority | A | C | R | — |
| Dispatch story to worker | I | — | R | — |
| Read story spec and ADRs | — | — | — | R |
| Implement code | — | — | — | R |
| Write tests | — | — | — | R |
| Write completion report | — | — | — | R |
| Validate completion evidence | A | — | R | — |
| Close GitHub issue | I | — | R | — |
| Flag escalation | I | C | R | — |
| Resolve architectural escalation | A | R | I | I |
| Generate state of play | I | C | R | — |
| Approve schema / contract changes | A | R | I | — |

### Infrastructure and Operations

| Activity | Boss | Arch. Advisor | PM Agent | Repo Agent |
|---|---|---|---|---|
| Maintain worklog repo docs | A | C | R | — |
| Maintain ADRs | A | R | R | — |
| Maintain CLAUDE.md standards | A | C | R | — |
| Monitor dashboard | R | — | — | — |
| Intervene via live monitor | R | — | — | — |
| K8s pod lifecycle management | I | — | — | — |
| Deploy infrastructure changes | A | C | — | — |
| Security review of agent permissions | A | C | — | — |

### Key Boundaries

| Boundary | Rule |
|---|---|
| Repo agents never | Make architecture decisions, change event schemas or adapter interfaces, modify worklog files beyond completions/ |
| PM agent never | Write application code, approve schema changes, override acceptance criteria |
| Architecture advisor never | Dispatch stories, close issues, access code repos directly |
| Boss always | Has final authority on escalations, schema changes, and strategic direction; is the only human who can intervene via the live monitor |
