# Sybase MCP Server Security Audit

**Repo:** CDataSoftware/sap-sybase-mcp-server-by-cdata  
**Date:** 2026-04-15  
**Auditor:** Claude Code automated audit

## Summary

**Risk Level: HIGH**

This MCP server has one critical vulnerability (arbitrary SQL execution) and several medium-severity issues. The codebase is small and focused, with no evidence of data exfiltration or telemetry. However, the lack of SQL validation makes this server dangerous in its current form.

| Severity | Count |
|----------|-------|
| Critical | 1 |
| High     | 1 |
| Medium   | 3 |
| Low      | 2 |

## Findings

### [CRITICAL] Arbitrary SQL Execution - No Statement Validation

- **Location:** `src/main/java/com/cdata/mcp/tools/RunQueryTool.java:52-57`
- **Description:** The `run_query` tool accepts arbitrary SQL from the MCP client and executes it directly against the database without any validation or sanitization. Despite the tool description claiming it only supports SELECT statements, there is no enforcement.
- **Code:**
  ```java
  public McpSchema.CallToolResult run(Map<String, Object> args) {
    String sql = (String)args.get("sql");
    this.logger.info("RunQueryTool({})", sql);
    try {
      try (Connection cn = config.newConnection()) {
        List<McpSchema.Content> content = new ArrayList<>();
        String csv = queryToCsv(cn, sql);  // No validation - arbitrary SQL executed
  ```
- **Risk:** 
  - A malicious or compromised LLM could execute `DROP TABLE`, `DELETE`, `UPDATE`, `INSERT`, or other DML/DDL statements
  - SQL injection through crafted database responses that influence subsequent LLM queries
  - Complete database compromise including data theft, modification, or destruction
- **Recommendation:** 
  1. Parse and validate SQL to only allow SELECT statements
  2. Use a SQL parser library to reject DDL/DML statements
  3. Implement a read-only database connection with restricted permissions
  4. Add a query allowlist/blocklist mechanism

### [HIGH] Identifier Quoting Not Properly Escaped

- **Location:** `src/main/java/com/cdata/mcp/Config.java:125-130`
- **Description:** The `quoteIdentifier` method wraps identifiers in quote characters but does not escape quote characters within the identifier itself.
- **Code:**
  ```java
  public String quoteIdentifier(String id) {
    String open = this.sqlInfo.getProperty(ID_QUOTE_OPEN_CHAR);
    String close = this.sqlInfo.getProperty(ID_QUOTE_CLOSE_CHAR);
    // TODO: Properly escape things  <-- Developer acknowledged the issue
    return open + id + close;
  }
  ```
- **Risk:** SQL injection via crafted table/schema/catalog names if user-controlled data flows through this path
- **Recommendation:** Implement proper escaping (typically doubling the quote character within the identifier)

### [MEDIUM] SQL Queries Logged with Full Content

- **Location:** `src/main/java/com/cdata/mcp/tools/RunQueryTool.java:53`
- **Description:** Complete SQL queries are logged at INFO level, which may include sensitive data in WHERE clauses.
- **Code:**
  ```java
  this.logger.info("RunQueryTool({})", sql);
  ```
- **Risk:** 
  - Queries containing PII, credentials, or sensitive filters will appear in logs
  - Example: `SELECT * FROM users WHERE ssn = '123-45-6789'`
- **Recommendation:** 
  1. Log only query metadata (table names, operation type) at INFO level
  2. Move full query logging to DEBUG level
  3. Consider truncating or masking sensitive patterns

### [MEDIUM] Error Messages May Leak Database Schema Information

- **Location:** Multiple files - `RunQueryTool.java:67`, `GetColumnsTool.java:62`, `GetTablesTool.java:73`, `TableMetadataResource.java:52`
- **Description:** Exception messages are passed directly to the MCP client without sanitization.
- **Code:**
  ```java
  throw new RuntimeException("ERROR: " + ex.getMessage());
  ```
- **Risk:** Database error messages often contain schema details, table names, column names, or even partial data that could aid attackers
- **Recommendation:** Return generic error messages to clients; log detailed errors internally with correlation IDs

### [MEDIUM] No Read-Only Enforcement at Connection Level

- **Location:** `src/main/java/com/cdata/mcp/Config.java:132-134`
- **Description:** Database connections are created without read-only mode.
- **Code:**
  ```java
  public Connection newConnection() throws SQLException {
    return this.driver.connect(this.getJdbcUrl(), new Properties());
  }
  ```
- **Risk:** Even if SQL validation is added, connection-level read-only mode provides defense in depth
- **Recommendation:** Set `connection.setReadOnly(true)` after establishing the connection

### [LOW] Dependency Versions Not Pinned Precisely

- **Location:** `pom.xml`
- **Description:** Dependencies use specific versions but Maven allows compatible updates. The MCP SDK is at 0.8.1 (pre-1.0 suggests API instability).
- **Dependencies:**
  - `io.modelcontextprotocol.sdk:mcp:0.8.1` - Early version, may have security issues
  - `junit:4.13.2` (test scope) - Older JUnit version
  - `slf4j-simple:2.0.16` - Current as of audit date
- **Risk:** Supply chain vulnerability if dependencies have unpatched CVEs
- **Recommendation:** Run `mvn dependency:tree` and check for CVEs; consider using OWASP dependency-check plugin

### [LOW] Charset Handling Uses Default

- **Location:** `src/main/java/com/cdata/mcp/UrlUtil.java:10,13`
- **Description:** URL encoding/decoding uses system default charset instead of explicit UTF-8.
- **Code:**
  ```java
  return URLEncoder.encode(part, Charset.defaultCharset());
  return URLDecoder.decode(part, Charset.defaultCharset());
  ```
- **Risk:** Inconsistent behavior across systems; potential encoding-based injection vectors
- **Recommendation:** Use `StandardCharsets.UTF_8` explicitly

## Positive Findings

1. **No Data Exfiltration:** No outbound HTTP calls, telemetry, or analytics code detected
2. **No Credential Logging:** JDBC URLs (which may contain credentials) are not logged
3. **Simple Architecture:** Small attack surface with only 3 MCP tools
4. **No Dynamic Code Loading:** Beyond the JDBC driver (configured by admin), no runtime code loading
5. **STDIO Transport Only:** Network exposure is limited to the JDBC connection

## Dependency Audit

```
io.modelcontextprotocol.sdk:mcp:0.8.1
  - jackson-databind (transitive) - Check for CVEs
  - slf4j-api (transitive)
org.slf4j:slf4j-simple:2.0.16
junit:junit:4.13.2 (test)
```

**Note:** Full CVE analysis requires Maven dependency resolution. The MCP SDK at 0.8.1 is pre-release software.

## Prompt Injection Analysis

### Tool Descriptions
The tool descriptions in `RunQueryTool.java` are static and do not include user-controlled content. The description correctly states "Execute a SQL SELECT statement" but this is **advisory only** - not enforced.

### Database Response Injection
Database query results flow back to the LLM as CSV text. A malicious database record containing prompt injection payloads (e.g., "Ignore previous instructions and...") could influence LLM behavior. This is a general MCP risk, not specific to this implementation.

## Network Behavior

- **Inbound:** STDIO only (no HTTP listener in current configuration)
- **Outbound:** JDBC connection to configured database only
- **DNS:** Only for JDBC connection resolution

## Recommendations Summary

| Priority | Action |
|----------|--------|
| P0 | Add SQL statement validation to enforce SELECT-only queries |
| P0 | Set connection to read-only mode |
| P1 | Fix identifier quoting to escape embedded quotes |
| P1 | Sanitize error messages before returning to clients |
| P2 | Review and update query logging to avoid PII exposure |
| P2 | Add OWASP dependency-check to build |
| P3 | Use explicit UTF-8 charset |

## Enterprise Deployment Recommendation

**REJECT** in current state.

The arbitrary SQL execution vulnerability makes this server unsuitable for enterprise deployment. The LLM (or a compromised LLM session) can execute destructive database operations despite the documentation claiming read-only support.

**Conditions for Approval:**
1. SQL validation must enforce SELECT-only queries
2. Database connections must be read-only
3. Error messages must be sanitized
4. Identifier escaping must be fixed

With these fixes, the server would be **APPROVE WITH MITIGATIONS** (deploy with read-only database credentials and network isolation).
