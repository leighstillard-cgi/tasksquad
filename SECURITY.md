# Security

TaskSquad is an internal agent-workflow framework. Treat it as a tool that can coordinate high-impact actions.

## Reporting Security Issues

Report security issues through the approved internal CGI security or project governance channel for this repository. Do not put vulnerability details in public issues, public chats, or broad mailing lists.

Include:

- Affected file or feature
- Impact
- Reproduction steps
- Suggested mitigation if known

## Safe Usage Requirements

- Run agents in an isolated environment.
- Use least-privilege credentials.
- Prefer read-only access to critical systems.
- Review changes that touch production systems before merge or deployment.
- Keep deployment approvals outside the agent loop unless explicitly approved by governance.
- Do not store secrets in this repository.

## Supported Scope

Security guidance applies to:

- TaskSquad scripts and hooks
- Claude Code project configuration
- Dashboard code
- Agent workflow files
- Documentation that influences agent behavior

Target application repositories remain responsible for their own application security controls.

## Secrets

Never commit:

- API keys
- Tokens
- Passwords
- Private keys
- Connection strings with credentials
- Customer data
- Production incident data unless explicitly approved and sanitized

Use approved internal secret management.

## Dashboard Exposure

The dashboard is for trusted local use. Do not expose it to untrusted networks without adding authentication, TLS, CSRF protection, and deployment hardening.
