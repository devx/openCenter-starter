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

## API examples

### Get cluster by id

```bash
curl -s http://localhost:8080/api/v1/clusters/cluster-id
```

Response (not found):

```json
{
  "error": {
    "code": "not_found",
    "message": "cluster not found"
  },
  "meta": {
    "request_id": "f3d0cb92-2a4c-4b3c-8c04-4a9a0d6b82b1"
  }
}
```
