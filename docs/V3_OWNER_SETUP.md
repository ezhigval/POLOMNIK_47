# Чеклист владельца к v3

Этапы **0–10** в коде. Без ключа любой адаптер `Configured()==false`, сайт жив. Чаты/бот, SMM, ИИ-фичи, UX-расписание и доки уже на проде; в эфир и в эквайринг ничего не уходит, пока адаптер `noop`.

Сайт: **https://tikhvin-palomnik.ru**  
API: **https://api.tikhvin-palomnik.ru**  
Секреты только в gitignored `.env.production`. После правок env — `make deploy` (без `compose down -v`).

v2-секреты (OAuth, smtp, sms.ru, Telegram-уведомления): [V2_OWNER_SETUP.md](V2_OWNER_SETUP.md).

---

## MessengerPort

Чат на сайте и команды бота — этап 5, на проде. Fan-out клиенту идёт только если `MESSENGER_ADAPTER` не noop и у кабинета есть identity этого канала. Без ключа сайт жив, уведомления Telegram не выключаются.

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

## PublisherPort

Календарь SMM — этап 6, на проде (`/management/smm`). Адаптер публикует текст из плана. По умолчанию `noop`: в эфир ничего не уходит.

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

## AIPort

Фичи поддержки/рекомендаций/дайджеста/watchdog — этап 7, на проде. Без ключа `Complete` не вызывается, сайт жив. По умолчанию `noop`. На проде ключа YandexGPT нет — оставлять `AI_ADAPTER=noop`. ИИ-звонки и ИИ-продавец — [V4_PLAN.md](V4_PLAN.md) этапы 7–8.

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
| MessengerPort | `MESSENGER_ADAPTER` (на проде `noop`; чаты/бот работают через webhook уведомлений) |
| PublisherPort | `PUBLISHER_ADAPTER` (на проде `noop`) |
| AIPort | `AI_ADAPTER` (на проде `noop`; ключа нет) |

Оплата Сбер/ЮKassa — этап 8 на проде как код адаптеров, **`PAYMENT_ADAPTER=noop`**. Живой эквайринг не включать, пока владелец не подтвердит статусы `AWAITING_PAYMENT` / `PAID` (их ещё нет в domain). Кнопки Pay в кабинете нет.

```bash
PAYMENT_ADAPTER=noop
```

| Адаптер | Env | Сумма |
|---------|-----|--------|
| `sber` | `SBER_USERNAME`, `SBER_PASSWORD`; опц. `SBER_API_BASE` (`https://securepayments.sberbank.ru`) | `booking.TotalPrice` × 100 (копейки), `register.do` |
| `yookassa` | `YOOKASSA_SHOP_ID`, `YOOKASSA_SECRET_KEY`; опц. `YOOKASSA_API_BASE` | `booking.TotalPrice` рубли, `POST /v3/payments` |

Return URL: `PAYMENT_RETURN_URL` или `{PUBLIC_SITE_URL}/account/trips`. Номера договоров не выдумываем. Возвратов нет. Чеклист: вписать ключ только после решения по статусам.

---

## Согласия и реквизиты оператора (#22)

Код и таблицы (goose 18–20) уже на проде. Новых адаптеров нет. Тексты в `legal_documents` **не сертифицированы юристом** — не считать публикацию проверкой юротдела.

Код читает `OPERATOR_*` (`backend/internal/legal/operator`) и `NEXT_PUBLIC_OPERATOR_*` (`frontend/src/lib/operator-config.ts`). С **v4 этапа 1** compose прокидывает эти переменные (пустые = плейсхолдеры `название` / «—»). Не выдумывать ИНН/ОГРН/адреса. Чеклист: [V4_OWNER_SETUP.md](V4_OWNER_SETUP.md).

## Инфра

Та же ВМ: снаружи только 22 / 80 / 443. Перед живым ботом / GPT / SMM в эфир — **4 vCPU / 8 ГБ**, диск ≥80 ГБ. Не `compose down -v`. После деплоя можно `docker builder prune`.
