# Documentation Source Checklist

Use this checklist before changing a README or architecture doc.

## Source of truth order

1. Application entrypoints
   Check `cmd/`, service bootstraps, background jobs, and task runners.
2. Routing and interfaces
   Check routes, handlers/controllers, RPC definitions, CLI commands, webhooks, or API specs.
3. Core behavior
   Check service, domain, workflow, and business-logic packages.
4. Data layer
   Check repositories/stores, schemas, migrations, and persistence models.
5. Runtime configuration
   Check config loaders, required environment variables, Compose files, deployment manifests, and startup commands.
6. Tests
   Use tests to confirm supported behavior, constraints, and edge cases.
7. Existing docs
   Reconcile older docs with the implementation instead of trusting them blindly.

## Questions to answer from the code

- What does the project do for a user or operator?
- What major entrypoints exist today?
- Which flows are complete, partial, or placeholder?
- Where does request or job orchestration live?
- What external systems does the project depend on?
- What commands and environment variables are actually needed for local work?
- What files should a new developer read first?
