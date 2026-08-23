# Чеклист владельца к v3 (этап 4)

Код адаптеров можно вписать ключом. Чаты, бот-команды, SMM и ИИ — следующие этапы; без ключа сайт жив (`Configured()==false`).

Сайт: **https://tikhvin-palomnik.ru**  
API: **https://api.tikhvin-palomnik.ru**  
Секреты только в gitignored `.env.production`. После правок env — `make deploy` (без `compose down -v`).

v2-секреты (OAuth, smtp, sms.ru, Telegram-уведомления): [V2_OWNER_SETUP.md](V2_OWNER_SETUP.md).

---

## MessengerPort (этап 4)

Чат на сайте и ответы в боте — этап 5. Сейчас адаптер только отправляет текст, если его вызвать.

В `.env.production` один из вариантов. Compose (`docker-compose.yml` и `docker-compose.prod.yml`) прокидывает `MESSENGER_ADAPTER` (по умолчанию `noop`), `WHATSAPP_*` и `MAX_BOT_TOKEN` / `MAX_API_BASE` — без ключа контейнер остаётся noop.

```bash
MESSENGER_ADAPTER=noop
```

| Адаптер | Env | Адрес в Send |
|---------|-----|----------------|
| `telegram` | тот же `TELEGRAM_BOT_TOKEN`; в проде `TELEGRAM_API_BASE` = URL Cloudflare Worker, не `api.telegram.org` | chat_id |
| `max` | `MAX_BOT_TOKEN`; опц. `MAX_API_BASE` (по умолчанию `https://platform-api2.max.ru`) | `user_id` или `chat:<chat_id>` |
| `whatsapp` | официальный Cloud API: `WHATSAPP_TOKEN`, `WHATSAPP_PHONE_NUMBER_ID`; опц. `WHATSAPP_GRAPH_BASE` | номер в международном виде |

WhatsApp — только Cloud API. Неофициальные клиенты не подключаем.

Проверка в админке: `/management/integrations` → карточка MessengerPort (`noop` / ждёт credentials / credentials OK).

---

## PublisherPort (этап 4)

Календарь SMM — этап 6. Сейчас адаптер только публикует, если его вызвать. По умолчанию `noop`: в эфир ничего не уходит.

```bash
PUBLISHER_ADAPTER=noop
```

`PUBLISHER_ADAPTER=live` включает те каналы, у которых уже есть ключ и id. Можно указать один канал: `site_news` / `telegram_channel` / `vk_wall` / `max_feed`.

| Канал в Publish | Env | Куда пишет |
|-----------------|-----|------------|
| `site_news` | достаточно `PUBLISHER_ADAPTER=live` или `site_news` (своя БД) | новость на сайте, `is_published=true`; slug/excerpt из title/body |
| `telegram_channel` | `TELEGRAM_BOT_TOKEN` + `TELEGRAM_CHANNEL_ID`; в проде `TELEGRAM_API_BASE` = Worker | канал Telegram. Бот должен быть админом канала; webhook сам не ставится |
| `vk_wall` | `VK_WALL_TOKEN` + `VK_WALL_OWNER_ID` (id сообщества **со знаком минус**, например `-123456`) | официальный `wall.post`. Не путать с `VK_OAUTH_*` для входа |
| `max_feed` | `MAX_BOT_TOKEN` + `MAX_FEED_CHAT_ID`; опц. `MAX_API_BASE` | `POST https://platform-api2.max.ru/messages?chat_id=` |

Текст поста = title + body + URL из вызова, без дописанных слоганов.

Проверка: `/management/integrations` → карточка PublisherPort.

---

## AIPort (этап 4)

Фичи поддержки/рекомендаций/дайджеста — этап 7. Сейчас адаптер только вызывает `Complete(system, user)`, если его вызвать. По умолчанию `noop`.

```bash
AI_ADAPTER=noop
```

| Адаптер | Env |
|---------|-----|
| `yandexgpt` | `YANDEXGPT_API_KEY` (ключ сервисного аккаунта, заголовок `Api-Key`) + `YANDEXGPT_FOLDER_ID`. Опц. `YANDEXGPT_MODEL` (по умолчанию `yandexgpt-lite`) или полный `modelUri` `gpt://…`. Опц. `YANDEXGPT_API_BASE` (`https://llm.api.cloud.yandex.net`) |

IAM-токен на 12 часов не используем — только долгоживущий API-ключ. Промпты продукта не зашиты: system/user приходят от вызывающего кода этапа 7. Цены и богословие адаптер сам не дописывает.

Проверка: `/management/integrations` → карточка AIPort.

---

## Уже в коде (ключи по желанию)

| Порт | Env |
|------|-----|
| SmartCaptcha | `CAPTCHA_ADAPTER=smartcaptcha`, `SMARTCAPTCHA_SERVER_KEY`, `SMARTCAPTCHA_CLIENT_KEY` |
| Object Storage | `BACKUP_STORAGE_ADAPTER=s3`, `S3_*` |
| Telegram уведомления заявок/поддержки | `NOTIFICATION_ADAPTER=telegram` (уже на проде) |
| MessengerPort | `MESSENGER_ADAPTER` (на проде `noop`) |
| PublisherPort | `PUBLISHER_ADAPTER` (на проде `noop`) |
| AIPort | `AI_ADAPTER` (пока не на проде) |

Оплата Сбер/ЮKassa — отдельный PR (этап 8).

---

## Инфра перед живым ботом / ИИ (этапы 5–7)

Та же ВМ: **4 vCPU / 8 ГБ**, диск ≥80 ГБ. Порты снаружи только 22 / 80 / 443.
