# Чеклист владельца к v3 (этап 4)

Код адаптеров можно вписать ключом. Чаты, бот-команды, SMM и ИИ — следующие этапы; без ключа сайт жив (`Configured()==false`).

Сайт: **https://tikhvin-palomnik.ru**  
API: **https://api.tikhvin-palomnik.ru**  
Секреты только в gitignored `.env.production`. После правок env — `make deploy` (без `compose down -v`).

v2-секреты (OAuth, smtp, sms.ru, Telegram-уведомления): [V2_OWNER_SETUP.md](V2_OWNER_SETUP.md).

---

## MessengerPort (этап 4)

Чат на сайте и ответы в боте — этап 5. Сейчас адаптер только отправляет текст, если его вызвать.

В `.env.production` один из вариантов:

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

## Уже в коде (ключи по желанию)

| Порт | Env |
|------|-----|
| SmartCaptcha | `CAPTCHA_ADAPTER=smartcaptcha`, `SMARTCAPTCHA_SERVER_KEY`, `SMARTCAPTCHA_CLIENT_KEY` |
| Object Storage | `BACKUP_STORAGE_ADAPTER=s3`, `S3_*` |
| Telegram уведомления заявок/поддержки | `NOTIFICATION_ADAPTER=telegram` (уже на проде) |

PublisherPort, YandexGPT, оплата Сбер/ЮKassa — отдельные PR.

---

## Инфра перед живым ботом / ИИ (этапы 5–7)

Та же ВМ: **4 vCPU / 8 ГБ**, диск ≥80 ГБ. Порты снаружи только 22 / 80 / 443.
