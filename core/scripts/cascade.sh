#!/usr/bin/env bash
#
# cascade.sh — Propagates a central file to all repos listed in a manifest.
# Verifies symlinks are intact; recreates if broken.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WORKLOG_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
MANIFEST="$WORKLOG_ROOT/core/repos.manifest"

usage() {
  echo "Usage: $(basename "$0") FILE"
  echo ""
  echo "FILE is relative to core/ (e.g., 'Code-SOP.md' or 'docs/standards/testing.md')"
  echo ""
  echo "Reads repo paths from: $MANIFEST"
  echo "Format: one absolute or sibling-relative repo path per line."
  exit 1
}

[[ $# -lt 1 ]] && usage

FILE="$1"
SOURCE="$WORKLOG_ROOT/core/$FILE"

if [[ ! -e "$SOURCE" ]]; then
  echo "ERROR: Source file does not exist: $SOURCE"
  exit 1
fi

if [[ ! -f "$MANIFEST" ]]; then
  echo "ERROR: Manifest not found: $MANIFEST"
  echo "Create it with one repo path per line."
  exit 1
fi

UPDATED=0
SKIPPED=0

while IFS= read -r line || [[ -n "$line" ]]; do
  [[ -z "$line" || "$line" =~ ^# ]] && continue

  if [[ "$line" = /* ]]; then
    REPO_DIR="$line"
  else
    REPO_DIR="$(dirname "$WORKLOG_ROOT")/$line"
  fi

  if [[ ! -d "$REPO_DIR" ]]; then
    echo "  SKIP: $line (directory not found)"
    ((SKIPPED++)) || true
    continue
  fi

  LINK_PATH="$REPO_DIR/$FILE"

  # Code-SOP.md maps to CLAUDE.md in target repos
  if [[ "$FILE" = "Code-SOP.md" ]]; then
    LINK_PATH="$REPO_DIR/CLAUDE.md"
  fi

  LINK_DIR="$(dirname "$LINK_PATH")"
  if [[ ! -d "$LINK_DIR" ]]; then
    mkdir -p "$LINK_DIR"
  fi

  if [[ -L "$LINK_PATH" ]]; then
    CURRENT="$(readlink "$LINK_PATH")"
    if [[ "$CURRENT" = "$SOURCE" ]]; then
      echo "  OK: $line ($FILE symlink intact)"
      continue
    else
      echo "  FIX: $line (symlink was -> $CURRENT, updating)"
      ln -sf "$SOURCE" "$LINK_PATH"
      ((UPDATED++)) || true
    fi
  elif [[ -e "$LINK_PATH" ]]; then
    echo "  WARN: $line ($FILE exists but is not a symlink — skipping)"
    ((SKIPPED++)) || true
  else
    ln -s "$SOURCE" "$LINK_PATH"
    echo "  NEW: $line ($FILE symlink created)"
    ((UPDATED++)) || true
  fi
done < "$MANIFEST"

echo ""
echo "Summary: $UPDATED updated, $SKIPPED skipped"
