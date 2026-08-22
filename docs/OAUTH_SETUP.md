# Настройка соцвхода (владелец)

Канон сайта: **https://tikhvin-palomnik.ru**.  
Callback’и — на **сайт** (Next.js), не на `api.tikhvin-palomnik.ru`.

Пути зафиксированы в коде: `frontend/src/lib/auth/social-paths.ts`.

Секреты кладите **только** в gitignored `.env.production` на сервере (или локальный `.env`). Не коммитьте. После заполнения — `make deploy` по вашей просьбе.

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

1. Откройте [https://oauth.yandex.ru/](https://oauth.yandex.ru/) → войдите → **Создать новое приложение** (или «Зарегистрировать новое»).
2. Тип: веб-сервис / доступ к данным пользователя (логин, имя, email — по минимуму).
3. **Callback URI** (точно):  
   `https://tikhvin-palomnik.ru/api/auth/social/yandex/callback`
4. Сохраните приложение. Скопируйте **ClientID** и **Client secret**.
5. В `.env.production`:

```bash
YANDEX_OAUTH_CLIENT_ID=...
YANDEX_OAUTH_CLIENT_SECRET=...
```

6. Пока переменных нет — кнопка Яндекс в UI будет «Пока что недоступно…» (когда кнопка появится в релизе II.3).

---

## 2. VK ID

1. Откройте [https://id.vk.com/about/business/go](https://id.vk.com/about/business/go) (или кабинет VK ID / «Мои приложения»).
2. Создайте приложение типа **Веб-сайт** / VK ID для сайта.
3. Укажите сайт: `https://tikhvin-palomnik.ru`
4. **Redirect URL** (точно):  
   `https://tikhvin-palomnik.ru/api/auth/social/vk/callback`
5. Скопируйте **App ID** (client id) и **защищённый ключ** (secret).
6. В `.env.production`:

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

На момент записи в коде заложены placeholder-URL (консоль Max может менять названия полей).

1. Откройте кабинет разработчика Max (официальная консоль OAuth/Login для партнёров Max — тот URL, который даёт поддержка Max на момент регистрации приложения).
2. Создайте приложение типа **веб / OAuth для сайта**.
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
3. Smoke: кнопка провайдера → редирект → кабинет `/account/trips`.
