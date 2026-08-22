# Telegram: один бот на уведомления и поддержку

Один токен (`TELEGRAM_BOT_TOKEN`) на всё: заявки, поддержка, webhook `/health`, позже login widget. Второй `TELEGRAM_LOGIN_BOT_TOKEN` не нужен — login env берёт тот же токен и username, если они не заданы отдельно.

`TELEGRAM_CHAT_ID` — только запасной адрес, если списки в настройках пустые.

## Получатели

В админке **Настройки** (`/management/settings`):

1. В «Новые заявки и смена статуса» добавьте username, например `ezhigval` (без обязательного `@`).
2. В «Сообщения в поддержку» — кому писать о клиентских обращениях «в поддержку».
3. Сохраните.

Человек из списка должен **один раз написать боту** (`/start` или любое сообщение). Telegram не умеет писать тому, кто бота не открывал. Пока этого не было, адрес в списке показывается как «ждёт /start» и при отправке пропускается (сервис не падает).

Команда `/health` в боте отвечает, кто ещё не писал.

Отдельного экрана поддержки в админке нет: менеджер видит текст обращения в Telegram.

## Env (VM)

С ВМ **нельзя** ходить на `api.telegram.org`. Исходящие вызовы — через Worker:

```env
NOTIFICATION_ADAPTER=telegram
TELEGRAM_BOT_TOKEN=...
TELEGRAM_API_BASE=https://polomnik-telegram-api.chemical-red.workers.dev
TELEGRAM_WEBHOOK_URL=https://api.tikhvin-palomnik.ru/api/v1/webhooks/telegram
```

Эти переменные должны быть и у **api**, и у **worker**. Worker повторяет `sendMessage` через тот же `TELEGRAM_API_BASE`.

`TELEGRAM_LOGIN_BOT_TOKEN` / `TELEGRAM_LOGIN_BOT_USERNAME` необязательны: если пусто, используются `TELEGRAM_BOT_TOKEN` и `TELEGRAM_BOT_USERNAME`.

## Архитектура

```text
NOTIFICATION_ADAPTER=noop     -> noop
NOTIFICATION_ADAPTER=telegram -> telegram adapter + recording -> outbox on failure
```

События outbox:

- `notification.booking.created`
- `notification.booking.status_changed`
- `notification.support.message`

Входящий webhook: `POST /api/v1/webhooks/telegram` — пишет `username → chat_id`.

Management API: `GET|PATCH /api/v1/management/telegram-settings`.
