# TaskSquad

TaskSquad is a framework for coordinating software delivery with AI agents. It gives agents a shared operating model: structured project context, standards, story dispatch, completion evidence, audit trails, and generated knowledge graphs.

The goal is not to replace engineering judgement. The goal is to make agent-assisted work repeatable, reviewable, and safer.

## Start Here

New adopters should read these in order:

1. [Onboarding Guide](docs/onboarding.md)
2. [Permissions and Safety](docs/permissions-and-safety.md)
3. [Troubleshooting](docs/troubleshooting.md)
4. [FAQ](docs/faq.md)

## Strong Safety Recommendation

Run TaskSquad and Claude Code in a container, VM, devcontainer, cloud workspace, or another similarly isolated environment.

Agentic coding tools are much faster when allowed to run commands, edit files, and use broad tool permissions, including dangerous or permission-bypass modes. Those same permissions can damage local files, push unwanted code, spend money, or modify systems if credentials are too powerful. Isolation gives you speed without handing the agent your whole machine.

For critical systems:

- Give agents read-only access wherever possible.
- Use least-privilege credentials and scoped tokens.
- Do not give write access to production infrastructure by default.
- Review every change that touches production systems before it is merged or deployed.
- Keep secrets out of the repository, prompts, logs, completion reports, and screenshots.

## Quick Start

Prerequisites:

- Git
- Bash-compatible shell, such as Linux, macOS, WSL, Git Bash, or a devcontainer
- Node.js for Claude Code and plugins
- Python 3 with pip
- Claude Code CLI
- Optional: Go for the dashboard, Rust for RTK

Install and verify:

```bash
bash core/scripts/install.sh
bash core/scripts/post-setup.sh
```

`install.sh` checks prerequisites, installs or verifies the expected agent tooling, and attempts the first graphify knowledge graph generation when a shell runner is available. `post-setup.sh` reruns validation after you add project content.

If graphify is installed but `graphify-out/GRAPH_REPORT.md` is still missing, open Claude Code in this repository and run:

```text
/graphify data/wiki core/docs/standards core/templates docs --update
```

See [Troubleshooting](docs/troubleshooting.md) for graphify, plugin, Windows/WSL, and permission issues.

## First Project Setup

1. Clone TaskSquad into an isolated workspace.
2. Run the install and post-setup commands above.
3. Edit [data/project/data/canonical-facts.md](data/project/data/canonical-facts.md) with non-secret project facts.
4. Add project-specific wiki pages under `data/wiki/`.
5. Add initial stories to [backlog.md](backlog.md) or to a gitignored `.client/backlog-client.md`.
6. Bootstrap a target code repository:

```bash
bash core/scripts/bootstrap-repo.sh /path/to/target-repo
bash core/scripts/check-repo-health.sh /path/to/target-repo
```

7. Start Claude Code in this TaskSquad worklog repository for PM work, or in the target repository for implementation work.

## Operating Model

1. The PM agent reads [backlog.md](backlog.md), chooses ready work, and writes dispatch files under `data/dispatches/`.
2. Worker agents read the dispatch, standards, project context, and target repository context.
3. Workers implement the change and write completion reports under `data/completions/`.
4. The PM validates completion evidence, updates `data/dispatch-log.md`, and archives completed work.
5. Important decisions and reusable knowledge are documented in `data/wiki/`.

Used workflow records should not stay in active directories. Move them into `data/archive/` or document the durable learning in the wiki.

## Repository Layout

```text
.
|-- CLAUDE.md                 # Agent standards loaded into every session
|-- PM_INSTRUCTIONS.md         # PM agent operating instructions
|-- backlog.md                # Active framework backlog
|-- docs/                     # Human onboarding and support docs
|-- core/                     # Reusable framework
|   |-- dashboard/            # Go monitoring dashboard
|   |-- docs/standards/       # Pillar standards
|   |-- templates/            # Story, completion, ADR, wiki templates
|   `-- scripts/              # Install, lint, bootstrap, health scripts
|-- data/                     # Project overlay and workflow state
|   |-- archive/              # Completed or historical work records
|   |-- dispatch-log.md       # Active assignment log
|   |-- project/              # Project-specific config and canonical facts
|   |-- wiki/                 # Structured project documentation
|   |-- dispatches/           # Active dispatch files
|   |-- completions/          # Unprocessed completion reports
|   |-- escalations/          # Open issues needing human review
|   `-- state-of-play/        # Generated status reports
|-- graphify-out/             # Generated, ignored knowledge graph output
`-- .claude/                  # Claude Code skills, hooks, agents, settings
```

## Permissions and Data Handling

The project-level Claude settings allow common Git operations and run safety hooks before tool use. This is intentional: the PM workflow needs to create dispatches, update logs, and commit changes.

Important boundaries:

- `gh-guard` limits GitHub issue mutations.
- `canonical-infra-inject` injects known infrastructure facts when relevant patterns are detected.
- `.client/` is gitignored for client-sensitive content.
- `graphify-out/`, generated lint reports, and generated manuals are gitignored.
- The local dashboard is intended for trusted local use. Do not expose it to the public internet without adding authentication, TLS, and deployment hardening.

Read [Permissions and Safety](docs/permissions-and-safety.md) before running agents with broad permissions.

## Knowledge Graph

`graphify-out/` is generated and intentionally ignored by Git. A fresh clone will not have the graph until setup runs.

Use the generated graph when available:

- `graphify-out/GRAPH_REPORT.md` for a text summary
- `graphify-out/graph.html` for visual exploration

Regenerate after changing wiki pages, standards, templates, or onboarding docs:

```bash
bash core/scripts/post-setup.sh
```

If the shell script cannot generate the graph directly, run the `/graphify ... --update` command from Claude Code as shown in Quick Start.

## Further Reading

- [CONTRIBUTING.md](CONTRIBUTING.md)
- [SECURITY.md](SECURITY.md)
- [SUPPORT.md](SUPPORT.md)
- [CHANGELOG.md](CHANGELOG.md)
- [core/docs/standards/workflow-discipline.md](core/docs/standards/workflow-discipline.md)
- [core/docs/standards/mcp-safety.md](core/docs/standards/mcp-safety.md)
