---
title: "Completion: STORY-03.4"
id: "COMPLETION-STORY-03.4"
status: draft
page_class: agent-write
page_type: story
tags: [completion]
created: "2026-04-15T00:00:00Z"
last_updated: "2026-04-15T00:00:00Z"
supersedes: []
superseded_by:
inbound_links: []
outbound_links: []
github_issue: ""
labels: [completion]
phase: "03"
parent_epic: "EPIC-03"
depends_on: []
repos: ["tasksquad"]
---

# Completion: STORY-03.4

**Story:** STORY-03.4 - Hooks and Safety Controls
**Repo:** tasksquad
**Branch:** feature/STORY-03.4-hooks-and-safety
**Agent:** claude-opus-4-5-20251101
**Timestamp:** 2026-04-15T00:00:00Z

## Summary

Created project-level hook infrastructure for TaskSquad including gh-guard to restrict GitHub issue mutations to the PM agent, canonical-infra-inject template for infrastructure fact injection, and comprehensive documentation. Evaluated lasso-security/claude-hooks and documented integration path for prompt injection defense. RTK rewrite hook documented as user-level configuration.

## Sub-task Evidence

- [x] `.claude/settings.json` created with project-level hook configuration -- File created at `.claude/settings.json` with PreToolUse hooks for gh-guard and canonical-infra-inject, plus git/gh permissions
- [x] Canonical infrastructure injection hook created (template) -- `.claude/hooks/canonical-infra-inject.sh` and `.claude/hooks/canonical-infra-inject.config` created; template facts file at `data/project/data/canonical-facts.md`
- [x] RTK rewrite hook documented (user-level) -- Documented in CLAUDE.md under "User-Level Hooks" section with configuration example
- [x] gh-guard adapted for GitHub issue tracking -- `.claude/hooks/gh-guard.sh` blocks issue mutations except from PM repo (non-worktree); allows view/list/comment
- [x] Hook documentation in CLAUDE.md -- Added "Hooks" section explaining each hook's purpose and configuration
- [x] lasso-security/claude-hooks evaluated and integrated -- Evaluated PostToolUse prompt injection defender; documented installation via `git clone` + `./install.sh` in CLAUDE.md

## Verification

```bash
# gh-guard blocks mutations from worktrees
$ echo '{"tool_name":"Bash","tool_input":{"command":"gh issue close 123"}}' | PROJECT_DIR=...worktree... .claude/hooks/gh-guard.sh
{
  "hookSpecificOutput": {
    "hookEventName": "PreToolUse",
    "permissionDecision": "deny",
    "permissionDecisionReason": "Only the PM agent may close/edit/create GitHub issues..."
  }
}

# gh-guard allows read operations
$ echo '{"tool_name":"Bash","tool_input":{"command":"gh issue view 123"}}' | PROJECT_DIR=...worktree... .claude/hooks/gh-guard.sh
# (no output = allow)

# canonical-infra-inject fires on infrastructure commands
$ echo '{"tool_name":"Bash","tool_input":{"command":"psql -h localhost"}}' | PROJECT_DIR=... .claude/hooks/canonical-infra-inject.sh
{
  "hookSpecificOutput": {
    "hookEventName": "PreToolUse",
    "permissionDecision": "allow",
    "permissionDecisionReason": "Infrastructure operation detected (call 1 of session)..."
  }
}

# settings.json validates
$ jq . .claude/settings.json > /dev/null && echo "Valid JSON"
Valid JSON
```

## Deviations from Spec

**lasso-security/claude-hooks integration approach:** Documented installation as external dependency rather than vendoring or submodule. Rationale: (1) The project has its own installer script with interactive configuration, (2) keeping it external allows independent updates, (3) STORY-03.8 (bootstrap scripts) is the appropriate place for automated installation. The CLAUDE.md documentation provides clear installation instructions.

## Architectural Escalations

None.

## New Patterns Discovered

**Hook testing pattern:** Created test commands that pipe JSON input to hooks and verify output structure. This pattern should be documented for future hook development:
```bash
echo '{"tool_name":"X","tool_input":{...}}' | PROJECT_DIR=... ./hook.sh
```

## Files Changed

- `.claude/settings.json` -- Created with hook configuration and permissions
- `.claude/hooks/canonical-infra-inject.sh` -- Created PreToolUse hook for infrastructure fact injection
- `.claude/hooks/canonical-infra-inject.config` -- Created pattern configuration for canonical injection
- `.claude/hooks/gh-guard.sh` -- Created PreToolUse hook to guard GitHub issue mutations
- `data/project/data/canonical-facts.md` -- Created template for canonical infrastructure facts
- `CLAUDE.md` -- Added Hooks documentation section
