---
name: bootstrap
description: "Bootstrap a sibling repo to use the TaskSquad framework"
---

# Bootstrap Repo

Bootstrap a sibling repo so it inherits the TaskSquad framework structure.

## Arguments

The user provides:
- `REPO_NAME` -- the name of the sibling repo directory (or absolute path)

## Workflow

1. **Run the bootstrap script**:
   ```bash
   core/scripts/bootstrap-repo.sh REPO_NAME
   ```

2. **Verify the result**:
   ```bash
   core/scripts/check-repo-health.sh REPO_NAME
   ```

3. **Report the outcome**:
   - If health check passes: confirm the repo is bootstrapped and list what was created
   - If health check fails: list the failing checks and suggest remediation

## Rules

- Never overwrite existing files in the target repo without user confirmation
- The bootstrap script handles all file creation -- do not manually create files
- If the target repo doesn't exist, report the error rather than creating it
