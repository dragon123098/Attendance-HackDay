# Attendance-HackDay

Attendance Quest is a Go server-rendered attendance rewards app for students,
teachers, and admins. Students can log in, mark attendance, earn coins, buy shop
items, and customize an avatar with unlocked cosmetics.

## Current Capabilities

- SQL-backed login with in-memory session tokens and role-aware routing.
- SQL-backed student dashboard with attendance status, coin balance, the current avatar, and a Sunday-through-Saturday assignment calendar.
- Classroom assignment templates in `dbo.WeeklyAssignmentTemplates` recur by weekday; the dashboard derives their due dates for the current server-local week on each request. This is proof-of-concept mock data, not an assignment editing or completion workflow.
- Student shop with SQL-backed catalog, ownership, coin validation, and atomic purchases.
- Avatar customization with free base avatars, owned cosmetic unlocks, layered visual preview, and SQL-backed saves.
- Student pages include persistent light/dark controls, free background colors, and unlocked special background themes.
- Manual coin adjustments are stored in `dbo.ManualCoinAdjustments` without creating transaction records.
- The admin dashboard, User Settings, Add Student, Add Teacher, and classroom create/edit flows use SQL Server; dashboard and edit classroom pages load rosters from `ClassroomStudents`.
- Teacher and admin dashboard scaffolding plus classroom management routes.

Some teacher/admin reporting and schedule-management flows are still in progress;
see `todo.md` for the remaining project checklist.

## Codebase Map

- `cmd/webserver/main.go` starts the HTTP server on `localhost:4000`.
- `internal/web` contains routes, handlers, in-memory session helpers, and server-rendered student/admin flows.
- `internal/store` contains all SQL Server data access, including atomic attendance rewards and shop purchases.
- `internal/domain` contains persisted application models.
- `internal/view` contains embedded templates, static CSS, and images.

SQL Server is the application's only runtime data store. The browser cookie contains an opaque token; its short-lived session record remains in application memory and references the SQL `Users.UserID`.

## Local Development

Run automated validation with:

```sh
go test ./...
```

To run the app manually:

```sh
go run ./cmd/webserver
```

To manually add or subtract coins, insert or update the student's amount in
`dbo.ManualCoinAdjustments`. That amount is added to the starting balance and
the sum of `dbo.Transactions`.

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

3. Apply the idempotent delta seed for the latest student records, image path metadata, and recurring weekly assignment templates. This includes the final records migrated from the retired JSON store.

    ```powershell
    docker run --rm -v "${PWD}\Seed_DataBase2.sql:/tmp/Seed_DataBase2.sql" --entrypoint "/opt/mssql-tools/bin/sqlcmd" mcr.microsoft.com/mssql-tools -S host.docker.internal,1433 -U SA -P "Password123!" -i /tmp/Seed_DataBase2.sql
    ```

    Bash / WSL:

    ```bash
    docker run --rm -v "$(pwd)/Seed_DataBase2.sql:/tmp/Seed_DataBase2.sql" --entrypoint "/opt/mssql-tools/bin/sqlcmd" mcr.microsoft.com/mssql-tools -S host.docker.internal,1433 -U SA -P "Password123!" -i /tmp/Seed_DataBase2.sql
    ```
