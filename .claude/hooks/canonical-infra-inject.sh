#!/usr/bin/env bash
# canonical-infra-inject: Inject canonical facts before infrastructure operations
# Prevents hallucination by re-injecting verbatim facts every Nth matching call.
#
# TEMPLATE: Populate canonical-infra-inject.config with environment-specific patterns
#           and create data/project/data/canonical-facts.md with your infrastructure values.
#
# Fires on: psql, pg_dump, aws, kubectl, terraform, pulumi, ssh, rsync
#           file paths under /infra/, /infrastructure/, /deploy/, /ops/
#           commands/paths containing canonical hostnames/URIs

set -euo pipefail

# ─────────────────────────────────────────────────────────────────────────────
# Environment validation
# ─────────────────────────────────────────────────────────────────────────────

if ! command -v jq &>/dev/null; then
    exit 0
fi

# PROJECT_DIR is set by Claude Code to the project root
if [[ -z "${PROJECT_DIR:-}" ]]; then
    exit 0
fi

# Look for canonical facts in the project
CANONICAL_FILE="${PROJECT_DIR}/data/project/data/canonical-facts.md"

if [[ ! -f "$CANONICAL_FILE" ]]; then
    # No canonical facts configured yet - skip silently
    exit 0
fi

# Config file with patterns to match
CONFIG_FILE="${PROJECT_DIR}/.claude/hooks/canonical-infra-inject.config"

if [[ ! -f "$CONFIG_FILE" ]]; then
    exit 0
fi

# ─────────────────────────────────────────────────────────────────────────────
# Parse input and extract target
# ─────────────────────────────────────────────────────────────────────────────

input=$(cat)
tool_name=$(echo "$input" | jq -r '.tool_name // empty')

if [[ -z "$tool_name" ]]; then
    exit 0
fi

case "$tool_name" in
    Bash)
        target=$(echo "$input" | jq -r '.tool_input.command // empty')
        ;;
    Edit|Write|Read)
        target=$(echo "$input" | jq -r '.tool_input.file_path // empty')
        ;;
    *)
        exit 0
        ;;
esac

if [[ -z "$target" ]]; then
    exit 0
fi

# ─────────────────────────────────────────────────────────────────────────────
# Pattern matching
# ─────────────────────────────────────────────────────────────────────────────

# Load patterns from config (skip comments and empty lines)
patterns=$(grep -v '^\s*#' "$CONFIG_FILE" | grep -v '^\s*$')

matched=false
while IFS= read -r pattern; do
    if echo "$target" | grep -qE "$pattern"; then
        matched=true
        break
    fi
done <<< "$patterns"

if [[ "$matched" != "true" ]]; then
    exit 0
fi

# ─────────────────────────────────────────────────────────────────────────────
# Counter management
# ─────────────────────────────────────────────────────────────────────────────

# Cleanup stale counter files (older than 7 days)
find /tmp -name '.canonical-inject-*' -mtime +7 -delete 2>/dev/null || true

COUNTER_FILE="/tmp/.canonical-inject-${SESSION_KEY:-$$}"
INJECTION_CADENCE="${CANONICAL_INJECTION_CADENCE:-10}"

count=$(cat "$COUNTER_FILE" 2>/dev/null || echo 0)
count=$((count + 1))
echo "$count" > "$COUNTER_FILE"

# Inject on first match OR every Nth match
if (( count != 1 && count % INJECTION_CADENCE != 0 )); then
    exit 0
fi

# ─────────────────────────────────────────────────────────────────────────────
# Extract canonical facts
# ─────────────────────────────────────────────────────────────────────────────

# Extract content between markers
canonical_excerpt=$(sed -n '/<!-- canonical-facts-start -->/,/<!-- canonical-facts-end -->/p' "$CANONICAL_FILE" | grep -v '<!-- canonical-facts')

# Sanity check
if [[ -z "$canonical_excerpt" || ${#canonical_excerpt} -lt 50 ]]; then
    reason="Warning: Canonical facts extraction failed - read $CANONICAL_FILE manually before proceeding."
else
    reason="Infrastructure operation detected (call $count of session).

CANONICAL FACTS - use these values verbatim:

$canonical_excerpt

Do not generate hostnames, URIs, or account IDs from memory.
If a value is not listed above, surface the gap explicitly."
fi

# ─────────────────────────────────────────────────────────────────────────────
# Emit JSON output
# ─────────────────────────────────────────────────────────────────────────────

jq -n --arg reason "$reason" '{
  "hookSpecificOutput": {
    "hookEventName": "PreToolUse",
    "permissionDecision": "allow",
    "permissionDecisionReason": $reason
  }
}'

exit 0
