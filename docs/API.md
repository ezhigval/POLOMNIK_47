# API Specification MVP

## 1. Base URL

```text
/api/v1
```

## 2. Response principles

Successful single-resource response:

```json
{
  "data": {}
}
```

Successful list response:

```json
{
  "data": [],
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "has_next": true
  }
}
```

Error response:

```json
{
  "error": {
    "code": "STRING_CODE",
    "message": "Human readable"
  }
}
```

## 3. Health

### GET /health

Returns service health outside API versioning.

### GET /api/v1/health

Returns API health inside versioned API.

### GET /health/ready and GET /api/v1/health/ready

Readiness probe: PostgreSQL. Redis is optional: if it is down the site stays up and ready still returns 200 (`checks.cache` may be `degraded`). Use for uptime monitors and load balancer checks.

Responses include `X-Request-ID`. Error bodies may include `request_id`. Rate-limited responses send `Retry-After`.

Public write forms (`/bookings`, auth, support) accept a hidden honeypot field `website`. If filled, API returns a generic 422. Captcha token is verified only when `CAPTCHA_ADAPTER=smartcaptcha` and keys exist; availability is in `GET /auth/methods` → `captcha`.

Outbox worker has no HTTP health endpoint. Docker healthcheck uses heartbeat file
`WORKER_HEARTBEAT_PATH` (default `/tmp/palomnik-worker-heartbeat`), refreshed after each successful poll.

### GET /api/v1/management/system-info

Admin-only ops snapshot: adapter modes, outbox counts (`pending` / `failed` / `processed`), request latency, last backup marker (`last_backup`). Not a public pprof.

## 4. Public Tours API

### GET /api/v1/tours

Query params:

```text
page
limit
date_from
date_to
price_min
price_max
location
is_hot
```

Rules:

- default `page = 1`;
- default `limit = 20`;
- max `limit = 100`;
- dates use `YYYY-MM-DD`;
- public endpoint returns only `is_active = true`;
- sorting defaults to nearest `date_start`.

Response item:

```json
{
  "id": "uuid",
  "slug": "string",
  "title": "string",
  "description": "string",
  "price": 15000,
  "currency": "RUB",
  "date_start": "2026-05-10",
  "date_end": "2026-05-15",
  "slots_total": 20,
  "slots_left": 5,
  "location": "string",
  "images": ["url"],
  "is_hot": false
}
```

### GET /api/v1/tours/{id}

Returns one active tour.

If tour is inactive or missing, return `404`.

### GET /api/v1/tours/popular

Returns active tours where `is_hot = true`.

Query params:

```text
limit
```

Default `limit = 10`.

### GET /api/v1/tours/{id}/reviews

Returns approved reviews for a specific tour.

Query params:

```text
page
limit
```

## 5. Public Reviews API

### GET /api/v1/reviews

Query params:

```text
tour_id
rating
page
limit
```

Rules:

- public endpoint returns only `is_approved = true`;
- if `tour_id` is present, return reviews for one tour;
- default `page = 1`;
- default `limit = 20`;
- max `limit = 100`.

Response item:

```json
{
  "id": "uuid",
  "tour_id": "uuid",
  "client_name": "string",
  "rating": 5,
  "text": "string",
  "created_at": "2026-05-20T10:00:00Z"
}
```

## 6. Public Bookings API

### POST /api/v1/bookings

Creates a booking request in local backend storage.

Bitrix24 is not called in MVP 1. Future sync uses integration ports.

Request:

```json
{
  "tour_id": "uuid",
  "name": "Иван Иванов",
  "phone": "+79999999999",
  "email": "mail@test.com",
  "people_count": 2,
  "comment": "string"
}
```

Rules:

- `tour_id` is required;
- `name` is required;
- `phone` is required;
- `email` is optional but must be valid if present;
- `people_count` must be greater than 0;
- `comment` is optional;
- client does not send `use_overbooking`;
- client does not send `payment_method` in MVP 1;
- payment response fields are not returned in MVP 1.

Success response:

```json
{
  "data": {
    "status": "ok",
    "booking_id": "uuid",
    "booking_status": "NEW",
    "total_price": 30000,
    "integration_status": "not_configured"
  }
}
```

Possible errors:

- `TOUR_NOT_FOUND`;
- `TOUR_INACTIVE`;
- `INSUFFICIENT_SLOTS`;
- `VALIDATION_ERROR`;
- `INVALID_PEOPLE_COUNT`.

## 7. Management API

Management API is part of MVP backend logic so data can be created and processed before Bitrix24 exists.

Until full auth is implemented, these endpoints must be either:

- available only locally;
- protected by reverse proxy;
- protected by `X-Admin-Token`.

This is not a replacement for future user auth.

### Tours management

```text
GET    /api/v1/management/tours
POST   /api/v1/management/tours
GET    /api/v1/management/tours/{id}
PATCH  /api/v1/management/tours/{id}
DELETE /api/v1/management/tours/{id}
POST   /api/v1/management/uploads
```

### Image upload

`POST /api/v1/management/uploads` — multipart form field `file` (JPEG/PNG/WebP/GIF).

Public files are served at `GET /uploads/{path}`.

Response:

```json
{
  "data": {
    "url": "http://localhost:8080/uploads/2026/07/uuid.jpg",
    "path": "2026/07/uuid.jpg"
  }
}
```

Returned `url` is stored in tour `images[]`.

Create/update request fields:

```json
{
  "slug": "string",
  "title": "string",
  "description": "string",
  "price": 15000,
  "currency": "RUB",
  "date_start": "2026-05-10",
  "date_end": "2026-05-15",
  "slots_total": 20,
  "slots_left": 20,
  "location": "string",
  "images": ["url"],
  "is_active": true,
  "is_hot": false,
  "overbooking_enabled": false
}
```

### Bookings management

```text
GET   /api/v1/management/bookings
GET   /api/v1/management/bookings/{id}
PATCH /api/v1/management/bookings/{id}/status
```

`GET /bookings` query (optional): `status` (`NEW` / `CONTACTED` / `CONFIRMED` / `COMPLETED` / `CANCELLED`), `date_from`, `date_to` (`ГГГГ-ММ-ДД`, `date_to` inclusive), `tour_id`, `page`, `limit`. `format=csv` — UTF-8 CSV тех же полей заявки, без выдуманных колонок; выгрузка всех совпадений фильтра, не только текущей страницы.

Status update request:

```json
{
  "status": "CONTACTED"
}
```

Allowed MVP statuses:

```text
NEW
CONTACTED
CONFIRMED
COMPLETED
CANCELLED
```

### Roles management

```text
GET    /api/v1/management/roles
GET    /api/v1/management/roles/templates
POST   /api/v1/management/roles
PATCH  /api/v1/management/roles/{id}
DELETE /api/v1/management/roles/{id}
GET    /api/v1/management/roles/{id}/assignments
POST   /api/v1/management/roles/{id}/assignments
DELETE /api/v1/management/roles/{id}/assignments/{userId}
```

`GET /roles/templates` — именованные пресеты формы создания роли (этап 2). Не создают строку в БД. `manage_roles` в набор не входит.

| id | Имя в UI | Права |
|----|----------|--------|
| `booking_manager` | Менеджер заявок | `manage_bookings`, `manage_support` |
| `advertiser` | Рекламщик | `view_stats` |
| `smm` | Сммщик | `manage_content` |
| `director` | Директор | все, кроме `manage_roles` |
| `developer` | Разработчик | `manage_integrations`, `view_stats` |

### Reviews management

```text
GET    /api/v1/management/reviews
POST   /api/v1/management/reviews
PATCH  /api/v1/management/reviews/{id}/approve
PATCH  /api/v1/management/reviews/{id}/reject
DELETE /api/v1/management/reviews/{id}
```

Create review request:

```json
{
  "tour_id": "uuid",
  "client_name": "string",
  "rating": 5,
  "text": "string",
  "is_approved": false
}
```

### Integration references (internal)

Read-only endpoint for operators to inspect sync state before/after Bitrix24 and 1C go live.

```text
GET /api/v1/management/integration-references
```

Query params:

```text
page
limit
external_system
local_entity_type
sync_status
```

Response item:

```json
{
  "id": "uuid",
  "local_entity_type": "booking",
  "local_entity_id": "uuid",
  "external_system": "bitrix24",
  "external_entity_type": "deal",
  "external_entity_id": "local:uuid-or-external-id",
  "sync_status": "pending",
  "last_sync_at": "2026-06-10T12:00:00Z",
  "last_error": "",
  "created_at": "2026-06-10T12:00:00Z",
  "updated_at": "2026-06-10T12:00:00Z"
}
```

Records appear when `CRM_ADAPTER=bitrix` or `ACCOUNTING_ADAPTER=onec` is enabled and adapter methods are invoked.

### Outbox events (internal)

Retry queue for integration events. Populated when sync returns `pending` or `failed`. Worker processing is Stage C.

```text
GET /api/v1/management/outbox-events
```

Query params:

```text
page
limit
status
entity_type
event_type
```

Response item:

```json
{
  "id": "uuid",
  "event_type": "crm.push_booking",
  "entity_type": "booking",
  "entity_id": "uuid",
  "payload": { "integration_status": "pending" },
  "status": "pending",
  "attempts": 0,
  "last_error": "",
  "created_at": "2026-06-10T12:00:00Z",
  "updated_at": "2026-06-10T12:00:00Z"
}
```

## 8. Webhooks

### Bitrix24 deal update (inbound)

```text
POST /api/v1/webhooks/bitrix/deal?token={BITRIX_INBOUND_TOKEN}
```

Bitrix24 outbound webhook sends deal stage changes. Backend resolves `ORIGIN_ID` → booking UUID and updates status using standard stage mapping.

## 8a. Cabinet identities

User session: `Authorization: Bearer`. Internal OAuth completion (Next.js → API) also requires `X-Internal-Secret`.

```text
GET   /api/v1/me
PATCH /api/v1/me
GET   /api/v1/me/identities
GET   /api/v1/me/passengers
POST  /api/v1/me/passengers
PATCH /api/v1/me/passengers/:id
DELETE /api/v1/me/passengers/:id
GET   /api/v1/me/bookings
POST  /api/v1/auth/oauth
```

`PATCH /me` (session required): `name`, `email`, `phone`, optional `phone_check_id`, honeypot `website`. Имя обязательно. Почта и телефон могут быть пустыми (oauth-only). Смена телефона при включённом callcheck требует тот же звонок, что регистрация. Уникальность почты и телефона как при регистрации. Пароль здесь не меняется (есть восстановление).

`GET /me/identities` returns `{ "data": [ { "provider", "subject", "created_at" } ] }` for the current user. Providers in product: `yandex`, `vk`, `max`, `telegram`.

Passengers (session required): body `name`, `phone`, `birth_date` (`YYYY-MM-DD`), `passport`, honeypot `website`. All four fields required. Phone normalized like the profile. No SNILS. Date of birth cannot be in the future. Passport is stored as given text — no invented series/number format. Cabinet returns full fields to the owner. Texts for messengers use `Passenger.MessengerLine()`: name readable, phone/passport last two characters, date of birth `**.**.****`.

`POST /auth/oauth` (internal) body: `provider`, `subject`, `email`, `name`, `phone`. If the request also has a valid user `Authorization: Bearer`:

- unused identity is linked to the current user;
- identity of the same user is a no-op;
- identity of another user merges that account **into the current** (bookings, favorites, support threads, admin role assignments, passengers). Profile fields that already differ are not overwritten.

Response extra fields (omitted when false/empty): `linked`, `merged`, `kept_fields` (conflict field names only: `name` / `email` / `phone` — no values). Invalid Bearer is ignored and the call is treated as login.

Invalid internal secret → 401, as before.

## 9. Future APIs

Not in MVP 1 (частично уже есть в коде; этот список устарел — см. ROADMAP):

```text
POST /auth/register
POST /auth/login
GET  /auth/methods
POST /auth/phone/start      # sms.ru callcheck
GET  /auth/phone/status
POST /auth/phone/complete
POST /auth/refresh
GET  /users/me
GET  /users/bookings
POST /payments
POST /webhooks/bitrix
POST /webhooks/onec
```

Подтверждение телефона — **звонок** (callcheck), не SMS. Без `PHONE_ADAPTER=smsru` + `SMSRU_API_ID` методы телефона отвечают «недоступно».

These endpoints require separate specification before implementation.

