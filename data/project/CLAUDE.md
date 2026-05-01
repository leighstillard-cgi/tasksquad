# Project Configuration

## Tooling Reference

When a `tooling.md` exists, load it at session start:

@data/project/tooling.md

## Domain Standards

Read these before any conversion work:

- `data/standards/ase-to-mssql.md` — conversion rules and validation requirements
- `data/guides/conversion-patterns-guide.md` — known divergence patterns with before/after examples

## Dispatch Rules

- Schema, event contract, or interface changes require human approval before dispatch.
- Completion reports are verified before issues are closed.
- Maximum concurrent dispatches: configurable (default 3).

## Escalation Triggers

Flag for the user (do not resolve yourself):
- Architectural escalations in completion reports
- Deviations affecting schema or data contracts
- New ADR proposals from coding agents
- Performance regressions > 5x baseline
