# Project Standards

> **PM Agent:** If your working directory is this worklog repository, read `PM_INSTRUCTIONS.md` — those are your primary instructions. The sections below are shared standards for repo worker agents.

## Framework

This repo uses a **core + project overlay** pattern:
- `core/` — Reusable PM framework (templates, scripts, standards, wiki schema). Generic across projects.
- `project/` — Engagement-specific content (tooling config, domain knowledge). Gitignored content lives in `.client/`.
- `wiki/` — Structured documentation with YAML frontmatter, linted by `core/scripts/lint-wiki.sh`.

## Defaults

- Prefer Go. Python 3 with strict type hints for scripting/tooling. TypeScript strict for frontend. .NET/C# for target platform work.
- Open source over proprietary at comparable quality.
- OpenTelemetry for observability. Structured JSON logging only.
- All database writes: `INSERT ON CONFLICT DO UPDATE` for idempotency unless stated otherwise.

## Behavioural Rules

- Propose test cases before writing implementation. Get approval.
- Before modifying existing code, run existing tests to confirm a passing baseline.
- Never skip error handling, logging, validation, or tests — regardless of time pressure.
- If a required external dependency is missing, stop and flag it. No placeholders or local substitutes.
- Conventional commits. Small, reviewable changesets. Touch only what the task requires.
- State assumptions before acting. When ambiguous, ask — do not guess silently.
- Minimum code that solves the problem. No speculative features or premature abstractions.
- Close the loop: define success criteria before starting, run tests to confirm, do not mark done until verified.

## Critical Rules

**Security** — No hardcoded secrets, ever. Parameterised queries only. Input validation at every boundary (allowlists over denylists). Centralised auth middleware, never per-handler.

**Data privacy** — All data access scoped by tenant. PII annotated in schema, masked in logs. Soft delete with retention hooks, no hard deletes without explicit policy.

**Testing** — No feature code without tests. Cover: happy path, edge cases, error conditions, auth enforcement, tenant isolation. Test that user A cannot access user B's data.

**Error handling** — No swallowed exceptions. User-facing errors: generic and safe. Internal errors: specific and actionable. Every log entry: structured JSON, correlation ID, timestamp, service name.

**MCP safety** — Read-only by default. Write/mutate/destroy requires explicit user approval and logging. Never create/modify/delete cloud resources without confirmation of environment and estimated cost.

## Knowledge Graph

A knowledge graph is maintained at `graphify-out/` indexing wiki, standards, guides, and templates.

**Before architecture questions:** Read `graphify-out/GRAPH_REPORT.md` to understand how concepts connect across the codebase. The God Nodes section shows core abstractions. The Surprising Connections section reveals non-obvious relationships.

**Interactive exploration:** Open `graphify-out/graph.html` in a browser to visually navigate the knowledge graph.

**Rebuild after edits:** When wiki, standards, or guides are modified, rebuild the graph:
```bash
/graphify wiki core/docs/standards guides core/templates --update
```

## Workflow: Before Starting Any Story

1. Read the story spec — check `story-specs/` first, then `.client/` if applicable
2. Read `adrs/` — follow locked decisions, do not re-derive
3. Read the relevant `core/docs/standards/` pillar doc if working in that domain
4. Read any domain-relevant documentation in `guides/`
5. Consult `graphify-out/GRAPH_REPORT.md` for cross-cutting concerns and related concepts
6. Branch: `feature/<story-id>-short-description`

## Workflow: Completion Reports

When a story is complete, create `completions/STORY-XX.X-completion.md` using the template at `core/templates/story-completion.md`. The report MUST include evidence for every acceptance criterion, any deviations from spec with rationale, and any architectural escalations. The PM agent validates reports and closes issues.

## Workflow: Boundaries

Do not: make architecture decisions (escalate in the completion report), skip acceptance criteria (document why if you can't meet one), or modify worklog files other than completion reports.

## Standards

Detailed standards are in `core/docs/standards/`:

@core/docs/standards/code-quality.md
@core/docs/standards/testing.md
@core/docs/standards/error-handling.md
@core/docs/standards/security.md
@core/docs/standards/database.md
@core/docs/standards/workflow-discipline.md

## Skills

Available skills (invoke via `/skill-name`):

- **dispatch** — Dispatch stories to subagents with worktree isolation, retry logic, and session logging
- **process-completion** — Validate and archive completion reports, update backlog and dispatch log
- **state-of-play** — Generate a comprehensive status report across all tracking files
- **audit** — Cross-reference backlog, dispatch log, completions, and escalations for drift
- **end-session** — Write session handoff state to `.claude/session-state.md`
- **lint-wiki** — Run wiki linter and report structural issues
- **bootstrap** — Bootstrap a sibling repo to use the TaskSquad framework
- **cascade** — Propagate a core framework file to all sibling repos

## Agents

- **backlog-auditor** — Background agent (Haiku) that audits story states across tracking files. Read-only.

## Tracking Directories

- `dispatches/` — Dispatch files for in-progress stories
- `completions/` — Completion reports from subagents (move to `completions/archive/` when processed)
- `escalations/` — Escalation reports for blocked/failed stories
- `session-logs/` — Audit trail of subagent dispatch sessions
- `state-of-play/` — Generated status reports

## Hooks

This project uses Claude Code hooks for safety controls. Hooks are configured in `.claude/settings.json`.

### Project-Level Hooks

**gh-guard** (PreToolUse on Bash)
Guards GitHub issue mutations. Only the PM agent (running in the tasksquad project root) can close, edit, or create issues. Coding agents dispatched to worktrees can view and comment but must write completion reports rather than closing issues directly.

**canonical-infra-inject** (PreToolUse on Bash, Edit, Write, Read)
Injects canonical infrastructure facts when operations matching infrastructure patterns are detected. Prevents hallucination of hostnames, account IDs, and connection details. Configure patterns in `.claude/hooks/canonical-infra-inject.config` and populate facts in `project/data/canonical-facts.md`.

### User-Level Hooks (Not Project-Level)

**rtk-rewrite** (PreToolUse on Bash)
Rewrites shell commands to use RTK for token savings. This is a user-level hook - users need RTK installed (`cargo install rtk` or see https://github.com/rtk-ai/rtk). Configure in `~/.claude/settings.json`:

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          {
            "type": "command",
            "command": "~/.claude/hooks/rtk-rewrite.sh"
          }
        ]
      }
    ]
  }
}
```

### Prompt Injection Defense

For environments processing untrusted content (files from external sources, web scraping, user uploads), consider adding the lasso-security/claude-hooks prompt injection defender:

```bash
# Clone and install
git clone https://github.com/lasso-security/claude-hooks.git
cd claude-hooks && ./install.sh /path/to/project
```

This adds PostToolUse hooks that scan tool outputs for injection patterns (instruction overrides, role-playing/DAN attempts, encoding obfuscation, context manipulation). See the lasso-security repo for configuration options.

## Project-Specific Configuration

@project/CLAUDE.md
