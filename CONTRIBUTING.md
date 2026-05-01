# Contributing

Thank you for improving TaskSquad. Keep contributions small, reviewable, and safe for internal reuse.

## Before You Start

1. Read [README.md](README.md).
2. Read [docs/permissions-and-safety.md](docs/permissions-and-safety.md).
3. Run setup validation:

```bash
bash core/scripts/post-setup.sh --check-only
```

## Contribution Rules

- Keep framework changes generic.
- Put project-specific content under `data/` or `.client/`.
- Do not commit secrets, credentials, customer data, or proprietary client details.
- Do not add broad write access to external systems.
- Keep generated artifacts out of Git unless explicitly required.
- Update documentation when behavior changes.

## Pull Request Checklist

- [ ] Change is scoped to one logical purpose.
- [ ] Docs updated where needed.
- [ ] Wiki lint passes if wiki content changed.
- [ ] Dashboard tests pass if dashboard code changed.
- [ ] No secrets or client-sensitive content included.
- [ ] Any production-impacting behavior is clearly called out.

## Verification Commands

```bash
bash core/scripts/lint-wiki.sh
cd core/dashboard && go test ./...
```

Run the dashboard tests only when Go is installed and dashboard code changed.
