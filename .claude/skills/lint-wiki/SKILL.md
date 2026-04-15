---
name: lint-wiki
description: "Lint wiki pages for frontmatter, naming, and link issues"
---

# Lint Wiki

Run the wiki linter to check for structural issues in wiki pages.

## Workflow

1. **Run the lint script**:
   ```bash
   core/scripts/lint-wiki.sh
   ```

2. **Parse the output** and categorize findings:
   - **Errors** (must fix): missing required frontmatter, broken internal links, invalid filenames
   - **Warnings** (should fix): missing optional fields, orphan pages

3. **Report the results**:
   - If no issues: confirm the wiki is clean
   - If errors: list each error with file path and what needs fixing
   - If only warnings: list them and suggest fixes but note they're non-blocking

4. **Suggest fixes** for each error:
   - Missing frontmatter: show the expected template from `core/templates/`
   - Broken links: suggest the correct target
   - Invalid filenames: show the naming convention from `core/docs/wiki/STRUCTURE.md`

## Rules

- Errors block merges; warnings do not
- Do not auto-fix errors without user confirmation
- The lint script defines what constitutes an error vs warning
