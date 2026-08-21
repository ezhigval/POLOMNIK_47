# Deploy

Прод-запуск **без Bitrix/1C** — backend остаётся source of truth, интеграции `noop`.

## 1. Сервер

- VPS (2 vCPU, 4 GB RAM) или managed container.
- Docker + Docker Compose v2.
- DNS A-записи: `SITE_DOMAIN` и `API_DOMAIN` → IP сервера.

## 2. Переменные окружения

Скопируйте шаблон и заполните:

```bash
cp .env.production.example .env.production
```

| Переменная | Назначение |
|------------|------------|
| `SITE_DOMAIN` | Публичный домен сайта (`polomnik47.ru`) |
| `API_DOMAIN` | Домен API (`api.polomnik47.ru`) |
| `ACME_EMAIL` | Email для Let's Encrypt (Caddy) |
| `POSTGRES_PASSWORD` | Сильный пароль БД (не `polomnik`) |
| `ADMIN_TOKEN` | Пароль для `/management/login` и management API |
| `JWT_SECRET` | ≥32 символов, для пользовательских сессий |
| `INTERNAL_API_SECRET` | ≥16 символов, для Google OAuth → API |
| `CORS_ALLOW_ORIGINS` | `https://ваш-домен.ru` |
| `NEXT_PUBLIC_SITE_URL` | `https://ваш-домен.ru` |
| `NEXT_PUBLIC_API_URL` | `https://api.ваш-домен.ru/api/v1` |
| `UPLOAD_PUBLIC_BASE_URL` | Публичный origin API (`https://api.ваш-домен.ru`) для URL загруженных фото |
| `NEXT_PUBLIC_CONTACT_*` | Телефон и email на сайте |
| `NEXT_PUBLIC_YM_ID` | Яндекс.Метрика (опционально) |
| `NEXT_PUBLIC_GA_ID` | Google Analytics (опционально) |
| `NOTIFICATION_ADAPTER` | `telegram` при настройке бота |
| `TELEGRAM_BOT_TOKEN` / `TELEGRAM_CHAT_ID` | Stage E |
| `MANAGEMENT_BASE_URL` | `https://домен.ru/management/bookings` |

Интеграции (Stage E):

- `CRM_ADAPTER=bitrix`, `BITRIX_WEBHOOK_URL=...`, `BITRIX_INBOUND_TOKEN=...`
- `ACCOUNTING_ADAPTER=onec`, `ONEC_BASE_URL=...`

При `APP_ENV=production` API **не стартует** с dev-секретами — проверка в `config.Validate`.

## 3. Запуск (production)

```bash
docker compose --env-file .env.production \
  -f docker-compose.yml \
  -f docker-compose.prod.yml \
  up -d --build
```

Что делает prod-overlay:

- Caddy на `:80`/`:443` с автоматическим HTTPS
- Postgres/Redis/API/Frontend **без** публичных портов
- `APP_ENV=production`, без dev-seed
- Worker с restart policy

Проверка:

- `https://api.домен/health/ready` → `status: ok`
- Сайт открывается, туры грузятся
- `/management/login` → admin token
- Тестовая заявка → management

## 4. Локальный dev (как раньше)

```bash
docker compose up --build -d
# http://localhost:3000
```

## 5. Бэкапы PostgreSQL

```bash
make backup-db
# или
./scripts/backup-postgres.sh
```

Дампы в `backups/`, хранятся 7 дней. Для cron:

```cron
0 3 * * * cd /opt/polomnik && ./scripts/backup-postgres.sh
```

Восстановление:

```bash
gunzip -c backups/polomnik-YYYYMMDD-HHMMSS.sql.gz | \
  docker exec -i polomnik_47-postgres-1 psql -U polomnik -d polomnik
```

## 5.1 Мониторинг

Uptime (внешний мониторинг → алерт при падении):

- `https://api.домен/health/ready`
- `make check-ops` локально / на сервере: readiness + `system-info` outbox.failed + worker healthy

Worker:

- Compose `healthcheck` по heartbeat-файлу
- Management → Integrations: блок «Outbox здоровье»

Cron ops check (пример каждый 5 мин):

```cron
*/5 * * * * cd /opt/polomnik && CHECK_WORKER=1 ADMIN_TOKEN=... ./scripts/check-ops.sh
```

## 6. CI / тесты перед деплоем

```bash
make smoke-docker   # API smoke
make e2e-docker     # Playwright E2E
make check-ops      # API ready + worker healthy + outbox failed=0
```

## 7. Telegram (Stage E)

Код готов. Подключать после контента и smoke. См. `docs/TELEGRAM_SETUP.md`.

## 8. Чеклист перед открытием

- [ ] Реальные туры и фото в management
- [ ] Контакты в env
- [ ] Политика ПДн согласована с юристом
- [ ] `ADMIN_TOKEN`, `JWT_SECRET`, `INTERNAL_API_SECRET` — не dev-значения
- [ ] `POSTGRES_PASSWORD` сменён
- [ ] HTTPS работает (Caddy + DNS)
- [ ] Бэкап cron настроен
- [ ] Telegram (если нужен) проверен тестовой заявкой
