# Integration Discovery (Stage C0)

Документ фиксирует решения по интеграциям Bitrix24 и 1C.

**Статус:** заполнено дефолтной схемой (assumption-based). Строки с пометкой *подтвердить* нужно сверить с владельцем проекта перед продом.

## 1. Роли систем

| Вопрос | Ответ |
|--------|-------|
| Какую роль играет Bitrix24? | CRM для менеджеров: контакты, сделки по заявкам, смена стадий. *подтвердить* |
| Какую роль играет 1C? | Учёт: экспорт подтверждённых заявок, контрагенты. Оплаты — позже. *подтвердить* |
| Кто source of truth для туров? | **PostgreSQL (backend)** |
| Кто source of truth для заявок/bookings? | **PostgreSQL (backend)** |
| Кто source of truth для отзывов? | **PostgreSQL (backend)** |

## 2. Bitrix24

| Вопрос | Ответ |
|--------|-------|
| Тип сущности для туров | **crm.product** (каталог Bitrix), sync при create/update тура |
| `entityTypeId` для туров | — |
| Pipeline / category для сделок | `BITRIX_DEAL_CATEGORY_ID` (0 = default). *подтвердить* |
| Маппинг `STAGE_ID` ↔ статусы backend | Через env: `BITRIX_STAGE_NEW`, `BITRIX_STAGE_CONTACTED`, … *подтвердить после настройки воронки* |
| Где хранятся отзывы в Bitrix24 | **Timeline comment** на crm.product при approve |
| Формат изображений туров | — (туры не в Bitrix) |
| Auth | **Incoming webhook** (`BITRIX_WEBHOOK_URL`) |
| Webhook payload для обратной синхронизации | `POST /api/v1/webhooks/bitrix/deal?token=...` — обновление статуса booking из STAGE_ID сделки |

### Техническая связка booking ↔ deal

- `ORIGINATOR_ID=polomnik` (env `BITRIX_ORIGINATOR_ID`)
- `ORIGIN_ID=<booking UUID>` — поиск и идемпотентность без кастомных UF-полей
- Опционально: `BITRIX_UF_TOUR_ID`, `BITRIX_UF_PEOPLE_COUNT` после создания полей в Bitrix

Инструкция: [BITRIX_SETUP.md](BITRIX_SETUP.md)

## 3. 1C

| Вопрос | Ответ |
|--------|-------|
| Версия и конфигурация 1C | **На стороне заказчика** (Бухгалтерия / УТ / ERP). *подтвердить* |
| Протокол обмена | **HTTP JSON** (`ONEC_PROTOCOL=http`) или **OData** (`ONEC_PROTOCOL=odata`, стандартный `/odata/standard.odata`) |
| Тип документа для экспорта booking | Определяет 1C-специалист (см. [ONEC_INTEGRATOR_TZ.md](ONEC_INTEGRATOR_TZ.md)) |
| Обязательные поля контрагента | `name`, `phone`, `email` в JSON |
| Нужна ли синхронизация оплат в MVP интеграций | **Нет** |
| Auth | Basic (`ONEC_USERNAME` / `ONEC_PASSWORD`) или `ONEC_API_KEY` header |
| Idempotency key | `booking_id` в теле запроса |

Endpoints по умолчанию:

```text
POST {ONEC_BASE_URL}/hs/polomnik/booking
POST {ONEC_BASE_URL}/hs/polomnik/counterparty
```

## 4. Sync policy

| Вопрос | Ответ |
|--------|-------|
| Направление sync для туров | Backend only (Bitrix skip) |
| Направление sync для bookings | **Push** backend → Bitrix24 / 1C |
| Conflict resolution | Backend wins; внешние системы обновляются по `ORIGIN_ID` / `booking_id` |
| Retry policy | Outbox worker: `OUTBOX_WORKER_MAX_ATTEMPTS=5`, poll 30s |
| Кто мониторит failed sync | Management UI `/management/integrations` + логи worker |

## 5. Текущая техническая реализация

```text
CRMPort          → noop | bitrix (REST webhook)
AccountingPort   → noop | onec (HTTP JSON)
integration_references + outbox_events
recording wrapper
cmd/worker — retry pending events
GET /api/v1/management/integration-references
GET /api/v1/management/outbox-events
```

После подтверждения таблиц — обновить env в проде и прогнать smoke с реальным Bitrix webhook.
