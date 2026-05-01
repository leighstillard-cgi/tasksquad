# Wiki Structure Definition

Canonical schema for all wiki pages in `data/wiki/`. This file is the
single source of truth for page types, frontmatter requirements, naming
conventions, and lint rules.

## Page Classes

| Class | Who writes | Purpose |
|---|---|---|
| `pm-only` | PM agent only | Canonical content: ADRs, epics, stories, standards |
| `agent-write` | Any coding agent | Scratch pages: drafts, completion reports, lint reports |
| `canonical-facts` | PM only | Verified infrastructure facts that hooks inject near attention cursor |

## Page Types

| Type | Class (default) | Description |
|---|---|---|
| `adr` | pm-only | Architecture Decision Records |
| `epic` | pm-only | Project epics |
| `story` | pm-only | Implementation stories |
| `concept` | pm-only | Domain knowledge concepts |
| `component` | pm-only | System components |
| `risk` | pm-only | Risk register entries |
| `standard` | pm-only | Coding/process standards |
| `runbook` | pm-only | Operational runbooks |
| `infrastructure` | canonical-facts | Verified infrastructure facts with provenance metadata |

## Required Frontmatter (all page types)

```yaml
---
title: <human title>
id: <stable id, e.g. ADR-035, EPIC-22, STORY-05.3>
status: draft|accepted|deprecated|superseded
page_class: pm-only|agent-write|canonical-facts
page_type: adr|epic|story|concept|component|risk|standard|runbook|infrastructure
tags: [...]
created: <ISO date>
last_updated: <ISO date>
supersedes: []
superseded_by: null
inbound_links: []
outbound_links: []
---
```

### Field semantics

- **`id`** --- Stable identifier. Must be unique across all wiki pages. Format
  depends on page type (see Naming Conventions below).
- **`status`** --- Lifecycle state: `draft` (work in progress), `accepted`
  (canonical), `deprecated` (no longer relevant), `superseded` (replaced by
  another page).
- **`page_class`** --- Access control. `pm-only` pages can only be written by
  the PM. `agent-write` pages can be created/updated by coding agents.
- **`supersedes`** / **`superseded_by`** --- Tracks replacement lineage between
  pages.
- **`inbound_links`** / **`outbound_links`** --- Cross-page references, kept in
  sync by the linter.

## Type-Specific Frontmatter

### `adr`

```yaml
decision: <one-line decision statement>
rationale: <why this decision was made>
constraints: <constraints that shaped the decision>
```

### `epic`

```yaml
github_issue: <issue number or URL>
labels: [...]
phase: <project phase>
```

### `story`

```yaml
github_issue: <issue number or URL>
labels: [...]
phase: <project phase>
parent_epic: <epic ID, e.g. EPIC-05>
depends_on: [<story IDs>]
repos: [<repo names>]
```

### `concept`

No additional fields.

### `component`

```yaml
repo: <repo name>
language: <primary language>
```

### `risk`

No additional fields beyond the common set.

### `standard`

No additional fields.

### `runbook`

No additional fields beyond the common set.

### `infrastructure` (page_class: canonical-facts)

```yaml
staleness_threshold_days: <integer, default 90>
certainty_basis: told-by-user|verified-in-system|inferred-from-code|inferred-from-docs
last_verified_at: <ISO date>
```

**Required fields for `canonical-facts` pages:**
- `certainty_basis` --- page-level default for certainty of facts
- `last_verified_at` --- date the page was last verified as accurate

Infrastructure pages use in-body tables with per-fact metadata columns:
- **Fact** --- the canonical fact being recorded
- **Value** --- the verified value
- **Certainty Basis** --- `told-by-user`, `verified-in-system`, `inferred-from-code`, `inferred-from-docs`
- **Last Verified** --- ISO date of last verification
- **Staleness** --- days until the fact should be re-verified

## Naming Conventions

| Type | Pattern | Example |
|---|---|---|
| `adr` | `ADR-NNN-slug.md` | `ADR-001-use-subagents-for-dispatch.md` |
| `epic` | `EPIC-NN.md` | `EPIC-03.md` |
| `story` | `STORY-XX.X.md` | `STORY-05.1.md` |
| `concept` | `lowercase-with-hyphens.md` | `stored-procedure-validation.md` |
| `component` | `lowercase-with-hyphens.md` | `dispatch-mechanism.md` |
| `risk` | `RISK-NNN-slug.md` | `RISK-001-rate-limiting.md` |
| `standard` | `slug.md` | `database-conversion.md` |
| `runbook` | `slug.md` | `batch-validation-restart.md` |
| `infrastructure` | `canonical.md` or `CANONICAL-slug.md` | `canonical.md` |

## Link Conventions

- Internal cross-page references use `[[wikilink]]` syntax.
- All paths are relative from `data/wiki/`.
- No absolute URLs for internal pages.
- Example: `[[ADR-001-use-subagents-for-dispatch]]` links to
  `data/wiki/adrs/ADR-001-use-subagents-for-dispatch.md`.

## Lint Rules

Enforced by `core/scripts/lint-wiki.sh`:

1. Every page must have all required frontmatter fields.
2. All `outbound_links` must resolve to existing files in `data/wiki/`.
3. Every page must be indexed in `wiki.md`.
4. Naming convention enforcement per page type (see table above).
5. No duplicate `id` values across all wiki pages.
6. `canonical-facts` pages must have `certainty_basis`, `last_verified_at`, and `staleness_threshold_days`.
