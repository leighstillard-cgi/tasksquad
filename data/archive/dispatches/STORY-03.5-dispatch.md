---
story_id: STORY-03.5
dispatched_at: 2026-04-15T02:26:00Z
dispatched_by: pm-agent
attempt: 1
max_retries: 5
---

# STORY-03.5: Graphify Knowledge Graph Setup

## Story Spec

**Status:** ready → in-progress
**Repo:** tasksquad (this repo)
**Depends on:** STORY-03.1 (done)
**Priority:** Medium

**Description:** Set up graphify to index the TaskSquad wiki and generate a knowledge graph. Configure CLAUDE.md to reference graphify output. Create the maintenance workflow (rebuild after edits).

**Acceptance criteria:**
- [ ] graphify installed and runnable in the environment
- [ ] Initial graph generated from wiki + standards + guides
- [ ] `graphify-out/GRAPH_REPORT.md` generated with god nodes and community structure
- [ ] `graphify-out/graph.html` interactive visualization available
- [ ] CLAUDE.md updated with graphify instructions (read GRAPH_REPORT before architecture questions)
- [ ] Rebuild command documented for post-edit maintenance

## Context

- STORY-03.1 (Wiki Structure and Lint Tooling) is complete — wiki directory exists with content
- Wiki content is in `data/wiki/` directory
- Standards are in `core/docs/standards/`
- Guides are in `data/guides/` (if exists) or `core/docs/guides/`
- graphify skill is available at `~/.claude/skills/graphify/SKILL.md`

## graphify Info

graphify is a skill that builds a knowledge graph from inputs such as code, docs, papers, and images, outputting HTML visualization and JSON data.

Key outputs:
- `graphify-out/GRAPH_REPORT.md` — god nodes, community structure, insights
- `graphify-out/graph.html` — interactive visualization

## Completion Output

Write completion report to: `data/completions/STORY-03.5-completion.md`
Use template at: `core/templates/story-completion.md`
