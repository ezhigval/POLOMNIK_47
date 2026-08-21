# Telegram: уведомления менеджеру

Код готов. **Подключение откладывается на финальный этап** вместе с Bitrix24 и 1С.

## Роль

- Уведомление в Telegram-чат при **новой заявке**
- Уведомление при **смене статуса** заявки в management
- Повтор через **outbox** при ошибке API Telegram

По умолчанию: `NOTIFICATION_ADAPTER=noop` — уведомления не отправляются, сайт работает как раньше.

## Архитектура

```text
NOTIFICATION_ADAPTER=noop     -> noop (default)
NOTIFICATION_ADAPTER=telegram -> telegram adapter + recording -> outbox on failure
```

События outbox:

- `notification.booking.created`
- `notification.booking.status_changed`

Worker (`cmd/worker`) повторяет доставку без повторной записи в outbox.

## Быстрый старт (когда дойдёте до подключения)

### 1. Создать бота

1. Открой [@BotFather](https://t.me/BotFather) в Telegram
2. `/newbot` → имя и username
3. Сохрани **token**

### 2. Узнать chat_id

1. Напиши боту любое сообщение (или добавь в группу менеджеров)
2. Открой в браузере:

```text
https://api.telegram.org/bot<TOKEN>/getUpdates
```

3. Найди `"chat":{"id": ...}` — это `TELEGRAM_CHAT_ID`

### 3. Env

```env
NOTIFICATION_ADAPTER=telegram
TELEGRAM_BOT_TOKEN=123456:ABC...
TELEGRAM_CHAT_ID=-1001234567890
MANAGEMENT_BASE_URL=https://ваш-домен.ru/management/bookings
```

### 4. Перезапуск

```bash
docker compose up --build -d api worker
```

### 5. Проверка

1. Отправь тестовую заявку на сайте
2. Проверь сообщение в чате
3. Смени статус в management — должно прийти второе уведомление
4. При сбое смотри **Outbox — Telegram** в `/management/integrations`

## Management

`GET /api/v1/management/system-info` — статус адаптеров без секретов:

- `notification_adapter`
- `telegram_configured` (token + chat_id заданы)

## Ограничения MVP

- Только исходящие сообщения (нет inline-кнопок / команд бота)
- Один chat_id на все уведомления
- HTML-разметка в тексте

См. также [INTEGRATIONS.md](INTEGRATIONS.md) и [DEPLOY.md](DEPLOY.md).
