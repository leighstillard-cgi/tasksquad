#!/usr/bin/env python3
#
# lint-wiki-helper.py — Python helper for lint-wiki.sh.
# Validates wiki page frontmatter, naming conventions, wikilinks, and structure.
# Outputs JSON to stdout. No external dependencies.

import json
import os
import re
import sys
from datetime import date, datetime


REQUIRED_FIELDS = [
    "title", "id", "status", "page_class", "page_type",
    "tags", "created", "last_updated",
]

VALID_PAGE_CLASSES = {"pm-only", "agent-write", "canonical-facts"}

VALID_PAGE_TYPES = {
    "adr", "epic", "story", "concept", "component",
    "risk", "standard", "runbook", "infrastructure",
}

CANONICAL_FACTS_REQUIRED = ["certainty_basis", "last_verified_at", "staleness_threshold_days"]

NAMING_PATTERNS = {
    "adr": re.compile(r"^ADR-\d{3}-.+\.md$"),
    "epic": re.compile(r"^EPIC-([A-Z]+-)?[0-9]+.*\.md$"),
    "story": re.compile(r"^STORY-[A-Za-z0-9]+[-.][\d.]+[a-z]?(-[a-z]+)?\.md$"),
    "concept": re.compile(r"^[a-z0-9]+(-[a-z0-9]+)*\.md$"),
    "component": re.compile(r"^[a-z0-9]+(-[a-z0-9]+)*\.md$"),
    "risk": re.compile(r"^RISK-\d{3}-.+\.md$"),
    "standard": re.compile(r"^[a-z0-9]+(-[a-z0-9]+)*\.md$"),
    "runbook": re.compile(r"^[a-z0-9]+(-[a-z0-9]+)*\.md$"),
    "infrastructure": re.compile(r"^(canonical|CANONICAL-.+)\.md$"),
}


def parse_frontmatter(filepath):
    """Parse YAML frontmatter between --- delimiters.
    Returns (dict, body_text, error) or (None, full_text, error)."""
    try:
        with open(filepath, "r", encoding="utf-8") as f:
            content = f.read()
    except Exception as e:
        return None, "", str(e)

    lines = content.split("\n")
    if not lines or lines[0].strip() != "---":
        return None, content, None

    end_idx = None
    for i in range(1, len(lines)):
        if lines[i].strip() == "---":
            end_idx = i
            break

    if end_idx is None:
        return None, content, None

    fm = {}
    fm_lines = lines[1:end_idx]
    current_key = None
    current_list = None

    for line in fm_lines:
        stripped = line.strip()
        if not stripped or stripped.startswith("#"):
            continue

        if stripped.startswith("- ") and current_key and current_list is not None:
            val = stripped[2:].strip().strip("\"'")
            current_list.append(val)
            fm[current_key] = current_list
            continue

        match = re.match(r"^([a-z_]+)\s*:\s*(.*)", line)
        if match:
            key = match.group(1)
            value = match.group(2).strip().strip("\"'")
            current_key = key

            if value == "" or value == "[]":
                current_list = []
                fm[key] = current_list
            elif value.startswith("[") and value.endswith("]"):
                items = [v.strip().strip("\"'") for v in value[1:-1].split(",") if v.strip()]
                fm[key] = items
                current_list = None
            elif value == "null" or value == "~":
                fm[key] = None
                current_list = None
            else:
                fm[key] = value
                current_list = None
        else:
            current_list = None

    body = "\n".join(lines[end_idx + 1:])
    return fm, body, None


def extract_wikilinks(body):
    return re.findall(r"\[\[([^\]]+)\]\]", body)


def extract_markdown_link_paths(body):
    return re.findall(r"\[[^\]]*\]\(([^)]+\.md)\)", body)


def lint_wiki(wiki_dir):
    errors = []
    warnings = []
    all_ids = {}
    all_files = []
    wikilinks = {}
    pages_by_type = {}

    for root, _dirs, files in os.walk(wiki_dir):
        for fname in files:
            if not fname.endswith(".md"):
                continue

            filepath = os.path.join(root, fname)
            relpath = os.path.relpath(filepath, wiki_dir)
            all_files.append(relpath)

            fm, body, parse_error = parse_frontmatter(filepath)

            if parse_error:
                errors.append({"file": relpath, "message": f"Read error: {parse_error}"})
                continue

            if fm is None:
                errors.append({"file": relpath, "message": "Missing or invalid YAML frontmatter"})
                continue

            for field in REQUIRED_FIELDS:
                val = fm.get(field)
                if field == "tags":
                    if val is not None and isinstance(val, list) and len(val) == 0:
                        warnings.append({"file": relpath, "message": "Empty tags list"})
                    elif val is None:
                        errors.append({"file": relpath, "message": f"Missing required field: {field}"})
                elif val is None or (isinstance(val, str) and not val):
                    errors.append({"file": relpath, "message": f"Missing required field: {field}"})

            page_class = fm.get("page_class", "")
            if page_class and page_class not in VALID_PAGE_CLASSES:
                errors.append({
                    "file": relpath,
                    "message": f"Invalid page_class '{page_class}' (valid: {', '.join(sorted(VALID_PAGE_CLASSES))})",
                })

            page_type = fm.get("page_type", "")
            if page_type and page_type not in VALID_PAGE_TYPES:
                errors.append({
                    "file": relpath,
                    "message": f"Invalid page_type '{page_type}' (valid: {', '.join(sorted(VALID_PAGE_TYPES))})",
                })

            if page_class == "canonical-facts":
                for field in CANONICAL_FACTS_REQUIRED:
                    if field not in fm or not fm.get(field):
                        errors.append({
                            "file": relpath,
                            "message": f"canonical-facts page missing required field: {field}",
                        })
                if "## " not in body or "|" not in body:
                    warnings.append({
                        "file": relpath,
                        "message": "canonical-facts page appears to have no fact tables",
                    })
                last_verified = fm.get("last_verified_at")
                threshold_days = fm.get("staleness_threshold_days")
                if last_verified and threshold_days:
                    try:
                        threshold = int(threshold_days)
                        if isinstance(last_verified, str):
                            verified_date = datetime.fromisoformat(last_verified).date()
                        elif isinstance(last_verified, date):
                            verified_date = last_verified
                        else:
                            raise ValueError(f"Unknown date format: {type(last_verified)}")
                        age_days = (date.today() - verified_date).days
                        if age_days > threshold:
                            warnings.append({
                                "file": relpath,
                                "message": f"Stale canonical-facts: last verified {age_days} days ago (threshold: {threshold} days)",
                            })
                    except (ValueError, TypeError) as e:
                        warnings.append({
                            "file": relpath,
                            "message": f"Could not check staleness: {e}",
                        })

            if page_type in NAMING_PATTERNS:
                pattern = NAMING_PATTERNS[page_type]
                if not pattern.match(fname):
                    errors.append({
                        "file": relpath,
                        "message": f"Filename '{fname}' does not match naming convention for page_type '{page_type}'",
                    })

            pages_by_type[page_type] = pages_by_type.get(page_type, 0) + 1

            page_id = fm.get("id", "")
            if page_id:
                if page_id in all_ids:
                    errors.append({
                        "file": relpath,
                        "message": f"Duplicate id '{page_id}' (also in {all_ids[page_id]})",
                    })
                else:
                    all_ids[page_id] = relpath

            links = extract_wikilinks(body)
            wikilinks[relpath] = links

    # Check wiki.md exists and indexes pages
    wiki_index = os.path.join(wiki_dir, "wiki.md")
    indexed_pages = set()
    path_to_id = {v: k for k, v in all_ids.items()}

    if os.path.exists(wiki_index):
        with open(wiki_index, "r", encoding="utf-8") as f:
            wiki_content = f.read()
        wiki_links = extract_wikilinks(wiki_content)
        md_link_paths = extract_markdown_link_paths(wiki_content)
        for link_path in md_link_paths:
            if link_path in path_to_id:
                wiki_links.append(path_to_id[link_path])
        indexed_pages = set(wiki_links)
    else:
        errors.append({"file": "wiki.md", "message": "wiki.md does not exist"})

    # Check wikilink resolution
    for filepath, links in wikilinks.items():
        for link in links:
            if link not in all_ids:
                warnings.append({
                    "file": filepath,
                    "message": f"Wikilink [[{link}]] does not resolve to any known page id",
                })

    # Check for orphan pages
    all_referenced_ids = set()
    for links in wikilinks.values():
        all_referenced_ids.update(links)
    all_referenced_ids.update(indexed_pages)

    for page_id, relpath in all_ids.items():
        if relpath == "wiki.md":
            continue
        if page_id not in all_referenced_ids:
            warnings.append({
                "file": relpath,
                "message": f"Orphan page: '{page_id}' is not referenced by wiki.md or any other page",
            })

    summary = {
        "total_pages": len(all_files),
        "error_count": len(errors),
        "warning_count": len(warnings),
        "pages_by_type": pages_by_type,
    }

    return {
        "errors": errors,
        "warnings": warnings,
        "summary": summary,
    }


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: lint-wiki-helper.py WIKI_DIR", file=sys.stderr)
        sys.exit(1)

    result = lint_wiki(sys.argv[1])
    json.dump(result, sys.stdout, indent=2)
