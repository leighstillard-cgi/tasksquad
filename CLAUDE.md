# Project Standards

> **PM Agent:** If your working directory is this worklog repository, read `PM_INSTRUCTIONS.md` — those are your primary instructions. The sections below are shared standards for repo worker agents.

## Defaults

- Prefer Go. Python 3 with strict type hints for scripting/tooling. TypeScript strict for frontend.
- AWS serverless over self-managed. Open source over proprietary at comparable quality.
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

## Workflow: Before Starting Any Story

1. Read the story spec — check `docs/story-specs/` first
2. Read `docs/adrs/` — follow locked decisions, do not re-derive
3. Read any domain-relevant documentation in `docs/` (including `docs/guides/` for conversion work)
4. Branch: `feat/<story-id>-short-description`

## Workflow: Completion Reports

When a story is complete, create `../worklog/completions/STORY-XX.X-completion.md` using the template at `../worklog/templates/story-completion.md`. The report MUST include evidence for every acceptance criterion, any deviations from spec with rationale, and any architectural escalations. The PM agent validates reports and closes issues.

## Workflow: Boundaries

Do not: make architecture decisions (escalate in the completion report), skip acceptance criteria (document why if you can't meet one), or modify worklog files other than completion reports.
