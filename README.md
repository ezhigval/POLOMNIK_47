# POLOMNIK_47

Платформа паломнической службы: backend API, frontend-сайт и будущие интеграции с Bitrix24 и 1C.

## Быстрый старт (Docker)

```bash
docker compose up --build -d
```

Сервисы:

| URL | Описание |
|-----|----------|
| http://localhost:3000 | Frontend |
| http://localhost:8080 | Backend API |
| http://localhost:8080/api/v1/tours | Публичный API |
| http://localhost:3000/management | Management UI (дашборд) |

Outbox worker (`worker` service) автоматически поднимается в compose и ретраит pending-события интеграций.

Management UI доступен только при `ADMIN_TOKEN` в env frontend-контейнера (по умолчанию `dev-admin-token`).

Остановка:

```bash
docker compose down
```

## Структура

```text
POLOMNIK_47/
  docs/       — документация
  backend/    — Go API (hexagonal architecture)
  frontend/   — Next.js
```

## Локальная разработка без Docker

Backend:

```bash
cd backend
ADMIN_TOKEN=dev-admin-token go run ./cmd/api
```

Outbox worker (отдельный процесс):

```bash
cd backend
DATABASE_URL=postgres://... go run ./cmd/worker
```

Frontend:

```bash
cd frontend
cp .env.example .env.local
npm install
npm run dev
```

## Документация

- [docs/PROJECT_SPEC.md](docs/PROJECT_SPEC.md)
- [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)
- [docs/API.md](docs/API.md)
- [docs/ROADMAP.md](docs/ROADMAP.md)
- [docs/INTEGRATION_DISCOVERY.md](docs/INTEGRATION_DISCOVERY.md) — решения по интеграциям
- [docs/BITRIX_SETUP.md](docs/BITRIX_SETUP.md) — Bitrix24 за 10 минут
- [docs/ONEC_INTEGRATOR_TZ.md](docs/ONEC_INTEGRATOR_TZ.md) — ТЗ для 1С-специалиста

## Текущий статус

- **Stage A (Backend MVP)** — завершён
- **Stage B (Frontend MVP + marketing UX)** — ~95%
- **Stage C (Bitrix24 / 1C code)** — адаптеры готовы; live подключение — Stage E
- **Stage D (Production readiness)** — ~65%: SEO, about, CI smoke; нужны реальный контент и деплой
- **Stage E (Live integrations)** — отложено: Telegram → Bitrix → 1C

## Проверки

```bash
cd backend && go test ./... && go vet ./...
cd frontend && npm run build && npm run lint && npm run test:smoke
make smoke-docker   # полный smoke через Docker
make e2e-docker     # Playwright E2E (guest booking + auth)
make check-ops      # API ready + worker heartbeat + outbox summary
make integration-smoke-docker  # Bitrix + 1C mocks end-to-end
make docker-prod    # prod stack (нужен .env.production)
make backup-db      # pg_dump → backups/
```
