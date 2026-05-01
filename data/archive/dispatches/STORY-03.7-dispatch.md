---
story_id: STORY-03.7
dispatched_at: 2026-04-15T02:53:00Z
dispatched_by: pm-agent
attempt: 1
max_retries: 5
---

# STORY-03.7: Full Gap Analysis and Port Recommendations

## Story Spec

**Status:** ready → in-progress
**Repo:** tasksquad (this repo)
**Depends on:** none
**Priority:** Low — produces a checklist for review

**Description:** Systematic comparison of reference harness processes vs TaskSquad. Document every process, pattern, and way of working present in the reference but absent in TaskSquad. Tag each as: must-port, nice-to-have, or not-applicable. Present for decision.

**Acceptance criteria:**
- [ ] Gap analysis document created covering: skills, hooks, templates, standards, wiki features, scripts, agents, dispatch patterns, completion workflows, escalation handling
- [ ] Each gap tagged with recommendation and rationale
- [ ] Document reviewed
- [ ] Decisions recorded (which gaps to close, which to defer, which to skip)

## Context

This is a research/documentation task, not a coding task. The goal is to identify what processes, patterns, and tooling from the reference harness have NOT yet been ported to TaskSquad.

### What Has Been Ported (EPIC-03 completed stories)

**STORY-03.1** (Wiki Structure and Lint Tooling):
- wiki/ directory with STRUCTURE.md
- lint-wiki.sh and lint-wiki-helper.py
- generate-manuals.sh

**STORY-03.2** (Core Framework Separation):
- core/ + project/ overlay pattern
- Standards moved to core/docs/standards/
- Templates in core/templates/
- bootstrap-repo.sh, cascade.sh, check-repo-health.sh

**STORY-03.3** (Skills and Agents):
- .claude/skills/: dispatch, process-completion, state-of-play, audit, end-session, lint-wiki, bootstrap, cascade
- .claude/agents/: backlog-auditor

**STORY-03.4** (Hooks and Safety Controls):
- .claude/settings.json with hook configuration
- gh-guard.sh (GitHub issue mutation guard)
- canonical-infra-inject.sh (infrastructure fact injection)
- lasso-security/claude-hooks documented for prompt injection defense

**STORY-03.5** (Graphify Knowledge Graph):
- graphify-out/ with GRAPH_REPORT.md and graph.html
- CLAUDE.md updated with knowledge graph workflow

**STORY-03.6** (Claude-Mem Cross-Session Memory):
- project/tooling.md with full claude-mem documentation
- CLAUDE.md updated with query workflow

### Reference Harness Location

The reference harness patterns are documented in the user's global CLAUDE.md and skills at:
- ~/.claude/skills/ — user-level skills
- ~/.claude/CLAUDE.md — user-level instructions
- The superpowers-extended-cc skills

### Output Format

Create `data/gap-analysis/EPIC-03-gap-analysis.md` with:

1. **Comparison Table**: Each category (skills, hooks, templates, etc.) with columns: Item | Reference | TaskSquad | Gap? | Recommendation | Rationale

2. **Recommendations Summary**: Grouped by must-port, nice-to-have, not-applicable

3. **Decision Log**: Space for the user to record decisions on each gap

## Completion Output

Write completion report to: `data/completions/STORY-03.7-completion.md`
Use template at: `core/templates/story-completion.md`
