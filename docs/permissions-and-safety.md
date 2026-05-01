# Permissions and Safety

TaskSquad is designed for capable agents operating quickly. That means the environment must be designed defensively.

## Recommended Isolation

Run agents in one of these:

- Devcontainer
- Disposable VM
- Dedicated WSL distro
- Cloud workstation
- Container with only the needed repositories mounted

Do not run broad-permission agents from a shell that has access to unrelated repositories, personal files, unrestricted cloud credentials, or production write credentials.

## Why Dangerous Modes Need Isolation

Agent workflows can be slow when every command, file edit, and Git operation requires manual confirmation. Teams often enable broader permissions to keep delivery moving.

Broader permissions are useful, but they change the risk profile. An agent can:

- Edit or delete files.
- Run commands with side effects.
- Commit and push changes.
- Use credentials already available in the shell.
- Make incorrect assumptions quickly.

Isolation limits the blast radius when mistakes happen.

## Critical System Rules

- Prefer read-only access to production and critical systems.
- Do not give agents default write access to production infrastructure.
- Use scoped, short-lived credentials.
- Confirm the active cloud account, region, subscription, tenant, and environment before any infrastructure operation.
- Review every code change that touches production systems before merge or deployment.
- Require human approval before applying infrastructure changes.
- Keep deployment and release approvals outside the agent loop unless your governance model explicitly permits it.

## Repository Permissions

This project's `.claude/settings.json` allows common Git commands so the PM agent can maintain dispatch logs, completion reports, and backlog state.

The safety hooks are guardrails, not guarantees:

- `gh-guard` restricts GitHub issue mutations.
- `canonical-infra-inject` adds known facts to reduce hallucinated infrastructure values.

You are still responsible for reviewing commits and pull requests.

## Dashboard Safety

The dashboard is intended for trusted local use.

Do not expose it outside a trusted development environment unless you add:

- Authentication
- TLS
- CSRF protection
- Network restrictions
- Deployment logging and monitoring

## Secrets

Never put secrets in:

- `data/project/data/canonical-facts.md`
- `.client/`
- Prompts
- Completion reports
- Dispatch files
- Wiki pages
- Screenshots
- Graph outputs

Use secret managers, environment variables, or approved internal credential workflows.
