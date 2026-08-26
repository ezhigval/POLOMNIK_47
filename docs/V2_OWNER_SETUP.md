# Чеклист владельца к v2.0

Один файл: что сделать **вам** после выкладки кода. Кодить не нужно — только кабинеты, DNS/env и проверки.

Сайт: **https://tikhvin-palomnik.ru**  
API: **https://api.tikhvin-palomnik.ru**  
Секреты только в gitignored `.env.production` на сервере. После правок env — попросить агента/`make deploy` (без `compose down -v`).

Подробности OAuth также в [OAUTH_SETUP.md](./OAUTH_SETUP.md). SEO/реклама — [SEO_ADS.md](./SEO_ADS.md).

---

## 0. Что уже в коде (не ждёт вашей разработки)

**Код v2 готов.** Ниже — только ваши кабинеты, DNS и env.

- Админка: настройки сайта, получатели уведомлений (канал+адрес), роли (хеш пароля, назначение по UUID из `/account`), экран поддержки `/management/support`
- Полный админ = только `ADMIN_TOKEN` из env; право `manage_support` для менеджеров
- Телефон: sms.ru **callcheck** (не SMS)
- Соцвход: Яндекс / VK / Max / Telegram Login — кнопки; без env → «Пока что недоступно…»
- Google в UI снят
- Mailer: порт + SMTP; без `MAIL_ADAPTER=smtp` регистрация жива, письмо не шлётся; восстановление пароля показывает «пока недоступно»
- Max/WhatsApp-уведомления: адаптер Max — noop до ключей; WhatsApp не подключён live
- SEO/Метрика/GA/CSP/sitemap — код готов, нужны ваши ID
- Синие цвета бренда
- Telegram: deep link на тред поддержки в админке

---

## 1. OAuth (Яндекс, VK, Max, Telegram)

Redirect URI — на **сайт**, не на `api.*`. Пути в коде: `frontend/src/lib/auth/social-paths.ts`.

| Провайдер | Redirect / домен | Env |
|-----------|------------------|-----|
| Яндекс ID | `https://tikhvin-palomnik.ru/api/auth/social/yandex/callback` | `YANDEX_OAUTH_CLIENT_ID`, `YANDEX_OAUTH_CLIENT_SECRET` |
| VK | `https://tikhvin-palomnik.ru/api/auth/social/vk/callback` | `VK_OAUTH_CLIENT_ID`, `VK_OAUTH_CLIENT_SECRET` |
| Max | `https://tikhvin-palomnik.ru/api/auth/social/max/callback` | `MAX_OAUTH_CLIENT_ID`, `MAX_OAUTH_CLIENT_SECRET`, `MAX_OAUTH_AUTHORIZE_URL`, `MAX_OAUTH_TOKEN_URL`, `MAX_OAUTH_USERINFO_URL` |
| Telegram Login | домен виджета: `tikhvin-palomnik.ru`; handler `…/api/auth/social/telegram` | тот же бот: `TELEGRAM_BOT_TOKEN` (+ опц. `TELEGRAM_LOGIN_BOT_USERNAME`) |

### Шаги кратко

1. **Яндекс:** [oauth.yandex.ru](https://oauth.yandex.ru/) → приложение → callback как в таблице → ClientID/secret в env.
2. **VK ID:** кабинет VK ID → веб-сайт → redirect как в таблице → App ID + защищённый ключ.
3. **Max:** кабинет Max для разработчиков → URL authorize/token/userinfo (когда выдадут) → env.
4. **Telegram:** BotFather → тот же `@Tikhvinpalomnik_bot` → Domain = `tikhvin-palomnik.ru` (без https).
5. На сервер: прописать env → `make deploy`.
6. Проверка: `/account/login` — доступные кнопки активны; без env — disabled с текстом «пока недоступно».

Также нужен `INTERNAL_API_SECRET` на **frontend** (Next) — уже должен быть для OAuth → backend `/auth/oauth`.

---

## 2. Телефон (sms.ru callcheck)

1. Кабинет [sms.ru](https://sms.ru/) → API ID.
2. В `.env.production`:

```bash
PHONE_ADAPTER=smsru
SMSRU_API_ID=...
```

3. Деплой. В формах входа/регистрации появится звонок-проверка (не SMS OTP).
4. Без ключа: «Пока что недоступно, используйте другой вариант.»

---

## 3. Почта (Яндекс 360 / SMTP)

Сейчас MX у вас: **mail.tikhvin-palomnik.ru** — оставьте как решите с Яндекс 360 / reg.ru (MX, SPF, DKIM).

1. Ящик исходящих, например `info@tikhvin-palomnik.ru`.
2. В `.env.production`:

```bash
MAIL_ADAPTER=smtp
SMTP_HOST=smtp.yandex.ru
SMTP_PORT=587
SMTP_USERNAME=info@tikhvin-palomnik.ru
SMTP_PASSWORD=...
SMTP_FROM=info@tikhvin-palomnik.ru
# запасной список пересылки (если не задан в админке «Настройки»):
MAIL_FORWARD_TO=smailikin70@yandex.ru
```

3. В админке **Настройки** → поле «Пересылка почты» (не-секретные адреса).
4. Входящая пересылка на личные ящики обычно настраивается **в Яндекс 360 / панели DNS**, не в коде сайта.
5. Без SMTP: регистрация работает, письмо подтверждения не уходит; «Забыли пароль?» показывает «пока недоступно» (это норма).
6. С SMTP: `/account/forgot-password` шлёт ссылку на `/account/reset-password`.

---

## 4. Метрика / GA / Вебмастер / реклама

См. [SEO_ADS.md](./SEO_ADS.md).

```bash
NEXT_PUBLIC_YM_ID=...
NEXT_PUBLIC_GA_ID=...          # опционально
NEXT_PUBLIC_YM_WEBVISOR=1      # только после согласия на запись сессий
# Вебмастер/GSC: коды уже в фронте; env нужен только чтобы сменить
# NEXT_PUBLIC_YANDEX_VERIFICATION=...
# NEXT_PUBLIC_GOOGLE_SITE_VERIFICATION=...
```

Дальше: после деплоя — «Проверить» в Яндекс.Вебмастере и Google Search Console, sitemap `https://tikhvin-palomnik.ru/sitemap.xml`.

---

## 5. Telegram уведомления

1. Получатели в **Настройки** (канал telegram + username).
2. Каждый из списка один раз пишет боту `/start`.
3. Webhook: если в BotFather/getWebhookInfo ошибка — проверить `TELEGRAM_WEBHOOK_URL` и Worker `TELEGRAM_API_BASE`.
4. Smoke: новая заявка → сообщение; сообщение в поддержку → сообщение со ссылкой на `/management/support/{id}`.

---

## 6. Роли админки

1. Войти полным админом: **пустая роль** + `ADMIN_TOKEN`.
2. Создать роли (manager и т.п.), пароли ≥ 8 символов, права галочками.
3. Пользователь смотрит UUID в `/account` (Профиль) → вставить в «Назначить» у роли.
4. Вход менеджера: имя роли + пароль роли (не ADMIN_TOKEN).

---

## 7. Что вам **не** нужно кодить

- Адаптеры Bitrix/1С/оплата — **не включать live** без отдельной просьбы.
- Юридические тексты, цены туров, тексты новостей — только ваши материалы.
- Счётчики «визитов» в админке — не выдумываем; смотрите Метрику.
- Экран поддержки в админке с ответами в чат — **уточнить продукт** (сейчас тред в БД + пуш в Telegram). Не делаем фантазийный UI без правил.
- WhatsApp-уведомления — ждать вашего решения и credentials.

---

## 8. После заполнения секретов

1. Обновить `.env.production` на ВМ (не коммитить).
2. `make deploy` (Postgres volume не трогать).
3. Hard-refresh сайта (Cmd+Shift+R).
4. Проверить: главная 200, `/management`, заявка → Telegram, вход по телефону/соцсети если включили.
