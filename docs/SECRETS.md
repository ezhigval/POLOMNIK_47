# Секреты: откуда взять и куда класть

Имена переменных — в [`.env.production.example`](../.env.production.example). **Значения в git, в чат и в PR не писать.**

Почта DNS (MX/SPF/DKIM) — отдельный файл: [MAIL_DNS.md](MAIL_DNS.md).  
Соцвход по шагам кабинетов: [OAUTH_SETUP.md](OAUTH_SETUP.md).  
Telegram-бот: [TELEGRAM_SETUP.md](TELEGRAM_SETUP.md).

После правок `.env.production` на ВМ нужен **`make deploy`** (пересборка фронта). Пока GitHub Actions падает на SSH — деплой только вручную, по вашей просьбе. Не `compose down -v`.

---

## Два места хранения

| Куда | Что |
|------|-----|
| ВМ **`/opt/polomnik/.env.production`** | Почти все секреты сайта и API. Файл gitignored. |
| GitHub → репозиторий → **Settings → Secrets and variables → Actions** | Только деплой: `DEPLOY_SSH_KEY`, опционально `DEPLOY_SSH_HOST`. **Не** копировать сюда весь `.env.production`. |

Админка **Настройки** — не секреты: имя сайта, телефон, email, пересылка, получатели уведомлений (канал+адрес), роли. OAuth-ключи и SMTP-пароль туда **не** дублировать.

---

## 1. Уже на проде — не ротировать без нужды

Эти значения уже лежат в `.env.production` на ВМ. Если «потеряли» — **открыть файл на сервере**, не генерировать вторую копию параллельно (иначе отвалятся сессии, админка и Telegram).

| Имя | Зачем | Откуда когда-то взяли |
|-----|--------|------------------------|
| `ADMIN_TOKEN` | Вход в `/management` (пустая роль + этот пароль) | Сгенерировали при установке (≥11 символов). Смотреть на ВМ. |
| `JWT_SECRET` | Сессии кабинета паломника | Сгенерировали (≥32 символов). Смотреть на ВМ. |
| `INTERNAL_API_SECRET` | Next (сайт) → API, в том числе OAuth | Сгенерировали (≥16). **Один и тот же** у контейнеров `frontend` и `api`. Без него соцвход не завершится, даже если OAuth-ключи есть. |
| `POSTGRES_PASSWORD` | Пароль роли БД | Сгенерировали при установке. **Не переименовывать** `POSTGRES_USER` / `POSTGRES_DB` в уже живом томе. |
| `TELEGRAM_BOT_TOKEN` | Один бот `@Tikhvinpalomnik_bot`: заявки, поддержка, webhook, позже Login | [@BotFather](https://t.me/BotFather) → `/mybots` → бот → **API Token**. Уже в проде. |
| `TELEGRAM_BOT_USERNAME` | `Tikhvinpalomnik_bot` | Тот же BotFather. |
| `TELEGRAM_API_BASE` | Обход блокировки `api.telegram.org` с ВМ (Cloudflare Worker) | URL воркера **уже в** `.env.production` на ВМ. Не выдумывать. Токен самого воркера — в кабинете Cloudflare, не в git. |
| `TELEGRAM_WEBHOOK_URL` | Входящие апдейты бота | Уже: `https://api.tikhvin-palomnik.ru/api/v1/webhooks/telegram`. Не менять ради Login. |
| `NOTIFICATION_ADAPTER` | `telegram` | Не выключать. |
| `NEXT_PUBLIC_YM_ID` | Метрика (публичный номер) | Кабинет Метрики; прод **111985266**. |
| `ACME_EMAIL` | Let's Encrypt у Caddy | Уже `smailikin70@yandex.ru`. |

`TELEGRAM_CHAT_ID` — только запас, если в админке пустой список получателей. Живые получатели: **Настройки** (telegram + username); человек один раз пишет боту `/start`.

---

## 2. OAuth (соцвход) — откуда взять

Кнопки на `/account/login` и `/account/register` уже есть. Пока переменных нет — «Пока что недоступно…». Google в UI скрыт, приложение Google **не** создавать.

**Redirect URI — на сайт, не на `api.*`.** Compose уже прокидывает OAuth-переменные и в `frontend`, и в `api`.

### Общая таблица

| Провайдер | Кабинет | Что скопировать | Redirect / домен | Переменные |
|-----------|---------|-----------------|------------------|------------|
| Яндекс ID | [oauth.yandex.ru](https://oauth.yandex.ru/) | ClientID, Client secret | `https://tikhvin-palomnik.ru/api/auth/social/yandex/callback` | `YANDEX_OAUTH_CLIENT_ID`, `YANDEX_OAUTH_CLIENT_SECRET` |
| VK | [id.vk.com/about/business/go](https://id.vk.com/about/business/go) или [vk.com/apps?act=manage](https://vk.com/apps?act=manage) | App ID, защищённый ключ | `https://tikhvin-palomnik.ru/api/auth/social/vk/callback` | `VK_OAUTH_CLIENT_ID`, `VK_OAUTH_CLIENT_SECRET` |
| Telegram Login | [@BotFather](https://t.me/BotFather), **тот же** бот | домен виджета; токен уже есть | домен `tikhvin-palomnik.ru` (без `https://`) | токен тот же; опц. `TELEGRAM_LOGIN_BOT_USERNAME=Tikhvinpalomnik_bot` |
| Max | консоль Max, если выдадут OAuth2 | Client ID/Secret + три URL | `https://tikhvin-palomnik.ru/api/auth/social/max/callback` | `MAX_OAUTH_CLIENT_ID`, `MAX_OAUTH_CLIENT_SECRET`, `MAX_OAUTH_AUTHORIZE_URL`, `MAX_OAUTH_TOKEN_URL`, `MAX_OAUTH_USERINFO_URL` |

Локально (если проверяете): тот же путь, origin `http://localhost:3000`.

### Яндекс ID — клики

Официально: [регистрация приложения для авторизации](https://yandex.ru/dev/id/doc/ru/register-auth).

1. Войдите тем Яндекс-аккаунтом, к которому не потеряете доступ: [https://oauth.yandex.ru/](https://oauth.yandex.ru/).
2. **Создать приложение** → тип **Для авторизации пользователей** (прямая ссылка: [https://oauth.yandex.ru/client/new/id/](https://oauth.yandex.ru/client/new/id/)).
3. Название сервиса (как на сайте), контактная почта → **Продолжить**.
4. Платформа: **Веб-сервисы**. Redirect URI **точно**:

   `https://tikhvin-palomnik.ru/api/auth/social/yandex/callback`

   Без этого адреса вход после разрешения у Яндекса сломается (`redirect_uri mismatch`).
5. Доступ к данным — **минимум**, что реально читает код: логин / имя / фамилия; **адрес электронной почты**. Телефон, дата рождения, портрет — не обязательны.
6. Сохраните. Со страницы приложения скопируйте **ClientID** и **Client secret** в `.env.production` на ВМ:

```bash
YANDEX_OAUTH_CLIENT_ID=
YANDEX_OAUTH_CLIENT_SECRET=
```

Значения **не** коммитить и не класть в GitHub Secrets. Пара, переданная владельцем в чат, пишется только в этот файл на сервере.

7. Верификация приложения в кабинете Яндекса снижает предупреждение «неизвестное приложение»; на работу callback не влияет.

### VK — клики

Код ходит на классический `oauth.vk.com` (code + `client_id` / `client_secret`), не VK ID SDK.

1. Кабинет VK ID: [https://id.vk.com/about/business/go](https://id.vk.com/about/business/go) → создать приложение, платформа **Web**.  
   Или старый список: [https://vk.com/apps?act=manage](https://vk.com/apps?act=manage) → **Веб-сайт**.
2. Базовый домен: `tikhvin-palomnik.ru` (или `https://tikhvin-palomnik.ru` — как просит форма).
3. Доверенный **Redirect URL точно**:

   `https://tikhvin-palomnik.ru/api/auth/social/vk/callback`

4. **App ID** (id приложения) → `VK_OAUTH_CLIENT_ID`.  
   **Защищённый ключ** (не сервисный ключ, если показывают оба) → `VK_OAUTH_CLIENT_SECRET`.
5. Если есть переключатель «почта» / email — включите: код запрашивает `scope=email`.

### Telegram Login — клики

Отдельный бот и отдельный токен **не нужны**.

1. Telegram → [@BotFather](https://t.me/BotFather) → `/mybots` → **@Tikhvinpalomnik_bot**.
2. **Bot Settings → Domain** (команда `/setdomain`): **`tikhvin-palomnik.ru`** без `https://` и без пути.
3. Webhook уведомлений (`/api/v1/webhooks/telegram` на `api.*`) **не менять**.
4. Опционально в env:

```bash
TELEGRAM_LOGIN_BOT_USERNAME=Tikhvinpalomnik_bot
# TELEGRAM_LOGIN_BOT_TOKEN пустой = берётся TELEGRAM_BOT_TOKEN
```

### Max — клики

Кнопка Max **не заработает**, пока нет всех пяти переменных: id, secret и трёх URL. В коде **нет** URL по умолчанию.

1. Redirect заранее пропишите в консоли Max, когда она это умеет:

   `https://tikhvin-palomnik.ru/api/auth/social/max/callback`

2. Если консоль показывает Authorize / Token / UserInfo — вставьте их в `MAX_OAUTH_AUTHORIZE_URL`, `MAX_OAUTH_TOKEN_URL`, `MAX_OAUTH_USERINFO_URL`.
3. **Не путать** с `MAX_BOT_TOKEN` из [business.max.ru](https://business.max.ru) — это токен **бота** (уведомления / Publisher), не OAuth входа на сайт. Пока адаптеры Max noop — бот-токен не обязателен.

Если кабинета OAuth2 у Max ещё нет — оставьте `MAX_OAUTH_*` пустыми. Кнопка останется «недоступно».

### После заполнения OAuth

Вписать в `.env.production` на ВМ → `make deploy` → hard-refresh → `/account/login`: рабочие кнопки ведут на провайдера, затем в `/account/trips`. Если уже были в кабинете, привязка: `/account?linked=1`.

---

## 3. Почта SMTP (после MX/SPF/DKIM)

Пошагово DNS и ящик: [MAIL_DNS.md](MAIL_DNS.md). Кратко, откуда секрет:

| Имя | Откуда |
|-----|--------|
| `MAIL_ADAPTER=smtp` | Включает отправку с сайта |
| `SMTP_HOST=smtp.yandex.ru` | Фиксировано Яндексом |
| `SMTP_PORT=587` | Как в коде (STARTTLS). Не 465, пока код не сменят |
| `SMTP_USERNAME` / `SMTP_FROM` | `info@tikhvin-palomnik.ru` |
| `SMTP_PASSWORD` | [Пароль приложения](https://id.yandex.ru/security/app-passwords) **учётки ящика info@**, тип «Почта». Не пароль входа на Яндекс и **не** Client secret OAuth-приложения |
| `MAIL_FORWARD_TO` | Не секрет; лучше админка «Настройки» |

OAuth ClientID/secret Яндекса и пароль SMTP — **разные** вещи из разных кабинетов.

В ящике включить IMAP и «Пароли приложений» (настройки Почты → Почтовые программы).

---

## 4. GitHub Actions (автодеплой)

Сейчас workflow **красный**: `Permission denied (publickey)`.

1. GitHub → этот репозиторий → **Settings → Secrets and variables → Actions**.
2. Секрет **`DEPLOY_SSH_KEY`**: вставить **приватный** ключ OpenSSH целиком (`-----BEGIN … KEY-----` …), не `.pub` и не passphrase-загадку без ключа.
3. Ключ должен открывать `smailikin70@93.77.165.81`. На сервере уже есть pubkey комментария `github-actions-deploy-palomnik`:

   `ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIA/Pf0QxoAH4bbrcC5bRW9Y7lN0OPigypFRn+0KI0HQZ`

4. Если **приватная** половина этого ключа есть у вас в сейфе — вставляете её. Если потеряна: сгенерировать **новую** пару, **публичную** дописать в `authorized_keys` на ВМ (старую можно оставить), приватную — в секрет.
5. Опционально `DEPLOY_SSH_HOST`, если host не `smailikin70@93.77.165.81`.

Пока секрет неверный — релиз только `make deploy`.

---

## 5. По желанию (формы входа / звонок)

| Имя | Откуда | Зачем |
|-----|--------|--------|
| `PHONE_ADAPTER=smsru` + `SMSRU_API_ID` | Кабинет [sms.ru](https://sms.ru/) → API ID | Звонок-проверка номера, **не** SMS. Без ключа кнопка «пока недоступно». |

---

## 6. Пока не трогаем (оставить noop / пусто)

Код есть, в эфир ничего не уходит. Ключи **не** заводить, пока не попросите live.

| Адаптер | Переменные | Откуда были бы |
|---------|------------|----------------|
| `CRM_ADAPTER=noop` | `BITRIX_WEBHOOK_URL`, `BITRIX_INBOUND_TOKEN` | Входящий вебхук Bitrix24 |
| `ACCOUNTING_ADAPTER=noop` | `ONEC_*` | 1С / интегратор |
| `CAPTCHA_ADAPTER=noop` | `SMARTCAPTCHA_SERVER_KEY`, `SMARTCAPTCHA_CLIENT_KEY` | [Яндекс SmartCaptcha](https://cloud.yandex.ru/services/smartcaptcha) |
| `BACKUP_STORAGE_ADAPTER=noop` | `S3_*` | Yandex Object Storage: бакет + статический ключ |
| `MESSENGER_ADAPTER=noop` | `WHATSAPP_TOKEN`, `WHATSAPP_PHONE_NUMBER_ID`, `MAX_BOT_TOKEN` | WhatsApp Cloud API / Max bot |
| `PUBLISHER_ADAPTER=noop` | `TELEGRAM_CHANNEL_ID`, `VK_WALL_TOKEN`, `VK_WALL_OWNER_ID`, `MAX_FEED_CHAT_ID` | Канал/стена; `VK_WALL_*` **не** те же, что `VK_OAUTH_*` |
| `AI_ADAPTER=noop` | `YANDEXGPT_API_KEY`, `YANDEXGPT_FOLDER_ID` | Yandex Cloud: сервисный аккаунт, API-ключ (не IAM на 12 ч) |
| `PAYMENT_ADAPTER=noop` | `SBER_*`, `YOOKASSA_*` | Договор эквайринга; **не включать** до решения по живой оплате |
| `OPERATOR_*` / `NEXT_PUBLIC_OPERATOR_*` | ИНН, ОГРН, адреса | Только из ваших документов / юриста, **не выдумывать** |
| `NEXT_PUBLIC_GA_ID`, коды Вебмастера/GSC | кабинеты Google/Яндекс | SEO; отдельно от этой инструкции |
| `GOOGLE_CLIENT_ID` / `SECRET` | — | Не нужно |

---

## 7. Как править файл на ВМ

```bash
# с вашей машины, тем же ключом, что make deploy
ssh -i ~/.ssh/palomnik_yc smailikin70@93.77.165.81
sudo -e /opt/polomnik/.env.production   # или nano
```

Добавить строки из таблиц выше. Файл не коммитить, не класть в GitHub Secrets целиком. Затем попросить **`make deploy`**.
