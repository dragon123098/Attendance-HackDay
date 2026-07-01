# Attendance-HackDay

Attendance Quest is a Go server-rendered attendance rewards app for students,
teachers, and admins. Students can log in, mark attendance, earn coins, buy shop
items, and customize an avatar with unlocked cosmetics.

## Current Capabilities

- Session-based login and logout with role-aware routing.
- Student dashboard with attendance status, coin balance, schedule data, and the current avatar.
- Student shop with seeded visual cosmetics, purchasable pastel background themes, coin validation, purchase persistence, and owned-item display.
- Avatar customization with free base avatars, owned cosmetic unlocks, layered visual preview, and persisted saves.
- Student pages include persistent light/dark controls, free background colors, and unlocked special background themes.
- Manual coin adjustments can be added in `data/data.json` without writing transaction records.
- Teacher and admin dashboard scaffolding plus classroom management routes.

Some teacher/admin reporting and schedule-management flows are still in progress;
see `todo.md` for the remaining project checklist.

## Codebase Map

- `cmd/webserver/main.go` starts the HTTP server on `localhost:4000`.
- `internal/web` contains routes, handlers, auth/session helpers, persistence, and student feature logic.
- `internal/domain` contains persisted application models.
- `internal/view` contains embedded templates, static CSS, and images.
- `data/data.json` is the local JSON data store loaded and saved by the app.

## Local Development

Run automated validation with:

```sh
go test ./...
```

To run the app manually:

```sh
go run ./cmd/webserver
```

To manually add or subtract coins, edit the `manual_coin_adjustments` map in
`data/data.json`. For example, `"student1": 25` adds 25 coins on top of the
student's starting coins and transaction history.

Added CodeQL
