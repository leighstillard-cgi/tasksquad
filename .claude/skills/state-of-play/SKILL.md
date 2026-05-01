---
name: state-of-play
description: "Generate a state-of-play document -- queries backlog, completions, dispatch log, and session logs. Use when asked for project status."
---

# Generate State-of-Play

Generate a comprehensive state-of-play document for the current project.

## Workflow

1. **Query current state from multiple sources**:

   a. **Backlog status**: Parse `backlog.md` and `.client/backlog-client.md` (if exists) for story statuses:
   ```bash
   grep -E '^\*\*Status:\*\*' backlog.md
   ```

   b. **Dispatch log**: Parse `data/dispatch-log.md` for active dispatches.

   c. **Completions**: List recent completions in `data/archive/completions/` (last 7 days).

   d. **Escalations**: List unresolved items in `data/escalations/` and recently resolved items in `data/archive/escalations/`.

   e. **Session logs**: Parse `data/session-logs/` for recent activity, pass/fail rates.

   f. **GitHub issues** (if available):
   ```bash
   gh issue list --state open --limit 100 --json number,title,labels 2>/dev/null || echo "GitHub not available"
   gh issue list --state closed --limit 10 --json number,title,closedAt 2>/dev/null || true
   ```

   g. **Database queries** (if MCP servers available):
   Run project-specific queries to gather metrics (e.g., validation counts, conversion progress).

2. **Read the template** at `core/templates/state-of-play.md` and fill in all sections.

3. **Save** to `data/state-of-play/YYYY-MM-DD.md` using today's date.

4. **Output a summary** with key highlights:
   - What completed since last state-of-play
   - What's in progress
   - What's blocked or needs attention
   - Batch progress (if applicable: N/total validated, M failures)
   - Next actions

## Rules

- Use whatever data sources are available -- degrade gracefully if GitHub or DB is not accessible
- Always save the report to data/state-of-play/ directory
- Keep metrics factual -- do not estimate or extrapolate without stating assumptions
