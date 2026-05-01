---
story_id: STORY-03.6
dispatched_at: 2026-04-15T02:41:00Z
dispatched_by: pm-agent
attempt: 1
max_retries: 5
---

# STORY-03.6: Claude-Mem Cross-Session Memory

## Story Spec

**Status:** ready → in-progress
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

## Context

- `.claude/settings.json` already exists with hook configuration (from STORY-03.4)
- CLAUDE.md already has a Hooks section that can be extended
- claude-mem is an MCP plugin that provides cross-session memory via observation recording

## claude-mem Info

claude-mem is a Claude Code plugin that:
- Records observations during sessions (PostToolUse hook)
- Generates session summaries on stop (Stop hook)
- Injects recent context at session start (SessionStart hook)
- Provides MCP tools for querying: search, timeline, get_observations

The three-layer query workflow:
1. `search` — find relevant observations by keyword
2. `timeline` — see chronological activity
3. `get_observations` — retrieve full observation details

## Completion Output

Write completion report to: `data/completions/STORY-03.6-completion.md`
Use template at: `core/templates/story-completion.md`
