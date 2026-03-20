---
story: STORY-XX.X
status: complete
repo: <org>/<repo-name>
branch: <branch-name>
agent: <agent-identifier>
timestamp: <ISO-8601>
---

## Summary

<One paragraph: what was built, key design choices made during implementation>

## Sub-task Evidence

- [x] <Sub-task from story spec> — <evidence: query result, test output, or file reference>
- [x] <Sub-task from story spec> — <evidence>
- [ ] <Any incomplete sub-task> — <reason>

## Verification

<Include test output, query results, or other evidence that acceptance criteria are met.
Format depends on the story — SQL results for data stories, test suite output for code stories,
screenshot or log snippet for infrastructure stories.>

## Deviations from Spec

<Any changes made during implementation that differ from the story spec.
If none, write "None — implemented as specified.">

## Architectural Escalations

<Any design decisions that should be reviewed by the architecture advisor.
If none, write "None.">

## Files Changed

- `path/to/file.go` — <what changed>
- `path/to/file_test.go` — <what changed>

## Git Reference

- Branch: `<branch-name>`
- Final commit: `<short sha>`
