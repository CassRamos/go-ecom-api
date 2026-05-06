# ecom-go

Simple Go e-commerce API example. This repository contains a small web API, database queries generated with `sqlc`, and migration SQL files for a PostgreSQL backend.

## Project structure

* `cmd/` - application entrypoint
* `internal/` - application internals (adapters, services, handlers)
* `internal/adapters/postgresql/migrations/` - SQL migrations
* `internal/adapters/postgresql/sqlc/` - generated and helper SQL code

## Requirements

* Go 1.20+ (or current stable Go)
* Docker (for sqlc generation on Windows without native install)

## Setup

1. Install Go modules:

```bash
go mod download
```

2. Ensure you have a PostgreSQL instance available and the `DATABASE_URL` environment variable set. Example DSN:

```
postgres://username:password@localhost:5433/your_database
```

## Running locally

* POSIX shells (Git Bash / WSL / macOS / Linux):

```bash
DATABASE_URL="postgres://username:password@localhost:5433/your_database" go run cmd/*.go
```

* PowerShell:

```powershell
$env:DATABASE_URL = "postgres://username:password@localhost:5433/your_database"; go run cmd/*.go
```

* Windows cmd.exe:

```cmd
set DATABASE_URL=postgres://username:password@localhost:5433/your_database & go run cmd/*.go
```

Note: The code will panic if `DATABASE_URL` is not present — see the entrypoint in [cmd/main.go](cmd/main.go#L1-L200).

---

## sqlc (generate) on Windows

If you prefer not to install `sqlc` natively on Windows, you can run it via Docker.

### PowerShell (recommended)

```powershell
docker run --rm -v ${PWD}:/src -w /src sqlc/sqlc generate
```

---

### Git Bash / MSYS (⚠️ requires fix)

Git Bash automatically converts Unix-style paths (e.g. `/src`) into Windows paths, which breaks Docker volume mounting.

#### ✅ Correct command:

```bash
MSYS_NO_PATHCONV=1 docker run --rm -v "$(pwd -W):/src" -w /src sqlc/sqlc generate
```

#### 💡 Why this is needed

Without `MSYS_NO_PATHCONV=1`, Git Bash rewrites:

```
/src → C:/Program Files/Git/src ❌
```

This causes the error:

```
the working directory 'C:/Program Files/Git/src' is invalid
```

Setting `MSYS_NO_PATHCONV=1` disables this behavior and allows Docker to correctly mount:

```
C:/Users/...:/src ✅
```

---

### WSL (best experience)

If you are using WSL with Docker Desktop:

```bash
docker run --rm -v $(pwd):/src -w /src sqlc/sqlc generate
```

No additional fixes are required.

---

## Troubleshooting

### **DATABASE_URL not set**

* Symptom: the program panics with `DATABASE_URL is not set` during startup.
* Cause: `cmd/main.go` checks for the env var and exits when empty. See [cmd/main.go](cmd/main.go#L1-L60).
* Fix: Set the environment variable before running (examples above).

---

### **Database connection errors (pgx.Connect)**

* Symptom: `pgx.Connect` returns an error and the server fails to start.

* Common causes:

  * Database not running or listening on the specified host/port
  * Wrong username/password or database name
  * Firewall or Docker port mapping issues

* Quick checks:

  * Try connecting with `psql` or a DB client using the same DSN
  * Verify the DB is listening on the expected port (5433 in examples)
  * If using Docker for Postgres, ensure the container maps the port (e.g. `-p 5433:5432`)
