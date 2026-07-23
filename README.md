# Attendance-HackDay

Attendance Quest is a Go server-rendered attendance rewards app for students,
teachers, and admins. Students can log in, mark attendance, earn coins, buy shop
items, and customize an avatar with unlocked cosmetics.

## Current Capabilities

- SQL-backed login with in-memory session tokens and role-aware routing.
- SQL-backed student dashboard with attendance status, coin balance, the current avatar, and a Sunday-through-Saturday assignment calendar.
- Classroom assignment templates in `WeeklyAssignmentTemplates` recur by weekday; the dashboard derives their due dates for the current server-local week on each request. This is proof-of-concept mock data, not an assignment editing or completion workflow.
- Student shop with SQL-backed catalog, ownership, coin validation, and atomic purchases.
- Avatar customization with free base avatars, owned cosmetic unlocks, layered visual preview, and SQL-backed saves.
- Student pages include persistent light/dark controls, free background colors, and unlocked special background themes.
- Manual coin adjustments are stored in `ManualCoinAdjustments` without creating transaction records.
- The admin dashboard, User Settings, Add Student, Add Teacher, and classroom create/edit flows use PostgreSQL. `ClassroomMemberships` is the normalized roster source after `Seed_DataBase3.sql`; compatibility writes continue maintaining the legacy classroom columns and tables.
- Teacher and admin dashboard scaffolding plus classroom management routes.
- Provider-neutral roster-source and attendance-destination contracts, encrypted connection persistence, external identity mappings, and versioned attendance-export records. No external provider is connected yet and no integration routes are exposed.

Some teacher/admin reporting and schedule-management flows are still in progress;
see `todo.md` for the remaining project checklist.

## Codebase Map

- `cmd/webserver/main.go` starts the HTTP server on `localhost:4000`.
- `internal/web` contains routes, handlers, in-memory session helpers, and server-rendered student/admin flows.
- `internal/store` contains all PostgreSQL data access, including atomic attendance rewards and shop purchases.
- `internal/domain` contains persisted application models.
- `internal/integrations` contains provider-neutral contracts, capability metadata, provider registration, and AES-GCM credential encryption.
- `internal/view` contains embedded templates, static CSS, and images.

PostgreSQL is the application's only runtime data store. The browser cookie contains an opaque token; its short-lived session record remains in application memory and references the SQL `Users.UserID`.

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
`ManualCoinAdjustments`. That amount is added to the starting balance and
the sum of `Transactions`.

Added CodeQL


## Database setup

1. Start the PostgreSQL container. The Compose environment creates the
   `attendancehackday` database and runs `init.sql` automatically on the first
   start of a new database volume:

```powershell
docker compose up -d
```

2. Verify PostgreSQL and the application database are ready:

```powershell
docker compose exec postgres pg_isready -U attendance -d attendancehackday
```

The application defaults to
`postgres://attendance:Password123!@localhost:5433/attendancehackday?sslmode=disable`.
Set `DATABASE_URL` to override that connection string.

`INTEGRATION_CREDENTIAL_KEY` optionally enables encrypted provider credential
storage. Set it to a base64-encoded 32-byte AES-256 key. The main application
continues to start when the value is absent or invalid, but integration
credential reads and writes remain disabled. Do not change or discard a key
after credentials have been stored unless a future key-rotation process has
re-encrypted those records.


## Check Database in DBeaver

1. Open a new connection

2. Select PostgreSQL

3. Fill in the Info
 Host: localhost
 Port: 5433
 Database: attendancehackday
 User name: attendance
 Password: Password123!

4. Click test connection. If it passes then click finish


## Seed the DataBase

1. Make sure the database is up

    ```powershell
    docker compose up -d
    ```

2. Seed the base data with the following command.

    ```powershell
    Get-Content -Raw .\Seed_DataBase.sql | docker exec -i attendance-postgres psql --set=ON_ERROR_STOP=1 --username=attendance --dbname=attendancehackday
    ```

    Bash / WSL:

    ```bash
    docker exec -i attendance-postgres psql --set=ON_ERROR_STOP=1 --username=attendance --dbname=attendancehackday < Seed_DataBase.sql
    ```

3. Apply the idempotent delta seed for the latest student records, image path metadata, and recurring weekly assignment templates. This includes the final records migrated from the retired JSON store.

    ```powershell
    Get-Content -Raw .\Seed_DataBase2.sql | docker exec -i attendance-postgres psql --set=ON_ERROR_STOP=1 --username=attendance --dbname=attendancehackday
    ```

    Bash / WSL:

    ```bash
    docker exec -i attendance-postgres psql --set=ON_ERROR_STOP=1 --username=attendance --dbname=attendancehackday < Seed_DataBase2.sql
    ```

4. Apply the additive integration foundation migration before deploying code
   that reads `ClassroomMemberships`. It backfills existing classroom and
   attendance data while retaining `Users.ClassroomID`, `Classrooms.TeacherID`,
   `ClassroomStudents`, and `AttendanceRecords` for compatibility.

    ```powershell
    Get-Content -Raw .\Seed_DataBase3.sql | docker exec -i attendance-postgres psql --set=ON_ERROR_STOP=1 --username=attendance --dbname=attendancehackday
    ```

    Bash / WSL:

    ```bash
    docker exec -i attendance-postgres psql --set=ON_ERROR_STOP=1 --username=attendance --dbname=attendancehackday < Seed_DataBase3.sql
    ```

The application does not apply `Seed_DataBase3.sql` automatically. Part 1 only
adds the integration foundation; Canvas, SIS, and attendance-export adapters
remain unimplemented.
