---
title: "Completion: STORY-<XX.X>"
id: "COMPLETION-STORY-<XX.X>"
status: draft
page_class: agent-write
page_type: story
tags: [completion]
created: "<ISO-DATE>"
last_updated: "<ISO-DATE>"
supersedes: []
superseded_by:
inbound_links: []
outbound_links: []
github_issue: "<ISSUE-NUMBER>"
labels: [completion]
phase: "<PHASE>"
parent_epic: "EPIC-<NN>"
depends_on: []
repos: ["<REPO-NAME>"]
---

# Completion: STORY-<XX.X>

**Story:** STORY-<XX.X>
**Repo:** <REPO-NAME>
**Branch:** <BRANCH-NAME>
**Agent:** <AGENT-IDENTIFIER>
**Timestamp:** <ISO-8601>

## Summary

<One paragraph: what was built, key design choices made during implementation.>

## Sub-task Evidence

- [x] <Sub-task from story spec> -- <evidence: query result, test output, or file reference>
- [x] <Sub-task from story spec> -- <evidence>
- [ ] <Any incomplete sub-task> -- <reason>

## Verification

```sql
-- <Query from acceptance criteria>
-- Expected: <value>
-- Actual: <value>
```

## Deviations from Spec

<Any changes from the story spec. If none, write "None -- implemented as specified.">

## Architectural Escalations

<Design decisions needing review. If none, write "None.">

## New Patterns Discovered

<Any reusable implementation, testing, operational, or workflow patterns discovered. If none, write "None.">

## Files Changed

- `path/to/file` -- <what changed>
