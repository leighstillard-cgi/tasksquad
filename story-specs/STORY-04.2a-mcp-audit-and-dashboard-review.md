# STORY-04.2a + Dashboard Review — Parallel Dispatch

**Purpose:** Run three parallel workstreams: dashboard PR review, Sybase MCP server security audit, MSSQL MCP server security audit. Time-boxed to ~50 minutes.

---

## Instructions

Run all three workstreams in parallel using the Agent tool. Each workstream is independent. Use `mode: "bypassPermissions"` or run via a session with `--dangerously-skip-permissions` since this is unattended.

---

## Workstream 1: Dashboard PR and Review

The dashboard (STORY-00.1 + STORY-00.2) was built directly on `main`. Create a PR for review, then run the review skills.

### Steps

1. `cd ~/workspace/tasksquad`
2. Identify the dashboard commits:
   ```bash
   git log --oneline --all | grep -i dashboard
   ```
   The relevant commits start from `4db1331 Add TaskSquad monitoring dashboard (STORY-00.1)` through the latest dashboard commit.

3. Create a feature branch from the commit before the first dashboard commit, cherry-pick dashboard commits onto it, and push:
   ```bash
   git checkout -b feature/dashboard-review main
   git push -u origin feature/dashboard-review
   ```
   Actually — the dashboard commits are already on main. Create the PR comparing the dashboard file changes. Use:
   ```bash
   gh pr create --base main --head feature/dashboard-review --title "Dashboard: Subagent Monitor + Live Updates (STORY-00.1, 00.2)" --body "Review PR for dashboard implementation"
   ```
   If that doesn't work because both point to the same commits, create a review branch from before the dashboard work and compare.

4. Run `/simplify` against the dashboard code:
   - Get the diff of all dashboard files: `git diff <pre-dashboard-commit>..HEAD -- dashboard/`
   - Launch 3 parallel review agents per the simplify skill (Code Reuse, Code Quality, Efficiency)
   - Fix any findings of substance

5. Run `/superpowers-extended-cc:requesting-code-review` against the dashboard code.

6. Run a security review of the dashboard code:
   - Path traversal in `/api/completion/:filename` and `/api/session-log/:filename` endpoints
   - WebSocket origin validation
   - XSS in markdown rendering (goldmark configuration)
   - File serving restricted to expected directories only
   - No sensitive data exposure in API responses
   - OWASP top 10 relevant checks for a Go HTTP server
   - Fix any medium+ findings

7. Commit and push all fixes. Update the PR.

### Dashboard location
- Code: `~/workspace/tasksquad/dashboard/`
- Main files: `main.go`, `internal/handlers/handlers.go`, `internal/watcher/`, `static/`, `templates/`
- Dependencies: gorilla/websocket, fsnotify, goldmark

---

## Workstream 2: Sybase MCP Server Security Audit

**Repo:** https://github.com/CDataSoftware/sap-sybase-mcp-server-by-cdata
**Goal:** Security audit, fork, implement fixes for medium+ findings, create PR, run reviews.

### Steps

1. Fork the repo:
   ```bash
   gh repo fork CDataSoftware/sap-sybase-mcp-server-by-cdata --clone --remote
   ```
   This forks to `leighstillard-cgi/sap-sybase-mcp-server-by-cdata` and clones locally.

2. `cd` into the cloned repo.

3. **Security Audit** — Review the entire codebase for:
   - **Data exfiltration**: Any HTTP calls, telemetry, analytics, phone-home behaviour beyond the target database connection
   - **Credential handling**: Are credentials logged? Stored insecurely? Exposed in error messages?
   - **Prompt injection**: Do MCP tool descriptions or parameter schemas contain injection vectors? Could crafted DB responses inject into the LLM context?
   - **Dependency audit**: Run `npm audit` / `pip audit` / equivalent. Check for known CVEs.
   - **Permission model**: Verify read-only capability exists. Check for hidden write/DDL/DML paths. Can a crafted MCP call execute arbitrary SQL?
   - **Input validation**: Are SQL parameters properly sanitised? Can tool parameters be used for SQL injection?
   - **Network behaviour**: What outbound connections does the server make? DNS, HTTP, etc.
   - **Supply chain**: Check dependency tree for suspicious or unmaintained packages.

4. **Document findings** in `security-audits/mcp-server-audit-sybase.md` (in the tasksquad repo, not the fork):
   ```markdown
   # Sybase MCP Server Security Audit
   **Repo:** CDataSoftware/sap-sybase-mcp-server-by-cdata
   **Date:** 2026-04-15
   **Auditor:** Claude Code automated audit

   ## Summary
   Risk level: [Low/Medium/High/Critical]

   ## Findings
   ### [SEVERITY] Finding title
   - **Location:** file:line
   - **Description:** what the issue is
   - **Risk:** what could happen
   - **Recommendation:** how to fix
   ```

5. **Implement fixes** for any Medium or above findings:
   ```bash
   git checkout -b security/audit-fixes
   ```
   Make the fixes. Commit with clear messages explaining each fix.

6. **Create PR** on the fork:
   ```bash
   git push -u origin security/audit-fixes
   gh pr create --title "Security: audit fixes for MCP server" --body "Findings and fixes from security audit"
   ```

7. **Run reviews on the PR:**
   - Run the simplify workflow (3 parallel agents: reuse, quality, efficiency) against the diff
   - Run requesting-code-review
   - Run a security-focused re-review to verify fixes are correct and complete
   - Fix any findings, push updates

8. **Write results** to `~/workspace/tasksquad/security-audits/mcp-server-audit-sybase.md`

### Risk assessment output
End with a clear **APPROVE / APPROVE WITH MITIGATIONS / REJECT** recommendation for deploying this MCP server in the Macquarie environment.

---

## Workstream 3: MSSQL MCP Server Security Audit

**Repo:** https://github.com/aadversteeg/mssqlclient-mcp-server
**Goal:** Identical to Workstream 2 but for the MSSQL server.

### Steps

1. Fork the repo:
   ```bash
   gh repo fork aadversteeg/mssqlclient-mcp-server --clone --remote
   ```

2. `cd` into the cloned repo.

3. **Security Audit** — same checklist as Workstream 2:
   - Data exfiltration / phone-home
   - Credential handling
   - Prompt injection vectors
   - Dependency audit
   - Permission model (read-only verification)
   - Input validation / SQL injection
   - Network behaviour
   - Supply chain

4. **Document findings** in `~/workspace/tasksquad/security-audits/mcp-server-audit-mssql.md`

5. **Implement fixes** for Medium+ findings on `security/audit-fixes` branch.

6. **Create PR** on the fork.

7. **Run reviews** (simplify, code-review, security re-review). Fix findings.

8. **Write final assessment** with APPROVE / APPROVE WITH MITIGATIONS / REJECT.

---

## Completion

When all three workstreams finish, write a summary to `~/workspace/tasksquad/completions/STORY-04.2a-completion.md` using the template at `core/templates/story-completion.md`. Include:

- Dashboard review: findings, fixes applied, PR URL
- Sybase MCP: risk assessment, findings count by severity, fixes applied, PR URL
- MSSQL MCP: risk assessment, findings count by severity, fixes applied, PR URL
- Overall recommendation: are both MCP servers safe to deploy?

Also create `~/workspace/tasksquad/security-audits/` directory if it doesn't exist.

Notify Leigh via Slack when complete:
```bash
cc-connect send --channel <appropriate-channel> --text "STORY-04.2a complete: dashboard reviewed, both MCP servers audited. See completions/STORY-04.2a-completion.md for results. <@U0AL7H3HQ56>"
```

---

## Context Files

The agent should read these before starting:
- `~/workspace/tasksquad/CLAUDE.md` — project standards
- `~/workspace/tasksquad/.client/backlog-client.md` — client backlog with STORY-04.2a spec
- `~/workspace/tasksquad/core/docs/standards/security.md` — security standards to apply
