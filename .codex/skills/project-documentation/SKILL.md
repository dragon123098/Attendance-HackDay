---
name: project-documentation
description: Update or create project README files, onboarding docs, architecture docs, and Mermaid diagrams so they stay accurate to the current implementation. Use when Codex needs to explain what an application does, how it is structured, how requests or data flow through it, where a new developer should start reading, or when code changes require documentation refreshes.
---

# Project Documentation

## Overview

Use this skill to keep project documentation accurate, onboarding-friendly, and grounded in the code that actually exists. Start by verifying current behavior from source files, then write concise docs that help a new developer understand the system quickly.

## Quick Routing

- For what files count as source of truth, read [references/doc-sources.md](references/doc-sources.md).
- For README structure and onboarding expectations, read [references/readme-checklist.md](references/readme-checklist.md).
- For architecture and flow diagrams, read [references/mermaid-patterns.md](references/mermaid-patterns.md).

## Workflow

1. Inspect the implementation before writing.
   Read the app entrypoints, routes, handlers/controllers, services, stores/repositories, config, migrations or schemas, and any existing docs that describe the same area.
2. Separate current behavior from planned behavior.
   If a route, job, or feature is incomplete, document it as partial, placeholder, or in-progress instead of implying full support.
3. Write for a new developer first.
   Explain what the project does, the main execution paths, where important logic lives, and what files to read next.
4. Keep the docs compact and navigable.
   Prefer a high-level overview in `README.md`, then link to deeper docs for architecture, operations, or migration detail.
5. Use Mermaid when it reduces reading time.
   Add a diagram only when a request flow, data flow, subsystem map, or migration story is easier to understand visually.

## Working Rules

- Treat code, configuration, schemas, and tests as the source of truth over stale prose.
- Update claims about setup, runtime behavior, jobs, and architecture when related code changes.
- Prefer exact names for packages, commands, environment variables, routes, and jobs.
- Use repository-relative Markdown links for repo files and directories. Never write machine-specific absolute filesystem links such as `/Users/...`, `/home/...`, or `file://...` in generated documentation.
- Mark limitations and unfinished flows clearly.
- Avoid copying large amounts of implementation detail into README files; summarize and link outward.
- Keep diagrams stable by using behavior-level labels instead of low-level function trivia.
- If the repo already has a docs structure, extend it instead of inventing a parallel structure.

## Deliverables

- A README that explains:
  - what the project does,
  - what exists today,
  - how it is structured,
  - how to run or work on it,
  - where a new developer should start.
- Optional deeper docs for architecture, operations, migrations, or feature-specific flows.
- Mermaid diagrams when they improve comprehension enough to justify their maintenance cost.
