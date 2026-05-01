---
story_id: STORY-03.4
dispatched_at: 2026-04-15T02:13:00Z
dispatched_by: pm-agent
attempt: 1
max_retries: 5
---

# STORY-03.4: Hooks and Safety Controls

## Story Spec

**Status:** ready → in-progress
**Repo:** tasksquad (this repo)
**Depends on:** STORY-03.2 (done)
**Priority:** Medium

**Description:** Port applicable hooks from the reference harness. Create canonical-infra-inject for the target environment (DB connection details, environment facts). Adapt gh-guard if GitHub is used. Create project-level settings.json with hook configuration.

**Additional requirement from PM:** Evaluate and integrate lasso-security/claude-hooks for prompt injection defense. See https://github.com/lasso-security/claude-hooks

**Acceptance criteria:**
- [ ] `.claude/settings.json` created with project-level hook configuration
- [ ] Canonical infrastructure injection hook created (template — populated when environment details are known)
- [ ] RTK rewrite hook documented (user-level, not project-level — users will need RTK installed)
- [ ] gh-guard adapted if GitHub issues are used for tracking
- [ ] Hook documentation in CLAUDE.md explaining what each hook does and why
- [ ] lasso-security/claude-hooks evaluated and integrated (prompt injection defense via PostToolUse hooks)

## Context

- STORY-03.2 (Core Framework Separation) is complete — `.claude/` directory structure exists
- STORY-03.3 (Skills and Agents) is complete — skills are in `.claude/skills/`
- Existing `.claude/` structure: `skills/`, `agents/`
- No `settings.json` exists yet in `.claude/`

## lasso-security/claude-hooks Info

From https://github.com/lasso-security/claude-hooks:
- Prompt injection defender that scans tool outputs (files, web pages, command results)
- Detects 50+ injection patterns across 5 attack categories
- Uses PostToolUse hooks to warn Claude about suspicious content
- Written in TypeScript
- Can be installed via their installer script or manually configured

## Completion Output

Write completion report to: `data/completions/STORY-03.4-completion.md`
Use template at: `core/templates/story-completion.md`
