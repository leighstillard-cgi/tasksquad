# STORY-02.XXX · Convert <object_type> <object_name> to MS SQL

**Status:** ready
**Repo:** tasksquad/ase-mssql-conversion
**Depends on:** STORY-01.2 (patterns guide)
**Severity:** <Low / Medium / High>

## Source

- **Database:** <ASE database name>
- **Schema:** <schema>
- **Object:** <object name> (<query / stored procedure / function / trigger>)

## SSMA Status

<What SSMA produced: partial conversion with errors / skipped / converted but incorrect output>

## Known Issues

<Specific ASE constructs identified, e.g.:>
- Uses @@error pattern (needs TRY/CATCH conversion)
- ASE-specific CONVERT style 140
- Cross-database reference to <other_db>

## Source SQL

```sql
-- Paste the original ASE SQL here
```

## Baseline Output

- **Result set:** <row count> rows, checksum <value>
- **Sample rows:** <first 5 rows for visual comparison>
- **Execution time:** <seconds on ASE test instance>
- **Output parameters:** <if stored procedure, list output param values>
- **Affected rows:** <if DML, count of rows affected>

## Acceptance Criteria

- [ ] Converted SQL compiles on MS SQL without errors
- [ ] Result set matches ASE baseline (row count, values, ordering)
- [ ] Execution time within 2x of ASE baseline
- [ ] All known issues addressed (list specific patterns converted)
- [ ] No ASE-specific constructs remain in converted SQL
- [ ] New patterns documented in completion report (if any encountered)

## MCP Connection

- **Environment:** migration_target (MS SQL test instance)
- **Schemas:** <list allowed schemas for this conversion>

## Required Reading

- `data/standards/ase-to-mssql.md`
- `data/guides/conversion-patterns-guide.md`
- `CLAUDE.md`
