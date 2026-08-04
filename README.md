# Devfolio API

Devfolio API is the backend service for a fullstack personal portfolio and blog platform.

This project was built to support a real portfolio/blog use case while also serving as a hands-on backend, cloud infrastructure, and CI/CD practice project. It demonstrates how a Go backend API, PostgreSQL database, cloud object storage, secret management, and automated deployment pipeline work together in a production-like environment.

---

## Production Overview

- Frontend: `https://kraiwit.dev`
- Frontend hosting: Google Cloud Run
- Backend hosting: Google Cloud Run
- Database: Neon (serverless PostgreSQL) — previously Google Cloud SQL, migrated over
- Upload storage: Google Cloud Storage
- Container registry: Google Artifact Registry
- Secrets: Google Secret Manager
- DNS / Domain management: Cloudflare
- CI/CD: GitHub Actions

---

## Stack

### Backend

- Go
- Fiber
- Clean Architecture
- GORM
- JWT authentication
- Swagger / OpenAPI

### Database

- PostgreSQL
- Neon (serverless Postgres, autoscaling/autosuspend compute)
- golang-migrate for schema migrations

### Infrastructure / DevOps

- Docker
- Google Cloud Run
- Neon
- Google Cloud Storage
- Google Secret Manager
- Google Artifact Registry
- GitHub Actions
- Workload Identity Federation

---

## Features

### Public Features

- Health check endpoint
- Public profile API
- Featured projects API
- Tags API
- Published blog posts API
- Blog post detail by slug

### Admin Features

- Admin login with JWT
- Admin post management
- Create, edit, delete, and list posts
- Tag management
- Profile update
- Cover image upload
- Cloud Storage-based image hosting

---

## API Endpoints

### Public Endpoints

- `GET /api/v1/health`
- `GET /api/v1/profile`
- `GET /api/v1/projects`
- `GET /api/v1/tags`
- `GET /api/v1/posts`
- `GET /api/v1/posts/:slug`

### Admin Endpoints

- `POST /api/v1/auth/login`
- `POST /api/v1/auth/logout`
- `GET /api/v1/admin/posts`
- `GET /api/v1/admin/posts/:id`
- `POST /api/v1/admin/posts`
- `PUT /api/v1/admin/posts/:id`
- `DELETE /api/v1/admin/posts/:id`
- `POST /api/v1/admin/tags`
- `PUT /api/v1/admin/profile`
- `POST /api/v1/admin/uploads/cover`

---

## Project Structure

```text
cmd/api
internal/
  config/
  database/
  router/
  delivery/http/
    handlers/
    middleware/
    response/
  domain/
    entities/
    repositories/
  usecase/
  infrastructure/persistence/
    gormmodel/
    repository/
pkg/
  auth/
  utils/
migrations/
docs/
```

---

## High-Level Architecture

```text
User Browser
    |
Cloudflare DNS
    |
Frontend Cloud Run
    |
Backend Cloud Run
    |
Neon PostgreSQL (serverless)

Backend Cloud Run
    |
Google Cloud Storage
```

---

## Architecture Diagram

```mermaid
flowchart LR
    U[User Browser]
    CF[Cloudflare DNS]
    FE[Cloud Run<br/>Next.js Frontend]
    BE[Cloud Run<br/>Go Fiber API]
    SQL[(Neon<br/>Serverless PostgreSQL)]
    GCS[Google Cloud Storage<br/>Uploads]
    SM[Secret Manager]
    AR[Artifact Registry<br/>Backend Image]
    GHA[GitHub Actions<br/>CI/CD]
    OBS[Cloud Logging / Monitoring<br/>Observability & Alerts]

    U --> CF --> FE
    FE --> BE

    BE --> SQL
    BE --> GCS
    BE --> SM

    GHA --> AR
    AR --> BE

    FE -. logs / metrics .-> OBS
    BE -. logs / metrics .-> OBS
    SQL -. metrics .-> OBS
```

---

## Deployment Architecture

The backend is deployed on Google Cloud Run.

### Runtime Components

- Cloud Run service: `devfolio-api-cr`
- Region: `us-central1`
- Runtime service account: `devfolio-backend-run`
- Database: Neon (serverless PostgreSQL), accessed over the public internet via the pooled connection endpoint (no Cloud SQL Unix socket, no VPC connector needed)
- Neon project: `devfolio` (org `Kraiwit`), branch `production`

### Cloud Run Settings

```text
Port: 8080
CPU: 1
Memory: 512Mi
Min instances: 0
Max instances: 2
Ingress: All
Authentication: Allow unauthenticated
```

---

## Environment Variables

### Normal Environment Variables

```env
APP_ENV=production
APP_PORT=8080
AUTO_MIGRATE=false

ADMIN_EMAIL=admin_kraiwit@gmail.com
ADMIN_DISPLAY_NAME=Kraiwit

GCS_BUCKET_NAME=kraiwit-devfolio-uploads
GCS_PUBLIC_BASE_URL=https://storage.googleapis.com/kraiwit-devfolio-uploads
```

Production connects to Neon using discrete `DB_HOST`/`DB_PORT`/`DB_USER`/`DB_NAME`/`DB_SSLMODE` env vars, **not** a single `DATABASE_URL` secret. (A single-`DATABASE_URL`-secret approach was tried first but broke: the Neon password contains characters that aren't safe unescaped inside a `postgresql://` URI, which made the driver silently fall back to an empty/default DSN instead of failing loudly. Discrete fields avoid that URI-parsing problem entirely.) `DATABASE_URL` is still supported by `buildDSN` in `internal/database/postgres.go` and takes priority when set, but it is left unset in production.

```env
DB_HOST=ep-round-king-ajnprccc-pooler.c-3.us-east-2.aws.neon.tech
DB_PORT=5432
DB_USER=neondb_owner
DB_NAME=neondb
DB_SSLMODE=require
```

### Secret Environment Variables

These values are stored in Google Secret Manager and exposed to Cloud Run as environment variables.

```env
DB_PASSWORD=Neo_PG_Password:latest
JWT_SECRET=devfolio-jwt-secret:latest
ADMIN_PASSWORD=devfolio-admin-password:latest
```

Do not commit real secret values into the repository.

---

## Secret Management

This project uses Google Secret Manager for sensitive runtime values.

Current secrets:

```text
Neo_PG_Password
devfolio-jwt-secret
devfolio-admin-password
```

Cloud Run reads these secrets through the runtime service account:

```text
devfolio-backend-run@famous-crossing-351110.iam.gserviceaccount.com
```

Required IAM role:

```text
Secret Manager Secret Accessor
```

---

## Database (Neon)

The production database runs on [Neon](https://neon.tech) — serverless PostgreSQL with autoscaling compute and autosuspend (scale-to-zero) after 5 minutes of inactivity on the Free plan.

### Project

```text
Org: Kraiwit
Project: devfolio
Branch: production
Database: neondb
```

### Connecting from Cloud Run

Cloud Run connects directly over the internet using Neon's **pooled** connection endpoint (hostname ending in `-pooler`), which fronts PgBouncer and copes much better with Cloud Run's many short-lived serverless connections than a direct (unpooled) endpoint would. There is no Cloud SQL-style Unix socket or VPC connector involved. The connection is configured as discrete fields rather than a single connection-string secret:

```text
DB_HOST=ep-round-king-ajnprccc-pooler.c-3.us-east-2.aws.neon.tech
DB_PORT=5432
DB_USER=neondb_owner
DB_NAME=neondb
DB_SSLMODE=require
DB_PASSWORD=<from Secret Manager secret Neo_PG_Password>
```

Discrete fields are used instead of a single `DATABASE_URL` because the Neon password contains characters (e.g. `+`, `/`, `=`) that aren't safe unescaped inside a `postgresql://user:password@host/db` URI — when it broke, the driver didn't error loudly, it silently fell back to an empty/default DSN (`host=/tmp`, `user=root`, empty dbname), which took a while to diagnose. Get the exact current host/user from the Neon console: project `devfolio` → Connect → make sure **Connection pooling** is toggled on.

### Free plan limits to watch

- **100 compute hours (CU-hrs) per month.** This is wall-clock active compute time, not query count — if something pings the API/DB often enough to prevent the 5-minute autosuspend from ever completing, the compute effectively never sleeps and the monthly allowance can be exhausted in days rather than the full month. Check Neon's Monitoring tab periodically for a "should be spiky, not a flat line" sanity check.
- Once the monthly allowance is used up, compute stops starting until the next billing cycle (or until the plan is upgraded) — the API will run in degraded mode (see `db_guard.go` / `RequireDB`) rather than crash, but every DB-backed route returns 503 until then.

---

## Connecting to Database from Local Machine

Connect directly with `psql` (or any Postgres client) using the Neon connection string — no proxy needed, since Neon is reachable over the public internet with TLS (`sslmode=require`):

```bash
psql "postgresql://neondb_owner:<password>@ep-round-king-ajnprccc-pooler.c-3.us-east-2.aws.neon.tech/neondb?sslmode=require"
```

### DBeaver / GUI client connection

```text
Host: ep-round-king-ajnprccc-pooler.c-3.us-east-2.aws.neon.tech
Port: 5432
Database: neondb
Username: neondb_owner
Password: <from Neon console>
SSL mode: require
```

Prefer the **unpooled** endpoint (no `-pooler` suffix, shown in Neon's Connect dialog under "Direct connection") for one-off admin/GUI sessions, and reserve the pooled endpoint for the running Cloud Run service.

---

## Upload Handling

Uploads are stored in Google Cloud Storage.

Bucket:

```text
kraiwit-devfolio-uploads
```

Public base URL:

```text
https://storage.googleapis.com/kraiwit-devfolio-uploads
```

Upload benefits:

- Avoids mixed-content issues
- Removes dependency on local VM disk
- Makes the backend stateless
- Fits Cloud Run deployment model
- Supports production-like file hosting

Older blog posts may still contain legacy local image paths from the previous VM deployment. New uploads should use Google Cloud Storage URLs.

---

## Local Development

### 1. Clone repository

```bash
git clone https://github.com/poppbsfs4za/devfolio-api.git
cd devfolio-api
```

### 2. Copy environment file

```bash
cp .env.example .env
```

### 3. Start local services

```bash
docker compose up -d
```

### 4. Install dependencies

```bash
go mod tidy
```

### 5. Run migrations

```bash
make migrate-up
```

or:

```bash
migrate -path migrations \
  -database "postgres://postgres:postgres@localhost:5432/devfolio?sslmode=disable" \
  up
```

### 6. Run API locally

```bash
go run ./cmd/api
```

Local API:

```text
http://localhost:8080
```

---

## Useful Commands

### Run locally

```bash
make up
```

### Stop local services

```bash
make down
```

### Run database migrations

```bash
make migrate-up
```

### Build Docker image

```bash
make docker-build
```

### Run tests

```bash
go test ./...
```

### Generate Swagger docs

```bash
swag init -g cmd/api/main.go --parseDependency --parseInternal
```

---

## Example Requests

### Health Check

```bash
curl -i http://localhost:8080/api/v1/health
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"changeme123"}'
```

### Get Published Posts

```bash
curl -X GET http://localhost:8080/api/v1/posts
```

### Create Post

```bash
curl -X POST http://localhost:8080/api/v1/admin/posts \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "CI/CD Basics for Backend Developers",
    "summary": "A practical intro to CI/CD",
    "content": "# Hello\nThis is my first post.",
    "status": "published",
    "tags": ["Go", "CI/CD"]
  }'
```

### Upload Cover Image

```bash
curl -X POST http://localhost:8080/api/v1/admin/uploads/cover \
  -H "Authorization: Bearer <TOKEN>" \
  -F "file=@/path/to/image.png"
```

---

## CI/CD

This project uses GitHub Actions for CI and Cloud Run deployment.

### CI Workflow

The CI workflow runs on every branch push and pull request to `main`.

CI steps:

1. Start PostgreSQL service in GitHub Actions
2. Install migration tool
3. Install Swagger generator
4. Download Go dependencies
5. Verify `go.mod` and `go.sum`
6. Wait for PostgreSQL
7. Run database migrations
8. Generate Swagger docs
9. Run tests
10. Build application

### Deployment Workflow

The deployment workflow runs after CI succeeds on `main`.

Deployment steps:

1. Authenticate to Google Cloud using Workload Identity Federation
2. Configure Docker authentication for Artifact Registry
3. Build Docker image
4. Push image to Artifact Registry
5. Deploy backend to Cloud Run
6. Configure runtime environment variables (including discrete `DB_HOST`/`DB_PORT`/`DB_USER`/`DB_NAME`/`DB_SSLMODE` for Neon)
7. Mount secrets from Secret Manager, including `DB_PASSWORD` (Neon password, secret `Neo_PG_Password`)
8. Run post-deployment health check

### Deployment Target

```text
Cloud Run service: devfolio-api-cr
Region: us-central1
Artifact Registry repository: devfolio
Image name: devfolio-api
```

---

## GitHub Actions Authentication

GitHub Actions authenticates to Google Cloud using Workload Identity Federation.

This avoids storing long-lived Google Cloud service account keys in GitHub Secrets.

GitHub secret required:

```text
GCP_WORKLOAD_IDENTITY_PROVIDER
```

Deployer service account:

```text
devfolio-github-deployer@famous-crossing-351110.iam.gserviceaccount.com
```

Runtime service account:

```text
devfolio-backend-run@famous-crossing-351110.iam.gserviceaccount.com
```

---

## Runtime Health Check

Production health check endpoint:

```text
GET /api/v1/health
```

Example:

```bash
curl -i https://devfolio-api-cr-294009483204.us-central1.run.app/api/v1/health
```

Expected response:

```json
{
  "data": {
    "status": "ok"
  }
}
```

---

## Security Notes

- The database host/user/name/sslmode are plain env vars; only the password is a secret, stored in Secret Manager as `DB_PASSWORD` (secret `Neo_PG_Password`)
- JWT secret is stored in Secret Manager
- Admin password is stored in Secret Manager
- Cloud Run uses a dedicated runtime service account
- GitHub Actions uses Workload Identity Federation
- Neon connections always use `sslmode=require`; there is no IP-allowlist on the Free plan, so treat the password itself as the sole secret boundary — rotate the Neon password immediately if it's ever exposed
- Do not commit `.env` or real secret values

---

## Migration History

The backend was originally deployed on a GCP Compute Engine VM with PostgreSQL running in Docker on the same VM.

The project was later migrated to a more production-like cloud architecture:

```text
Before:
Frontend Cloud Run → Backend VM → PostgreSQL Docker

After:
Frontend Cloud Run → Backend Cloud Run → Cloud SQL PostgreSQL
                                  ↓
                           Google Cloud Storage
```

The migration included:

- Moving PostgreSQL data from Docker PostgreSQL to Cloud SQL
- Deploying the backend to Cloud Run
- Updating frontend environment variables to point to the Cloud Run backend
- Moving sensitive values to Secret Manager
- Updating GitHub Actions deployment from VM SSH to Cloud Run
- Keeping GCS as the upload storage layer

### Cloud SQL → Neon (August 2026)

The database was migrated a second time, from Cloud SQL to Neon:

```text
Before:
Frontend Cloud Run → Backend Cloud Run → Cloud SQL PostgreSQL

After:
Frontend Cloud Run → Backend Cloud Run → Neon PostgreSQL (serverless)
```

The actual data migration (schema + all rows: `users`, `profiles`, `projects`, `tags`, `posts`, `post_tags`) had already been done by hand earlier, directly against Neon's `devfolio` project (branch `production`, database `neondb`) — the Cloud SQL instance (`devfolio-postgres`) was deleted afterwards and no longer exists. What was actually out of sync was the **deploy pipeline**, which still had `deploy.yml` configured to attach a Cloud SQL instance that no longer existed — every Cloud Run revision was failing to connect at startup (`Internal error looking up Cloud SQL instance ...`) and running in degraded mode (`db_available=false`, all DB-backed routes returning 503) until this was fixed.

The fix included:

- Removing the `--set-cloudsql-instances` Cloud SQL connection attachment from the Cloud Run deploy step (replaced with `--clear-cloudsql-instances`) and dropping the old Cloud SQL env vars/secret
- Pointing Cloud Run at Neon using discrete `DB_HOST`/`DB_PORT`/`DB_USER`/`DB_NAME`/`DB_SSLMODE` env vars plus a `DB_PASSWORD` secret (reusing the existing but previously-unused `Neo_PG_Password` secret in Secret Manager, after resetting the Neon password and adding a new secret version)
- A single-`DATABASE_URL`-secret approach was tried first and abandoned: it produced a connection error where the driver fell back to an empty/default DSN (`host=/tmp`, `user=root`, empty dbname) rather than a clear auth failure, traced to unescaped special characters (`+`, `/`, `=`) in the freshly-reset Neon password breaking `postgresql://` URI parsing. Discrete fields sidestep URI parsing entirely.
- No application code changes were needed — `buildDSN` in `internal/database/postgres.go` already supported both `DATABASE_URL` and the discrete fields
- No data export/import was needed — the data was already on Neon

Separately: this same Neon project (`devfolio`) had its monthly 100 CU-hr compute allowance exhausted within the first few days of the billing cycle, before this fix — traced to the compute apparently never getting a full 5-minute idle window to autosuspend (Monitoring showed continuous usage rather than the expected spiky pattern). Even with the deploy pipeline fixed, Cloud Run will keep hitting 503s until the Neon compute allowance resets (next billing cycle) or the plan is upgraded — fixing the pipeline alone does not restore live service while the quota is exhausted. Keep an eye on Neon's Monitoring tab going forward for that "flat line, never suspends" pattern.

---

## Operational Notes

### View Cloud Run Logs

```bash
gcloud run services logs read devfolio-api-cr \
  --region us-central1 \
  --project famous-crossing-351110 \
  --limit 100
```

### Back Up the Neon Database

Neon keeps continuous point-in-time restore history automatically (retention depends on plan), but for an explicit snapshot before risky changes, create a child branch (instant, copy-on-write) from the Neon console or CLI:

```bash
neonctl branches create --project-id <project-id> --name backup-$(date +%Y%m%d)
```

Or take a manual `pg_dump` for an offline copy:

```bash
pg_dump "postgresql://neondb_owner:<password>@ep-round-king-ajnprccc-pooler.c-3.us-east-2.aws.neon.tech/neondb?sslmode=require" \
  --format=custom --file=devfolio-backup-$(date +%Y%m%d).dump
```

### Test Production API

```bash
curl -i https://devfolio-api-cr-294009483204.us-central1.run.app/api/v1/health
curl -i https://devfolio-api-cr-294009483204.us-central1.run.app/api/v1/posts
```

---

## Known Notes

- Some legacy uploaded image URLs may still point to old local VM paths.
- New uploads are stored in Google Cloud Storage.
- Cloud Run should remain stateless.
- VM infrastructure can be kept temporarily as rollback backup after migration.
- Admin password rotation requires updating the stored password hash in the database if the admin user already exists.

---

## Next Improvements

- Add centralized monitoring and alerting
- Add structured logging and request tracing
- Add automated database backup policy review
- Add rollback strategy using Cloud Run revisions
- Add more unit and integration tests
- Add alerting on Neon compute-hour usage so a runaway consumer is caught early, not after the monthly quota is already exhausted
- Clean up legacy local upload URLs
- Remove old VM infrastructure after migration is stable

---

## Why This Project Exists

This project was created as a practical backend and infrastructure learning project.

The goal was not only to build API features, but also to understand how these parts work together in a real production-like system:

- Backend API design
- Clean Architecture in Go
- PostgreSQL database design
- Cloud SQL and Neon (serverless Postgres) migrations
- Object storage for uploads
- Secret management
- Container deployment
- CI/CD automation
- Runtime health checks
- Cloud infrastructure operations
