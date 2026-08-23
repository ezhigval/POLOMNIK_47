# Backend

Go backend for PALOMNIK 47.

## Run locally

With in-memory storage:

```bash
go run ./cmd/api
```

With PostgreSQL and Redis:

```bash
cp ../.env.example ../.env
export $(grep -v '^#' ../.env | xargs)
go run ./cmd/api
```

Default address:

```text
:8080
```

## Docker Compose

From repository root:

```bash
docker compose up --build
```

API:

```text
http://localhost:8080
```

Management API header:

```text
X-Admin-Token: dev-admin-token
```

## Checks

```bash
go test ./...
go vet ./...
```

## Migrations

Migrations use goose-compatible SQL files in `migrations/`.

Install goose if needed:

```bash
go install github.com/pressly/goose/v3/cmd/goose@v3.27.1
```

Run migrations:

```bash
DATABASE_URL=postgres://palomnik:palomnik@localhost:5432/palomnik?sslmode=disable make migrate-up
```

Check status:

```bash
DATABASE_URL=postgres://palomnik:palomnik@localhost:5432/palomnik?sslmode=disable make migrate-status
```

Seed local dev data:

```bash
DATABASE_URL=postgres://palomnik:palomnik@localhost:5432/palomnik?sslmode=disable make seed-dev
```

## PostgreSQL Integration Tests

Repository integration tests are skipped unless `TEST_DATABASE_URL` is set:

```bash
TEST_DATABASE_URL=postgres://palomnik:palomnik@localhost:5432/palomnik_test?sslmode=disable go test ./internal/adapters/repository/postgres
```
