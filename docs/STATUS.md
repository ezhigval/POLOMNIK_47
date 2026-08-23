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

## Дальше: v3 (в работе)

Код линейки v3 начат: этапы **0–1** в `main`; этап **2** — меню по правам, фильтры заявок/CSV, шаблоны ролей (пресеты формы). Этап **3** начат: identity и слияние кабинета (профиль PATCH и пассажиры — ещё нет).  
Этапы 4–10 и v4 не начаты. План: **[V3_PLAN.md](V3_PLAN.md)**. Кратко: [ROADMAP.md](ROADMAP.md).

Без секретов владельца E2E OAuth/SMTP не прогнать — [V2_OWNER_SETUP.md](V2_OWNER_SETUP.md).
