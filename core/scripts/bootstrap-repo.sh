#!/usr/bin/env bash
#
# bootstrap-repo.sh — Sets up a sibling repo to use the TaskSquad framework.
# Creates CLAUDE.md symlink, .claude/repo-context.md skeleton, and
# docs/upstream/ symlink to core docs. Idempotent.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WORKLOG_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

usage() {
  echo "Usage: $(basename "$0") REPO_PATH"
  echo ""
  echo "REPO_PATH can be:"
  echo "  - A bare name (looked up as sibling of this worklog)"
  echo "  - An absolute path to the repo"
  exit 1
}

[[ $# -lt 1 ]] && usage

REPO_NAME="$1"

if [[ "$REPO_NAME" = /* ]]; then
  REPO_DIR="$REPO_NAME"
else
  REPO_DIR="$(dirname "$WORKLOG_ROOT")/$REPO_NAME"
fi

if [[ ! -d "$REPO_DIR" ]]; then
  echo "ERROR: Repo directory does not exist: $REPO_DIR"
  exit 1
fi

echo "Bootstrapping: $REPO_DIR"

# 1. CLAUDE.md symlink -> core/Code-SOP.md
CLAUDE_MD="$REPO_DIR/CLAUDE.md"
TARGET_SOP="$WORKLOG_ROOT/core/Code-SOP.md"

if [[ -L "$CLAUDE_MD" ]]; then
  CURRENT_TARGET="$(readlink "$CLAUDE_MD")"
  if [[ "$CURRENT_TARGET" = "$TARGET_SOP" ]]; then
    echo "  CLAUDE.md symlink: OK (already correct)"
  else
    echo "  CLAUDE.md symlink: updating (was -> $CURRENT_TARGET)"
    ln -sf "$TARGET_SOP" "$CLAUDE_MD"
  fi
elif [[ -e "$CLAUDE_MD" ]]; then
  echo "  WARNING: CLAUDE.md exists but is not a symlink — skipping (back up and remove manually)"
else
  ln -s "$TARGET_SOP" "$CLAUDE_MD"
  echo "  CLAUDE.md symlink: created"
fi

# 2. .claude/repo-context.md skeleton
REPO_CONTEXT="$REPO_DIR/.claude/repo-context.md"
mkdir -p "$REPO_DIR/.claude"

if [[ -e "$REPO_CONTEXT" ]]; then
  echo "  .claude/repo-context.md: already exists"
else
  cat > "$REPO_CONTEXT" << 'SKELETON'
# Repo Context

## Purpose

<!-- What does this repo do? One paragraph. -->

## Key Paths

<!-- Important directories and files for agents working in this repo. -->

## Dependencies

<!-- Other repos or services this repo depends on. -->

## Notes

<!-- Anything an agent should know before working here. -->
SKELETON
  echo "  .claude/repo-context.md: created skeleton"
fi

# 3. docs/upstream/ symlink -> core/docs/
UPSTREAM_DIR="$REPO_DIR/docs/upstream"
TARGET_DOCS="$WORKLOG_ROOT/core/docs"

mkdir -p "$REPO_DIR/docs"

if [[ -L "$UPSTREAM_DIR" ]]; then
  CURRENT_TARGET="$(readlink "$UPSTREAM_DIR")"
  if [[ "$CURRENT_TARGET" = "$TARGET_DOCS" ]]; then
    echo "  docs/upstream symlink: OK (already correct)"
  else
    echo "  docs/upstream symlink: updating (was -> $CURRENT_TARGET)"
    ln -sfn "$TARGET_DOCS" "$UPSTREAM_DIR"
  fi
elif [[ -d "$UPSTREAM_DIR" ]]; then
  echo "  WARNING: docs/upstream/ exists as a directory — skipping (remove manually)"
else
  ln -s "$TARGET_DOCS" "$UPSTREAM_DIR"
  echo "  docs/upstream symlink: created"
fi

echo "Done."
