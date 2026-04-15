---
story_id: STORY-03.8
dispatched_at: 2026-04-15T03:05:00Z
dispatched_by: pm-agent
attempt: 1
max_retries: 5
---

# STORY-03.8: Environment Bootstrap Install Script

## Story Spec

**Status:** ready → in-progress
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
- [ ] Post-setup checklist printed after install (and written to `SETUP_COMPLETE.md`)
- [ ] Checklist includes: rebuild graphify, run wiki lint, populate canonical-facts.md, optional lasso-security hooks
- [ ] `scripts/post-setup.sh` created — runs rebuild steps (graphify, wiki lint) automatically after user adds content

## Context from Completed Stories

**STORY-03.4** (Hooks):
- `.claude/settings.json` exists with hook configuration
- Hooks: gh-guard, canonical-infra-inject
- lasso-security/claude-hooks documented for prompt injection defense

**STORY-03.5** (Graphify):
- graphify output in `graphify-out/`
- Rebuild command: `/graphify wiki core/docs/standards guides core/templates --update`

**STORY-03.6** (Claude-Mem):
- claude-mem plugin documented in `project/tooling.md`
- Plugin: `thedotmack/claude-mem` from marketplace

**STORY-03.7** (Gap Analysis):
- Must-port: `poormansadvisor` skill (can be follow-up story)
- All nice-to-haves declined
- Core framework fully ported

## Implementation Notes

1. **scripts/install.sh** should:
   - Check prerequisites (Node.js, Python, Rust/cargo for RTK)
   - Install claude-mem plugin via Claude Code CLI or marketplace
   - Install graphify (check how it's installed - likely npm or pip)
   - Install RTK via `cargo install rtk`
   - Configure ~/.claude/settings.json with plugins enabled
   - Detect .NET SDK and configure LSP if present
   - Print summary and write SETUP_COMPLETE.md

2. **scripts/post-setup.sh** should:
   - Run wiki lint: `./core/scripts/lint-wiki.sh`
   - Run graphify: invoke the skill or run the underlying command
   - Remind about canonical-facts.md population

3. Script should be bash, work without sudo, be idempotent

## Completion Output

Write completion report to: `completions/STORY-03.8-completion.md`
Use template at: `core/templates/story-completion.md`
