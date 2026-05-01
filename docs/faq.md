# FAQ

## Is TaskSquad an application framework?

No. TaskSquad is a delivery coordination framework for agent-assisted software work. Your application code lives in target repositories.

## Do I have to use the dashboard?

No. The file-based workflow works without the dashboard. The dashboard is a convenience for viewing active work, completions, escalations, and session logs.

## Why are generated graph files ignored?

Graph outputs can be rebuilt from source docs and can be large or environment-specific. Keep the source documentation in Git and regenerate `graphify-out/` after setup.

## Where should client-sensitive work go?

Use `.client/`. It is ignored by Git. Do not put secrets anywhere in the repository, including `.client/`.

## Where should completed dispatches and completion reports go?

Move processed records into `data/archive/`. If the record contains durable learning, summarize that learning in `data/wiki/`.

## Can agents push code?

They can if your environment and permissions allow it. TaskSquad includes guardrails, but you should still review all commits and pull requests, especially anything touching production systems.

## Should this run in dangerous or permission-bypass mode?

Only inside an isolated environment with scoped credentials. Broad permissions make agents faster, but they also increase the impact of mistakes.

## What belongs in canonical facts?

Non-secret facts that agents often need and should not guess: environment names, non-secret endpoints, repository names, service names, account identifiers where policy allows, and known operational constraints.

## What should I do when an agent is uncertain?

Tell it to stop, write an escalation, or document the assumption in the completion report. Do not let agents invent missing infrastructure, credentials, or business rules.
