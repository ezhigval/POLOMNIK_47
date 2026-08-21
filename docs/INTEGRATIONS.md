# Integrations: Bitrix24, 1C, Telegram

## 1. Current decision

**Все внешние сервисы по умолчанию выключены (noop).** Код адаптеров готов; для локальной проверки — mock-серверы в Docker.

```bash
make integration-smoke-docker
# или
docker compose -f docker-compose.yml -f docker-compose.integrations.yml up -d --build
cd frontend && npm run test:integration
```

Mock-сервисы:

| Сервис | URL | Назначение |
|--------|-----|------------|
| `mock-bitrix` | `http://localhost:8091` | Bitrix24 REST (сделки, контакты, товары) |
| `mock-onec` | `http://localhost:8092/accounting` | HTTP JSON обмен 1C |

Живое подключение — см. `BITRIX_SETUP.md`, `ONEC_INTEGRATOR_TZ.md`.

MVP wires integration through ports:

```text
CRM_ADAPTER=noop           -> noop adapter (default)
CRM_ADAPTER=bitrix         -> Bitrix24 REST (incoming webhook)

ACCOUNTING_ADAPTER=noop    -> noop adapter (default)
ACCOUNTING_ADAPTER=onec    -> 1C HTTP JSON / OData

NOTIFICATION_ADAPTER=noop      -> noop (default)
NOTIFICATION_ADAPTER=telegram  -> Telegram Bot API (manager alerts)
```

Application services always call ports. Integration and notification failures do not break public booking/tour flows.

Setup guides (подключать на финальном этапе):

- [BITRIX_SETUP.md](BITRIX_SETUP.md)
- [ONEC_INTEGRATOR_TZ.md](ONEC_INTEGRATOR_TZ.md)
- [TELEGRAM_SETUP.md](TELEGRAM_SETUP.md)
- [INTEGRATION_DISCOVERY.md](INTEGRATION_DISCOVERY.md)

When `CRM_ADAPTER=bitrix` or `ACCOUNTING_ADAPTER=onec` is selected, sync attempts are persisted in `integration_references` through the recording wrapper around adapter ports.

Events with `pending` or `failed` status are also enqueued to `outbox_events` for future retry worker (Stage C). Management API: `GET /api/v1/management/outbox-events`.

`cmd/worker` polls pending outbox events and retries integration calls through the same `CRMPort` / `AccountingPort` stack. Configure via:

```text
OUTBOX_WORKER_POLL_INTERVAL=30s
OUTBOX_WORKER_BATCH_SIZE=10
OUTBOX_WORKER_MAX_ATTEMPTS=5
```

Docker Compose includes a `worker` service. Pending events are deduplicated per `(event_type, entity_id)`.

## 2. Integration stage goals

Stage C must answer and implement:

- what system becomes source of truth for tours;
- how bookings move between local backend and Bitrix24;
- how reviews are imported/exported;
- how 1C receives accounting/order/payment data;
- how conflicts are resolved;
- how retries and failures are monitored.

## 3. Bitrix24 future scope

Potential responsibilities:

- tours storage or sync;
- contacts;
- deals/bookings;
- reviews;
- manager workflow;
- CRM status tracking.

Potential adapter methods:

```text
PullTours()
PushTour(tour)
PullReviews()
PushBooking(booking)
UpdateBookingStatus(booking)
FindOrCreateContact(contact)
```

Do not implement these against real Bitrix24 until the following are known:

- entity type for tours: Smart Process or products;
- Bitrix24 `entityTypeId`;
- exact user field codes;
- deal pipeline id;
- `STAGE_ID` mapping;
- review entity structure;
- image storage format;
- webhook payload format;
- auth method and credentials.

## 4. 1C future scope

Potential responsibilities:

- order export;
- payment/accounting sync;
- counterparty sync;
- service/tour nomenclature sync;
- reconciliation reports.

Potential adapter methods:

```text
ExportBooking(booking)
ExportPayment(payment)
SyncCounterparty(contact)
SyncNomenclature(tour)
GetAccountingStatus(entity)
```

Do not implement these against real 1C until the following are known:

- 1C version/configuration;
- exchange protocol: HTTP, file exchange, OData, SOAP, custom API;
- auth method;
- document types;
- required fields;
- status mapping;
- retry strategy;
- idempotency keys;
- who is source of truth for accounting statuses.

## 5. Integration references

External links should not pollute domain entities with provider-specific fields.

Use `integration_references` or adapter-specific mapping tables:

```text
local booking id -> Bitrix24 deal id
local booking id -> 1C order/document id
local tour id -> Bitrix24 product/smart process id
```

## 6. Source of truth options for Stage C

The source of truth must be decided before real integration:

### Option 1 - Backend source of truth

Backend owns tours/bookings/reviews and pushes to Bitrix24/1C.

Pros:

- stable product API;
- easier frontend behavior;
- less dependency on CRM availability.

Cons:

- managers must work in backend/admin UI or sync must be very reliable.

### Option 2 - Bitrix24 source of truth

Bitrix24 owns tours/bookings/reviews, backend caches and normalizes.

Pros:

- managers work in CRM;
- matches original idea.

Cons:

- CRM model and availability affect product behavior;
- more complex caching and sync.

### Option 3 - Hybrid

Backend owns public product behavior, Bitrix24 owns manager workflow.

Pros:

- flexible;
- gradual migration.

Cons:

- requires clear conflict resolution.

No option is selected yet. Agent must ask before implementing real integrations.

