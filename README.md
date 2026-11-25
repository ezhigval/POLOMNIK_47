# Тихвинский путь

Сайт и API паломнической службы Тихвинской епархии: каталог поездок, заявки, кабинет менеджера. Backend — источник правды; Bitrix24 и 1С подключаются адаптерами (сейчас `noop`).

Репозиторий: [github.com/ezhigval/POLOMNIK_47](https://github.com/ezhigval/POLOMNIK_47)

Публичный сайт: **https://tikhvin-palomnik.ru**  
API: **https://api.tikhvin-palomnik.ru**

Отдельного preview-стенда нет.

## Правила (для людей и агентов)

Полный канон: **[AGENTS.md](AGENTS.md)**. Кратко:

1. Правки в **отдельной ветке** и PR в `main`. Деплой: GitHub Actions после успешного CI, резервный путь — `SSH_KEY=<путь> make deploy`. Актуальные настройки: [docs/RELEASE.md](docs/RELEASE.md).
2. Не выдумывать бизнес-логику. Неясное поведение — вопрос владельцу.
3. Backend — гексагональная архитектура. Логика не в HTTP-handlers и не во frontend.
4. Секреты (`.env.production`, ключи, токены) не коммитить и не светить в чат.
5. Commit / push / force-push — только по просьбе владельца. В `main` не force-push.
6. Юридические тексты и контент туров не выдумывать. Секреты OAuth/SMTP/Метрики — владелец по [docs/V2_OWNER_SETUP.md](docs/V2_OWNER_SETUP.md).

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
  backend/                  Go API (module palomnik)
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
ветка → PR → CI → merge в main → CI(main) → GitHub Actions Deploy production → проверка /health/ready
```

Production размещён на **Selectel**, Ubuntu 24.04, Docker Compose + Caddy. Фактический target в `.github/workflows/deploy.yml`: `deploy@92.53.100.23`, рабочий каталог `/opt/polomnik`, файл конфигурации `.env.production.test`. SSH-ключ хранится только в GitHub Actions secrets как `DEPLOY_SSH_KEY`; значения секретов в репозитории не публикуются. Cloudflare-компоненты Telegram не следует отключать при изменениях инфраструктуры.

Ручной резервный вариант: `SSH_KEY=<путь-к-приватному-ключу> make deploy`. Нельзя запускать два деплоя одновременно или выполнять `docker compose down -v`. После публикации сверять `/release.txt` с `main`.

См. [docs/RELEASE.md](docs/RELEASE.md), [docs/DEPLOY.md](docs/DEPLOY.md), [AGENTS.md](AGENTS.md).

## Документация

- [docs/STATUS.md](docs/STATUS.md) — что на проде сейчас
- [docs/V4_PLAN.md](docs/V4_PLAN.md) — линейка v4
- [docs/V4_OWNER_SETUP.md](docs/V4_OWNER_SETUP.md) — чеклист секретов и кабинетов
- [docs/RELEASE.md](docs/RELEASE.md) — конвейер релиза
- [docs/DEPLOY.md](docs/DEPLOY.md) — сервер, env, HTTPS, бэкапы
- [docs/ROADMAP.md](docs/ROADMAP.md)
- [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)
- [docs/API.md](docs/API.md)
- [docs/OAUTH_SETUP.md](docs/OAUTH_SETUP.md)
- [docs/SEO_ADS.md](docs/SEO_ADS.md)
- [docs/TELEGRAM_SETUP.md](docs/TELEGRAM_SETUP.md)

Остальное: [docs/README.md](docs/README.md).

## Статус

По изменениям в GitHub на **2026-11-25**: PR #54–64 охватывают режимы главной страницы, разделение пространств CMS, формы в модальных окнах и предпросмотр неопубликованных правок. PR #65–68 (merged) добавили email-получателей, SMTP через Selectel и обработку постоянной SMTP-ошибки, чтобы не отправлять повторно Telegram-уведомления по тому же событию. Проверка фактического состояния production и независимых backup/restore — отдельные операционные задачи, не подменяемые фактом merge.

Текущий перечень подтверждённых изменений и непроверенных задач: [docs/STATUS.md](docs/STATUS.md). Старое состояние от 2026-08-27 не является актуальным.
