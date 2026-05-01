# Support

TaskSquad is intended to be self-serve for routine onboarding and usage.

## First Checks

1. Read [docs/onboarding.md](docs/onboarding.md).
2. Read [docs/troubleshooting.md](docs/troubleshooting.md).
3. Run:

```bash
bash core/scripts/post-setup.sh --check-only
bash core/scripts/lint-wiki.sh
```

## When To Ask For Help

Ask for support when:

- Setup fails after following troubleshooting.
- A script would need elevated permissions.
- You need write access to a critical system.
- An agent proposes a production-impacting change.
- You find a possible security issue.
- The docs conflict with internal CGI policy.

## What To Include

- Operating system or container image
- Shell used
- Command run
- Relevant error output
- Whether this is a worklog repo or target repo
- Whether credentials are read-only or write-capable

Do not include secrets or customer-sensitive content in support requests.
