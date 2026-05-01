# ASE to MS SQL — Conversion Patterns Guide

**Status:** Seed version — enriched by agents as conversions are completed
**Last updated:** 2026-03-20
**Mandatory reading for:** Every worker agent assigned a STORY-02.X conversion story

---

## How This Guide Works

This guide documents known ASE-to-MS-SQL divergence patterns with before/after examples. Worker agents read this before starting any conversion. When an agent encounters a new pattern not covered here, they document it in the completion report as a deviation. The PM agent adds new patterns to this guide after processing the completion.

After 20–30 conversions, this guide typically covers 90%+ of encountered patterns.

---

## Error Handling: @@error to TRY/CATCH

**This is the single most common conversion pattern.** ASE uses per-statement `@@error` checking. MS SQL uses structured `TRY/CATCH`.

### ASE

```sql
INSERT INTO orders (order_id, customer_id, amount)
VALUES (@order_id, @customer_id, @amount)

IF @@error <> 0
BEGIN
    RAISERROR 50001 'Order insert failed'
    RETURN 1
END

UPDATE inventory SET qty = qty - @qty WHERE product_id = @product_id

IF @@error <> 0
BEGIN
    RAISERROR 50002 'Inventory update failed'
    RETURN 1
END
```

### MS SQL

```sql
BEGIN TRY
    INSERT INTO orders (order_id, customer_id, amount)
    VALUES (@order_id, @customer_id, @amount)

    UPDATE inventory SET qty = qty - @qty WHERE product_id = @product_id
END TRY
BEGIN CATCH
    THROW 50001, 'Operation failed', 1
    -- Or for more detail:
    -- DECLARE @msg NVARCHAR(4000) = ERROR_MESSAGE()
    -- RAISERROR(@msg, 16, 1)
    RETURN 1
END CATCH
```

### Notes

- THROW requires SQL Server 2012+. For older targets, use `RAISERROR(@msg, 16, 1)` with MS SQL syntax.
- ASE `RAISERROR <number> '<message>'` → MS SQL `RAISERROR('<message>', 16, 1)` or `THROW <number>, '<message>', 1`
- When converting multi-statement procedures with many `@@error` checks, wrap the entire logical unit in a single TRY/CATCH rather than converting each check individually.
- Preserve the original error numbers where possible for downstream error handling compatibility.

---

## RAISERROR Syntax Differences

### ASE

```sql
RAISERROR 50001 'Something failed'
RAISERROR 50001, 'Something failed with %1!', @detail
```

### MS SQL

```sql
-- Option A: THROW (preferred, SQL Server 2012+)
THROW 50001, 'Something failed', 1

-- Option B: RAISERROR (MS SQL syntax)
RAISERROR('Something failed', 16, 1)
RAISERROR('Something failed with %s', 16, 1, @detail)
```

### Notes

- ASE uses `%1!`, `%2!` for parameter substitution. MS SQL uses `%s`, `%d`, `%i`.
- THROW does not support parameter substitution — build the message string first with `CONCAT()` or `FORMATMESSAGE()`.

---

## Row Limiting: SET ROWCOUNT to TOP

### ASE

```sql
SET ROWCOUNT 1000
DELETE FROM staging_table WHERE processed = 1
SET ROWCOUNT 0
```

### MS SQL

```sql
DELETE TOP (1000) FROM staging_table WHERE processed = 1
```

### Notes

- `SET ROWCOUNT` is deprecated in MS SQL for INSERT, UPDATE, DELETE. Always convert to `TOP(n)`.
- `SET ROWCOUNT` still works for SELECT in MS SQL but `TOP` is preferred.
- If `SET ROWCOUNT` is used in a loop pattern (delete in batches), preserve the loop structure but replace with `DELETE TOP(n)`.

---

## Cursor Syntax

### ASE

```sql
DECLARE mycursor CURSOR FOR
    SELECT col1, col2 FROM mytable WHERE status = 'active'

OPEN mycursor
FETCH mycursor INTO @val1, @val2

WHILE @@sqlstatus = 0
BEGIN
    -- process row
    FETCH mycursor INTO @val1, @val2
END

CLOSE mycursor
DEALLOCATE CURSOR mycursor
```

### MS SQL

```sql
DECLARE mycursor CURSOR LOCAL FAST_FORWARD FOR
    SELECT col1, col2 FROM mytable WHERE status = 'active'

OPEN mycursor
FETCH NEXT FROM mycursor INTO @val1, @val2

WHILE @@FETCH_STATUS = 0
BEGIN
    -- process row
    FETCH NEXT FROM mycursor INTO @val1, @val2
END

CLOSE mycursor
DEALLOCATE mycursor
```

### Notes

- ASE `@@sqlstatus = 0` → MS SQL `@@FETCH_STATUS = 0`
- ASE `FETCH cursor_name` → MS SQL `FETCH NEXT FROM cursor_name`
- ASE `DEALLOCATE CURSOR name` → MS SQL `DEALLOCATE name` (no CURSOR keyword)
- Always add `LOCAL FAST_FORWARD` to MS SQL cursors for performance unless the cursor needs to be updatable or scrollable.
- If possible, evaluate whether the cursor can be replaced with a set-based operation entirely.

---

## CONVERT Style Codes

Most CONVERT style codes are compatible between ASE and MS SQL. The following are known differences:

| Style | ASE Meaning | MS SQL Meaning | Action |
|---|---|---|---|
| 0 | mon dd yyyy hh:miAM | mon dd yyyy hh:miAM | Compatible |
| 1 | mm/dd/yy | mm/dd/yy | Compatible |
| 100-113 | Various date formats | Various date formats | Check each — mostly compatible |
| 140 | yyyy-mm-dd hh:mi:ss:mmmAM | Not standard in MS SQL | Use style 121 or FORMAT() |

### Notes

- When in doubt, test the CONVERT with sample data on the MS SQL target and compare output.
- MS SQL has additional styles (126 for ISO 8601, 127 for ISO 8601 with timezone) that can be used as modern replacements.

---

## String Functions

| ASE | MS SQL | Notes |
|---|---|---|
| `CHAR_LENGTH(str)` | `LEN(str)` | LEN trims trailing spaces; use DATALENGTH for byte count |
| `STR_REPLACE(str, old, new)` | `REPLACE(str, old, new)` | Direct swap |
| `SUBSTRING(str, start, len)` | `SUBSTRING(str, start, len)` | Compatible |
| `CHARINDEX(substr, str)` | `CHARINDEX(substr, str)` | Compatible |
| `LTRIM(RTRIM(str))` | `TRIM(str)` | TRIM available in SQL Server 2017+; otherwise use LTRIM(RTRIM()) |
| `REPLICATE(str, n)` | `REPLICATE(str, n)` | Compatible |

---

## Date Functions

| ASE | MS SQL | Notes |
|---|---|---|
| `GETDATE()` | `GETDATE()` | Compatible |
| `DATEADD(dd, n, date)` | `DATEADD(DAY, n, date)` | Use full part name for clarity |
| `DATEDIFF(dd, d1, d2)` | `DATEDIFF(DAY, d1, d2)` | Compatible |
| `DATEPART(yy, date)` | `DATEPART(YEAR, date)` | Use full part name |

### Notes

- ASE abbreviated date part names (dd, mm, yy) work in MS SQL but are discouraged. Convert to full names (DAY, MONTH, YEAR) for clarity.
- MS SQL 2022+ has DATETRUNC() which can simplify some patterns.

---

## Temporary Tables

Both ASE and MS SQL use `#temp` for local temp tables and `##temp` for global temp tables. The syntax is mostly compatible, but there are scope and lifetime differences:

- In ASE, a #temp table created in a stored procedure is visible to called sub-procedures. This is also true in MS SQL.
- ASE allows `SELECT INTO #temp` with an existing #temp table name (it drops and recreates). MS SQL requires an explicit `DROP TABLE IF EXISTS #temp` first.
- Always add `DROP TABLE IF EXISTS #temp` before `SELECT INTO #temp` in MS SQL conversions.

---

## Cross-Database References

### ASE

```sql
SELECT * FROM other_db..table_name
SELECT * FROM other_db.dbo.table_name
```

### MS SQL

```sql
SELECT * FROM other_db.dbo.table_name  -- explicit schema always
-- Or if on different server:
SELECT * FROM [linked_server].other_db.dbo.table_name
```

### Notes

- ASE `db..table` (double dot) is shorthand for `db.dbo.table`. MS SQL supports this but always convert to explicit three-part names for clarity.
- If the referenced database doesn't exist on the MS SQL target, this needs an architectural decision (linked server, synonym, or schema merge).

---

## Identity Columns

| ASE | MS SQL | Notes |
|---|---|---|
| `IDENTITY(1,1)` in column def | `IDENTITY(1,1)` in column def | Compatible |
| `@@identity` | `SCOPE_IDENTITY()` | Always use SCOPE_IDENTITY() in MS SQL — @@IDENTITY can return wrong value with triggers |
| `SET IDENTITY_INSERT ON/OFF` | `SET IDENTITY_INSERT ON/OFF` | Compatible |

---

## Dynamic SQL

### ASE

```sql
EXEC('SELECT * FROM ' + @tablename + ' WHERE id = ' + @id)
```

### MS SQL

```sql
-- Preferred: parameterised dynamic SQL
DECLARE @sql NVARCHAR(MAX) = N'SELECT * FROM ' + QUOTENAME(@tablename) + N' WHERE id = @id_param'
EXEC sp_executesql @sql, N'@id_param INT', @id_param = @id
```

### Notes

- Always convert `EXEC()` string concatenation to `sp_executesql` with parameters where possible. This is both a security improvement (SQL injection prevention) and a performance improvement (plan caching).
- Use `QUOTENAME()` for dynamic table/column names.
- If the dynamic SQL is too complex to parameterise, document it as an escalation in the completion report.

---

## Adding New Patterns

When you encounter a pattern not covered here, include it in the **Deviations from Spec** section of your completion report with:

1. ASE pattern (before)
2. MS SQL equivalent (after)
3. Any gotchas or edge cases
4. Complexity rating (Low / Medium / High)

The PM agent will add validated patterns to this guide.
