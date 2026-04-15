---
name: end-session
description: "Write session handoff summary before ending"
---

# End Session

Write a session state summary to `.claude/session-state.md` so the next session can pick up where this one left off.

## Workflow

1. **Gather session context**:
   - What task(s) were worked on this session
   - Current status of each task (complete, in-progress, blocked)
   - Any pending issues or open questions
   - Key decisions made during the session
   - What should be picked up next

2. **Write the handoff file** to `.claude/session-state.md`:
   ```markdown
   # Session State -- YYYY-MM-DD

   ## Current Task
   <description of what was being worked on>

   ## Status
   <complete | in-progress | blocked>

   ## Pending Issues
   - <list of unresolved items>

   ## Key Decisions
   - <decisions made this session>

   ## Next Steps
   - <what to pick up next>

   ## Relevant Paths
   - <file paths relevant to the work>

   ## Workflow IDs
   - <any active workflow or issue IDs>
   ```

3. **Commit the file**:
   ```bash
   git add .claude/session-state.md
   git commit -m "session: update handoff state"
   ```

## Rules

- Always overwrite the previous session-state.md (it's a snapshot, not a log)
- Include absolute file paths, not relative ones
- Include story IDs and issue numbers where relevant
- Keep it concise -- bullet points, not prose
