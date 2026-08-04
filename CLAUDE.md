# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Development
make up          # Start PostgreSQL + API via Docker Compose
make down        # Stop all services
make tidy        # Tidy go.mod and go.sum

# Code quality
make fmt         # go fmt ./...
make lint        # golangci-lint run
make test        # go test ./...

# Database
make migrate-up      # Run all pending migrations
make migrate-down    # Roll back last migration
make db-connect      # Connect to the running database

# Docs
make docgen      # Regenerate Swagger docs (swag init)
```

To run a single test:
```bash
go test ./internal/usecase/... -run TestAuthUseCase_Login -v
```

To run the API locally without Docker (requires a running Postgres):
```bash
cp .env.example .env   # edit DB_HOST=localhost
go run ./cmd/api
```

## Architecture

Clean Architecture with four layers — dependencies point inward:

```
delivery/http  →  usecase  →  domain  →  infrastructure/persistence
(handlers)        (business)   (entities,   (GORM repos)
                               interfaces)
```

| Layer | Path | Role |
|-------|------|------|
| Delivery | `internal/delivery/http/handlers/` | Fiber handlers; parse request → call usecase → format response |
| Usecase | `internal/usecase/` | Business rules; accepts/returns domain entities |
| Domain | `internal/domain/` | Pure Go structs (`entities/`) and repository interfaces (`repositories/`) |
| Infrastructure | `internal/infrastructure/persistence/` | GORM models (`gormmodel/`) and repository implementations (`repository/`) |
| Config | `internal/config/config.go` | Loads all env vars (with defaults) via `godotenv` |
| Router | `internal/router/router.go` | Wires Fiber routes to handlers; applies JWT middleware on `/admin/*` |

**Key pattern**: GORM models live in `gormmodel/` and are separate from domain entities in `entities/`. Repositories translate between the two. Do not expose GORM models outside the persistence layer.

## Configuration

All configuration is env-based. Copy `.env` and adjust; defaults suit local Docker Compose:

| Variable | Default | Notes |
|----------|---------|-------|
| `APP_ENV` | `local` | Swagger UI is disabled in non-local envs |
| `APP_PORT` | `8080` | |
| `DB_HOST` | `localhost` | Use `postgres` inside Docker Compose |
| `AUTO_MIGRATE` | `true` | GORM AutoMigrate on startup |
| `JWT_SECRET` | `super-secret-change-me` | **Change in production** |
| `GCS_BUCKET_NAME` | *(empty)* | Leave empty to use local storage (`UPLOAD_DIR`) |
| `CORS_ALLOW_ORIGINS` | `http://localhost:3000` | Comma-separated list |
| `DATABASE_URL` | *(empty)* | Full Postgres DSN; when set, takes priority over `DB_HOST`/`DB_PORT`/etc. Not used in production — see note below. |

Production uses the discrete `DB_HOST`/`DB_PORT`/`DB_USER`/`DB_NAME`/`DB_SSLMODE` env vars (pointed at Neon) plus a `DB_PASSWORD` secret, **not** `DATABASE_URL`. A single-`DATABASE_URL`-secret approach was tried first but broke because the Neon password contains characters that aren't safe unescaped inside a `postgresql://` URI, which made the driver silently fall back to an empty/default DSN. Production secrets (`DB_PASSWORD`, `JWT_SECRET`, `ADMIN_PASSWORD`) are pulled from Google Secret Manager by Cloud Run at runtime; `DB_PASSWORD` maps to the Secret Manager secret `Neo_PG_Password`.

## Database & Migrations

Schema is managed two ways:
- **Local development**: GORM AutoMigrate (`AUTO_MIGRATE=true`) creates/alters tables on startup.
- **CI & production**: `golang-migrate` SQL files in `migrations/` are run explicitly before the binary starts.

The `migrations/` directory currently contains a placeholder `000001_init` pair. Add new migration files as `000002_<name>.up.sql` / `000002_<name>.down.sql`.

Admin user is seeded on every startup — idempotent, only creates if absent (see `cmd/api/main.go`).

## File Uploads

- If `GCS_BUCKET_NAME` is set, uploads go to Google Cloud Storage.
- Otherwise, files are written to `UPLOAD_DIR` (default `/app/storage/uploads/covers`) and served at `/uploads/*`.

## API Routes

- **Public**: `GET /api/v1/health`, `/profile`, `/posts`, `/posts/:slug`, `/tags`, `/projects`
- **Auth**: `POST /api/v1/auth/login` (sets `HttpOnly` JWT cookie), `POST /api/v1/auth/logout`
- **Admin** (JWT cookie required): everything under `/api/v1/admin/`
- **Swagger UI**: `/swagger/` — only rendered when `APP_ENV=local`

## CI/CD

- **CI** (`.github/workflows/ci.yml`): runs on every push; spins up Postgres service, runs migrations, generates Swagger, runs tests, builds binary.
- **Deploy** (`.github/workflows/deploy.yml`): triggers after CI passes on `main`; builds Docker image, pushes to Artifact Registry (`us-central1`), deploys to Cloud Run service `devfolio-api-cr` in GCP project `famous-crossing-351110`. Production database is **Neon** (serverless Postgres), connected via discrete `DB_HOST`/`DB_PORT`/`DB_USER`/`DB_NAME`/`DB_SSLMODE` env vars plus the `DB_PASSWORD` secret (`Neo_PG_Password`) — not Cloud SQL, and not a single `DATABASE_URL` secret.
- Authentication to GCP uses **Workload Identity Federation** — no long-lived service account keys in CI.
