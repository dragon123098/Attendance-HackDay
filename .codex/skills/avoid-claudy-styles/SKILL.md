---
name: avoid-claudy-styles
description: make focused CSS and template styling changes that avoid generic AI-looking layouts, preserve the app's visual identity, and keep diffs small by editing existing stylesheets first
---

# avoid claudy styles

## Use this skill
Use this skill when the user asks for styling, layout, visual polish, page redesign, or CSS changes in this project. The goal is to avoid generic "Claudey" AI styling while keeping changes focused, readable, and consistent with the current app.

## Design direction
- Start from the app's actual identity and audience before changing colors or layout. Student pages should feel like a saturated, kid-focused learning app with clear contrast, playful but controlled color, and integrated surfaces. Admin and teacher pages should feel calmer, flatter, and more professional.
- Prefer Duolingo-style clarity for student UI: deep dark backgrounds, bright green/cyan/yellow/coral/purple accents, strong readable labels, visible avatars/icons, and content integrated into the page instead of isolated white bubbles.
- Prefer Codecademy-style restraint for structure: flatter content regions, practical rows, bands, dividers, and focused forms instead of decorative card stacks.
- Use the established font direction: Geist Pixel for headings, page titles, nav labels, big numbers, and major buttons; Rubik for body text, labels, metadata, form fields, and normal UI copy.
- Keep text readable in every state. Check dark mode, light mode, active, hover, disabled, locked, owned, selected, and error states when the touched component has them.
- Use existing CSS variables and local design tokens when possible. Add new colors only when the current tokens cannot express the needed state.
- Keep page content feeling continuous. Avoid turning every section, sentence, or stat into a separate floating white card.

## Avoid
- Beige/cream/rusty-orange dominance unless the product specifically calls for it.
- Large italic serif headlines, tracked-out letter spacing, faux editorial ticker bars, generic SaaS hero layouts, neon/glow-heavy outlines, and decorative gradients that do not serve the UI.
- Rounded rectangle overload. Cards should be purposeful: repeated items, forms, modals, or framed tools. Avoid cards inside cards.
- One-note palettes dominated by one hue family. Student pages need enough color contrast to feel alive without becoming overstimulating.
- Landing-page copy or explanatory text inside app surfaces unless the user explicitly asks for marketing content.
- Inline styles unless unavoidable.
- New CSS files when the existing stylesheet can handle the change cleanly.

## Minimal CSS rules
- Check whether the needed style already exists before adding anything new.
- Prefer updating the current stylesheet over creating a new css file.
- Add a new stylesheet only when the style is clearly page-specific or would make the existing file cluttered.
- Keep selectors narrow and changes minimal.
- Match the current visual language before introducing new spacing, color, or component patterns.
- Avoid redesigning components when a small adjustment solves the issue.
- Put overrides near the related existing rules when practical; append a short focused override block only when the file already has layered overrides.
- Prefer semantic class hooks in templates over brittle descendant selectors, but only add template markup when CSS alone would be fragile.
- Maintain responsive behavior. If a layout changes on desktop, check whether mobile needs a small fallback.
- Avoid text overflow and overlap. Use stable dimensions, wrapping, or smaller type inside compact controls.

## Workflow
1. Inspect the current templates and css files involved in the request.
2. Identify the local visual pattern: student, admin, teacher, auth, shop, avatar, dashboard, or popup.
3. Reuse existing class names, tokens, fonts, surfaces, and component patterns that already solve part of the problem.
4. Edit the current stylesheet first.
5. Add new markup only when it creates a cleaner, more stable hook or improves semantics.
6. Add new CSS only for the missing piece.
7. If the user asks for minimal/no tests, honor that. Otherwise, validate with the smallest relevant automated test command.
8. Return the exact files changed and the reason each change was needed.

## Output style
Keep styling changes compact and practical. Favor small diffs over broad refactors. Call out when a change intentionally avoids a stale AI-design pattern, but do not over-explain obvious CSS.
