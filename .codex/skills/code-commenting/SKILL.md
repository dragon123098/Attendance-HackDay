---
name: code-commenting
description: Add or preserve clear code comments for larger functions, core workflow entrypoints, high-churn code hotspots, and constants that are critical to app-wide behavior in non-test files. Use when Codex is editing non-trivial functions, central request or job flows, shared integration boundaries, or critical constants whose purpose and usage would not be obvious without a brief comment, but skip test files and test-only helpers.
---

# Code Commenting

## Overview

Use this skill to keep important code self-explaining when behavior is added or changed. Focus on short, accurate comments that tell the reader why a larger or central function exists, or what a critical constant controls.

## Commenting Rules

- Do not use this skill for test files such as `_test.go`, `*.spec.*`, `*.test.*`, fixture-heavy test helpers, or other code whose main purpose is exercising behavior in tests.
- Do not add comments to every new function or method by default.
- Add or refresh a comment when the function is large enough, central enough, or churn-prone enough that future readers would benefit from extra guidance.
- Refresh the existing comment when a function's purpose or main behavior changes.
- Explain the function's role, inputs or outputs only when that detail helps, and the main work it performs.
- Keep comments high signal. Do not restate the function name line by line.
- Match the file's native comment style and naming conventions.

## Function Guidance

- Apply this guidance only to production code and shared runtime code, not test files.
- Prioritize comments for functions that are one or more of the following:
  - long or multi-branch functions that take time to mentally parse,
  - core workflow entrypoints such as handlers, controllers, jobs, orchestration methods, or integration adapters,
  - shared helpers whose behavior is relied on in many places,
  - hot spots in recent git history, when local history is available and repeated churn suggests the code is easy to misread or regress.
- Usually skip comments for tiny private helpers, thin getters, obvious constructors, and straightforward wrappers unless the intent is still non-obvious.

- Prefer one or two sentences that answer:
  - what the function is for,
  - what main job it performs,
  - and any important usage expectation when it is not obvious.
- For exported or public APIs, write comments that help another developer understand when to call them.
- For private helpers, keep comments brief and add them only when the intent would otherwise be easy to miss.

Examples:

- `// loadConfig reads environment-backed settings and normalizes defaults for app startup.`
- `// buildSummary groups raw migration rows into the totals shown on the dashboard.`

## Constant Guidance

- Apply this guidance only to constants used by the application at runtime, not test-only constants.
- Add a short comment for constants that are critical to app-wide behavior or whose meaning, origin, or intended use is not obvious from the name alone.
- Refresh the existing comment when a constant's meaning or intended usage changes.
- Explain what the constant represents and how callers or maintainers should use it.
- Prioritize comments for shared defaults, limits, status values, sentinel strings, route names, environment keys, auth- or security-related values, cache keys, feature flags, bitmasks, and enum-like groups that affect behavior across the app.
- Skip redundant comments for trivial local constants unless the surrounding code would still be unclear.

Examples:

- `// maxUploadBytes caps request bodies to the largest file size the import flow supports.`
- `// statusArchived marks records that should stay visible in history but not in active lists.`

## Editing Workflow

1. Inspect the changed code before writing comments.
2. First confirm the file is not a test file and the code is not test-only support code.
3. Decide whether the code is a bigger function, a core path, a churn hotspot, or a critical constant before adding a comment.
4. If local git history is available and relevant, use it as a signal for hotspot code, not as a requirement for every edit.
5. Add or update comments only where the behavior or intended use would benefit from explanation.
6. Keep the comment consistent with the final implementation after refactors.
7. If an old comment is now inaccurate, rewrite or remove it in the same change.
