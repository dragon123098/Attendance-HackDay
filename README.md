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
- The admin Add Student flow reads classrooms from SQL Server and creates the student there.
- Teacher and admin dashboard scaffolding plus classroom management routes.

Some teacher/admin reporting and schedule-management flows are still in progress;
see `todo.md` for the remaining project checklist.

## Codebase Map

- `cmd/webserver/main.go` starts the HTTP server on `localhost:4000`.
- `internal/web` contains routes, handlers, auth/session helpers, persistence, and student feature logic.
- `internal/store` contains SQL Server data access for flows that have moved off `data/data.json`.
- `internal/domain` contains persisted application models.
- `internal/view` contains embedded templates, static CSS, and images.
- `data/data.json` is still the local JSON data store for most app flows.

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


## Database setup

1. Start the SQL Server container:

```powershell
docker compose up -d
```

2. Run the initialization script from the repository root. Use the relative path to `init.sql`, not a local user path.

PowerShell:

```powershell
docker run --rm -v "${PWD}\init.sql:/tmp/init.sql" --entrypoint "/opt/mssql-tools/bin/sqlcmd" mcr.microsoft.com/mssql-tools -S host.docker.internal,1433 -U SA -P "Password123!" -i /tmp/init.sql
```

> Note: The SQL Server Linux image does not automatically execute files mounted into `/docker-entrypoint-initdb.d/`, so we run the script manually from a SQL tools container.

Bash / WSL:

```bash
docker run --rm -v "$(pwd)/init.sql:/tmp/init.sql" --entrypoint "/opt/mssql-tools/bin/sqlcmd" mcr.microsoft.com/mssql-tools -S host.docker.internal,1433 -U SA -P "Password123!" -i /tmp/init.sql
```

3. Verify the database exists:

```powershell
docker run --rm --entrypoint "/opt/mssql-tools/bin/sqlcmd" mcr.microsoft.com/mssql-tools -S host.docker.internal,1433 -U SA -P "Password123!" -Q "SELECT name FROM sys.databases WHERE name='AttendanceHackday';"
```

If the command succeeds, it should return `AttendanceHackday` and `(1 rows affected)`.


## Check Database in DBeaver

1. Open a new connection

2. Select SQL Server

3. Fill in the Info
 Host: localhost
 Port: 1433
 Database: Leave blank for now, or master
 User name: sa
 Password: Password123!

4. Click test connection. If it passes then click finish


## Seed the DataBase

1. Make sure the database is up

    ```powershell
    docker compose up -d
    ```

2. Seed the base data with the following command.

    ```powershell
    docker run --rm -v "${PWD}\Seed_DataBase.sql:/tmp/Seed_DataBase.sql" --entrypoint "/opt/mssql-tools/bin/sqlcmd" mcr.microsoft.com/mssql-tools -S host.docker.internal,1433 -U SA -P "Password123!" -i /tmp/Seed_DataBase.sql
    ```

    Bash / WSL:

    ```bash
    docker run --rm -v "$(pwd)/Seed_DataBase.sql:/tmp/Seed_DataBase.sql" --entrypoint "/opt/mssql-tools/bin/sqlcmd" mcr.microsoft.com/mssql-tools -S host.docker.internal,1433 -U SA -P "Password123!" -i /tmp/Seed_DataBase.sql
    ```

3. Apply the delta seed for newer JSON data and image path metadata.

    ```powershell
    docker run --rm -v "${PWD}\Seed_DataBase2.sql:/tmp/Seed_DataBase2.sql" --entrypoint "/opt/mssql-tools/bin/sqlcmd" mcr.microsoft.com/mssql-tools -S host.docker.internal,1433 -U SA -P "Password123!" -i /tmp/Seed_DataBase2.sql
    ```

    Bash / WSL:

    ```bash
    docker run --rm -v "$(pwd)/Seed_DataBase2.sql:/tmp/Seed_DataBase2.sql" --entrypoint "/opt/mssql-tools/bin/sqlcmd" mcr.microsoft.com/mssql-tools -S host.docker.internal,1433 -U SA -P "Password123!" -i /tmp/Seed_DataBase2.sql
    ```
