#!/usr/bin/env bash
#
# check-repo-health.sh — Validates framework symlinks in any repo.
# Returns non-zero on any failure.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WORKLOG_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

REPO_DIR="${1:-.}"
REPO_DIR="$(cd "$REPO_DIR" && pwd)"

FAIL=0

check_symlink() {
  local link_path="$1"
  local expected_target="$2"
  local label="$3"

  if [[ -L "$link_path" ]]; then
    local actual
    actual="$(readlink "$link_path")"
    if [[ "$actual" = "$expected_target" ]]; then
      echo "  OK: $label -> $expected_target"
    else
      echo "  FAIL: $label points to '$actual' (expected '$expected_target')"
      FAIL=1
    fi
  elif [[ -e "$link_path" ]]; then
    echo "  FAIL: $label exists but is not a symlink"
    FAIL=1
  else
    echo "  FAIL: $label does not exist"
    FAIL=1
  fi
}

check_exists() {
  local path="$1"
  local label="$2"

  if [[ -e "$path" ]]; then
    echo "  OK: $label exists"
  else
    echo "  FAIL: $label does not exist"
    FAIL=1
  fi
}

echo "Checking repo health: $REPO_DIR"
echo ""

check_symlink "$REPO_DIR/CLAUDE.md" "$WORKLOG_ROOT/core/Code-SOP.md" "CLAUDE.md"
check_symlink "$REPO_DIR/docs/upstream" "$WORKLOG_ROOT/core/docs" "docs/upstream/"
check_exists "$REPO_DIR/.claude/repo-context.md" ".claude/repo-context.md"

echo ""
if [[ "$FAIL" -ne 0 ]]; then
  echo "HEALTH CHECK FAILED"
  echo "Run bootstrap-repo.sh to fix symlinks."
  exit 1
else
  echo "All checks passed."
  exit 0
fi
