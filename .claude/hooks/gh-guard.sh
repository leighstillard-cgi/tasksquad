#!/usr/bin/env bash
# gh-guard: Blocks coding agents from mutating GitHub issues.
# Only the PM agent (tasksquad project root) may close/edit/create issues.
# All repos can: gh issue view, gh issue list, gh issue comment, gh pr *

set -euo pipefail

if ! command -v jq &>/dev/null; then
    exit 0
fi

INPUT=$(cat)
CMD=$(echo "$INPUT" | jq -r '.tool_input.command // empty')

if [[ -z "$CMD" ]]; then
    exit 0
fi

# Only guard gh issue commands
echo "$CMD" | grep -qE '^\s*gh\s+issue' || exit 0

# Allow read-only and comment operations
echo "$CMD" | grep -qE 'gh\s+issue\s+(view|list|comment)' && exit 0

# Allow if PROJECT_DIR is the PM repo (tasksquad project root, not a worktree)
# The PM agent can close/edit/create issues as part of story management
PM_REPO_PATTERN="tasksquad$"
if [[ "${PROJECT_DIR:-}" =~ $PM_REPO_PATTERN ]] && [[ ! "${PROJECT_DIR:-}" =~ worktrees ]]; then
    exit 0
fi

# Block mutations (close, edit, create, reopen, delete, transfer, lock, unlock)
jq -n '{
  "hookSpecificOutput": {
    "hookEventName": "PreToolUse",
    "permissionDecision": "deny",
    "permissionDecisionReason": "Only the PM agent may close/edit/create GitHub issues. Write a completion report instead."
  }
}'
