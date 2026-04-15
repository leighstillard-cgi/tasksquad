# EPIC-03 Gap Analysis: Reference Harness vs TaskSquad

> **Status:** DRAFT - Pending user review
> **Date:** 2026-04-15
> **Story:** STORY-03.7

This document compares the reference Claude Code harness (user-level `~/.claude/` configuration) with the TaskSquad project implementation. Each gap is tagged with a recommendation and rationale.

## Legend

| Tag | Meaning |
|-----|---------|
| **must-port** | Essential functionality that TaskSquad needs to operate correctly |
| **nice-to-have** | Would improve experience but not blocking |
| **not-applicable** | Reference feature not relevant to TaskSquad's use case |
| **already-ported** | Functionality already exists in TaskSquad |

---

## 1. Skills Comparison

### User-Level Skills (Reference Harness: `~/.claude/skills/`)

| Skill | In Reference? | In TaskSquad? | Gap? | Recommendation | Rationale |
|-------|---------------|---------------|------|----------------|-----------|
| `audit-tool-routing` | Y | N | Y | nice-to-have | Audits Bash tool usage compliance. Useful for maintaining discipline but not core to TaskSquad's PM workflow. |
| `graphify` | Y | Y (doc only) | N | already-ported | STORY-03.5 documented graphify; user-level skill is available. Output in `graphify-out/`. |
| `poormansadvisor` | Y | N | Y | must-port | Critical for escalation when agents are stuck. User CLAUDE.md references this. Auto-routing to Opus/Codex is valuable for PM workflows. |
| `find-skills` | Y (symlink) | N | Y | nice-to-have | Helps discover and install skills from `npx skills` ecosystem. Not core to TaskSquad but useful for extensibility. |
| `postgres` | Y (symlink) | N | Y | nice-to-have | Read-only PostgreSQL queries. Only relevant if TaskSquad workers need DB access. Project-specific, not framework-level. |

### Project-Level Skills (TaskSquad: `.claude/skills/`)

| Skill | In Reference? | In TaskSquad? | Gap? | Recommendation | Rationale |
|-------|---------------|---------------|------|----------------|-----------|
| `dispatch` | N | Y | N | already-ported | STORY-03.3: Subagent dispatch with retry logic |
| `process-completion` | N | Y | N | already-ported | STORY-03.3: Completion report verification |
| `state-of-play` | N | Y | N | already-ported | STORY-03.3: Project status generation |
| `audit` | N | Y | N | already-ported | STORY-03.3: Backlog drift detection |
| `end-session` | N | Y | N | already-ported | STORY-03.3: Session handoff summaries |
| `lint-wiki` | N | Y | N | already-ported | STORY-03.3: Wiki linting wrapper |
| `bootstrap` | N | Y | N | already-ported | STORY-03.3: Sibling repo bootstrapping |
| `cascade` | N | Y | N | already-ported | STORY-03.3: Framework file propagation |

### Plugin-Provided Skills (superpowers-extended-cc)

| Skill | In Reference? | In TaskSquad? | Gap? | Recommendation | Rationale |
|-------|---------------|---------------|------|----------------|-----------|
| `brainstorming` | Y (plugin) | N (plugin) | N | not-applicable | Available via user-level plugin. TaskSquad uses story specs, not freeform brainstorming. |
| `writing-plans` | Y (plugin) | N (plugin) | N | not-applicable | Available via user-level plugin. TaskSquad has its own planning via story specs. |
| `executing-plans` | Y (plugin) | N (plugin) | N | not-applicable | Available via user-level plugin. |
| `systematic-debugging` | Y (plugin) | N (plugin) | N | not-applicable | Available via user-level plugin. |
| `test-driven-development` | Y (plugin) | N (plugin) | N | not-applicable | Available via user-level plugin. Coding agents can use it. |
| `verification-before-completion` | Y (plugin) | N (plugin) | N | not-applicable | Available via user-level plugin. |
| `dispatching-parallel-agents` | Y (plugin) | N (plugin) | N | not-applicable | TaskSquad has its own dispatch skill with project-specific logic. |
| `using-git-worktrees` | Y (plugin) | N (plugin) | N | not-applicable | Available via user-level plugin. |
| `subagent-driven-development` | Y (plugin) | N (plugin) | N | not-applicable | Available via user-level plugin. |
| `requesting-code-review` | Y (plugin) | N (plugin) | N | not-applicable | Available via user-level plugin. |
| `receiving-code-review` | Y (plugin) | N (plugin) | N | not-applicable | Available via user-level plugin. |
| `finishing-a-development-branch` | Y (plugin) | N (plugin) | N | not-applicable | Available via user-level plugin. |
| `writing-skills` | Y (plugin) | N (plugin) | N | not-applicable | Available via user-level plugin. |

### Plugin-Provided Skills (claude-mem)

| Skill | In Reference? | In TaskSquad? | Gap? | Recommendation | Rationale |
|-------|---------------|---------------|------|----------------|-----------|
| `mem-search` | Y (plugin) | Y (doc) | N | already-ported | STORY-03.6 documented; available via user-level plugin. |
| `timeline-report` | Y (plugin) | Y (doc) | N | already-ported | STORY-03.6 documented. |
| `smart-explore` | Y (plugin) | Y (doc) | N | already-ported | STORY-03.6 documented. |
| `make-plan` | Y (plugin) | Y (doc) | N | already-ported | STORY-03.6 documented. |
| `do` | Y (plugin) | Y (doc) | N | already-ported | STORY-03.6 documented. |

### Plugin-Provided Skills (codex)

| Skill | In Reference? | In TaskSquad? | Gap? | Recommendation | Rationale |
|-------|---------------|---------------|------|----------------|-----------|
| `codex:rescue` | Y (plugin) | N (doc) | Y | nice-to-have | Delegates investigation to Codex. Useful for complex debugging but not core to PM workflow. |
| `codex:setup` | Y (plugin) | N (doc) | Y | nice-to-have | Codex CLI setup helper. |
| `gpt-5-4-prompting` | Y (plugin) | N | N | not-applicable | Internal plugin guidance. |
| `codex-cli-runtime` | Y (plugin) | N | N | not-applicable | Internal plugin helper. |
| `codex-result-handling` | Y (plugin) | N | N | not-applicable | Internal plugin guidance. |

---

## 2. Hooks Comparison

### User-Level Hooks (Reference Harness: `~/.claude/hooks/`)

| Hook | In Reference? | In TaskSquad? | Gap? | Recommendation | Rationale |
|------|---------------|---------------|------|----------------|-----------|
| `canonical-facts-inject.sh` | Y | N | Y | not-applicable | SessionStart hook for data-worklog only. TaskSquad uses canonical-infra-inject pattern instead. |
| `canonical-infra-inject.sh` | Y | Y | N | already-ported | STORY-03.4 ported. Prevents infrastructure hallucination. |
| `canonical-infra-inject.config` | Y | Y | N | already-ported | STORY-03.4 ported. |
| `gh-guard.sh` | Y | Y | N | already-ported | STORY-03.4 ported. Blocks coding agents from mutating GitHub issues. |
| `rtk-rewrite.sh` | Y | N | Y | not-applicable | RTK (Rust Token Killer) integration. User-level optimization, not TaskSquad-specific. |
| `slack-guard.sh` | Y | N | Y | nice-to-have | Blocks Slack write operations via MCP. Only needed if TaskSquad agents use Slack MCP tools. |
| `duncemode-detect.sh` | Y | N | Y | nice-to-have | UserPromptSubmit hook for skepticism mode. Quality-of-life feature, not core. |

### Project-Level Hooks (TaskSquad: `.claude/hooks/`)

| Hook | In Reference? | In TaskSquad? | Gap? | Recommendation | Rationale |
|------|---------------|---------------|------|----------------|-----------|
| `gh-guard.sh` | Y | Y | N | already-ported | STORY-03.4: Project-adapted version for TaskSquad. |
| `canonical-infra-inject.sh` | Y | Y | N | already-ported | STORY-03.4: Template-ready for project configuration. |
| `canonical-infra-inject.config` | Y | Y | N | already-ported | STORY-03.4: Pattern configuration file. |

---

## 3. Templates Comparison

| Template | In Reference? | In TaskSquad? | Gap? | Recommendation | Rationale |
|----------|---------------|---------------|------|----------------|-----------|
| `adr.md` | N | Y | N | already-ported | STORY-03.2: Architecture decision record template. |
| `component.md` | N | Y | N | already-ported | STORY-03.2: System component documentation. |
| `concept.md` | N | Y | N | already-ported | STORY-03.2: Domain concept template. |
| `draft.md` | N | Y | N | already-ported | STORY-03.2: Draft page template. |
| `epic.md` | N | Y | N | already-ported | STORY-03.2: Epic specification template. |
| `lint-report.md` | N | Y | N | already-ported | STORY-03.2: Wiki lint report template. |
| `pipeline-stage.md` | N | Y | N | already-ported | STORY-03.2: CI/CD stage documentation. |
| `runbook.md` | N | Y | N | already-ported | STORY-03.2: Operational procedure template. |
| `state-of-play.md` | N | Y | N | already-ported | STORY-03.2: Status report template. |
| `story-completion.md` | N | Y | N | already-ported | STORY-03.2: Completion report template. |
| `story.md` | N | Y | N | already-ported | STORY-03.2: Story specification template. |

---

## 4. Standards/Pillar Docs Comparison

| Standard | In Reference? | In TaskSquad? | Gap? | Recommendation | Rationale |
|----------|---------------|---------------|------|----------------|-----------|
| `code-quality.md` | N | Y | N | already-ported | STORY-03.2: Coding standards. |
| `testing.md` | N | Y | N | already-ported | STORY-03.2: Testing standards. |
| `error-handling.md` | N | Y | N | already-ported | STORY-03.2: Error handling and logging. |
| `security.md` | N | Y | N | already-ported | STORY-03.2: Security standards. |
| `database.md` | N | Y | N | already-ported | STORY-03.2: Database practices. |
| `workflow-discipline.md` | N | Y | N | already-ported | STORY-03.2: Agent workflow discipline. |

---

## 5. Wiki Features Comparison

| Feature | In Reference? | In TaskSquad? | Gap? | Recommendation | Rationale |
|---------|---------------|---------------|------|----------------|-----------|
| Wiki directory structure | Y (data-worklog) | Y | N | already-ported | STORY-03.1: Full wiki hierarchy. |
| Frontmatter validation | Y (data-worklog) | Y | N | already-ported | STORY-03.1: lint-wiki.sh validates YAML frontmatter. |
| Internal link checking | Y (data-worklog) | Y | N | already-ported | STORY-03.1: lint-wiki-helper.py checks links. |
| Wiki index page | Y (data-worklog) | Y | N | already-ported | STORY-03.1: `wiki/wiki.md` index. |
| Manual generation | Y (data-worklog) | Y | N | already-ported | STORY-03.1: generate-manuals.sh script. |

---

## 6. Scripts Comparison

| Script | In Reference? | In TaskSquad? | Gap? | Recommendation | Rationale |
|--------|---------------|---------------|------|----------------|-----------|
| `lint-wiki.sh` | Y | Y | N | already-ported | STORY-03.1: Wiki linting. |
| `lint-wiki-helper.py` | Y | Y | N | already-ported | STORY-03.1: Python-based frontmatter and link validation. |
| `generate-manuals.sh` | Y | Y | N | already-ported | STORY-03.1: Concatenates wiki sections into manuals. |
| `bootstrap-repo.sh` | N | Y | N | already-ported | STORY-03.2: Bootstraps sibling repos. |
| `cascade.sh` | N | Y | N | already-ported | STORY-03.2: Propagates framework files. |
| `check-repo-health.sh` | N | Y | N | already-ported | STORY-03.2: Validates repo structure. |

---

## 7. Agents Comparison

| Agent | In Reference? | In TaskSquad? | Gap? | Recommendation | Rationale |
|-------|---------------|---------------|------|----------------|-----------|
| `backlog-auditor` | N | Y | N | already-ported | STORY-03.3: Background Haiku agent for drift detection. |
| PM agent patterns | Y (data-worklog) | Y | N | already-ported | TaskSquad IS the PM framework; patterns embedded in skills. |

---

## 8. Dispatch Patterns Comparison

| Pattern | In Reference? | In TaskSquad? | Gap? | Recommendation | Rationale |
|---------|---------------|---------------|------|----------------|-----------|
| Worktree isolation | Y | Y | N | already-ported | Dispatch skill uses `isolation: "worktree"`. |
| Retry with context injection | Y | Y | N | already-ported | Dispatch skill implements ralph loop (max 5 retries). |
| Max concurrency limit | Y | Y | N | already-ported | Configurable (default 3). |
| Dispatch file tracking | Y | Y | N | already-ported | Writes to `dispatches/`. |
| Session logging | Y | Y | N | already-ported | Writes to `session-logs/`. |

---

## 9. Completion Workflows Comparison

| Workflow | In Reference? | In TaskSquad? | Gap? | Recommendation | Rationale |
|----------|---------------|---------------|------|----------------|-----------|
| Completion report template | Y | Y | N | already-ported | `core/templates/story-completion.md` |
| Evidence verification | Y | Y | N | already-ported | process-completion skill checks criteria. |
| Archival workflow | Y | Y | N | already-ported | Moves to `completions/archive/`. |
| Backlog status update | Y | Y | N | already-ported | Updates story status to done. |
| Dispatch log update | Y | Y | N | already-ported | Updates dispatch status to complete. |
| Pattern extraction | Y | Y | N | already-ported | Extracts new patterns to guides. |

---

## 10. Escalation Handling Comparison

| Pattern | In Reference? | In TaskSquad? | Gap? | Recommendation | Rationale |
|---------|---------------|---------------|------|----------------|-----------|
| Escalation file creation | Y | Y | N | already-ported | Writes to `escalations/`. |
| Architectural escalation flagging | Y | Y | N | already-ported | Completion template includes section. |
| Human approval gates | Y | Y | N | already-ported | Schema/contract changes require approval. |
| Retry exhaustion handling | Y | Y | N | already-ported | Dispatch skill marks `blocked` after max retries. |

---

## 11. Configuration Patterns Comparison

| Pattern | In Reference? | In TaskSquad? | Gap? | Recommendation | Rationale |
|---------|---------------|---------------|------|----------------|-----------|
| Project settings.json | Y | Y | N | already-ported | `.claude/settings.json` with schema ref. |
| Permission allowlists | Y | Y | N | already-ported | Git commands allowed in project settings. |
| Hook configuration | Y | Y | N | already-ported | PreToolUse hooks registered. |
| Environment variables | Y (user) | N (project) | Y | nice-to-have | User-level has many env vars. Project-level empty. |

---

## 12. Plugin Configuration

| Plugin | In Reference? | In TaskSquad? | Gap? | Recommendation | Rationale |
|--------|---------------|---------------|------|----------------|-----------|
| pyright-lsp | Y (enabled) | N | N | not-applicable | User-level plugin for Python LSP. |
| gopls-lsp | Y (enabled) | N | N | not-applicable | User-level plugin for Go LSP. |
| typescript-lsp | Y (enabled) | N | N | not-applicable | User-level plugin for TypeScript LSP. |
| superpowers-extended-cc | Y (enabled) | N | N | not-applicable | User-level plugin. Skills available to all projects. |
| claude-mem | Y (enabled) | Y (doc) | N | already-ported | STORY-03.6 documented configuration. |
| codex | Y (enabled) | N | N | not-applicable | User-level plugin for Codex integration. |

---

## Recommendations Summary

### Must-Port (1 item)

| Item | Category | Rationale |
|------|----------|-----------|
| `poormansadvisor` skill | Skills | Referenced in user CLAUDE.md for escalation. Auto-routes Sonnet/Haiku to Opus, Opus to Codex. Critical for agents that get stuck. |

### Nice-to-Have — DECLINED (Not Needed)

The following items were identified as potentially useful but have been explicitly declined as not needed for TaskSquad:

| Item | Category | Why Not Needed |
|------|----------|----------------|
| `audit-tool-routing` skill | Skills | Tool discipline enforced via standards, not automated auditing |
| `find-skills` skill | Skills | TaskSquad has a fixed skill set; discovery not required |
| `postgres` skill | Skills | Database access is project-specific, not framework-level |
| `slack-guard.sh` hook | Hooks | Slack write ops handled by cc-connect, not MCP tools |
| `duncemode-detect.sh` hook | Hooks | Skepticism mode not part of TaskSquad workflow |
| `codex:rescue` skill | Skills | poormansadvisor covers escalation needs |

### Not-Applicable (19 items)

Available at user/plugin level — no project-level porting needed:

**superpowers-extended-cc skills (12):**
- `brainstorming` — user-level plugin
- `writing-plans` — user-level plugin
- `executing-plans` — user-level plugin
- `systematic-debugging` — user-level plugin
- `test-driven-development` — user-level plugin
- `verification-before-completion` — user-level plugin
- `dispatching-parallel-agents` — TaskSquad has own dispatch skill
- `using-git-worktrees` — user-level plugin
- `subagent-driven-development` — user-level plugin
- `requesting-code-review` — user-level plugin
- `receiving-code-review` — user-level plugin
- `finishing-a-development-branch` — user-level plugin
- `writing-skills` — user-level plugin

**codex plugin internals (3):**
- `gpt-5-4-prompting` — internal plugin guidance
- `codex-cli-runtime` — internal plugin helper
- `codex-result-handling` — internal plugin guidance

**User-level plugins (3):**
- `pyright-lsp` — Python LSP, user-level
- `gopls-lsp` — Go LSP, user-level
- `typescript-lsp` — TypeScript LSP, user-level

**User-level hooks (1):**
- `canonical-facts-inject.sh` — data-worklog specific, TaskSquad uses canonical-infra-inject instead

### Already-Ported (38 items)

EPIC-03 successfully ported the complete PM framework:

**STORY-03.1 — Wiki Infrastructure (5):**
- Wiki directory structure
- Frontmatter validation (lint-wiki.sh)
- Internal link checking (lint-wiki-helper.py)
- Wiki index page (wiki/wiki.md)
- Manual generation (generate-manuals.sh)

**STORY-03.2 — Core Framework (17):**
- Core/project overlay pattern
- `adr.md` template
- `component.md` template
- `concept.md` template
- `draft.md` template
- `epic.md` template
- `lint-report.md` template
- `pipeline-stage.md` template
- `runbook.md` template
- `state-of-play.md` template
- `story-completion.md` template
- `story.md` template
- `code-quality.md` standard
- `testing.md` standard
- `error-handling.md` standard
- `security.md` standard
- `database.md` standard
- `workflow-discipline.md` standard
- `bootstrap-repo.sh` script
- `cascade.sh` script
- `check-repo-health.sh` script

**STORY-03.3 — Skills and Agents (9):**
- `dispatch` skill
- `process-completion` skill
- `state-of-play` skill
- `audit` skill
- `end-session` skill
- `lint-wiki` skill
- `bootstrap` skill
- `cascade` skill
- `backlog-auditor` agent

**STORY-03.4 — Hooks (3):**
- `gh-guard.sh` hook
- `canonical-infra-inject.sh` hook
- `canonical-infra-inject.config`

**STORY-03.5 — Knowledge Graph (1):**
- graphify documentation and output (`graphify-out/`)

**STORY-03.6 — Cross-Session Memory (5):**
- claude-mem documentation (`project/tooling.md`)
- `mem-search` skill reference
- `timeline-report` skill reference
- `smart-explore` skill reference
- `make-plan` / `do` skill references

---

## Decision Log

| Date | Item | Decision | Notes |
|------|------|----------|-------|
| 2026-04-15 | All nice-to-have items | DECLINED | Not needed for TaskSquad — poormansadvisor is the only must-port |
| 2026-04-15 | `poormansadvisor` skill | TO PORT | Will be included in STORY-03.8 or as follow-up story |

---

## Appendix: Reference Harness Locations

| Location | Purpose |
|----------|---------|
| `~/.claude/CLAUDE.md` | User-level instructions (advisor escalation, Slack patterns) |
| `~/.claude/skills/` | User-level custom skills |
| `~/.agents/skills/` | Shared skills (symlinked) |
| `~/.claude/hooks/` | User-level hooks |
| `~/.claude/settings.json` | User-level configuration |
| `~/.claude/plugins/` | Installed plugins (superpowers, claude-mem, codex) |
