# Attendance-HackDay
Attendance project for Optimization and Innovation project i,e secret combinations
(Sponsored by Shawn Mendix)
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
