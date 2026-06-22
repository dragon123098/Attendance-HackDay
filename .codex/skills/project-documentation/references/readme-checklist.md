# README Checklist

Use this structure by default unless the repo already has a stronger established format.

## Minimum sections

### 1. Project overview

- State what the app, service, library, or job is for.
- Name the primary audience or use case.
- Mention major constraints if they shape the design.

### 2. Current capabilities

- Summarize the major user-facing or operator-facing features that exist today.
- Call out important limitations, partial flows, or in-progress areas.

### 3. Architecture at a glance

- Explain the main subsystems and how they connect.
- Add a Mermaid diagram if that makes the flow faster to understand.

### 4. Codebase map

- Point new developers to the main entrypoints and packages.
- Explain where request handling, business logic, data access, templates/UI, and background jobs live.

### 5. Local development

- Document required tools, environment variables, databases, and startup commands.
- Keep run instructions exact and executable.

### 6. Deeper docs

- Link to architecture, migration, API, or operational docs instead of bloating the README.

## Writing rules

- Prefer present-tense descriptions of what exists now.
- Keep the overview high-level; move detailed procedures into `docs/`.
- Do not imply a feature is complete just because scaffolding exists.
- When the app is mid-migration or mid-rewrite, say so plainly.
