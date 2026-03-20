# ASE to MS SQL Conversion Standards

**Pillar:** Database Migration
**Scope:** SAP ASE (Adaptive Server Enterprise) to Microsoft SQL Server
**Mandatory reading for:** Every worker agent assigned a STORY-01.X or STORY-02.X conversion story

---

## 1. Objective

Convert SAP ASE database objects (stored procedures, queries, functions, triggers, views) to functionally equivalent Microsoft SQL Server implementations. Each converted object must produce identical results to its ASE original under the same input conditions.

**Non-objective:** This is not a rewrite or modernisation effort. Do not improve, optimise, or refactor business logic. The goal is functional equivalence — same inputs, same outputs, same side effects.

---

## 2. Required Reading Order

Before starting any conversion story, read in this order:

1. This document (conversion standards)
2. `guides/conversion-patterns-guide.md` (syntax-level transformation patterns)
3. The assigned story spec in `story-specs/`
4. `CLAUDE.md` (project-wide standards)

---

## 3. Conversion Process

For each object assigned via a story spec:

1. **Read** the source ASE SQL from the story spec. Understand its purpose, inputs, outputs, and side effects before writing any SQL.
2. **Identify** all ASE-specific constructs by comparing against the patterns guide. List them.
3. **Convert** mechanically using the patterns guide for known divergences. One pattern at a time, not all at once.
4. **Compile** the converted SQL on the MS SQL target. Fix syntax errors.
5. **Validate** against the baseline outputs provided in the story spec (see section 5).
6. **Document** all changes, new patterns, and deviations in the completion report.

---

## 4. Conversion Rules

### 4.1 Mandatory

These rules apply to every conversion without exception:

| Rule | Rationale |
|---|---|
| Convert `EXEC()` string concatenation to `sp_executesql` with parameters | Security (SQL injection prevention) and performance (plan caching) |
| Use `SCOPE_IDENTITY()`, never `@@IDENTITY` | `@@IDENTITY` returns wrong value when triggers fire |
| Wrap logical units in `TRY/CATCH`, not per-statement `@@error` | Cleaner error handling, matches MS SQL idiom |
| Convert `db..table` to explicit `db.dbo.table` three-part names | Clarity and maintainability |
| Add `DROP TABLE IF EXISTS #temp` before every `SELECT INTO #temp` | MS SQL does not auto-drop existing temp tables |
| Use `LOCAL FAST_FORWARD` on all cursors | Performance — unless cursor must be updatable or scrollable |
| Convert `SET ROWCOUNT n` to `TOP(n)` for all DML | `SET ROWCOUNT` is deprecated for INSERT/UPDATE/DELETE in MS SQL |
| Use `QUOTENAME()` for all dynamic object names | Security (prevents injection via object names) |
| Preserve original error numbers in `THROW`/`RAISERROR` | Downstream systems may depend on specific error numbers |

### 4.2 Preferred

Apply where practical. Document in the completion report if you deviate:

- **`THROW` over `RAISERROR`** for new error paths (SQL Server 2012+). Keep `RAISERROR` only when parameter substitution is needed.
- **Full date part names** (`DAY`, `MONTH`, `YEAR`) instead of abbreviations (`dd`, `mm`, `yy`).
- **`TRIM()` over `LTRIM(RTRIM())`** if targeting SQL Server 2017+.
- **Set-based over cursors** if the replacement is straightforward (< 20 lines). If non-trivial, keep the cursor and document why.

### 4.3 Prohibited

Do not do these, even if they would simplify the conversion:

- **Do not change business logic.** If you see a bug in the ASE original, document it in the completion report as an escalation — do not fix it.
- **Do not add indexes, constraints, or schema objects** unless required for the converted SQL to compile. Schema changes require escalation.
- **Do not use SQL Server features newer than the target version** without confirmation in the story spec.
- **Do not hard-code environment-specific values** (server names, database names, file paths). Use variables or configuration.
- **Do not rename objects, parameters, or variables.** The converted object must have the same interface as the original.

---

## 5. Validation Requirements

### 5.1 Functional Equivalence

Every converted object must be verified against its ASE baseline from the story spec:

| Check | Requirement |
|---|---|
| Result set row count | Exact match |
| Result set values | Exact match (see tolerance rules below) |
| Result set ordering | Exact match if original has `ORDER BY`; order-independent comparison otherwise |
| Output parameters | Exact match (stored procedures) |
| Affected row count | Exact match (DML statements) |
| Side effects | Inserts/updates/deletes to other tables must match |

### 5.2 Data Type Tolerances

| Data type | Tolerance |
|---|---|
| Integer types (`INT`, `SMALLINT`, `TINYINT`, `BIGINT`) | Exact match |
| `DECIMAL` / `NUMERIC` | Exact match (same scale and precision) |
| `FLOAT` / `REAL` | Within 1e-10 relative error |
| `DATETIME` | Within 3.33ms (ASE datetime resolution is 1/300th second) |
| `VARCHAR` / `NVARCHAR` | Exact match after trailing-space normalisation |
| `MONEY` / `SMALLMONEY` | Exact match |
| `BIT` | Exact match |
| `BINARY` / `VARBINARY` | Exact match (byte-for-byte) |

### 5.3 Performance Thresholds

| Metric | Acceptable | Escalate |
|---|---|---|
| Execution time | Within 2x of ASE baseline | Exceeds 5x |
| Tempdb usage | Within 2x of ASE temp table usage | Exceeds 5x |
| Query plan | No table scans on tables > 10,000 rows unless ASE original also scans | Full scan introduced where ASE used index seek |

Between 2x and 5x: document in the completion report but do not escalate. Above 5x: escalate.

---

## 6. Transaction Handling

- **Preserve original transaction boundaries.** If the ASE original uses `BEGIN TRAN` / `COMMIT`, the MS SQL version must use the same boundaries.
- **TRY/CATCH with transactions:** When wrapping multi-statement logic in TRY/CATCH, add `BEGIN TRAN` before the TRY block and `COMMIT` at the end of the TRY block. Add `ROLLBACK` in the CATCH block. Only do this if the ASE original implies transactional behaviour (multiple related DML statements).
- **Do not add transactions to single-statement logic.** If the ASE original has no explicit transaction and performs a single DML statement, do not wrap it in a transaction.
- **Nested transactions:** ASE and MS SQL handle nested transactions differently. ASE `@@trancount` semantics are similar but not identical. If the original uses nested transactions, test carefully and document behaviour.

---

## 7. NULL and ANSI Settings

### NULL Handling

- ASE and MS SQL handle NULLs identically in most cases.
- If the ASE original uses `= NULL` comparisons instead of `IS NULL`, convert them. MS SQL should run with `SET ANSI_NULLS ON` (the default and the only option in future SQL Server versions).
- Verify the ASE original doesn't depend on non-standard `SET STRING_RTRUNCATION` or `SET ANSINULL` settings. If it does, document in the completion report.

### ANSI Settings

The following settings should be assumed ON for all MS SQL target code:

```sql
SET ANSI_NULLS ON
SET QUOTED_IDENTIFIER ON
SET ANSI_PADDING ON
SET ANSI_WARNINGS ON
SET CONCAT_NULL_YIELDS_NULL ON
```

If the converted code behaves differently with these settings, that is a bug to be fixed — not a deviation to be documented.

---

## 8. Cross-Database References

When a converted object references another database:

| Scenario | Action |
|---|---|
| Same server, target database exists | Use explicit three-part name: `target_db.dbo.object` |
| Same server, target database does not exist | **Escalate.** Requires architectural decision (linked server, synonym, or schema merge). |
| Different server | **Escalate.** Requires linked server configuration. |

Do not create linked servers, synonyms, or other infrastructure objects. Escalate and wait for resolution.

---

## 9. Error Number Preservation

- If the ASE original uses `RAISERROR 50001`, the MS SQL version must use error number 50001.
- Custom error numbers (> 50000) must match between ASE and MS SQL.
- If the original uses a system error number that differs between ASE and MS SQL, document the mapping in the completion report.
- If `sp_addmessage` is needed to register custom error numbers on the MS SQL target, include the `sp_addmessage` call in the converted SQL and document it.

---

## 10. Testing Requirements

Every conversion must include the following evidence in the completion report:

### 10.1 Compilation

The converted SQL compiles without errors on the MS SQL target. Include the compilation output or confirmation.

### 10.2 Baseline Comparison

Run the converted object with the same inputs as the ASE baseline (provided in the story spec). Include:

- Row count comparison (expected vs actual)
- Value comparison for at least the first 5 rows and last 5 rows
- Output parameter comparison (if stored procedure)
- Affected row count (if DML)

### 10.3 Edge Cases

Test at minimum:

- **NULL inputs** — if the object accepts parameters, pass NULL for each nullable parameter
- **Empty result set** — if applicable, use inputs that produce zero rows
- **Boundary values** — MAX/MIN for date and numeric parameters if the object does range filtering

### 10.4 Error Path

If the object has error handling (@@error checks, RAISERROR):

- Trigger at least one error path and verify the MS SQL version produces the expected error behaviour (error number, message, rollback if applicable)

---

## 11. MCP Database Access

When using MCP to execute queries against the MS SQL target:

- **Read-only by default.** SELECT queries only unless the story spec explicitly authorises writes.
- **Test environment only.** Never execute against production. Verify the environment name in the MCP connection before every session.
- **Parameterised queries only.** No string concatenation in MCP queries.
- **Clean up.** Drop any temp tables or test data created during validation.
- **Log all queries.** Include significant queries in the completion report.

---

## 12. Escalation Triggers

Escalate in the completion report (do not proceed independently) when:

| Trigger | Why |
|---|---|
| Cross-database reference to a database that doesn't exist on the target | Requires architectural decision |
| Schema change required (new table, index, constraint, column) | Requires DBA / architect review |
| Business logic appears incorrect in the ASE original | Not our decision to fix |
| Performance exceeds 5x the ASE baseline after optimisation attempts | May require index or schema changes |
| ASE feature with no MS SQL equivalent | Requires architectural decision on replacement strategy |
| Dynamic SQL too complex to parameterise safely | Security review needed |
| Object depends on ASE-specific system tables or stored procedures | Requires replacement strategy |
| Object uses `WAITFOR`, `DBCC`, or other ASE admin commands | Requires operational review |
| Conversion requires changes to calling code (different parameter signature) | Affects other objects — scope creep |

---

## 13. Completion Report Requirements

In addition to the standard completion report template, every conversion report must include:

1. **Patterns applied:** List every ASE pattern converted and which pattern from the guide was used
2. **New patterns:** Any patterns not in the guide, documented with ASE before / MS SQL after / gotchas (the PM will add these to the guide)
3. **Full converted SQL** in the Files Changed section or as an appendix
4. **Baseline comparison results** showing expected vs actual for all checks in section 5.1
5. **Test evidence** for each item in section 10
6. **Deviations from this standards document** with rationale

---

## 14. Relationship to Other Documents

| Document | What it covers | When to read |
|---|---|---|
| This document (`standards/ase-to-mssql.md`) | Rules, process, quality gates, escalation triggers | Before every conversion |
| `guides/conversion-patterns-guide.md` | Syntax-level ASE → MS SQL transformation patterns | Before every conversion, during conversion |
| `templates/conversion-story.md` | Story spec format for individual conversions | When the PM creates new stories |
| `CLAUDE.md` | Project-wide standards (security, testing, error handling) | Before every story |
