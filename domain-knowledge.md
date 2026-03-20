# Domain Knowledge

**Status:** Business context for agents working on this project

This file supplements story specs with broader "why" context that agents need to make good implementation decisions.

---

## SAP ASE to MS SQL Migration

CGI is migrating a client's database workloads from SAP ASE (Adaptive Server Enterprise, formerly Sybase) to Microsoft SQL Server. The migration uses SSMA (SQL Server Migration Assistant) for bulk conversion, with manual/agent-assisted conversion for problem jobs that SSMA cannot handle automatically.

### Key Context

- The source database contains production stored procedures, queries, and jobs
- SSMA handles the majority of conversions automatically
- The remaining objects require manual intervention due to ASE-specific syntax and behaviour differences
- See `guides/conversion-patterns-guide.md` for documented ASE-to-MS-SQL divergence patterns
- Each converted object must produce identical output to the ASE original (verified by baseline comparison)
