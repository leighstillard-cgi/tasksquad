#!/usr/bin/env bash
#
# generate-manuals.sh — Generates human-readable concatenated manuals from wiki pages.
# Auto-discovers wiki subdirectories and generates one manual per directory.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WORKLOG_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

WIKI_DIR="$WORKLOG_ROOT/wiki"
MANUALS_DIR="$WORKLOG_ROOT/manuals"

if [[ ! -d "$WIKI_DIR" ]]; then
  echo "ERROR: Wiki directory not found: $WIKI_DIR"
  exit 1
fi

mkdir -p "$MANUALS_DIR"

# --- Helper functions ---

extract_field() {
  local file="$1"
  local field="$2"
  sed -n '/^---$/,/^---$/p' "$file" | grep "^${field}:" | head -1 | sed "s/^${field}:[[:space:]]*//" | sed 's/^["'"'"']//' | sed 's/["'"'"']$//'
}

extract_id() {
  extract_field "$1" "id"
}

extract_title() {
  extract_field "$1" "title"
}

strip_frontmatter() {
  local file="$1"
  awk 'BEGIN{fm=0} /^---$/{fm++; next} fm>=2{print}' "$file"
}

# Titlecase a directory name: "adrs" -> "ADRs", "stories" -> "Stories"
dir_to_title() {
  local dir="$1"
  case "$dir" in
    adrs) echo "Architecture Decision Records" ;;
    epics) echo "Epics" ;;
    stories) echo "Stories" ;;
    concepts) echo "Domain Concepts" ;;
    components) echo "Components" ;;
    risks) echo "Risk Register" ;;
    standards) echo "Standards" ;;
    runbooks) echo "Runbooks" ;;
    infrastructure) echo "Infrastructure" ;;
    *) echo "$dir" | sed 's/-/ /g; s/\b\(.\)/\u\1/g' ;;
  esac
}

# Build a manual from a wiki subdirectory
build_manual() {
  local source_dir="$1"
  local dir_name
  dir_name="$(basename "$source_dir")"
  local manual_title
  manual_title="$(dir_to_title "$dir_name")"
  local output="$MANUALS_DIR/${dir_name}.md"

  local file_list
  file_list="$(mktemp)"

  for f in "$source_dir"/*.md; do
    [[ -f "$f" ]] || continue
    local id
    id="$(extract_id "$f")"
    [[ -z "$id" ]] && id="$(basename "$f" .md)"
    printf '%s\t%s\n' "$id" "$f" >> "$file_list"
  done

  local sorted
  sorted="$(sort -t$'\t' -k1,1 "$file_list")"
  rm -f "$file_list"

  if [[ -z "$sorted" ]]; then
    return
  fi

  local tmpfile
  tmpfile="$(mktemp)"

  {
    echo "# $manual_title"
    echo ""
    echo "> Auto-generated from wiki/${dir_name}/. Do not edit directly."
    echo ""
    echo "## Table of Contents"
    echo ""

    local idx=1
    while IFS=$'\t' read -r id filepath; do
      local title
      title="$(extract_title "$filepath")"
      [[ -z "$title" ]] && title="$id"
      echo "${idx}. ${title} (${id})"
      ((idx++))
    done <<< "$sorted"

    echo ""
    echo "---"
    echo ""

    local first=true
    while IFS=$'\t' read -r id filepath; do
      if [[ "$first" = true ]]; then
        first=false
      else
        echo ""
        echo "---"
        echo ""
      fi
      strip_frontmatter "$filepath"
    done <<< "$sorted"
  } > "$tmpfile"

  if [[ -e "$output" ]]; then
    chmod u+w "$output" 2>/dev/null || true
  fi
  cp "$tmpfile" "$output"
  chmod 444 "$output"
  rm -f "$tmpfile"

  local count
  count="$(echo "$sorted" | wc -l)"
  echo "  Generated: $(basename "$output") ($count pages)"
}

echo "Generating manuals from: $WIKI_DIR"
echo "Output directory: $MANUALS_DIR"
echo ""

generated=0
for subdir in "$WIKI_DIR"/*/; do
  [[ -d "$subdir" ]] || continue
  # Skip if no .md files
  shopt -s nullglob
  md_files=("$subdir"*.md)
  shopt -u nullglob
  [[ ${#md_files[@]} -eq 0 ]] && continue

  build_manual "$subdir"
  ((generated++))
done

if [[ "$generated" -eq 0 ]]; then
  echo "No wiki subdirectories with .md files found."
else
  echo ""
  echo "Done. Generated $generated manuals in: $MANUALS_DIR"
fi
