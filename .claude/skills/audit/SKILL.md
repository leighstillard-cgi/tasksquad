---
name: audit
description: "Sync check across backlog, data/completions, dispatch log, and data/escalations. Flags drift. Use periodically or after processing completions."
---

# Backlog Audit

Check for drift between story states across backlog, dispatch log, data/completions, and data/escalations.

## Workflow

### Phase 1: Gather State

1. **Parse backlog**: Read `backlog.md` and `.client/backlog-client.md` (if exists). Extract all story IDs with their status.

2. **Parse dispatch log**: Read `data/dispatch-log.md`. Extract all dispatched stories with their status.

3. **Parse completions**: List all files in `data/completions/` and `data/archive/completions/`. Extract story IDs.

4. **Parse escalations**: List all files in `data/escalations/` (not archived). Extract story IDs.

5. **GitHub issues** (if available):
   ```bash
   gh issue list --state open --limit 100 --json number,title,state 2>/dev/null || echo "GitHub not available"
   gh issue list --state closed --limit 50 --json number,title,state 2>/dev/null || true
   ```

### Phase 2: Cross-Reference

6. **Diff the states**:
   - Story marked `done` in backlog but no completion report: flag as missing evidence
   - Story marked `in-progress` in backlog but not in dispatch log: flag as phantom dispatch
   - Story in dispatch log as `dispatched` but completion exists: flag as unprocessed completion
   - Completion in archive but story not marked `done`: flag as unprocessed
   - Escalation exists but story not marked `blocked`: flag as missed escalation

### Phase 3: Staleness Check

7. **Check for stalled stories**: Any story in dispatch log as `dispatched` for longer than 4 hours without a completion report.

### Phase 4: Report

8. **Output a drift report**:
   ```
   ## Audit Report -- YYYY-MM-DD

   ### Status Drift
   | Story | Backlog | Dispatch Log | Completion | Issue |
   |-------|---------|-------------|------------|-------|

   ### Stalled Stories (dispatched > 4 hours, no completion)
   - STORY-XX.X -- dispatched at <timestamp>

   ### Unprocessed Completions
   - data/completions/STORY-XX.X-completion.md

   ### Open Escalations
   - data/escalations/STORY-XX.X-escalation.md
   ```

9. **Fix the drift**: Update backlog and dispatch log to match actual state. Commit and push.

## Rules

- Completion reports are the source of truth for whether work was done
- Backlog is the source of truth for story specs and acceptance criteria
- When fixing drift, update documentation to match reality -- never the reverse
- If GitHub is available, include issue state in the audit
