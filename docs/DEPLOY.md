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
# каталог по умолчанию /opt/palomnik; на этой ВМ живой путь /opt/polomnik (скрипт подхватывает)
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

## 6. Российский VPS

Скрипт `deploy/vps/deploy.sh` переносит код на уже подготовленный VPS и запускает тот же Compose production overlay. Он не отправляет `.env.production`, uploads или `backups/`, не останавливает стек перед пересозданием и не удаляет Docker volumes. Deploy-пользователь должен уже иметь доступ к Docker; root SSH не требуется.

Файл `deploy/vps/compose.prod.yml` сохраняет PostgreSQL в отдельном volume `polomnik_postgres_vps_data`; повторные выкладки используют тот же volume. Файлы `.env.production.live` и `.env.mail` тоже исключены из синхронизации.

Почтовый контейнер принимает только письма на `info@tikhvin-palomnik.ru`, получает актуальный список адресатов из `site_settings.mail_forward_to` и отправляет сообщения через SMTP-реле Selectel. На VPS подготовьте `.env.mail` с `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `SMTP_RELAY_USERNAME` и `SMTP_RELAY_PASSWORD` из защищённых настроек сервера/Selectel; права файла — `600`. Порт TCP 25 требуется для приёма входящих писем. Не публикуйте порт базы данных.

Перед первым запуском проверьте fingerprint VPS у провайдера и внесите его в `~/.ssh/known_hosts`. Затем укажите проверенный SSH-target и ключ:

```bash
VPS_SSH_HOST='deploy@vps.example.ru' \
VPS_DEPLOY_DIR='/opt/polomnik' \
SSH_KEY="$HOME/.ssh/polomnik_vps" \
make deploy-vps
```

На VPS заранее создайте каталог, поместите туда `.env.production` и проверьте `docker compose version`. Скрипт завершается ошибкой, если каталога или файла секретов нет.

GitHub Actions содержит отдельный ручной workflow `Deploy to Russian VPS`. Он запускается только вручную из `main` и требует secrets `VPS_SSH_HOST`, `VPS_SSH_KEY`, `VPS_SSH_KNOWN_HOSTS` (проверенная строка из `known_hosts`) и `VPS_DEPLOY_DIR`. Production-секреты приложения остаются на VPS. Этот workflow не переключает DNS и не меняет действующий workflow Yandex Cloud.

После переноса сначала проверьте новый сервер по его IP/локальной маршрутизации и сохранённым сертификатам. DNS переключайте только после функциональной проверки нового сервера и согласования короткого окна записи; старую ВМ сохраняйте для отката.

## 7. CI / тесты перед деплоем

```bash
make smoke-docker   # API smoke
make e2e-docker     # Playwright E2E
make check-ops      # API ready + worker healthy + outbox failed=0
```

## 8. Telegram

Код готов. Подключать отдельно. См. `docs/TELEGRAM_SETUP.md`.

## 9. Чеклист перед открытием

- [ ] Реальные туры и фото в management
- [ ] Контакты в env
- [ ] Политика ПДн согласована с юристом
- [ ] `ADMIN_TOKEN`, `JWT_SECRET`, `INTERNAL_API_SECRET` — не dev-значения
- [ ] `POSTGRES_PASSWORD` сменён
- [ ] HTTPS работает (Caddy + DNS)
- [ ] Бэкап cron настроен
- [ ] Telegram (если нужен) проверен тестовой заявкой

## SMTP-уведомления на VPS Selectel

В `deploy/vps/compose.prod.yml` API и worker получают закрытый серверный `.env.mail` (не хранится в git) и используют уже действующий Selectel SMTP relay: `smtp.mail.selcloud.ru:1126`, отправитель `info@tikhvin-palomnik.ru`. Конфиг Go берёт `SMTP_USERNAME`/`SMTP_PASSWORD` либо `SMTP_RELAY_USERNAME`/`SMTP_RELAY_PASSWORD`. Канал Email считается настроенным только при наличии учётных данных. Изменения вступают в силу после пересоздания контейнеров через CI Deploy. При отсутствии `.env.mail` канал остаётся выключенным. Подтвердить фактическую доставку нужно тестовой заявкой. Не коммитить SMTP-секреты и не выводить их в логи.

## Ошибка SMTP 550 bounced check filter

По документации Selectel этот код означает, что адрес **получателя** занесён в локальный или глобальный стоп-лист. В кабинете Selectel откройте Почтовый сервис → ресурс → Стоп-листы, выполните поиск именно адреса, указанного в «Получателях уведомлений». Удалить из локального стоп-листа можно при ошибочном включении, из глобального самостоятельно нельзя; обратиться в поддержку Selectel. Не обходить стоп-лист новым отправителем и не отправлять повторно автоматически на заблокированные адреса.

Мультиканальный диспетчер не ставит в очередь полное уведомление при частичном успехе: иначе уже доставленный Telegram/MAX многократно повторится при SMTP-сбое. Ошибку частичной доставки сохраняют серверные логи. Окончательные SMTP 5xx не повторяются; временные ошибки при полном отказе могут попасть в outbox.
