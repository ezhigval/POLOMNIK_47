# ТЗ для 1С-специалиста: HTTP-обмен с PALOMNIK

Backend отправляет JSON на публикуемые HTTP-сервисы 1C. Формат фиксирован; реализация на стороне 1C — на усмотрение специалиста (Бухгалтерия 3.0, УТ, ERP и т.д.).

## Base URL

Заказчик передаёт публичный URL, например:

```text
https://1c.example.ru/accounting
```

В env backend:

```env
ACCOUNTING_ADAPTER=onec
ONEC_BASE_URL=https://1c.example.ru/accounting
ONEC_USERNAME=...
ONEC_PASSWORD=...
# или ONEC_API_KEY=...
```

## Endpoint 1: экспорт заявки

```http
POST /hs/palomnik/booking
Content-Type: application/json
```

### Request body

```json
{
  "booking_id": "uuid",
  "tour_id": "uuid",
  "name": "Иван Иванов",
  "phone": "+79990000000",
  "email": "client@example.com",
  "people_count": 2,
  "total_price": 24000,
  "status": "CONFIRMED",
  "comment": "комментарий",
  "source": "web",
  "overbooked": false
}
```

### Response (обязательно)

HTTP 200, JSON:

```json
{
  "external_id": "DOC-000123",
  "status": "ok",
  "message": "optional"
}
```

`external_id` — номер/ссылка документа в 1C (реализация, заказ покупателя и т.п.).

### Идемпотентность

Повторный запрос с тем же `booking_id` должен обновлять или возвращать существующий документ, не создавать дубликат.

## Endpoint 2: контрагент

```http
POST /hs/palomnik/counterparty
Content-Type: application/json
```

### Request body

```json
{
  "booking_id": "uuid",
  "name": "Иван Иванов",
  "phone": "+79990000000",
  "email": "client@example.com"
}
```

### Response

```json
{
  "external_id": "CP-000456",
  "status": "ok"
}
```

## Кастомные пути

Если в 1C другие URL:

```env
ONEC_EXPORT_BOOKING_PATH=/hs/exchange/booking
ONEC_COUNTERPARTY_PATH=/hs/exchange/counterparty
```

## Что нужно от 1C-специалиста

1. Опубликовать HTTP-сервис (IIS / встроенный публикация 1C)
2. Реализовать два метода POST по контракту выше
3. Создать/найти контрагента по телефону
4. Создать документ учёта для подтверждённых заявок (статус `CONFIRMED` / `COMPLETED` — согласовать с заказчиком)
5. Выдать credentials для backend

## Тест без 1C

Backend покрыт mock-тестами. Для ручной проверки:

```bash
make integration-smoke-docker
```

Поднимает `mock-onec` на `:8092` с эндпоинтами из этого ТЗ. Debug: `GET http://localhost:8092/accounting/debug/documents`.

## Контакт с backend-командой

При изменении обязательных полей — обновить `docs/INTEGRATION_DISCOVERY.md` и согласовать до деплоя.
