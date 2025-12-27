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

### List clusters (with pagination)

```bash
curl -s "http://localhost:8080/api/v1/clusters?limit=2&offset=0"
```

Response (success):

```json
{
  "data": [
    {
      "id": "5c7fbda1-2a80-4df0-8a35-5d4aa1f4b6f6",
      "name": "prod-cluster",
      "status": "ready",
      "created_at": "2025-01-09T12:15:11Z",
      "updated_at": "2025-01-10T09:02:54Z"
    }
  ],
  "meta": {
    "request_id": "5b2143f1-7ce1-4f2c-9d5c-6b4d29a4f2a1",
    "pagination": {
      "total": 1,
      "limit": 2,
      "offset": 0
    }
  }
}
```

### List clusters (filters)

```bash
curl -s "http://localhost:8080/api/v1/clusters?status=ready&name_prefix=prod&id_prefix=5c7f"
```

Query parameters:
- `status`: exact match on cluster status.
- `name_prefix`: prefix match on cluster name.
- `id_prefix`: prefix match on cluster id.

Response (success):

```json
{
  "data": {
    "id": "5c7fbda1-2a80-4df0-8a35-5d4aa1f4b6f6",
    "name": "prod-cluster",
    "status": "ready",
    "created_at": "2025-01-09T12:15:11Z",
    "updated_at": "2025-01-10T09:02:54Z"
  },
  "meta": {
    "request_id": "d9b9a046-1f31-45cc-8299-6f1dd903c702"
  }
}
```
