# Bitrix24: быстрый старт (≈10 минут)

## 1. Создать портал

1. Открой [bitrix24.ru](https://www.bitrix24.ru) → «Создать бесплатно»
2. Зарегистрируй компанию / портал
3. Зайди в CRM

## 2. Входящий вебхук

1. **Приложения** → **Разработчикам** → **Другое** → **Входящий вебхук**
2. Права: **CRM** (контакты, сделки — чтение и запись)
3. Сохрани URL вида:

```text
https://YOUR.bitrix24.ru/rest/1/xxxxxxxxxxxxxxxx/
```

4. В `.env` проекта:

```env
CRM_ADAPTER=bitrix
BITRIX_WEBHOOK_URL=https://YOUR.bitrix24.ru/rest/1/xxxxxxxxxxxxxxxx/
```

5. Перезапусти API и worker:

```bash
docker compose up --build -d api worker
```

## 3. Воронка сделок (рекомендуется)

Создай стадии под статусы backend или используй стандартную воронку.

После создания стадий скопируй их `STAGE_ID` из настроек CRM и задай в env:

```env
BITRIX_STAGE_NEW=NEW
BITRIX_STAGE_CONTACTED=UC_XXXXX
BITRIX_STAGE_CONFIRMED=UC_XXXXX
BITRIX_STAGE_COMPLETED=WON
BITRIX_STAGE_CANCELLED=LOSE
```

Если env пустые — Bitrix использует стадию по умолчанию для новой сделки.

Опционально — отдельная воронка:

```env
BITRIX_DEAL_CATEGORY_ID=1
```

## 4. Проверка

### Локально (mock Bitrix)

```bash
make integration-smoke-docker
```

### На живом портале

1. Создай заявку на сайте
2. Открой `/management/integrations` — sync ref со статусом `synced`, external id = ID сделки
3. В Bitrix24 → CRM → Сделки — должна появиться «Заявка: …»
4. Смени статус заявки в management → стадия сделки обновится (дефолтные стадии: NEW, PREPARATION, EXECUTING, WON, LOSE)

## 5. Опциональные поля

Если нужны tour_id и people_count в карточке сделки:

1. CRM → Настройки → Поля сделки → добавь пользовательские поля
2. Узнай коды (`UF_CRM_...`)
3. В env:

```env
BITRIX_UF_TOUR_ID=UF_CRM_1234567890
BITRIX_UF_PEOPLE_COUNT=UF_CRM_0987654321
```

## 6. Обратная синхронизация (исходящий webhook Bitrix → backend)

1. Bitrix24 → Роботы / Business Process → исходящий webhook на:

```text
POST https://YOUR-API/api/v1/webhooks/bitrix/deal?token=YOUR_BITRIX_INBOUND_TOKEN
```

2. В `.env`: `BITRIX_INBOUND_TOKEN=...`
3. При смене стадии сделки backend обновит статус заявки по `ORIGIN_ID`

## Troubleshooting

| Симптом | Решение |
|---------|---------|
| `not_configured` | Пустой `BITRIX_WEBHOOK_URL` |
| Outbox `failed` | Смотри `last_error` в management; проверь права вебхука |
| Сделка не находится при update | Проверь `BITRIX_ORIGINATOR_ID=polomnik` (дефолт) |
| HTTP timeout | Увеличь `INTEGRATION_HTTP_TIMEOUT=30s` |
