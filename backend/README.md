# Backend Service

This is the Go backend service for openCenter-base. It uses Fiber and follows a hexagonal (ports & adapters) layout.

## Local development

```bash
cd backend

go run ./cmd/server
```

## Database migrations

Install the `migrate` CLI (https://github.com/golang-migrate/migrate) and set `DATABASE_URL`.

```bash
export DATABASE_URL="postgres://user:pass@localhost:5432/opencenter?sslmode=disable"
migrate -path ./migrations -database "$DATABASE_URL" up
```

## API specification

The OpenAPI 3.1 spec for this service lives at `backend/openapi.yaml`.
