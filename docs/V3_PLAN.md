# План линейки v3

Канон: [ROADMAP.md](ROADMAP.md) · [STATUS.md](STATUS.md) · [AGENTS.md](../AGENTS.md).  
База прода: тег `v2.1.0`. Деплой и тег v3 — только по просьбе владельца.

Без ключа любой адаптер `Configured()==false`, сайт жив. Live Bitrix24 / 1С нет. ИИ-звонки и автономный ИИ-менеджер продаж — **v4**.

**Правило подключений:** код адаптера пишем до включения в UX, с mock-тестами и переменными в `.env.example`. Когда ключ появится — только `.env.production` + `make deploy`.

Оплата (Сбер / ЮKassa) — позже: договор эквайринга не блокирует остальной v3. Оба адаптера одного `PaymentPort`, выбор в env.

```text
0 ветка → 1 платформа → 2 админка и 3 кабинет
                         ↓
              4 адаптеры «вписать ключ»
           ↙        ↓         ↘
     5 чаты/бот   6 SMM      7 ИИ
           ↘        ↓         ↙
              8 оплата → 9 UX → 10 доки
```

| Этап | Зачем сейчас | Ключи владельца |
|------|----------------|-----------------|
| 0 Подготовка | Не смешать WIP с v3 | нет |
| 1 Платформа | Защитить прод; диск ≥80 ГБ, IP статическим | Object Storage — когда будет |
| 2 Админка | Сотрудники работают без внешних API | нет |
| 3 Кабинет | Слияние и пассажиры до мессенджеров | SMTP / sms.ru уже в V2_OWNER_SETUP |
| 4 Адаптеры | Код под ключ; перед живым ботом/ИИ — 4 vCPU / 8 ГБ | TG канал, Max, WABA, VK, SmartCaptcha, YandexGPT, S3 |
| 5 Чаты и бот | Продукт на готовом MessengerPort | те же, что этап 4 |
| 6 SMM | Календарь на готовых PublisherPort | те же |
| 7 ИИ | Фичи на готовом AIPort | YandexGPT |
| 8 Оплата | Когда договор эквайринга | Сбер и/или ЮKassa |
| 9 UX | Полировка один раз | нет |
| 10 Доки | Чеклист «вписать сюда» | — |

---

## Инфра: мощность, порты, DNS, IP

Одна ВМ (`93.77.165.81`), Docker Compose, Caddy. Снаружи только **22 / 80 / 443**. Postgres, Redis, API, frontend в проде без публикации портов. Сайт и API на одном IP (SNI). Исходящий Telegram — через Cloudflare Worker; второй IP из‑за `api.telegram.org` не нужен.

Вторая ВМ, второй белый IP и новые порты security group в v3 **не закладываем**. Мессенджеры, капча, Сбер/ЮKassa, WhatsApp Cloud, VK ходят HTTPS на `https://api.tikhvin-palomnik.ru/api/v1/webhooks/…` и return URL на сайт. Не открывать 5432 / 6379 / 8080.

Перед ресайзом: закрепить NAT IP как **статический**. Ресайз вертикальный (стоп → cores/RAM/диск → старт), A-записи те же.

| Когда | Что | Зачем |
|--------|-----|--------|
| До этапа 1 | Диск **≥80 ГБ**. RAM до **4 ГБ**, если меньше. Swap оставить | Дампы, слои Docker, Redis, логи |
| Перед этапами 4–7 (живой бот/SMM/GPT) | Та же ВМ: **4 vCPU / 8 ГБ** | Иначе swap на 2 vCPU |
| Этап 8 | Обычно без апсайза | Считает шлюз оплаты |
| После v3, watchdog >70% | Сначала диск/RAM. Вторая ВМ / Managed Postgres — только по факту | |

Сборка: `docker builder prune` после деплоя, чтобы не забивать диск.

DNS (REG.RU), тот же A `93.77.165.81`: `@`, `www`, `api` уже есть. MX/SPF/DKIM — Яндекс 360, не эта ВМ. Отдельные `pay.` / `bot.` / `wa.` не делаем. Опционально `media` CNAME на бакет — если фото уедем с диска ВМ. IPv6 не обязателен. SSH можно сузить на IP владельца (этап 1), без смены DNS.

Не покупать заранее: балансировщик, второй белый адрес, «IP для WhatsApp», открытый SMTP на ВМ.

---

## Этап 0 — подготовка

Ветка работы: `cursor/v3-platform-926b` от `main` (облачное имя; канон плана «ветка v3» сохранён по смыслу). В [ARCHITECTURE.md](ARCHITECTURE.md) порты: Cache, Messenger, Publisher, Captcha, AI, Payment, BackupStorage.

**Сделано в коде:** порты + noop; SmartCaptcha и S3-адаптеры с `Configured()==false` без ключа.

---

## Этап 1 — платформа

Уже есть: Redis, `CachePort` у туров (пока не вызывается), дамп в 03:00 на диск, slog, CSP, in-memory лимитер, секрет Telegram webhook.

- Кэш публичных туров и новостей; Redis down ≠ падение сайта.
- Лимиты auth/заявка/чат/сброс пароля через Redis, Retry-After.
- Логи: request_id наружу; не логировать токены, пароли, паспорт, сырые payment callback.
- Шаблон входящих webhook: HMAC/secret, constant-time, идемпотентность.
- `CaptchaPort` + Яндекс SmartCaptcha + honeypot. Без ключа — honeypot+лимит.
- `BackupStoragePort` + Object Storage; nightly dump + offsite когда ключ есть.
- Срез для разработчика: outbox, latency, last backup. Не публичный pprof.

**Сделано в коде этапа 1** (без деплоя, без live-ключей): Redis down не валит сайт и `/health/ready`; публичный кэш туров/новостей; лимиты auth/заявка/чат/сброс через Redis с fallback в память и `Retry-After`; `X-Request-ID` + `request_id` в ошибках; webhook secret constant-time + идемпотентность; honeypot на публичных формах; SmartCaptcha только если есть ключ; `cmd/backup-offsite` + отметка last-backup; system-info показывает latency и бэкап.

---

## Этап 2 — админка

- Шаблоны ролей: менеджер заявок, рекламщик, сммщик, директор, разработчик. `manage_roles` только у `ADMIN_TOKEN`.
- Меню и дашборд по правам.
- Фильтры заявок + CSV.
- `/management/support` — зеркало того же треда, что и бот.
- Рекламщик: Метрика, не выдуманные визиты.

**Сделано в коде:** меню/дашборд по permission; фильтры заявок и CSV. Шаблоны ролей — пресеты формы «Создать роль» (`GET /management/roles/templates`): не пишут в БД, `manage_roles` в набор не входит. Карточка «Визиты» — при `view_stats`, цифры визитов не рисуются.

---

## Этап 3 — кабинет паломника

- Identity `(user_id, provider, subject)`.
- Привязка → слияние в текущего (заявки, избранное, треды, пассажиры); конфликт полей не перезаписывать молча.
- Правка профиля.
- Пассажиры: ФИО, телефон, ДР, паспорт; автозаполнение; маска в мессенджерах. СНИЛС нет.
- Сброс пароля и письмо подтверждения — существующий Mailer.
- Политика ПД — только текст юриста.

**Сделано в коде (identity):** таблица `user_identities`, OAuth-логин через identity, привязка при сессии (`Authorization: Bearer` на внутреннем `POST /auth/oauth`) и слияние в текущего пользователя (заявки, избранное, треды поддержки, назначения ролей, пассажиры). Конфликт имени/почты/телефона не перезаписывается. СНИЛС нет.

**Сделано в коде (профиль):** `PATCH /me` — имя, почта, телефон текущего кабинета. Смена телефона при включённом callcheck — тот же звонок, что при регистрации. Пароль через существующее восстановление.

**Сделано в коде (пассажиры):** таблица `passengers` (ФИО, телефон, ДР, паспорт; СНИЛС нет). CRUD `/me/passengers`, кабинет `/account/passengers`, автозаполнение имени и телефона из профиля. В мессенджеры — маска телефона, паспорта и ДР. Слияние кабинета переносит пассажиров.

Политика ПД: код согласий влит владельцем (#22, goose 18–20). Тексты не сертифицированы юристом. См. [legal/README.md](legal/README.md).

---

## Этап 4 — адаптеры «вписать ключ»

| Порт | Адаптеры | Включение |
|------|----------|-----------|
| MessengerPort | telegram, max, whatsapp Cloud API | токен / WABA |
| PublisherPort | site_news, telegram_channel, vk_wall, max_feed | токен + id канала |
| CaptchaPort | smartcaptcha | ключ виджета |
| AIPort | yandexgpt, noop | API-ключ |
| BackupStoragePort | s3/yandex, noop | ключ бакета |
| Phone / Mailer | уже в v2 | sms.ru / SMTP |

WhatsApp — только официальный Cloud API. `ExportPayment` на AccountingPort — заготовка, 1С live нет.

**Сделано в коде (MessengerPort):** `MESSENGER_ADAPTER=telegram` / `max` / `whatsapp` (по умолчанию `noop`). Без ключа — noop, сайт жив. Telegram ходит в тот же Bot API / Worker, что уведомления. Max — `POST https://platform-api2.max.ru/messages`. WhatsApp — Graph Cloud API. Compose прокидывает `MESSENGER_ADAPTER` и WhatsApp/Max env. Чаты и бот-команды не включены (этап 5). Чеклист: [V3_OWNER_SETUP.md](V3_OWNER_SETUP.md).

**Сделано в коде (PublisherPort):** `PUBLISHER_ADAPTER=noop` (по умолчанию) или `live` / один канал `site_news` / `telegram_channel` / `vk_wall` / `max_feed`. Без ключа — noop. `live` вызывает только настроенные каналы. Telegram-канал — тот же бот и Worker; бот должен быть админом, webhook не трогаем. VK — официальный `wall.post`, `owner_id` сообщества отрицательный. Max — `POST /messages?chat_id=`. SMM-календарь не включён (этап 6). Чеклист: [V3_OWNER_SETUP.md](V3_OWNER_SETUP.md).

**Сделано в коде (AIPort):** `AI_ADAPTER=noop` (по умолчанию) или `yandexgpt`. Без ключа — noop. Официальный `POST https://llm.api.cloud.yandex.net/foundationModels/v1/completion`, заголовок `Authorization: Api-Key`. Нужны `YANDEXGPT_API_KEY` и `YANDEXGPT_FOLDER_ID`. Фичи поддержки/рекомендаций/watchdog не включены (этап 7). v4 (звонки, ИИ-продавец) сюда не входит. Чеклист: [V3_OWNER_SETUP.md](V3_OWNER_SETUP.md).

---

## Этап 5 — чаты и бот

Клиент на сайте → тред → менеджерам в каналы из Настроек. Ответ в боте → staff в тот же тред → сайт; дубль в мессенджер клиента, если привязан. Команды: заявки и туры (слоты, цена, вкл/выкл) с теми же Permission, что HTTP.

**Сделано в коде:** тот же webhook и `TELEGRAM_API_BASE` (на проде Worker). `/start` и `/health` не менялись. Реплай на уведомление с id диалога или `/reply <id> текст` пишет staff-сообщение в тот же тред, если username в получателях поддержки **или** у кабинета есть `manage_support` через identity `telegram` + назначение роли. Команды `/bookings`, `/booking`, `/tours`, `/tour` (slots / price / on|off) проверяют те же Permission, что HTTP. Fan-out ответа клиенту — `MessengerPort.Send` на привязанный identity; при `MESSENGER_ADAPTER=noop` это no-op. Телефон в тексте команд маскируется. Уведомления `NOTIFICATION_ADAPTER=telegram` не выключаются.

---

## Этап 6 — контент-план

Материал + слот времени + список publisher. Источник сначала `/management/smm`, позже тот же порт может читать таблицу. Страница `/news/[slug]`. Текст поста только из плана. Падение одного канала не откатывает остальные.

**Сделано в коде:** таблица `smm_posts` (goose 17), админка `/management/smm`, `POST /management/smm/:id/publish` и публикация due в worker. Текст = title/body/url из плана. Каналы — имена PublisherPort. `PUBLISHER_ADAPTER=noop` на проде ничего не постит. Публичная страница `/news/[slug]` для опубликованных статей.

---

## Этап 7 — ИИ

- Поддержка: первая линия и черновик менеджеру; эскалация человеку; не выдумывать цены и богословие.
- Рекомендации только из опубликованных туров.
- Дайджест метрик директору/рекламщику.
- Watchdog: health, диск, outbox, 5xx, просроченный бэкап → отчёт, **без** рестарта прода.

**Сделано в коде:** `POST /management/support/:id/draft` — черновик менеджеру, в тред и клиенту не пишется, эскалация человеку всегда. `GET /tours/:id/recommendations` — только `IsActive` туры (неактивные и выдуманные id отбрасываются; без ключа — остальные опубликованные). `GET /management/ai/metrics-digest` (`view_stats`) — заявки по статусам, активные туры, открытые диалоги, outbox; визитов нет. `GET /management/watchdog` и лог worker раз в 5 мин: БД, диск, outbox, 5xx с старта API, бэкап старше 26 ч; `restart_attempted` всегда false. На проде `AI_ADAPTER=noop` — фичи no-op, сайт жив. v4 не начат.

---

## Этап 8 — оплата

`PaymentPort`, адаптеры `sber` и `yookassa`, `PAYMENT_ADAPTER`. Сумма = `total_price`. Возвраты не делаем.

`AWAITING_PAYMENT` → `PAID` / `CANCELLED` (слоты назад) → `CONTACTED` → `CONFIRMED` → `COMPLETED`. Уведомление «новая заявка» на `PAID`. Без адаптера форма как в v2 (`NEW`).

**Сделано в коде (адаптеры, без живого эквайринга):** `PAYMENT_ADAPTER=noop` по умолчанию; `sber` (официальный `register.do`) и `yookassa` (`POST /v3/payments`). Сумма только `booking.TotalPrice`. Возвратов нет. **Статусы `AWAITING_PAYMENT` / `PAID` в domain нет** — живая машина статусов ждёт подтверждения владельца, в код их не добавляли. Checkout заявки не переключает статус. На проде ключей нет, оставлять noop.

---

## Этап 9 — UX

Каталог ближе к расписанию (даты, цена, длительность; тексты не копировать). Sticky CTA, карточки, кабинет (привязки, пассажиры, оплата), админка по ролям.

**Сделано в коде (без новых API и без новых полей тура):** `/search` — h1 «Расписание»; таблица/карточки из `date_start`, `date_end`, title, location, price, `slots_left`; длительность на фронте = число дней `date_end − date_start + 1`. Главная: CMS-блок популярных направлений читает `GET /tours?limit=8` (уже `ORDER BY date_start ASC`) и показывает «Ближайшие выезды». Sticky CTA на узком экране: `tel:` из контакта сайта и ссылка «Чат» на `/support/chat`. Карточка тура без выдержки description. Кабинет: счётчики привязок и пассажиров, блок оплаты «на сайте не подключена»; в поездках сумма = `total_price`, кнопки Pay нет. Админка туров — колонка «Длительность» тем же расчётом.

---

## Этап 10 — документы

Обновить STATUS, DECISIONS §12, DATA_MODEL, API, ARCHITECTURE.  
Чеклист ключей: [V3_OWNER_SETUP.md](V3_OWNER_SETUP.md).

**Сделано:** перечисленные файлы приведены к факту прода (этапы 0–9, goose 20, #22 на main, оплата noop). Юридические тела документов не переписывались.

Проверки этапа: `go test ./...`, `go vet`, `npm run lint` / `build`. Прод — `make deploy` по просьбе.

---

## Не входит в v3

- Live Bitrix24 / 1С
- Возвраты, рассрочка
- Неофициальный WhatsApp
- Посты «из нейросети»
- Фейковые визиты
- СНИЛС
- ИИ-звонки и ИИ-продавец (v4)
- Watchdog, который сам меняет прод
