# Deploy

Прод-запуск **без live Bitrix/1C** — backend остаётся source of truth, интеграции `noop`.

Конвейер: **[RELEASE.md](RELEASE.md)**, правила агентов: **[AGENTS.md](../AGENTS.md)**.

```bash
make deploy
```

Публичный сайт: `https://tikhvin-palomnik.ru`. API: `https://api.tikhvin-palomnik.ru`.

## 1. Сервер

Рекомендуемый прод: **Yandex Cloud Compute** (Ubuntu 24.04, 2 vCPU / ≥2 GB RAM + swap) + Docker Compose + Caddy.

Скрипты:

```bash
export PATH="$HOME/yandex-cloud/bin:$PATH"
yc init   # один раз, логин в браузере
./deploy/yandex/provision.sh          # VPC + VM + IP
./deploy/yandex/deploy.sh             # rsync + docker compose prod
# каталог по умолчанию /opt/palomnik; иначе DEPLOY_DIR=...
```

DNS на REG.RU (A-записи → публичный IP VM) для зоны **tikhvin-palomnik.ru**:

| Имя | Тип | Значение |
|-----|-----|----------|
| `@` | A | IP сервера |
| `www` | A | IP сервера |
| `api` | A | IP сервера |

- Docker + Docker Compose v2.
- Порты 22 / 80 / 443 открыты в security group.

## 2. Переменные окружения

Скопируйте шаблон и заполните:

```bash
cp .env.production.example .env.production
```

| Переменная | Назначение |
|------------|------------|
| `SITE_DOMAIN` | `tikhvin-palomnik.ru` |
| `API_DOMAIN` | `api.tikhvin-palomnik.ru` |
| `ACME_EMAIL` | Email для Let's Encrypt (Caddy) |
| `POSTGRES_PASSWORD` | Сильный пароль БД (не дефолт из local compose) |
| `ADMIN_TOKEN` | Пароль для `/management/login` и management API |
| `JWT_SECRET` | ≥32 символов, для пользовательских сессий |
| `INTERNAL_API_SECRET` | ≥16 символов, для Google OAuth → API |
| `CORS_ALLOW_ORIGINS` | `https://tikhvin-palomnik.ru` |
| `NEXT_PUBLIC_SITE_URL` | `https://tikhvin-palomnik.ru` |
| `NEXT_PUBLIC_API_URL` | `/api/v1` (same-origin через Next rewrite) |
| `UPLOAD_PUBLIC_BASE_URL` | `https://api.tikhvin-palomnik.ru` |
| `NEXT_PUBLIC_CONTACT_*` | Телефон и email на сайте |
| `NEXT_PUBLIC_YM_ID` | Яндекс.Метрика (прод: `111985266`; пусто = без скрипта) |
| `NEXT_PUBLIC_GA_ID` | Google Analytics (опционально) |
| `NEXT_PUBLIC_YM_WEBVISOR` | `1` — вебвизор Метрики (прод включён; нужен rebuild фронта) |
| `NEXT_PUBLIC_YM_CLICKMAP` | `0` — выключить карту кликов (иначе вкл. вместе с Метрикой) |
| `NEXT_PUBLIC_YANDEX_VERIFICATION` | Код Яндекс.Вебмастера (мета-тег, опц.) |
| `NEXT_PUBLIC_GOOGLE_SITE_VERIFICATION` | Код Google Search Console (опц.) |
| `NOTIFICATION_ADAPTER` | `telegram` при настройке бота |
| `TELEGRAM_BOT_TOKEN` / `TELEGRAM_CHAT_ID` | следующий этап |
| `MANAGEMENT_BASE_URL` | `https://tikhvin-palomnik.ru/management/bookings` |

Интеграции (позже):

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

- `https://tikhvin-palomnik.ru` — HTTPS 200
- `https://www.tikhvin-palomnik.ru` → 301 на apex
- `https://api.tikhvin-palomnik.ru/health/ready` → `status: ok` (нужна A-запись `api`)
- `/management/login` → admin token

На ВМ API с хоста не слушает `:8080`. Проверка изнутри:

```bash
API_VIA_DOCKER=1 ./scripts/check-ops.sh
```

## 4. Локальный dev

```bash
docker compose up --build -d
# http://localhost:3000
```

## 5. Бэкапы PostgreSQL

```bash
make backup-db
# или на ВМ:
./scripts/backup-postgres.sh
```

Дампы в `backups/`, хранятся 7 дней. Cron:

```cron
0 3 * * * cd /opt/palomnik && ./scripts/backup-postgres.sh
```

Восстановление (имя контейнера — из `docker compose ps`; пользователь и БД — из `.env.production`):

```bash
gunzip -c backups/palomnik-YYYYMMDD-HHMMSS.sql.gz | \
  docker exec -i "$(docker ps --format '{{.Names}}' | grep postgres-1 | head -1)" \
    psql -U "$POSTGRES_USER" -d "$POSTGRES_DB"
```

## 5.1 Мониторинг

- Снаружи: `https://tikhvin-palomnik.ru` и (если есть DNS) `https://api.tikhvin-palomnik.ru/health/ready`
- На ВМ: `API_VIA_DOCKER=1 ./scripts/check-ops.sh` — readiness + outbox.failed + worker healthy

Worker:

- Compose `healthcheck` по heartbeat-файлу
- Management → Integrations: блок «Outbox здоровье»

```cron
*/10 * * * * cd /opt/palomnik && API_VIA_DOCKER=1 ./scripts/check-ops.sh
```

`ADMIN_TOKEN` скрипт читает из `.env.production`, в crontab его писать не нужно.

## 6. CI / тесты перед деплоем

```bash
make smoke-docker   # API smoke
make e2e-docker     # Playwright E2E
make check-ops      # API ready + worker healthy + outbox failed=0
```

## 7. Telegram

Код готов. Подключать отдельно. См. `docs/TELEGRAM_SETUP.md`.

## 8. Чеклист перед открытием

- [ ] Реальные туры и фото в management
- [ ] Контакты в env
- [ ] Политика ПДн согласована с юристом
- [ ] `ADMIN_TOKEN`, `JWT_SECRET`, `INTERNAL_API_SECRET` — не dev-значения
- [ ] `POSTGRES_PASSWORD` сменён
- [ ] HTTPS работает (Caddy + DNS)
- [ ] Бэкап cron настроен
- [ ] Telegram (если нужен) проверен тестовой заявкой
