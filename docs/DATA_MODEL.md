# Data Model

## 1. Principles

- Domain model is independent from PostgreSQL schema.
- PostgreSQL is the active storage adapter in MVP 1.
- Bitrix24 and 1C identifiers are stored as external references only after integrations are added.
- User accounts are not required in MVP 1.
- Booking stores contact data directly.

## 2. Tour

Fields:

```text
id UUID
slug TEXT
title TEXT
description TEXT
price INT
currency TEXT
date_start DATE
date_end DATE
slots_total INT
slots_left INT
location TEXT
images TEXT[]
is_active BOOLEAN
is_hot BOOLEAN
overbooking_enabled BOOLEAN
created_at TIMESTAMPTZ
updated_at TIMESTAMPTZ
```

Rules:

- `price >= 0`;
- `slots_total >= 0`;
- `slots_left >= 0`;
- `slots_left <= slots_total`, unless overbooking creates a separate booking flag instead of negative slots;
- `date_start <= date_end`;
- public API returns only active tours.

## 3. Booking

Fields:

```text
id UUID
tour_id UUID
name TEXT
phone TEXT
email TEXT
people_count INT
status TEXT
total_price INT
comment TEXT
overbooked BOOLEAN
source TEXT
created_at TIMESTAMPTZ
updated_at TIMESTAMPTZ
```

MVP statuses:

```text
NEW
CONTACTED
CONFIRMED
COMPLETED
CANCELLED
```

Future statuses:

```text
AWAITING_PAYMENT
PAID
IN_TRIP
SYNC_PENDING
SYNC_FAILED
```

Rules:

- booking belongs to one tour;
- `people_count > 0`;
- `total_price = tour.price * people_count` in MVP 1;
- status transitions must be validated in domain/application layer;
- cancellation should release reserved slots if slots were already reserved.

## 4. Review

Fields:

```text
id UUID
tour_id UUID
client_name TEXT
rating INT
text TEXT
is_approved BOOLEAN
created_at TIMESTAMPTZ
updated_at TIMESTAMPTZ
```

Rules:

- `rating` is from 1 to 5;
- public API returns only approved reviews;
- management API can create, approve, reject and delete reviews.

## 5. IntegrationReference

Used in future integration stage.

Fields:

```text
id UUID
local_entity_type TEXT
local_entity_id UUID
external_system TEXT
external_entity_type TEXT
external_entity_id TEXT
sync_status TEXT
last_sync_at TIMESTAMPTZ
last_error TEXT
created_at TIMESTAMPTZ
updated_at TIMESTAMPTZ
```

Examples:

```text
external_system = bitrix24
external_system = onec
```

This table can be added in MVP 1 as a forward-compatible placeholder, but no active Bitrix24/1C calls are made in MVP 1.

## 6. Initial PostgreSQL schema direction

Tables:

```text
tours
bookings
reviews
integration_references
```

Optional tables:

```text
tour_images
outbox_events
```

Recommended indexes:

```text
tours(slug)
tours(is_active)
tours(is_hot)
tours(date_start, date_end)
bookings(tour_id)
bookings(status)
bookings(created_at)
reviews(tour_id)
reviews(is_approved)
integration_references(local_entity_type, local_entity_id)
integration_references(external_system, external_entity_id)
```

## 7. Outbox

Outbox is optional in MVP 1 and useful for future integrations.

If implemented, it stores future sync tasks without requiring Bitrix24/1C to be active:

```text
id UUID
event_type TEXT
entity_type TEXT
entity_id UUID
payload JSONB
status TEXT
attempts INT
last_error TEXT
created_at TIMESTAMPTZ
updated_at TIMESTAMPTZ
```

In MVP 1, outbox events can remain pending until `cmd/worker` processes them or adapters return `synced` / `not_configured`.

## 8. User and identity (v3 cabinet)

OAuth больше не хранится колонками на `users`. Одна пара `(provider, subject)` — одна строка `user_identities`; у пользователя может быть несколько входов.

```text
users
  id UUID
  email TEXT          # unique where email <> ''
  phone TEXT          # unique where not empty
  name TEXT
  password_hash TEXT  # nullable for oauth-only
  created_at TIMESTAMPTZ
  updated_at TIMESTAMPTZ

user_identities
  user_id UUID        # FK users ON DELETE CASCADE
  provider TEXT       # yandex / vk / max / telegram
  subject TEXT
  created_at TIMESTAMPTZ
  PRIMARY KEY (provider, subject)
```

Привязка соцсети при открытой сессии сливает другой кабинет **в текущего**: заявки (`bookings.user_id`), избранное, треды поддержки (открытый тред один: сообщения переносятся), назначения админ-ролей. Имя/почта/телефон целевого профиля не перезаписываются, если оба заполнены и различаются. Пассажиров в схеме ещё нет. СНИЛС нет.

