# Onboarding Guide

This guide is for someone adopting TaskSquad without a walkthrough from the original maintainer.

## What You Are Setting Up

TaskSquad is a worklog and control framework for AI-assisted delivery. It is not your application codebase. It coordinates work across one or more target repositories.

You will normally have:

- This TaskSquad repository as the worklog and agent control plane.
- One or more target code repositories where implementation work happens.
- Project-specific facts and docs under `data/`.
- Client-sensitive content in `.client/`, which is ignored by Git.

## Recommended Environment

Use an isolated environment:

- Devcontainer
- Disposable VM
- Cloud workstation
- WSL distro dedicated to this work
- Container with mounted target repositories

Avoid running broad-permission agents directly on a personal machine with access to unrelated source code, credentials, cloud CLIs, or production configuration.

## Prerequisites

Install these before running setup:

- Git
- Bash-compatible shell
- Node.js
- Python 3 and pip
- Claude Code CLI

Optional but useful:

- Go, for `core/dashboard`
- Rust/Cargo, for RTK
- Docker or devcontainer tooling, for isolation

## Step 1: Clone and Install

```bash
git clone <internal-task-squad-repo-url> tasksquad
cd tasksquad
bash core/scripts/install.sh
```

The installer checks prerequisites, verifies Claude Code plugin configuration, installs available supporting tools, and attempts the first graphify build when possible.

## Step 2: Validate Setup

```bash
bash core/scripts/post-setup.sh
bash core/scripts/lint-wiki.sh
```

Expected result:

- Wiki lint passes.
- `data/project/data/canonical-facts.md` exists.
- Graphify is either generated successfully or the script prints the exact Claude Code `/graphify` command to run.

## Step 3: Fill Project Facts

Edit [data/project/data/canonical-facts.md](../data/project/data/canonical-facts.md).

Include non-secret facts only:

- Service names
- Repository names
- Non-secret URLs
- Environment names
- Cloud account IDs where policy allows
- Read-only endpoints

Do not include:

- Passwords
- API keys
- Access tokens
- Private keys
- Production credentials

## Step 4: Add Project Documentation

Add project-specific docs under `data/wiki/`.

Use the structure in [core/docs/wiki/STRUCTURE.md](../core/docs/wiki/STRUCTURE.md). Run wiki lint after adding pages:

```bash
bash core/scripts/lint-wiki.sh
```

## Step 5: Create Initial Work

Add active work to [backlog.md](../backlog.md) or keep sensitive/client-specific work in `.client/backlog-client.md`.

Use [core/templates/story.md](../core/templates/story.md) for story specs and [core/templates/story-completion.md](../core/templates/story-completion.md) for completion reports.

## Step 6: Bootstrap a Target Repository

```bash
bash core/scripts/bootstrap-repo.sh /path/to/target-repo
bash core/scripts/check-repo-health.sh /path/to/target-repo
```

This links the target repository to the shared agent standards and creates `.claude/repo-context.md`, which you should fill in for that repository.

## Step 7: Run the Workflow

Use this repository for PM coordination:

- PM instructions: [PM_INSTRUCTIONS.md](../PM_INSTRUCTIONS.md)
- Active dispatch log: [data/dispatch-log.md](../data/dispatch-log.md)
- Active dispatches: `data/dispatches/`
- Completion reports: `data/completions/`
- Escalations: `data/escalations/`

Use target repositories for implementation work. Worker agents should write completion reports back to TaskSquad, not silently close work.

## Step 8: Archive Used Records

After completion reports and dispatches have been processed, move them to `data/archive/` or document the durable learning in `data/wiki/`.

Active directories should stay small and current.
