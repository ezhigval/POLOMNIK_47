# Тихвинский путь (POLOMNIK_47)

Сайт и API паломнической службы Тихвинской епархии: каталог поездок, заявки, кабинет менеджера. Backend — источник правды; Bitrix24 и 1С подключаются адаптерами (сейчас `noop`).

Репозиторий: [github.com/ezhigval/POLOMNIK_47](https://github.com/ezhigval/POLOMNIK_47)

Публичный сайт: **https://tikhvin-palomnik.ru**  
API: **https://api.tikhvin-palomnik.ru**

Старый адрес `tikhvin-polomnik.ru` (через «о») отдаёт 301 на palomnik. Отдельного preview-стенда нет.

## Правила (для людей и агентов)

Полный канон: **[AGENTS.md](AGENTS.md)**. Кратко:

1. Правки сначала **локально**, затем **GitHub**, затем деплой на palomnik — только если владелец попросил.
2. Не выдумывать бизнес-логику. Неясное поведение — вопрос владельцу.
3. Backend — гексагональная архитектура. Логика не в HTTP-handlers и не во frontend.
4. Секреты (`.env.production`, ключи, токены) не коммитить и не светить в чат.
5. Commit / push / force-push — только по просьбе владельца. В `main` не force-push.
6. Юридические тексты, контент туров и Telegram — отдельный этап, не трогать без просьбы.

## Быстрый старт (локально)

```bash
docker compose up --build -d
```

| URL | Что |
|-----|-----|
| http://localhost:3000 | Сайт |
| http://localhost:8080 | API |
| http://localhost:8080/api/v1/tours | Публичный список туров |
| http://localhost:3000/management | Админка (токен `dev-admin-token`) |

Остановка: `docker compose down`.

Без Docker:

```bash
# backend
cd backend && ADMIN_TOKEN=dev-admin-token go run ./cmd/api

# frontend
cd frontend && cp .env.example .env.local && npm install && npm run dev
```

## Структура

```text
.
  AGENTS.md                 правила для всех агентов
  docs/                     спецификация, API, деплой, релиз
  backend/                  Go API (module polomnik)
  frontend/                 Next.js
  deploy/                   Caddy + скрипты Yandex Cloud
  .cursor/rules/            правила Cursor
```

## Проверки

```bash
cd backend && go test ./... && go vet ./...
cd frontend && npm run lint && npm run build
make smoke-docker            # API smoke через Docker
make e2e-docker              # Playwright
make check-ops               # ready + worker + outbox
```

## Релиз и деплой

```text
локально → GitHub (если попросили commit/push) → make deploy
```

```bash
make deploy
```

ВМ: Ubuntu 24.04, пользователь `smailikin70`, IP `93.77.165.81`, каталог `/opt/polomnik`. SSH: `ssh polomnik-yc`. Compose: `docker-compose.yml` + `docker-compose.prod.yml`, env `.env.production` (не в git).

DNS на REG.RU для зоны **tikhvin-palomnik.ru**: `@` / `www` / `api` → A `93.77.165.81`. Старый `tikhvin-polomnik.ru` можно оставить с теми же A-записями ради 301. Подробности: [docs/DEPLOY.md](docs/DEPLOY.md), [docs/RELEASE.md](docs/RELEASE.md).

## Документация

- [docs/RELEASE.md](docs/RELEASE.md) — конвейер релиза
- [docs/DEPLOY.md](docs/DEPLOY.md) — сервер, env, HTTPS, бэкапы
- [docs/PROJECT_SPEC.md](docs/PROJECT_SPEC.md)
- [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)
- [docs/API.md](docs/API.md)
- [docs/ROADMAP.md](docs/ROADMAP.md)
- [docs/BITRIX_SETUP.md](docs/BITRIX_SETUP.md)
- [docs/ONEC_INTEGRATOR_TZ.md](docs/ONEC_INTEGRATOR_TZ.md)

## Статус

v1.0: сайт на https://tikhvin-palomnik.ru. Telegram (заявки и поддержка) — через настройки админки. Bitrix24 и 1С — адаптеры в коде, live выключен. Очередь: [docs/ROADMAP.md](docs/ROADMAP.md).
