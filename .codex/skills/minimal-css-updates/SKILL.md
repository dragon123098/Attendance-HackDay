---
name: minimal-css-updates
description: add the smallest useful css changes and prefer updating existing stylesheets instead of creating new ones; use when the user asks for styling, spacing, visibility, or layout tweaks in this project
---

# minimal css updates

## Use this skill
Use this skill when the user wants styling changes in this project and the goal is to keep css additions small, consistent, and placed in the existing stylesheet when possible.

## Rules
- Check whether the needed style already exists before adding anything new.
- Prefer updating the current stylesheet over creating a new css file.
- Add a new stylesheet only when the style is clearly page-specific or would make the existing file cluttered.
- Keep selectors narrow and changes minimal.
- Match the current visual language before introducing new spacing, color, or component patterns.
- Avoid redesigning components when a small adjustment solves the issue.

## Workflow
1. Inspect the current templates and css files involved in the request.
2. Reuse any existing class names or style patterns that already solve part of the problem.
3. Edit the current stylesheet first.
4. Add new css only for the missing piece.
5. Return the exact files changed and the reason each change was needed.

## Output style
Keep styling changes compact and practical. Favor small diffs over broad refactors.
