# Roadmap

Дорожная карта разделена на три глобальных этапа:

1. Backend Logic.
2. Frontend Logic.
3. Integrations: Bitrix24 and 1C.

Главное архитектурное правило: backend строится по гексагональной архитектуре. Business logic должна работать через порты, чтобы PostgreSQL, Bitrix24 и 1C были заменяемыми adapters.

## Stage A - Backend Logic

Цель:

Создать самостоятельный backend MVP без активных Bitrix24 и 1C, но с архитектурными заготовками под будущие интеграции.

### A0. Repository restructuring

Todo:

- [x] Создать структуру `docs/`, `backend/`, `frontend/`.
- [x] Перенести актуальную документацию в `docs/`.
- [x] Удалить устаревшие root-документы или заменить их ссылкой на `docs/`.
- [x] Перенести/заменить старые `go.mod` и `main.go` в `backend/`.
- [x] Создать `backend/go.mod` с module `polomnik`.
- [x] Убедиться, что в корне нет старого рабочего Go-кода, который сбивает агента.

Acceptance criteria:

- Документация находится в `docs/`.
- Backend-код находится в `backend/`.
- Frontend-код будет находиться в `frontend/`.
- Старый GoLand demo-код удален или заменен.

### A1. Backend skeleton

Todo:

- [x] Создать `backend/cmd/api/main.go`.
- [x] Подключить Fiber.
- [x] Добавить config loader.
- [x] Добавить structured logger.
- [x] Добавить request id middleware.
- [x] Добавить recover middleware.
- [x] Добавить graceful shutdown.
- [x] Добавить `GET /health`.
- [x] Добавить `GET /api/v1/health`.
- [x] Настроить `go test ./...`.
- [x] Настроить `go vet ./...`.

Acceptance criteria:

- Backend запускается локально.
- Health endpoints отвечают.
- `go test ./...` проходит.
- `go vet ./...` проходит.

### A2. Hexagonal structure

Todo:

- [x] Создать `internal/domain`.
- [x] Создать `internal/application`.
- [x] Создать `internal/ports`.
- [x] Создать `internal/adapters`.
- [x] Создать HTTP adapter `internal/adapters/http/fiber`.
- [x] Создать repository adapters:
  - [x] `postgres`;
  - [x] `memory`.
- [x] Создать cache adapters:
  - [x] `redis`;
  - [x] `noop`.
- [x] Создать integration adapters:
  - [x] `bitrix/noop`;
  - [x] `onec/noop`.
- [x] Запретить imports из domain в adapters.

Acceptance criteria:

- Domain не зависит от внешних технологий.
- Application работает через ports.
- Fiber используется только в HTTP adapter.

### A3. Domain model

Todo:

- [x] Реализовать `Tour`.
- [x] Реализовать `Booking`.
- [x] Реализовать `Review`.
- [x] Реализовать booking statuses:
  - [x] `NEW`;
  - [x] `CONTACTED`;
  - [x] `CONFIRMED`;
  - [x] `COMPLETED`;
  - [x] `CANCELLED`.
- [x] Реализовать validation для tour dates.
- [x] Реализовать validation для slots.
- [x] Реализовать validation для rating.
- [x] Реализовать calculation `total_price = tour.price * people_count`.
- [x] Реализовать status transition rules.
- [x] Реализовать domain errors.

Acceptance criteria:

- Domain tests работают без БД и HTTP.
- Invalid domain state невозможно создать без ошибки.

### A4. Ports

Todo:

- [x] Описать `TourRepository`.
- [x] Описать `BookingRepository`.
- [x] Описать `ReviewRepository`.
- [x] Описать `CachePort`.
- [x] Описать `CRMPort` для future Bitrix24.
- [x] Описать `AccountingPort` для future 1C.
- [x] Описать transaction boundary port, если понадобится.
- [x] Реализовать fake ports для tests.

Acceptance criteria:

- Application layer не знает про PostgreSQL, Redis, Bitrix24, 1C.
- Future ports есть, но реальные integrations не вызываются.

### A5. PostgreSQL adapter

Todo:

- [x] Выбрать migration tool.
- [x] Создать migrations.
- [x] Создать таблицу `tours`.
- [x] Создать таблицу `bookings`.
- [x] Создать таблицу `reviews`.
- [x] Создать таблицу `integration_references`.
- [x] Создать таблицу `outbox_events` (миграция + enqueue при pending/failed).
- [x] Реализовать PostgreSQL repositories.
- [x] Добавить repository tests.
- [x] Добавить seed/dev data mechanism.

Acceptance criteria:

- Миграции применяются на пустую БД.
- Repository tests проходят.
- Backend может работать с PostgreSQL adapter.

### A6. Application services

Todo:

- [x] Реализовать `TourService`.
- [x] Реализовать `BookingService`.
- [x] Реализовать `ReviewService`.
- [x] Реализовать `ManagementService` или отдельные use cases.
- [x] Реализовать list tours with filters.
- [x] Реализовать get tour.
- [x] Реализовать popular tours.
- [x] Реализовать create booking.
- [x] Реализовать slot reservation/release.
- [x] Реализовать booking status changes.
- [x] Реализовать reviews approval flow.
- [x] Подключить noop CRM and accounting ports.

Acceptance criteria:

- Все use cases покрыты tests с fake repositories.
- Booking logic не зависит от HTTP и PostgreSQL.

### A7. Public HTTP API

Todo:

- [x] Реализовать `GET /api/v1/tours`.
- [x] Реализовать `GET /api/v1/tours/{id}`.
- [x] Реализовать `GET /api/v1/tours/popular`.
- [x] Реализовать `GET /api/v1/tours/{id}/reviews`.
- [x] Реализовать `GET /api/v1/reviews`.
- [x] Реализовать `POST /api/v1/bookings`.
- [x] Реализовать DTO mapping.
- [x] Реализовать error format.
- [x] Реализовать request validation.
- [x] Добавить handler tests.

Acceptance criteria:

- Public MVP API соответствует `docs/API.md`.
- Невалидные requests возвращают понятные ошибки.

### A8. Management HTTP API

Todo:

- [x] Реализовать management routes group.
- [x] Добавить `X-Admin-Token` middleware или ограничить локальным доступом.
- [x] Реализовать tours CRUD.
- [x] Реализовать bookings list/detail.
- [x] Реализовать booking status update.
- [x] Реализовать reviews list/create/approve/reject/delete.
- [x] Добавить handler tests.

Acceptance criteria:

- Данные MVP можно создавать без Bitrix24.
- Management API не доступен публично без защиты.

### A9. Cache

Todo:

- [x] Реализовать cache port.
- [x] Реализовать Redis adapter.
- [x] Реализовать noop cache adapter.
- [x] Подключить cache-aside для tours list.
- [x] Подключить cache-aside для tour detail.
- [x] Подключить cache-aside для popular tours.
- [x] Добавить cache invalidation после management mutations.
- [x] Добавить tests на cache hit/miss.

Acceptance criteria:

- Redis можно отключить без поломки backend.
- Cache не меняет бизнес-логику.

### A10. Integration stubs

Todo:

- [x] Реализовать noop `CRMPort`.
- [x] Реализовать noop `AccountingPort`.
- [x] Добавить config flags:
  - [x] `CRM_ADAPTER=noop`;
  - [x] `ACCOUNTING_ADAPTER=noop`.
- [x] При создании booking вызывать порт только если adapter включен.
- [x] В MVP noop adapter возвращает `not_configured`.
- [x] Добавить tests, что отсутствие Bitrix24/1C не ломает booking.

Acceptance criteria:

- Backend явно готов к будущим integrations.
- MVP не требует Bitrix24/1C credentials.

### A11. Docker and local infrastructure

Todo:

- [x] Создать backend Dockerfile.
- [x] Создать docker-compose.
- [x] Добавить PostgreSQL.
- [x] Добавить Redis, если cache включен.
- [x] Добавить `.env.example`.
- [x] Добавить migrations command.
- [x] Добавить seed command.
- [x] Обновить README.

Acceptance criteria:

- `docker-compose up` поднимает backend dependencies.
- Новый разработчик может запустить backend по инструкции.

### A12. Backend DoD

Backend stage считается завершенным, если:

- backend находится в `backend/`;
- API соответствует `docs/API.md`;
- domain/application покрыты tests;
- `go test ./...` проходит;
- `go vet ./...` проходит;
- локальный запуск описан;
- Bitrix24 и 1C не нужны для MVP;
- порты под Bitrix24 и 1C существуют.

## Stage B - Frontend Logic

Цель:

Создать Next.js frontend, который использует только backend API.

### B0. Frontend skeleton

Todo:

- [x] Создать `frontend/`.
- [x] Создать Next.js app.
- [x] Подключить TypeScript.
- [x] Подключить TailwindCSS.
- [x] Настроить lint/test/build scripts.
- [x] Настроить API base URL через env.
- [x] Создать API client.

Acceptance criteria:

- Frontend запускается локально.
- Frontend build проходит.

### B1. Public pages

Todo:

- [x] Реализовать список туров.
- [x] Реализовать фильтры:
  - [x] даты;
  - [x] цена;
  - [x] location;
  - [x] hot/popular.
- [x] Реализовать карточку тура.
- [x] Реализовать страницу тура.
- [x] Реализовать блок отзывов.
- [x] Реализовать форму заявки.
- [x] Реализовать success state после заявки.
- [x] Реализовать error states.
- [x] Реализовать empty states.
- [x] Реализовать loading states.

Acceptance criteria:

- Пользователь может найти тур и отправить заявку.
- Frontend не обращается к Bitrix24/1C.
- Frontend не содержит business rules, кроме UX validation.

### B2. Management UI, optional for MVP

Todo:

- [x] Решить, нужен ли management UI в первом frontend.
- [x] Если нужен, реализовать внутренние страницы:
  - [x] tours management;
  - [x] bookings management;
  - [x] reviews moderation.
- [x] Защитить доступ через согласованный MVP механизм.

Acceptance criteria:

- Management UI не доступен публично без защиты.
- Если UI не нужен, management остается только API.

### B3. UX and quality

Todo:

- [x] Responsive layout.
- [x] Accessible forms.
- [x] Clear validation messages.
- [x] SEO metadata for tours.
- [x] Sitemap, robots.txt, JSON-LD, OG image.
- [x] Global 404 page.
- [x] API loading indicators.
- [x] Smoke/e2e test for main flow (Playwright: guest booking + auth).

Acceptance criteria:

- Main user flow работает на desktop and mobile.
- Форма заявки надежно обрабатывает успех и ошибки.

### B5. Marketing UX polish (post-MVP UI)

Todo:

- [x] Hero с фото, trust stats, продающий контент.
- [x] Секции: почему мы, отзывы, FAQ, CTA.
- [x] Типографика (display serif), mobile nav, тёмный footer.
- [x] Карточки туров с cover fallback, длительность, hover.
- [x] Страница тура: программа, «что включено», breadcrumbs.
- [x] Контакты через env (`NEXT_PUBLIC_CONTACT_*`).
- [x] Политика конфиденциальности (`/privacy`).
- [x] Favicon / app icon.
- [x] Загрузка URL фото туров в management UI.
- [x] Секция «О службе» (#about), site-config через env.
- [ ] Реальные фото и тексты от владельца продукта.
- [x] OG-image для шаринга в мессенджерах.
- [ ] Яндекс.Метрика / GA.

Acceptance criteria:

- Сайт выглядит как готовый продающий лендинг на демо-контенте.
- Юридический текст согласован перед публичным запуском.

### B4. Frontend DoD

Frontend stage завершен, если:

- Next.js app находится в `frontend/`;
- публичный пользовательский сценарий работает;
- frontend вызывает только backend API;
- build/lint проходят;
- базовые smoke tests проходят.

## Stage C - Integrations: Bitrix24 and 1C

Цель:

Подключить внешние системы через adapters, не переписывая domain/application.

### C0. Integration discovery

Todo:

- [x] Уточнить роль Bitrix24 (дефолт в discovery, *подтвердить*).
- [x] Уточнить роль 1C (дефолт в discovery, *подтвердить*).
- [x] Выбрать source of truth strategy (backend).
- [x] Получить Bitrix24 entity model (deals/contacts, без туров).
- [x] Получить 1C exchange format (HTTP JSON contract).
- [x] Утвердить sync direction (push).
- [x] Утвердить conflict resolution (backend wins).
- [x] Утвердить retry policy (outbox worker).
- [x] Заполнить `docs/INTEGRATION_DISCOVERY.md`.

Acceptance criteria:

- [x] Есть integration spec (assumption-based).
- [ ] Нет догадок по Bitrix24/1C fields — **требует подтверждения владельца перед продом**.

### C1. Bitrix24 adapter

Todo:

- [x] Реализовать real Bitrix24 adapter for `CRMPort`.
- [x] Настроить auth/credentials (incoming webhook).
- [x] Реализовать mapping tours (crm.product catalog).
- [x] Реализовать mapping bookings/deals (+ product rows).
- [x] Реализовать mapping contacts.
- [x] Реализовать mapping reviews (timeline comment on approve).
- [x] Реализовать sync jobs (outbox worker).
- [x] Реализовать webhooks (inbound deal stage → booking status).
- [x] Реализовать integration reference storage.
- [x] Реализовать outbox worker skeleton (retry через порты).
- [x] Добавить contract tests with fixtures (httptest mock servers).

Acceptance criteria:

- Bitrix24 adapter можно включить config-ом.
- Domain/application не меняются при включении adapter.

### C2. 1C adapter

Todo:

- [x] Реализовать real 1C adapter for `AccountingPort` (HTTP JSON + OData).
- [x] Настроить auth/credentials (basic / API key).
- [x] Реализовать export bookings.
- [ ] Реализовать payment/accounting sync, если payments уже есть.
- [x] Реализовать counterparty sync.
- [x] Реализовать retry and idempotency (outbox + booking_id).
- [x] Добавить contract tests with fixtures.

Acceptance criteria:

- 1C adapter можно включить config-ом.
- Ошибки 1C не ломают публичный API.

### C3. Source of truth switch

Todo:

- [x] Решить, остается ли PostgreSQL source of truth. **Да (реализован push).**
- [x] Если backend остается source of truth, реализовать push strategy.
- [ ] Если hybrid, реализовать conflict rules (не требуется для текущей стратегии).
- [x] Обновить API behavior documentation (`INTEGRATION_DISCOVERY.md`).

Acceptance criteria:

- [x] Источник правды явно описан.
- [x] Sync failures observable (management integrations + outbox).
- [ ] Нет расхождений между docs и behavior — **финальная сверка после live Bitrix/1C.**

### C4. Integration DoD

Integration stage **в коде завершён (~95%)**. Prod DoD:

- [x] Локальный e2e с mock Bitrix + mock 1C (`make integration-smoke-docker`).
- [ ] Bitrix24 подключен на живом портале и проверен end-to-end.
- [ ] 1C подключена интегратором и проверен export.
- [x] source of truth определен (PostgreSQL);
- [x] sync failures логируются и retry-ятся (outbox worker);
- [x] tests покрывают mappings and failure scenarios;
- [x] domain/application не переписаны под внешние системы.

## Stage D - Production readiness

Цель:

Запуск в прод без Bitrix/1C (автономный режим) + операционка для менеджеров.

### D0. Уведомления (код)

Todo:

- [x] `NotificationPort` + noop adapter.
- [x] Telegram adapter (`NOTIFICATION_ADAPTER=telegram`).
- [x] Уведомления: новая заявка + смена статуса.
- [x] Recording + outbox retry (как CRM).
- [x] `GET /management/system-info`.
- [x] `docs/TELEGRAM_SETUP.md`.
- [ ] **Live подключение Telegram** — Stage E.

Acceptance criteria:

- [x] Код готов, по умолчанию `noop`.
- [ ] Бот подключён и проверен end-to-end — финальный этап.

### D1. Контент и юридическое

Todo:

- [x] Базовая политика ПДн на сайте.
- [x] Ссылка на политику в форме заявки.
- [x] Seed с расширенными описаниями туров (dev).
- [ ] Реальные контакты, фото, тексты.
- [ ] Финальная политика / оферта с юристом.

### D2. Management UX

Todo:

- [x] Поле URL фото в формах тура.
- [x] Загрузка файлов (локальное хранилище) вместо URL-only.
- [ ] Полноценная авторизация (login/роли) вместо одного токена — **частично:** `/management/login` + session cookie.

### D3. Deploy and ops

Todo:

- [x] Docker Compose для локального стека.
- [x] CI: backend tests + frontend build.
- [x] CI: smoke с поднятым API (docker compose).
- [x] `docs/DEPLOY.md` — prod-запуск с Caddy, `.env.production.example`, backup script.
- [ ] Прод-сервер, домен, SSL, бэкапы БД — **инфра на стороне владельца;** скрипты и compose готовы.
- [x] CI: frontend lint + smoke + e2e с поднятым стеком.
- [x] Security hardening: OAuth lockdown, rate limits, prod secret validation, CSP.
- [x] Мониторинг / алерты (API `/health/ready`, worker heartbeat, `system-info` outbox summary, `scripts/check-ops.sh`).

### D4. Analytics

Todo:

- [x] Яндекс.Метрика / GA — env-gated (`NEXT_PUBLIC_YM_ID`, `NEXT_PUBLIC_GA_ID`).
- [x] События: `tour_view`, `booking_submit`.
- [ ] Подключить счётчики на prod и проверить в кабинете.

### D5. Production DoD (без внешних сервисов)

Автономный прод считается готовым, если:

- [ ] реальный контент и контакты на сайте;
- [x] политика ПДн опубликована (базовая версия);
- [ ] деплой на домен с HTTPS;
- [ ] бэкапы PostgreSQL;
- [x] код уведомлений и интеграций готов (adapters noop).

**Внешние подключения — Stage E, не блокируют автономный запуск.**

## Stage E - External connections (финальный этап)

Цель:

Подключить и отладить все внешние сервисы **после** наполнения контентом и внутренних тестов.

### E0. Telegram

Todo:

- [x] Adapter, outbox, worker, docs.
- [ ] Создать бота, задать env, smoke на реальной заявке.
- [ ] Проверить уведомление при смене статуса.

### E1. Bitrix24

Todo:

- [x] Adapter code (см. Stage C1).
- [ ] Живой портал, webhook, подгонка стадий/полей.

### E2. 1C

Todo:

- [x] Adapter code (см. Stage C2).
- [ ] HTTP-сервис у интегратора, smoke export.

### E3. Analytics (опционально)

Todo:

- [ ] Яндекс.Метрика / GA.
- [ ] События конверсии.

### E4. External connections DoD

- [ ] Telegram доставляет заявки менеджеру.
- [ ] Bitrix24 sync проверен (если нужен).
- [ ] 1C export проверен (если нужен).
- [ ] Документация совпадает с prod-поведением.

## Текущий статус (сводка)

| Область | Готовность |
|---------|------------|
| Backend MVP (Stage A) | **100%** |
| Frontend + marketing UX (B, B5) | **~95%** |
| Integrations **code** (C) | **~95%** |
| Content + legal + deploy prep (D) | **~75%** |
| **Live external services (E)** | **0%** (намеренно отложено) |

**Сейчас:** пишем контент, тестируем автономно.  
**Потом:** подключаем Telegram → Bitrix → 1C → аналитику и отлаживаем.

## Global Definition of Done

Любая задача считается выполненной, если:

- она соответствует документации;
- unclear requirements вынесены в вопросы;
- код находится в правильной зоне repo;
- tests добавлены для измененной логики;
- checks проходят;
- документация обновлена при изменении поведения;
- external integrations не добавлены в MVP без решения;
- business logic не находится во frontend or adapters.
