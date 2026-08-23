# Backend Architecture

Слои не менялись. Telegram на проде — адаптер уведомлений + webhook, не ядро домена. Bitrix24 и 1С в проде `noop`. Очередь: [ROADMAP.md](ROADMAP.md).

## 1. Обязательный стиль

Backend строится по гексагональной архитектуре.

Цель:

- отделить бизнес-логику от Fiber;
- отделить бизнес-логику от PostgreSQL;
- отделить бизнес-логику от Redis;
- отделить бизнес-логику от Bitrix24;
- отделить бизнес-логику от 1C;
- упростить замену adapters без переписывания domain и use cases.

## 2. Слои

```text
Domain
  - entities
  - value objects
  - domain errors
  - domain rules

Application
  - use cases
  - services
  - transactions
  - orchestration

Ports
  - repository ports
  - cache ports
  - integration ports
  - clock/id/logger abstractions where useful

Adapters
  - HTTP Fiber adapter
  - PostgreSQL adapter
  - Redis adapter
  - Bitrix24 adapter, future
  - 1C adapter, future
  - in-memory/fake adapters for tests
```

Dependency rule:

```text
adapters -> application -> domain
```

Domain must not import application, adapters, Fiber, SQL, Redis, Bitrix24 SDK, 1C packages or environment config.

## 3. Target backend structure

```text
backend/
  go.mod
  go.sum
    cmd/
    api/
      main.go
    worker/
      main.go
    backup-offsite/
      main.go
  internal/
    domain/
      tour.go
      booking.go
      review.go
      errors.go
    application/
      tour_service.go
      booking_service.go
      review_service.go
      management_service.go
    ports/
      repositories.go
      cache.go
      crm.go
      accounting.go
      transactions.go
    adapters/
      http/
        fiber/
          router.go
          handlers/
          middleware/
          dto/
      repository/
        postgres/
        memory/
      cache/
        redis/
        noop/
      integration/
        bitrix/
        onec/
        noop/
    config/
    logger/
    validation/
  migrations/
  tests/
```

`pkg/` не создается, пока нет реальной публичной библиотеки для внешнего использования.

## 4. Ports

### TourRepository

Отвечает за локальное хранение туров.

Минимальные методы:

```text
ListTours(filters, pagination)
GetTour(id)
CreateTour(tour)
UpdateTour(tour)
DeleteTour(id)
ReserveSlots(tourID, peopleCount)
ReleaseSlots(tourID, peopleCount)
```

### BookingRepository

Минимальные методы:

```text
CreateBooking(booking)
GetBooking(id)
ListBookings(filters, pagination)
UpdateBookingStatus(id, status)
MarkBookingOverbooked(id)
```

### ReviewRepository

Минимальные методы:

```text
ListReviews(filters, pagination)
CreateReview(review)
ApproveReview(id)
RejectReview(id)
DeleteReview(id)
```

### CachePort

Минимальные методы:

```text
Get(key)
Set(key, value, ttl)
Delete(key)
```

Cache must be optional. If Redis is unavailable, public API should continue through repository/service paths.

### MessengerPort (v3, этап 4)

```text
Configured()
Send(channel, address, text)
```

Адаптеры: telegram / max / WhatsApp Cloud API. Без ключа `Configured()==false`.

### PublisherPort (v3, этап 4)

```text
Configured()
Publish(channel, content)
```

Адаптеры: site_news / telegram_channel / vk_wall / max_feed.

### CaptchaPort (v3, этап 1)

```text
Configured()
Provider()
ClientKey()
Verify(token, ip)
```

Адаптер `smartcaptcha`. Без ключа — honeypot + Redis/memory rate limit, сайт жив.

### AIPort (v3, этап 4/7)

```text
Configured()
Complete(system, user)
```

Адаптер `yandexgpt` появится на этапе 4. v4 (звонки, ИИ-продавец) сюда не входит.

### PaymentPort (v3, этап 8)

```text
Configured()
CreatePayment(booking)  # сумма = booking.TotalPrice
```

Адаптеры `sber` / `yookassa`, выбор `PAYMENT_ADAPTER`. Возвраты не делаем. Без адаптера форма заявки как в v2 (`NEW`).

### BackupStoragePort (v3, этап 1)

```text
Configured()
Upload(name, body)
```

Адаптер S3-compatible (Yandex Object Storage). Ночной дамп на диск остаётся; offsite — когда есть ключ бакета.

### CRMPort

Future Bitrix24 integration port.

MVP 1 adapter: noop/fake.

Possible future methods:

```text
SyncTour(tour)
PullTours()
PushBooking(booking)
UpdateBookingStatus(booking)
PullReviews()
```

Methods must be finalized only when Bitrix24 model is known.

### AccountingPort

Future 1C integration port.

MVP 1 adapter: noop/fake.

Possible future methods:

```text
ExportBooking(booking)
ExportPayment(payment)
SyncCounterparties()
```

Methods must be finalized only when 1C exchange format is known.

## 5. Active adapters by stage

### Stage A - Backend Logic

```text
HTTP: Fiber
Repository: PostgreSQL
Cache: Redis or noop
CRM: noop/fake
Accounting: noop/fake
```

### Stage B - Frontend Logic

Backend adapters stay the same. Frontend uses only HTTP API.

### Stage C - Integrations

```text
CRM: Bitrix24 adapter
Accounting: 1C adapter
Repository: PostgreSQL remains local storage, cache, fallback or operational store
```

The exact source of truth can change in Stage C, but only through port-level decisions.

## 6. HTTP layer

Fiber is only an adapter.

HTTP handlers must:

- parse request;
- validate transport-level data;
- call application services;
- map application result to response DTO;
- map errors to API error format.

HTTP handlers must not:

- contain booking business rules;
- directly call PostgreSQL;
- directly call Redis;
- directly call Bitrix24;
- directly call 1C.

## 7. Application layer

Application services own workflows:

- list tours;
- create tour;
- update tour;
- create booking;
- change booking status;
- list reviews;
- approve review;
- future sync operations.

Application services work through ports.

## 8. Domain layer

Domain owns:

- entity invariants;
- status transitions;
- slot validation;
- price calculation basics;
- approval rules for reviews;
- domain errors.

Domain must be testable without database, HTTP server, Redis or external services.

## 9. Error strategy

Domain/application errors should be typed:

- validation error;
- not found;
- conflict;
- insufficient slots;
- invalid status transition;
- dependency unavailable;
- internal error.

HTTP adapter maps them to:

- `400 Bad Request`;
- `404 Not Found`;
- `409 Conflict`;
- `422 Unprocessable Entity`;
- `503 Service Unavailable`;
- `500 Internal Server Error`.

## 10. Configuration

Environment variables for backend MVP:

```text
APP_ENV=local
HTTP_ADDR=:8080
DATABASE_URL=postgres://...
REDIS_URL=redis://...
ADMIN_TOKEN=
LOG_LEVEL=info
CRM_ADAPTER=noop
ACCOUNTING_ADAPTER=noop
```

Future:

```text
BITRIX_BASE_URL=
BITRIX_WEBHOOK_URL=
BITRIX_CLIENT_ID=
BITRIX_CLIENT_SECRET=
ONEC_BASE_URL=
ONEC_USERNAME=
ONEC_PASSWORD=
```

Future secrets must not be required for MVP 1.

## 11. Testing architecture

Required tests:

- domain tests without adapters;
- application tests with fake ports;
- HTTP handler tests with fake application/services;
- repository tests for PostgreSQL;
- adapter tests for Redis if enabled;
- future contract tests for Bitrix24 and 1C.

Minimum backend checks:

```bash
go test ./...
go vet ./...
```

