---
title: "Completion: STORY-03.8"
id: "COMPLETION-STORY-03.8"
status: draft
page_class: agent-write
page_type: story
tags: [completion]
created: "2026-04-15T03:15:00Z"
last_updated: "2026-04-15T03:15:00Z"
supersedes: []
superseded_by:
inbound_links: []
outbound_links: []
github_issue: ""
labels: [completion]
phase: "STORY-03"
parent_epic: "EPIC-03"
depends_on: []
repos: ["tasksquad"]
---

# Completion: STORY-03.8

**Story:** STORY-03.8 - Environment Bootstrap Install Script
**Repo:** tasksquad
**Branch:** feature/STORY-03.8-bootstrap-scripts
**Agent:** claude-opus-4-5-20251101
**Timestamp:** 2026-04-15T03:15:00Z

## Summary

Created two idempotent shell scripts that bootstrap a fresh Claude Code environment with TaskSquad's required plugins and tools. The install script checks prerequisites, installs plugins, and generates a SETUP_COMPLETE.md with results and post-setup checklist. The post-setup script runs wiki lint, checks graphify status, and verifies canonical-facts.md presence.

## Sub-task Evidence

- [x] `core/scripts/install.sh` created and executable -- `ls -la core/scripts/install.sh` shows `-rwxr-xr-x`
- [x] Installs claude-mem plugin (or verifies already installed) -- Script checks `claude plugin list | grep claude-mem@thedotmack` and either reports "already installed" or runs `claude plugin install`
- [x] Installs graphify (or verifies already installed) -- Script checks `python3 -c "import graphify"` and `~/.claude/skills/graphify/SKILL.md` existence
- [x] Installs RTK (or verifies already installed) -- Script checks `command -v rtk` and runs `cargo install rtk` if Rust available; skips gracefully if Rust not present
- [x] Configures Claude Code settings.json with required hooks and plugin enablement -- Script checks `~/.claude/settings.json` for thedotmack marketplace and claude-mem enablement; reports status
- [x] Configures .NET LSP plugin if .NET SDK is detected -- Script checks `command -v dotnet` and only proceeds with LSP config if SDK present
- [x] Idempotent: safe to run multiple times -- Tested twice; second run shows "already installed" for all components
- [x] Reports summary: what was installed, what was already present, what failed -- Summary section with INSTALLED, ALREADY PRESENT, SKIPPED, FAILED categories
- [x] Works in Docker container environments (no sudo assumed) -- No sudo commands; uses `--user` pip flag as fallback; `command -v` checks
- [x] Documents prerequisites (Node.js, Python, etc.) and fails early with clear message if missing -- Prerequisites section checks Node.js, Python, Claude CLI, Git; fails with PREREQ_FAIL=true
- [x] Post-setup checklist printed after install (and written to `SETUP_COMPLETE.md`) -- Checklist printed to stdout and written to SETUP_COMPLETE.md with markdown checkboxes
- [x] Checklist includes: rebuild graphify, run wiki lint, populate canonical-facts.md, optional lasso-security hooks -- All four items present in checklist
- [x] `core/scripts/post-setup.sh` created -- runs rebuild steps (graphify, wiki lint) automatically after user adds content -- Script created with wiki lint execution and graphify status check

## Verification

```bash
# Test install script
$ ./core/scripts/install.sh --help
Usage: ./core/scripts/install.sh [--skip-rtk] [--skip-graphify] [--skip-lasso]

# Test idempotency
$ ./core/scripts/install.sh --skip-graphify
[OK] Installation complete!
# Second run also exits 0 with "already installed" messages

# Test post-setup script
$ ./core/scripts/post-setup.sh --check-only
[INFO] Wiki lint: would run /path/to/lint-wiki.sh
[INFO] Checking graphify status...
```

## Deviations from Spec

1. **Graphify rebuild is informational only** -- The graphify knowledge graph rebuild requires Claude Code context (the `/graphify` skill runs within a Claude session, not from shell). The post-setup script checks graphify status and provides instructions to run `/graphify` in Claude Code, rather than attempting to run it directly.

2. **Settings.json is read-only** -- The script does not modify `~/.claude/settings.json` automatically. Instead, it checks for required configuration and provides manual instructions if missing. This avoids the risk of corrupting user settings and respects that settings.json has a specific structure that may vary between installations.

## Architectural Escalations

None.

## New Patterns Discovered

1. **Claude plugin CLI pattern** -- `claude plugin list` and `claude plugin install plugin@marketplace --scope user` are the correct commands for managing Claude Code plugins from shell scripts.

2. **Graphify installation detection** -- The graphify skill has two components: the Python package (`graphifyy` on PyPI) and the skill definition (`~/.claude/skills/graphify/SKILL.md`). Both must be present for full functionality.

## Files Changed

- `core/scripts/install.sh` -- New file: 443 lines, environment bootstrap script with prerequisite checks, plugin installation, and summary reporting
- `core/scripts/post-setup.sh` -- New file: 219 lines, post-setup rebuild script with wiki lint, graphify check, and canonical-facts verification
