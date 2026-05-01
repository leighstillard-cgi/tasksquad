# TaskSquad

*A framework for managing software projects with AI agents.*

---

## What This Gives You

AI coding assistants are powerful but forgetful — each conversation starts fresh. They also lack the project context that makes humans effective: what decisions were made and why, what's already been tried, what patterns work here.

TaskSquad solves this by giving your AI agents a **shared brain**: structured documentation they can read, a knowledge graph they can navigate, and a workflow system that tracks what's been done and what's next. The result is AI that works more like a junior developer who's been properly onboarded — not one who just walked in the door.

Concretely, this means:
- **Agents remember context** across sessions and handoffs
- **Work is tracked** with the same discipline you'd apply to human teams
- **Standards are enforced** because agents read them before every task
- **Mistakes don't repeat** because patterns and decisions are documented

---

## Conceptual Foundation

TaskSquad builds on two ideas:

**Andrej Karpathy's LLM Wiki** — The insight that AI agents work better when they have structured, agent-readable documentation to consult. Rather than relying on training data or lengthy prompts, agents can look things up in a wiki designed for them. TaskSquad implements this via the `data/wiki/` directory, knowledge graph (`graphify-out/`), and pillar standards.

**Traditional PM discipline, adapted for agents** — Backlogs, story specs, completion reports, and dispatch logs aren't new. What's new is using them as the coordination layer between AI agents. A PM agent reads the backlog, dispatches work to coding agents, validates their completion reports, and tracks progress — the same workflow a human PM would run, but automated.

---

## Key Features

| Feature | What It Does |
|---|---|
| **Core + Project Overlay** | Reusable framework (`core/`) + engagement-specific content (`data/project/`). Bootstrap new repos without reinventing the wheel. |
| **Knowledge Graph** | `graphify-out/` contains a semantic map of all documentation. Agents consult it before architecture questions. |
| **Skills System** | Codified workflows (`.claude/skills/`) for dispatch, completion processing, auditing, and more. |
| **Safety Hooks** | `gh-guard` restricts GitHub mutations. `canonical-infra-inject` prevents hallucinated infrastructure details. Optional prompt injection defense via lasso-security. |
| **Dispatch Lifecycle** | Stories move through `ready` → `in-progress` → `done` with completion reports that require evidence for every acceptance criterion. |
| **Cross-Session Memory** | `claude-mem` integration preserves context across conversations. |

---

## Repo Structure

```
.
|-- CLAUDE.md                 # Agent standards loaded into every session
|-- backlog.md                # Product backlog with story status
|-- core/                     # Reusable framework
|   |-- dashboard/            # Go monitoring dashboard
|   |-- docs/standards/       # Pillar standards
|   |-- templates/            # Story, completion, ADR templates
|   `-- scripts/              # Lint, bootstrap, cascade scripts
|-- data/                     # Engagement-specific content and workflow state
|   |-- dispatch-log.md       # Assignment log and completion status
|   |-- project/              # Canonical facts and tooling config
|   |-- wiki/                 # Structured documentation
|   |-- guides/               # Pattern guides and how-tos
|   |-- standards/            # Project-specific standards
|   |-- dispatches/           # Active dispatch files
|   |-- completions/          # Completion reports from agents
|   |-- escalations/          # Issues needing human review
|   |-- session-logs/         # Audit trail of agent sessions
|   `-- story-specs/          # Detailed story specifications
|-- graphify-out/             # Knowledge graph output
`-- .claude/                  # Skills, hooks, agents, and settings
```

---

## How It Works

1. **PM agent** reads the backlog, finds ready stories, writes dispatch files
2. **Coding agents** receive dispatch context, read standards, do the work
3. **Completion reports** document what was done with evidence for each criterion
4. **PM validates** the report, updates tracking, closes the story
5. **Knowledge compounds** — decisions, patterns, and context persist for future work

---

## Getting Started

**Bootstrap a new repo:**
```bash
./core/scripts/bootstrap-repo.sh /path/to/target-repo
```

**After adding your content:**
```bash
./core/scripts/post-setup.sh  # Rebuilds graphify, runs wiki lint
```

**Or manually:**
```bash
/graphify data/wiki data/guides data/standards core/docs/standards core/templates  # Build knowledge graph
./core/scripts/lint-wiki.sh                                # Validate wiki structure
# Edit data/project/data/canonical-facts.md with your infrastructure details
```

---

## Further Reading

- `CLAUDE.md` — The standards agents follow
- `core/docs/standards/workflow-discipline.md` — How to approach work
- `graphify-out/GRAPH_REPORT.md` — Knowledge graph insights
- `backlog.md` — Current project status
