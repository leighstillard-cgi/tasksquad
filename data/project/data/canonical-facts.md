# Canonical Infrastructure Facts

> **Purpose:** This file contains the single source of truth for infrastructure values.
> The `canonical-infra-inject` hook reads this file and injects these facts into Claude's
> context when infrastructure operations are detected.
>
> **Format:** Values between the start and end markers are injected verbatim. Keep entries concise.

<!-- canonical-facts-start -->


## Database Connections

| Environment | Host | Port | Database |
|-------------|------|------|----------|
| Production  | TBD  | TBD  | TBD      |
| Staging     | TBD  | TBD  | TBD      |
| Development | TBD  | TBD  | TBD      |

## Cloud Resources

| Resource Type | Name/ID | Region | Notes |
|---------------|---------|--------|-------|
| TBD           | TBD     | TBD    | TBD   |

## Service Endpoints

| Service | URL | Auth Method |
|---------|-----|-------------|
| TBD     | TBD | TBD         |

<!-- canonical-facts-end -->

## Notes

Populate this file with your actual infrastructure values before enabling the
`canonical-infra-inject` hook. Values should be:

- Hostnames (databases, APIs, services)
- Account IDs (AWS, GCP, Azure)
- Bucket/container names
- Connection strings (without secrets)
- IP ranges for internal networks

Never include secrets (passwords, API keys, tokens) in this file.
