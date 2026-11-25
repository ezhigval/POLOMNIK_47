# Настройка соцвхода (владелец)

Канон сайта: **https://tikhvin-palomnik.ru**.  
Callback’и — на **сайт** (Next.js), не на `api.tikhvin-palomnik.ru`.

Пути зафиксированы в коде: `frontend/src/lib/auth/social-paths.ts`.

Секреты кладите **только** в gitignored `.env.production` на сервере (или локальный `.env`). Не коммитьте. После заполнения — `make deploy` по вашей просьбе.

Указатель «секрет → кабинет»: [SECRETS.md](SECRETS.md). Нужен уже существующий на проде `INTERNAL_API_SECRET` (один и тот же у frontend и api).

Google в UI **скрыт** — приложение Google не нужно.

---

## Общие Redirect URI (прод)

| Провайдер | Redirect URI / домен |
|-----------|----------------------|
| Яндекс ID | `https://tikhvin-palomnik.ru/api/auth/social/yandex/callback` |
| VK ID | `https://tikhvin-palomnik.ru/api/auth/social/vk/callback` |
| Max | `https://tikhvin-palomnik.ru/api/auth/social/max/callback` |
| Telegram Login | домен виджета: `tikhvin-palomnik.ru` (без `https://`); обработчик: `https://tikhvin-palomnik.ru/api/auth/social/telegram` |

Локально (если проверяете): замените origin на `http://localhost:3000`, те же пути.

---

## 1. Яндекс ID

Официально: [регистрация приложения для авторизации](https://yandex.ru/dev/id/doc/ru/register-auth).

1. Откройте [https://oauth.yandex.ru/](https://oauth.yandex.ru/) тем аккаунтом, к которому не потеряете доступ.
2. **Создать приложение** → **Для авторизации пользователей** (или сразу [https://oauth.yandex.ru/client/new/id/](https://oauth.yandex.ru/client/new/id/)).
3. Название сервиса, контактная почта → платформа **Веб-сервисы**.
4. **Redirect URI** (точно, иначе Яндекс вернёт `redirect_uri mismatch`):  
   `https://tikhvin-palomnik.ru/api/auth/social/yandex/callback`
5. Доступ: логин / имя / фамилия и **email**. Телефон и портрет не обязательны.
6. Сохраните. **ClientID** → `YANDEX_OAUTH_CLIENT_ID`, **Client secret** → `YANDEX_OAUTH_CLIENT_SECRET` — только в `.env.production` на ВМ, не в git.
7. Пока переменных нет на сервере — кнопка Яндекс: «Пока что недоступно…». После записи + `make deploy` кнопка должна открыть oauth.yandex.ru.

Это **не** SMTP и **не** пароль ящика `info@`. Пароль приложения для почты — [MAIL_DNS.md](MAIL_DNS.md).

---

## 2. VK ID

Код использует классический `oauth.vk.com` (не VK ID SDK).

1. Откройте [https://id.vk.com/about/business/go](https://id.vk.com/about/business/go) → приложение **Web**, или [https://vk.com/apps?act=manage](https://vk.com/apps?act=manage) → **Веб-сайт**.
2. Базовый домен: `tikhvin-palomnik.ru`.
3. **Redirect URL** (точно):  
   `https://tikhvin-palomnik.ru/api/auth/social/vk/callback`
4. **App ID** → `VK_OAUTH_CLIENT_ID`, **защищённый ключ** (не сервисный, если показывают оба) → `VK_OAUTH_CLIENT_SECRET`. Если есть опция email — включите (`scope=email` в коде).
5. В `.env.production`:

```bash
VK_OAUTH_CLIENT_ID=...
VK_OAUTH_CLIENT_SECRET=...
```

---

## 3. Telegram Login Widget

Используется **тот же бот**, что и для заявок: `@Tikhvinpalomnik_bot` (токен уже `TELEGRAM_BOT_TOKEN`). Отдельный бот не нужен.

1. Откройте Telegram → [@BotFather](https://t.me/BotFather).
2. `/mybots` → выберите `@Tikhvinpalomnik_bot`.
3. **Bot Settings** → **Domain** → `/setdomain` (или пункт Domain).
4. Укажите домен **без** схемы и пути:  
   `tikhvin-palomnik.ru`
5. Токен тот же, что в проде. Опционально (если хотите явно):

```bash
TELEGRAM_LOGIN_BOT_TOKEN=   # пусто = берётся TELEGRAM_BOT_TOKEN
TELEGRAM_LOGIN_BOT_USERNAME=Tikhvinpalomnik_bot
```

6. Виджет/кнопка бьёт в обработчик:  
   `https://tikhvin-palomnik.ru/api/auth/social/telegram`  
   (это путь нашего фронта, не webhook уведомлений).

Webhook уведомлений (`/api/v1/webhooks/telegram` на api.*) **не трогайте** — это другой контур.

---

## 4. Max

В коде **нет** URL authorize/token/userinfo по умолчанию: без пяти переменных кнопка останется «недоступно». Токен бота `MAX_BOT_TOKEN` с [business.max.ru](https://business.max.ru) — **другой** секрет (чат/публикации), не вход на сайт.

1. Консоль OAuth/Login Max — тот URL, который даёт Max партнёрам на момент регистрации.
2. Приложение типа **веб / OAuth для сайта**, если кабинет это умеет.
3. **Redirect URI** (точно, как в репо):  
   `https://tikhvin-palomnik.ru/api/auth/social/max/callback`
4. Скопируйте Client ID / Secret.
5. Если консоль показывает отдельные URL authorize / token / userinfo — впишите их явно:

```bash
MAX_OAUTH_CLIENT_ID=...
MAX_OAUTH_CLIENT_SECRET=...
MAX_OAUTH_AUTHORIZE_URL=   # если ещё не в коде по умолчанию
MAX_OAUTH_TOKEN_URL=
MAX_OAUTH_USERINFO_URL=
```

6. Если в кабинете Max пока нет готовых endpoint’ов — оставьте три URL пустыми и напишите разработчику: без них кнопка останется «недоступно». Redirect URI всё равно зарегистрируйте заранее.

---

## Что не нужно

| Не нужно | Почему |
|----------|--------|
| Google OAuth (`GOOGLE_CLIENT_*`) | Кнопка убрана с `/account/login` и `/account/register` (I.5). |
| Отдельный бот только для Login | Достаточно `@Tikhvinpalomnik_bot` + `/setdomain`. |
| Redirect на `api.tikhvin-palomnik.ru` | Соцвход идёт через Next на основном домене. |
| Коммит `.env.production` | Секреты только на сервере / локально, gitignored. |

---

## После заполнения env

1. Проверьте, что compose прокидывает переменные в `frontend` (и при необходимости в `api`).
2. Скажите агенту сделать `make deploy` (сами не деплоим без просьбы).
3. Smoke: кнопка провайдера → редирект → кабинет `/account/trips`. Если уже вошли, привязка с `/account` сливает входы в текущий кабинет (`/account?linked=1`).
