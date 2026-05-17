# Devfolio API

Devfolio API is the backend service for a fullstack personal portfolio and blog platform.

This project was built to support a real portfolio/blog use case while also serving as a hands-on backend, cloud infrastructure, and CI/CD practice project. It demonstrates how a Go backend API, PostgreSQL database, cloud object storage, secret management, and automated deployment pipeline work together in a production-like environment.

---

## Production Overview

- Frontend: `https://kraiwit.dev`
- Frontend hosting: Google Cloud Run
- Backend hosting: Google Cloud Run
- Database: Google Cloud SQL for PostgreSQL
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
- Cloud SQL for PostgreSQL
- golang-migrate for schema migrations

### Infrastructure / DevOps

- Docker
- Google Cloud Run
- Google Cloud SQL
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
Cloud SQL PostgreSQL

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
    DB[(Cloud SQL<br/>PostgreSQL)]
    GCS[Google Cloud Storage<br/>Uploads]
    AR[Artifact Registry<br/>Backend Image]
    GHA[GitHub Actions<br/>CI/CD]
    SM[Secret Manager]

    U --> CF --> FE
    FE --> BE
    BE --> DB
    BE --> GCS
    BE --> SM
    GHA --> AR
    AR --> BE
```

---

## Deployment Architecture

The backend is deployed on Google Cloud Run.

### Runtime Components

- Cloud Run service: `devfolio-api-cr`
- Region: `us-central1`
- Runtime service account: `devfolio-backend-run`
- Database: Cloud SQL PostgreSQL
- Cloud SQL connection name:

```text
famous-crossing-351110:us-central1:devfolio-postgres
```

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

DB_HOST=/cloudsql/famous-crossing-351110:us-central1:devfolio-postgres
DB_PORT=5432
DB_USER=postgres
DB_NAME=devfolio
DB_SSLMODE=disable
DB_TIMEZONE=Asia/Bangkok
AUTO_MIGRATE=false

ADMIN_EMAIL=admin_kraiwit@gmail.com
ADMIN_DISPLAY_NAME=Kraiwit

GCS_BUCKET_NAME=kraiwit-devfolio-uploads
GCS_PUBLIC_BASE_URL=https://storage.googleapis.com/kraiwit-devfolio-uploads
```

### Secret Environment Variables

These values are stored in Google Secret Manager and exposed to Cloud Run as environment variables.

```env
DB_PASSWORD=devfolio-db-password:latest
JWT_SECRET=devfolio-jwt-secret:latest
ADMIN_PASSWORD=devfolio-admin-password:latest
```

Do not commit real secret values into the repository.

---

## Secret Management

This project uses Google Secret Manager for sensitive runtime values.

Current secrets:

```text
devfolio-db-password
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

## Cloud SQL

The production database runs on Cloud SQL for PostgreSQL.

### Instance

```text
devfolio-postgres
```

### Database

```text
devfolio
```

### Connection Name

```text
famous-crossing-351110:us-central1:devfolio-postgres
```

### Connecting from Cloud Run

Cloud Run connects to Cloud SQL using the Cloud SQL connection integration.

The application uses:

```env
DB_HOST=/cloudsql/famous-crossing-351110:us-central1:devfolio-postgres
DB_PORT=5432
```

---

## Connecting to Database from Local Machine

Recommended method: Cloud SQL Auth Proxy.

### Start Cloud SQL Auth Proxy

```bash
cloud-sql-proxy famous-crossing-351110:us-central1:devfolio-postgres --port 5433
```

Keep this terminal open while using DBeaver or another database client.

### DBeaver Connection

```text
Host: localhost
Port: 5433
Database: devfolio
Username: postgres
Password: <DB_PASSWORD>
SSL: disable
```

This method does not require adding local public IP addresses to Cloud SQL Authorized Networks.

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
6. Attach Cloud SQL connection
7. Configure runtime environment variables
8. Mount secrets from Secret Manager
9. Run post-deployment health check

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

- Database password is stored in Secret Manager
- JWT secret is stored in Secret Manager
- Admin password is stored in Secret Manager
- Cloud Run uses a dedicated runtime service account
- GitHub Actions uses Workload Identity Federation
- Cloud SQL should be accessed locally through Cloud SQL Auth Proxy
- Avoid opening Cloud SQL public access with `0.0.0.0/0`
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

---

## Operational Notes

### View Cloud Run Logs

```bash
gcloud run services logs read devfolio-api-cr \
  --region us-central1 \
  --project famous-crossing-351110 \
  --limit 100
```

### Create Cloud SQL Backup

```bash
gcloud sql backups create \
  --instance=devfolio-postgres \
  --project famous-crossing-351110
```

### Start Cloud SQL Proxy Locally

```bash
cloud-sql-proxy famous-crossing-351110:us-central1:devfolio-postgres --port 5433
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
- Add least-privilege database user instead of using `postgres`
- Clean up legacy local upload URLs
- Remove old VM infrastructure after migration is stable

---

## Why This Project Exists

This project was created as a practical backend and infrastructure learning project.

The goal was not only to build API features, but also to understand how these parts work together in a real production-like system:

- Backend API design
- Clean Architecture in Go
- PostgreSQL database design
- Cloud SQL migration
- Object storage for uploads
- Secret management
- Container deployment
- CI/CD automation
- Runtime health checks
- Cloud infrastructure operations
