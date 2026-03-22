# Devfolio API

A clean-architecture backend starter for a personal developer platform. It supports profile data, projects, tags, blog posts, and admin authentication so you can publish blog content without redeploying the application for each new post.

## Stack
- Go + Fiber
- PostgreSQL
- GORM
- JWT auth
- Docker / Docker Compose
- GitHub Actions CI

## Architecture
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
```

## Features in this scaffold
- Health check endpoint
- Admin login endpoint
- Public profile endpoint
- Public featured projects endpoint
- Public tags endpoint
- Public published-post list and detail endpoints
- Admin create/update/delete post endpoints
- Admin create tag endpoint
- Admin upsert profile endpoint
- Automatic admin seed on startup

## Environment
Copy `.env.example` to `.env` and adjust values.

## Run locally
```bash
docker compose up -d
cp .env.example .env
go mod tidy
go run ./cmd/api
```

## Default admin
Uses the values from `.env`:
- email: `ADMIN_EMAIL`
- password: `ADMIN_PASSWORD`

## API overview
### Public
- `GET /api/v1/health`
- `GET /api/v1/profile`
- `GET /api/v1/projects`
- `GET /api/v1/tags`
- `GET /api/v1/posts`
- `GET /api/v1/posts/:slug`

### Admin
- `POST /api/v1/auth/login`
- `POST /api/v1/admin/posts`
- `PUT /api/v1/admin/posts/:id`
- `DELETE /api/v1/admin/posts/:id`
- `POST /api/v1/admin/tags`
- `PUT /api/v1/admin/profile`

## Example requests
### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"changeme123"}'
```

### Create post
```bash
curl -X POST http://localhost:8080/api/v1/admin/posts \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "title":"CI/CD Basics for Backend Developers",
    "summary":"A practical intro",
    "content":"# Hello\nThis is my first post.",
    "status":"published",
    "tags":["Go","CI/CD"]
  }'
```

## Next recommended steps
1. Add request validation.
2. Replace AutoMigrate with proper SQL migrations.
3. Add pagination and filtering to posts.
4. Add refresh tokens or session management.
5. Add Swagger or OpenAPI.
6. Add deploy workflow after CI.
