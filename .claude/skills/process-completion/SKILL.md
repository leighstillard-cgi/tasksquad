---
name: process-completion
description: "Process completion reports from completions/ -- verify evidence against acceptance criteria, archive file, commit and push. Use when new completion files appear."
---

# Process Completion Report

Process completion reports from the `completions/` directory.

## Workflow

1. **List new completions**: `ls completions/ | grep -v archive | grep -v gitkeep`
2. **For each completion file**, read it and extract:
   - `story` field from frontmatter -- maps to a story in `backlog.md` or `.client/backlog-client.md`
   - `status` field -- complete, partial, or escalation
   - Sub-task evidence -- checklist items with [x] or [ ]
   - Deviations from spec
   - Architectural escalations
   - New patterns discovered
3. **Find the matching story spec**: check `story-specs/` and `backlog.md`
4. **Verify evidence against acceptance criteria**:
   - Read the story spec or backlog entry for acceptance criteria
   - For each criterion, check if the completion report provides evidence
   - Evidence = query results, test outputs, file references, counts
5. **If ALL criteria have evidence AND status is complete**:
   - Update `backlog.md`: mark story as `done`
   - Update `dispatch-log.md`: set status to `complete`
   - Move the completion file to `completions/archive/`
   - If new patterns discovered: extract and add to `guides/conversion-patterns-guide.md`
   - `git add . && git commit -m "completion: STORY-XX.X processed" && git push`
6. **If ANY criterion lacks evidence OR status is partial**:
   - Write an escalation to `escalations/STORY-XX.X-escalation.md`
   - Do NOT mark the story as done
   - Flag: "STORY-XX.X completion is missing evidence for: <list>"
7. **If architectural escalations exist**:
   - Copy them to `escalations/` for human review
   - Do NOT mark done until escalation is resolved

## After processing

Report a summary of what was processed and any items needing attention.

## Rules

- Always commit and push after archiving
- If multiple completion files exist, process them all in one pass
- Never mark a story done if evidence is missing for any acceptance criterion
