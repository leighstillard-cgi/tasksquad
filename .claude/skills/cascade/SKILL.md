---
name: cascade
description: "Propagate a central framework file to all sibling repos"
---

# Cascade File

Propagate a central framework file from this repo to all sibling repos listed in `core/repos.manifest`.

## Arguments

The user provides:
- `FILE` -- path relative to core/ (e.g., `Code-SOP.md` or `docs/standards/testing.md`)

## Workflow

1. **Run the cascade script**:
   ```bash
   core/scripts/cascade.sh FILE
   ```

2. **Report the outcome**:
   - List which repos were updated
   - List which repos were skipped (and why)
   - Note any conflicts that need manual resolution

## Rules

- Only cascade files that originate from `core/`
- Never cascade project-specific files
- If a target repo has local modifications, report the conflict rather than overwriting
