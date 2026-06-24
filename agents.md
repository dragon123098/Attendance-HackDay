# AGENTS.md

### Do
- use **Go** for all backend logic
- use the **standard library first** (net/http, html/template, context, database/sql)
- use **HTMX** only as a backup tool for targeted interactivity when plain server-rendered HTML is not enough
- prefer existing reusable packages from **`.codex/skills`**  when they clearly fit
- if a skill is used, explicitly state which skill is being used
- keep handlers small and focused (one responsibility per handler)
- return full pages by default; return **HTML fragments** only when an HTMX enhancement is actually being used
- default to **small components/templates**
- default to **small diffs**
- validate changes with automated tests only; do not start the webserver or run the app for validation unless explicitly asked
- manual validation is performed by the user after Codex completes its changes
- if a dependency download fails, stop immediately and tell the user exactly what failed; do not continue trying to make it work

### Don't
- do not build a SPA or client-side state machine
- do not introduce frontend frameworks (React, Vue, etc.)
- do not fetch data directly in templates
- do not add heavy Go dependencies without approval
- do not start the app, run the webserver, or attempt browser/manual validation unless explicitly asked
- do not retry failed dependency downloads with workarounds, alternate install methods, or sandbox bypass attempts unless explicitly asked

---

### Safety and permissions

Allowed without prompt:
- read files, list files

Ask first:
- adding Go modules / new dependencies
- deleting files or changing permissions
- database migrations
- running the webserver, launching the app, or performing manual/browser validation
- moving or deleting folders

---

### HTMX conventions
- prefer normal server-rendered page flows first; add HTMX only as a fallback when it materially improves UX
- detect HTMX via `HX-Request: true`
- POST/PUT/DELETE:
  - validate input
  - enforce CSRF
  - return fragments or `204 No Content`
- redirects:
  - use `HX-Redirect` header for HTMX flows
- avoid client-side state; let the server be the source of truth

---

### Styling system
- use semantic HTML elements first (`main`, `section`, `article`, `nav`)
- no inline styles unless unavoidable

---

### API usage
- server-rendered HTML is the default interface
- JSON APIs are allowed only when necessary
- keep API logic in handlers or service layers
- document non-HTML endpoints in `./api/docs/*.md`

---

### PR checklist 
- HTMX responses verified manually when HTMX is used
- diff is small and focused with a clear summary

---

### When stuck
- ask one clarifying question
- propose a short plan before implementing
- prefer a small draft PR over a large speculative change

---

### Philosophy
- server-rendered first
- progressive enhancement, not JS-heavy UIs
- boring, readable Go over clever abstractions
- simple tools, explicit behavior

### Documentation
- use the `project-documentation` skill when updating `README.md` or architecture documentation
- keep `README.md` up to date when app behavior, setup, structure, or major workflows change
- document the current implemented behavior, and call out partial or placeholder flows explicitly
- document functions when function name is not enough
- document files to explain their purpose when created
- document in README how to run the app if process changed
- document in README any new significant changes to the app