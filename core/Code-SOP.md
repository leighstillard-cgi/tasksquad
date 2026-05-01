# Coding Agent Standards

These rules apply to all coding work across the project. Follow them in every session.

## Commit Discipline

- Atomic commits: one logical change per commit.
- Commit messages: imperative mood, under 72 chars for the subject line. Body explains WHY, not WHAT.
- Always commit and push before marking work as complete. Unpushed work does not count.
- Never force-push to main/master. Feature branches only.
- Never skip pre-commit hooks (`--no-verify`). If a hook fails, fix the issue.

## Branch Naming

- Feature: `feature/<issue-number>-<short-description>`
- Bugfix: `fix/<issue-number>-<short-description>`
- Chore: `chore/<short-description>`
- Always branch from the latest main. Rebase before PR if behind.

## Pull Requests

- Title: short, under 70 characters.
- Body: summary of changes, test plan, and link to the GitHub issue.
- One PR per issue. Do not bundle unrelated changes.
- All CI checks must pass before requesting review.

## Tool Usage

### RTK (Rust Token Killer)

RTK runs transparently via a PreToolUse hook. You do not need to invoke it manually.
If `rtk --version` fails or returns unexpected output, stop and report.

### LSP

Use LSP for code navigation: go-to-definition, find-references, hover info.
Enabled via `ENABLE_LSP_TOOL=1`. Available for Python (pyright), Go (gopls), TypeScript, .NET (OmniSharp/csharp-ls).

### jcodemunch

Use for structural code queries: symbol search, file outlines, blast radius analysis.
Must run `index_repo` before queries return results for a repo.

### jdocmunch

Use for navigating large markdown docs by section rather than reading entire files.
Must index the repo or directory first via `doc_index_repo` or `index_local`.

## Memory and Knowledge Graphs

Before asking clarifying questions, check these sources in order:

### claude-mem (cross-session memory)

Search prior session context via the `mcp__plugin_claude-mem_mcp-search__*` tools:
1. `search(query)` -- get index with IDs
2. `timeline(anchor=ID)` -- surrounding context
3. `get_observations([IDs])` -- full details

### graphify knowledge graphs

Check `./graphify-out/` for the repo-specific knowledge graph.
Start with `GRAPH_REPORT.md` for god nodes and community structure.

### When to ask the user

Only ask if the answer is genuinely not in memory or graphs. When you do ask, state what you already checked.

## Code Standards

Follow the standards in `core/docs/standards/`. Key files:
- `code-quality.md` -- structure, naming, dependencies
- `testing.md` -- test coverage requirements
- `error-handling.md` -- logging and error patterns
- `security.md` -- input validation, auth, TLS
- `workflow-discipline.md` -- how to approach work

Do not deviate from these standards without explicit user approval.

## Completion Protocol

When your task is done:
1. All tests pass.
2. All changes are committed and pushed.
3. Write a completion report to `data/completions/` in the worklog repo.
4. Include: what was done, evidence of verification, any architectural escalations.

## Project-Specific Context

Load repo-specific overrides, paths, and conventions:

```
@.claude/repo-context.md
```
