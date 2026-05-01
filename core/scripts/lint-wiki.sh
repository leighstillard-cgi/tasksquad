#!/usr/bin/env bash
#
# lint-wiki.sh — Validates wiki structure and content.
# Uses a Python helper for YAML parsing and structural checks.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WORKLOG_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

WIKI_DIR="${1:-$WORKLOG_ROOT/data/wiki}"
REPORT_DIR="$WORKLOG_ROOT/data/lint-reports"
TODAY="$(date +%Y-%m-%d)"
REPORT_FILE="$REPORT_DIR/lint-${TODAY}.md"

if [[ ! -d "$WIKI_DIR" ]]; then
  echo "ERROR: Wiki directory not found: $WIKI_DIR"
  exit 1
fi

mkdir -p "$REPORT_DIR"

HELPER="$SCRIPT_DIR/lint-wiki-helper.py"
if [[ ! -f "$HELPER" ]]; then
  echo "ERROR: Python helper not found: $HELPER"
  exit 1
fi

RESULT="$(python3 "$HELPER" "$WIKI_DIR")"

ERRORS="$(echo "$RESULT" | python3 -c "import sys,json; d=json.load(sys.stdin); print(len(d.get('errors',[])))")"
WARNINGS="$(echo "$RESULT" | python3 -c "import sys,json; d=json.load(sys.stdin); print(len(d.get('warnings',[])))")"

{
  echo "# Wiki Lint Report — $TODAY"
  echo ""
  echo "**Wiki directory:** \`$WIKI_DIR\`"
  echo ""

  echo "## Summary"
  echo ""
  echo "$RESULT" | python3 -c "
import sys, json
d = json.load(sys.stdin)
s = d.get('summary', {})
print(f\"- **Pages scanned:** {s.get('total_pages', 0)}\")
print(f\"- **Errors:** {s.get('error_count', 0)}\")
print(f\"- **Warnings:** {s.get('warning_count', 0)}\")
print(f\"- **Pages by type:** {s.get('pages_by_type', {})}\")
"
  echo ""

  if [[ "$ERRORS" -gt 0 ]]; then
    echo "## Errors"
    echo ""
    echo "$RESULT" | python3 -c "
import sys, json
d = json.load(sys.stdin)
for e in d.get('errors', []):
    print(f\"- **{e['file']}**: {e['message']}\")
"
    echo ""
  fi

  if [[ "$WARNINGS" -gt 0 ]]; then
    echo "## Warnings"
    echo ""
    echo "$RESULT" | python3 -c "
import sys, json
d = json.load(sys.stdin)
for w in d.get('warnings', []):
    print(f\"- **{w['file']}**: {w['message']}\")
"
    echo ""
  fi
} > "$REPORT_FILE"

echo "Wiki Lint — $TODAY"
echo "  Pages: $(echo "$RESULT" | python3 -c "import sys,json; print(json.load(sys.stdin).get('summary',{}).get('total_pages',0))")"
echo "  Errors: $ERRORS"
echo "  Warnings: $WARNINGS"
echo "  Report: $REPORT_FILE"

if [[ "$ERRORS" -gt 0 ]]; then
  echo ""
  echo "ERRORS found — see report for details."
  exit 1
else
  echo ""
  echo "No errors. Clean."
  exit 0
fi
