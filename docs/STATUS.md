# Статус продукта

Дата сверки: **2026-08-23**.  
Репозиторий: [github.com/ezhigval/POLOMNIK_47](https://github.com/ezhigval/POLOMNIK_47) · ветка `main`.  
Прод: **https://tikhvin-palomnik.ru** · API: **https://api.tikhvin-palomnik.ru**

## Версии

| Тег | Смысл |
|-----|--------|
| `v1.0.0` | Первый публичный запуск |
| `v1.2.x` | Настройки, RBAC, callcheck, SEO, синяя палитра |
| `v1.3.0` / `v2.0.0` | Mailer, OAuth-кнопки, чеклист владельца |
| `v2.0.1` | Hotfix compose + документация |
| `v2.1.0` | Freeze v2 code-complete: поддержка в админке, password reset, security P0/P1 |

Откат: `git checkout <тег>` / деплой с известного тега. Не `compose down -v` — данные Postgres сохранять.

## Что в коде (готово)

- Сайт: туры, заявки, отзывы, новости, CMS главной, кабинет `/account`, поддержка (чат в БД)
- Админка: Главная, новости, туры, заявки, **поддержка** (`/management/support`), отзывы, синхронизация, **Настройки** (сайт, получатели `канал:адрес`, роли; право `manage_support`)
- Восстановление пароля: `/account/forgot-password` → письмо (нужен SMTP) → `/account/reset-password`
- Полный админ = `ADMIN_TOKEN` в env; пароли ролей — хеш в БД; UUID пользователя в кабинете
- Telegram: один бот, webhook, исходящие через Cloudflare Worker; в уведомлении поддержки — ссылка на тред в админке
- Телефон: sms.ru **callcheck** (без ключа — «пока недоступно»)
- Соцвход: Яндекс / VK / Max / Telegram (без env — кнопки недоступны); Google в UI снят
- Mailer SMTP/noop; SEO/Метрика/GA — код готов
- Bitrix24 / 1С — адаптеры есть, live **выключен** (`noop`)
- CMS create/delete page в management API закрыты

**Код v2 готов.** Владельцу остаются только секреты / DNS / контент / Telegram `/start` — см. [V2_OWNER_SETUP.md](V2_OWNER_SETUP.md).

## Что делает владелец (не код)

Единый чеклист: **[V2_OWNER_SETUP.md](V2_OWNER_SETUP.md)**  
OAuth подробно: [OAUTH_SETUP.md](OAUTH_SETUP.md) · реклама: [SEO_ADS.md](SEO_ADS.md) · Telegram: [TELEGRAM_SETUP.md](TELEGRAM_SETUP.md)

## v3 (2026-08-23)

Этапы **0–10** в `main`. На проде после деплоя этого среза: goose **20**, сайт и `/health/ready` **200**.

| Этап | На проде | Заметка |
|------|----------|---------|
| 0–3 | да | платформа, админка, кабинет (identity / профиль / пассажиры) |
| 4 | да | адаптеры paste-key; без ключа `Configured()==false` |
| 5 | да | реплай и команды бота на том же Telegram webhook |
| 6 | да | `smm_posts` (goose 17), `/management/smm`, `/news/[slug]` |
| 7 | да | черновик поддержки, рекомендации, дайджест, watchdog |
| 8 | да | PaymentPort `sber` / `yookassa`; **`PAYMENT_ADAPTER=noop`** |
| 9 | да | каталог как расписание, sticky CTA, кабинет: привязки / пассажиры / заглушка оплаты |
| 10 | этот срез | STATUS, DECISIONS §12, DATA_MODEL, API, ARCHITECTURE, V3_OWNER_SETUP |

**Live на проде:** `NOTIFICATION_ADAPTER=telegram` (один бот, исходящие через Cloudflare Worker).  

**Noop / не live:** Messenger, Publisher, AI, Payment, Bitrix24, 1С. Живой эквайринг не включать: статусов `AWAITING_PAYMENT` / `PAID` в domain нет.

**UX (этап 9), без новых полей тура:** `/search` заголовок «Расписание»; таблица из `date_start` / `date_end` / title / location / price / slots; длительность = `date_end − date_start + 1` день; на главной блок «Ближайшие выезды» (`GET /tours?limit=8`, `ORDER BY date_start ASC`); sticky CTA на мобильном — телефон из настроек и «Чат»; в кабинете счётчики привязок и пассажиров, текст «на сайте не подключена»; кнопок Pay нет.

**Согласия (#22, влил владелец, `768ddf3`):** goose **18–20** — `legal_documents`, `consents`, `user_photos`, `reviews.allow_distribution`. Публичные `/legal`, кабинет `/account/consents` и `/account/photos`, cookie-баннер, админка `/management/legal`. Тексты документов **не сертифицированы юристом**. Реквизиты оператора — placeholders в коде; compose `OPERATOR_*` / `NEXT_PUBLIC_OPERATOR_*` пока не прокидывает.

v4 не начат. План: **[V3_PLAN.md](V3_PLAN.md)**. Ключи: [V3_OWNER_SETUP.md](V3_OWNER_SETUP.md). Юридическая папка: [legal/README.md](legal/README.md).

Без секретов владельца E2E OAuth/SMTP не прогнать — [V2_OWNER_SETUP.md](V2_OWNER_SETUP.md).
