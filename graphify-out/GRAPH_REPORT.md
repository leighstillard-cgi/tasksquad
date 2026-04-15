# Graph Report - wiki + core/docs/standards + guides + core/templates  (2026-04-15)

## Corpus Check
- 25 files · ~9,296 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 79 nodes · 72 edges · 17 communities detected
- Extraction: 89% EXTRACTED · 11% INFERRED · 0% AMBIGUOUS · INFERRED: 8 edges (avg confidence: 0.82)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_API Design Patterns|API Design Patterns]]
- [[_COMMUNITY_Workflow & Agent Discipline|Workflow & Agent Discipline]]
- [[_COMMUNITY_Data Privacy & Secrets|Data Privacy & Secrets]]
- [[_COMMUNITY_Observability & Error Handling|Observability & Error Handling]]
- [[_COMMUNITY_ASE-to-MSSQL Conversion|ASE-to-MSSQL Conversion]]
- [[_COMMUNITY_Database Practices|Database Practices]]
- [[_COMMUNITY_MCP Server Safety|MCP Server Safety]]
- [[_COMMUNITY_Infrastructure Scaffolding|Infrastructure Scaffolding]]
- [[_COMMUNITY_Testing Standards|Testing Standards]]
- [[_COMMUNITY_Code Quality|Code Quality]]
- [[_COMMUNITY_Story Templates|Story Templates]]
- [[_COMMUNITY_ADR Template|ADR Template]]
- [[_COMMUNITY_Epic Template|Epic Template]]
- [[_COMMUNITY_Concept Template|Concept Template]]
- [[_COMMUNITY_Component Template|Component Template]]
- [[_COMMUNITY_Runbook Template|Runbook Template]]
- [[_COMMUNITY_State of Play Template|State of Play Template]]

## God Nodes (most connected - your core abstractions)
1. `ASE to MS SQL Conversion Patterns Guide` - 7 edges
2. `Security Standards` - 6 edges
3. `Workflow Discipline` - 6 edges
4. `Error Handling & Logging` - 5 edges
5. `API Design` - 5 edges
6. `Data Privacy & Tenant Isolation` - 5 edges
7. `Test Categories` - 5 edges
8. `Database Practices` - 4 edges
9. `Testing Standards` - 4 edges
10. `Observability` - 4 edges

## Surprising Connections (you probably didn't know these)
- `@@error to TRY/CATCH Pattern` --semantically_similar_to--> `Error Handling & Logging`  [INFERRED] [semantically similar]
  guides/conversion-patterns-guide.md → core/docs/standards/error-handling.md
- `Dynamic SQL Parameterisation` --semantically_similar_to--> `Parameterised Queries`  [INFERRED] [semantically similar]
  guides/conversion-patterns-guide.md → core/docs/standards/database.md
- `Database Conversion Standards` --references--> `ASE to MS SQL Conversion Patterns Guide`  [EXTRACTED]
  wiki/standards/database-conversion.md → guides/conversion-patterns-guide.md
- `Workflow Discipline` --references--> `Test-Driven Development`  [INFERRED]
  core/docs/standards/workflow-discipline.md → core/docs/standards/testing.md
- `Tenant Isolation` --semantically_similar_to--> `Tenant Boundary Tests`  [INFERRED] [semantically similar]
  core/docs/standards/database.md → core/docs/standards/testing.md

## Hyperedges (group relationships)
- **Pillar Standards** — database_practices, testing_standards, security_standards, error_handling_logging, code_quality, workflow_discipline, observability, api_design, data_privacy, scaffolding, mcp_safety, secrets_management [EXTRACTED 1.00]
- **Wiki Templates** — story_template, story_completion_template, adr_template, epic_template, concept_template, component_template, runbook_template, state_of_play_template [EXTRACTED 1.00]
- **Tenant Security Concepts** — tenant_isolation, data_privacy, tenant_boundary_tests, soft_delete, data_classification [INFERRED 0.85]

## Communities

### Community 0 - "API Design Patterns"
Cohesion: 0.2
Nodes (10): API Design, API Versioning, CSRF Protection, Error Response Schema, Idempotency Keys, Input Validation, Pagination, Rate Limiting (+2 more)

### Community 1 - "Workflow & Agent Discipline"
Cohesion: 0.25
Nodes (9): Agent-Ergonomic Design, Close the Loop, Gap Analysis, Simplicity First, Surgical Changes, Test-Driven Development, Testing Standards, Think Before Coding (+1 more)

### Community 2 - "Data Privacy & Secrets"
Cohesion: 0.22
Nodes (9): Data Classification, Data Privacy & Tenant Isolation, detect-secrets, Environment Variables, PII Masking, Retention Policy, Secret Redaction, Secrets Management (+1 more)

### Community 3 - "Observability & Error Handling"
Cohesion: 0.25
Nodes (9): Correlation ID, @@error to TRY/CATCH Pattern, Error Handling & Logging, Four Golden Signals, Health Checks, Log Levels, Observability, OpenTelemetry (+1 more)

### Community 4 - "ASE-to-MSSQL Conversion"
Cohesion: 0.29
Nodes (7): ASE to MS SQL Conversion Patterns Guide, Cursor Syntax Conversion, Database Conversion Standards, Identity Columns, RAISERROR Conversion, SET ROWCOUNT to TOP, Wiki Index

### Community 5 - "Database Practices"
Cohesion: 0.29
Nodes (7): Connection Pooling, Database Practices, Dynamic SQL Parameterisation, Parameterised Queries, Tenant Boundary Tests, Tenant Isolation, Versioned Migrations

### Community 6 - "MCP Server Safety"
Cohesion: 0.33
Nodes (6): AWS MCP Safety, CloudFormation Templates, Cost Estimation, GitHub MCP Safety, MCP Server Safety, Terraform MCP Safety

### Community 7 - "Infrastructure Scaffolding"
Cohesion: 0.4
Nodes (5): CI/CD Configuration, Containerisation, Infrastructure as Code, SCAFFOLD/FLAG Pattern, Scaffolding & External Dependencies

### Community 8 - "Testing Standards"
Cohesion: 0.4
Nodes (5): Auth Enforcement Tests, Centralised Auth Middleware, Edge Case Tests, Happy Path Tests, Test Categories

### Community 9 - "Code Quality"
Cohesion: 0.5
Nodes (4): Code Quality, Dependency Injection, Linting Configuration, Separation of Concerns

### Community 10 - "Story Templates"
Cohesion: 1.0
Nodes (2): Story Completion Template, Story Template

### Community 11 - "ADR Template"
Cohesion: 1.0
Nodes (1): ADR Template

### Community 12 - "Epic Template"
Cohesion: 1.0
Nodes (1): Epic Template

### Community 13 - "Concept Template"
Cohesion: 1.0
Nodes (1): Concept Template

### Community 14 - "Component Template"
Cohesion: 1.0
Nodes (1): Component Template

### Community 15 - "Runbook Template"
Cohesion: 1.0
Nodes (1): Runbook Template

### Community 16 - "State of Play Template"
Cohesion: 1.0
Nodes (1): State of Play Template

## Knowledge Gaps
- **50 isolated node(s):** `Wiki Index`, `Story Template`, `Story Completion Template`, `ADR Template`, `Epic Template` (+45 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **Thin community `Story Templates`** (2 nodes): `Story Completion Template`, `Story Template`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `ADR Template`** (1 nodes): `ADR Template`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Epic Template`** (1 nodes): `Epic Template`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Concept Template`** (1 nodes): `Concept Template`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Component Template`** (1 nodes): `Component Template`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Runbook Template`** (1 nodes): `Runbook Template`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `State of Play Template`** (1 nodes): `State of Play Template`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Test Categories` connect `Testing Standards` to `Workflow & Agent Discipline`, `Database Practices`?**
  _High betweenness centrality (0.178) - this node is a cross-community bridge._
- **Why does `Security Standards` connect `API Design Patterns` to `Testing Standards`, `Database Practices`?**
  _High betweenness centrality (0.175) - this node is a cross-community bridge._
- **Why does `Tenant Isolation` connect `Database Practices` to `Data Privacy & Secrets`?**
  _High betweenness centrality (0.174) - this node is a cross-community bridge._
- **What connects `Wiki Index`, `Story Template`, `Story Completion Template` to the rest of the system?**
  _50 weakly-connected nodes found - possible documentation gaps or missing edges._