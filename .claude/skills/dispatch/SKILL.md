---
name: dispatch
description: "Dispatch a story to a subagent -- reads story spec, writes dispatch file, launches subagent with full context. Terminal-native, no Slack."
---

# Dispatch Story to Subagent

Dispatch a story to a Claude Code subagent for implementation. Handles single dispatch or batch dispatch with configurable concurrency.

## Arguments

The user provides either:
- A story ID (e.g., `STORY-05.1`)
- A batch request (e.g., "dispatch all ready STORY-02.X stories")
- A free-text description of work to dispatch

## Single Dispatch Workflow

1. **Read the story spec**: Check `data/story-specs/`, `backlog.md`, and `.client/backlog-client.md`.

2. **Verify the story is dispatchable**:
   - Status must be `ready` (not `in-progress`, `blocked`, or `done`)
   - All dependencies must be `done`
   - Not already in `data/dispatch-log.md` as `dispatched`

3. **Write the dispatch file** to `data/dispatches/STORY-XX.X-dispatch.md`:
   ```yaml
   ---
   story_id: STORY-XX.X
   dispatched_at: <ISO-8601>
   dispatched_by: pm-agent
   attempt: 1
   max_retries: 5
   ---
   ```
   Include: story title, spec location, relevant context from recent completions, dependencies completed.

4. **Update tracking**:
   - Update `data/dispatch-log.md`: add row with story ID, timestamp, status `dispatched`
   - Update `backlog.md`: set story status to `in-progress`

5. **Launch the subagent** using the Agent tool:
   ```
   Agent({
     description: "STORY-XX.X: <title>",
     subagent_type: "general-purpose",
     isolation: "worktree",
     prompt: <full context below>
   })
   ```

6. **Build the subagent prompt** with full context:
   - Story spec (full text)
   - Required reading: `CLAUDE.md`, relevant standards, patterns guide
   - Completion template path: `core/templates/story-completion.md`
   - Completion output path: `data/completions/STORY-XX.X-completion.md`
   - MCP connection details (from `data/project/tooling.md` if available)
   - Previous failure context (if this is a retry)

7. **Capture the result**:
   - On success: verify completion report was written, log to `data/session-logs/`
   - On failure: log failure reason, check retry count

8. **Handle retries** (ralph loop):
   - If subagent fails and attempt < max_retries:
     - Increment attempt counter in dispatch file
     - Inject previous failure context into next prompt
     - Re-launch subagent
   - If attempt >= max_retries:
     - Write escalation to `data/escalations/STORY-XX.X-escalation.md`
     - Update backlog: set story to `blocked`
     - Update dispatch log: set status to `failed`
     - Log: "STORY-XX.X exhausted retries -- escalation written"

9. **Write session log** to `data/session-logs/STORY-XX.X-<timestamp>.md`:
   - Dispatch file contents
   - Subagent prompt (summarized)
   - Subagent result summary
   - Completion report path or escalation path
   - Duration, exit status, attempt number

## Batch Dispatch Workflow

1. **Parse the batch request**: identify which stories to dispatch.
2. **Filter to dispatchable stories**: ready status, dependencies met, not already dispatched.
3. **Dispatch in parallel** up to the concurrency limit (default 3):
   - Launch multiple Agent calls in a single message
   - Each agent gets its own worktree isolation
4. **Process results** as they return: log, verify completions, handle retries.
5. **Report batch summary**: dispatched N stories, M succeeded, K failed, J retrying.

## Subagent Prompt Template

```
You are a coding agent working on STORY-{id}: {title}.

## Your Task
{story spec full text}

## Required Reading (already loaded)
- CLAUDE.md standards apply to all your work
- {relevant standards file}
- {patterns guide if conversion work}

## MCP Database Access
{connection details from tooling.md}

## Completion
When done, write your completion report to:
  data/completions/STORY-{id}-completion.md

Use the template at core/templates/story-completion.md.
Include evidence for EVERY acceptance criterion.

## Rules
- Do not modify backlog.md or data/dispatch-log.md
- Do not dispatch other stories
- If you hit an architectural decision, note it in "Architectural Escalations"
- If you discover new patterns, note them in "New Patterns Discovered"
- Commit and push all changes before writing the completion report
```

## Rules

- Never dispatch two stories to the same repo simultaneously
- Never dispatch stories out of dependency order
- Schema/contract changes require human approval -- write to data/escalations/ and wait
- Maximum concurrent dispatches: configurable (default 3)
- Always write dispatch file before launching subagent
- Always write session log after subagent completes
