# Аудит проекта: персональные данные и согласия

**Дата:** 2026-08-23  
**Репозиторий:** palomnik (Тихвинский путь)  
**Публичный сайт:** https://tikhvin-palomnik.ru  
**API:** https://api.tikhvin-palomnik.ru  
**Статус:** этап 0 — предварительный аудит (код не изменялся)

---

## 1. Архитектура проекта

| Слой | Путь | Технологии |
|------|------|------------|
| Backend | `backend/` | Go 1.25, модуль `palomnik`, гексагон (`adapters → application → domain`) |
| HTTP API | `backend/internal/adapters/http/fiber/` | Fiber v2, JWT, Goose migrations |
| Frontend | `frontend/` | Next.js 16, React 19, TypeScript, Tailwind CSS 4 |
| Инфраструктура | `docker-compose.yml`, `docker-compose.prod.yml`, `deploy/` | PostgreSQL 16, Redis 7, Caddy 2, worker |

Документация: `docs/ARCHITECTURE.md`, `AGENTS.md`, `docs/RELEASE.md`.

---

## 2. Формы сбора персональных данных

### 2.1. Публичные формы

| Форма | Файл(ы) | Собираемые данные | Согласие / юридический UI |
|-------|---------|-------------------|---------------------------|
| **Заявка на тур** | `frontend/src/components/booking-form.tsx` → `POST /api/bookings` | имя, телефон, email (опц.), кол-во человек, комментарий | Текст «Нажимая кнопку, вы соглашаетесь…» + ссылка на `/privacy`. **Нет checkbox, нет фиксации версии** |
| **Регистрация** | `frontend/src/components/auth/register-form.tsx` → `POST /api/auth/register` | имя, телефон, email, пароль; верификация телефона (sms.ru callcheck) | **Согласия нет** |
| **Вход** | `frontend/src/components/auth/login-form.tsx` | login (телефон/email), пароль; OAuth (Яндекс, VK, Max, Telegram) | **Согласия нет** |
| **Восстановление пароля** | `forgot-password-form.tsx`, `reset-password-form.tsx` | email / token + новый пароль | **Согласия нет** |
| **Профиль ЛК** | `frontend/src/components/account/profile-form.tsx` | имя, телефон, email | **Согласия нет** |
| **Пассажиры** | `frontend/src/components/account/passenger-form.tsx` | ФИО, телефон, дата рождения, **паспорт** | **Согласия нет** — чувствительные идентификационные данные |
| **Чат поддержки** | `frontend/src/components/support-chat.tsx` → `POST /api/support` | текст сообщения (требует авторизации) | **Согласия нет** |
| **Отзывы (публично)** | `frontend/src/app/(site)/reviews/page.tsx` | — | **Публичной формы отправки отзыва нет** — только чтение одобренных отзывов |
| **Оплата** | — | — | **Не реализована** («Без оплаты на сайте») |
| **Загрузка фото пользователем** | — | — | **Не реализована** — загрузка только в админке (`POST /api/v1/management/uploads`) |

### 2.2. Административные формы (с PII)

| Форма | Файл | Данные |
|-------|------|--------|
| Создание отзыва | `frontend/src/components/management/create-review-form.tsx` | имя клиента, текст, рейтинг, тур |
| Экспорт заявок CSV | `backend/.../handlers/bookings_csv.go` | имя, телефон, email, комментарий |
| Настройки Telegram | `frontend/src/components/management/telegram-settings-form.tsx` | chat_id, username |

### 2.3. Антиспам (не согласие)

- Honeypot (`HoneypotField`) на большинстве форм
- Backend: `handlers/formguard.go` — `rejectHoneypot()`
- Yandex SmartCaptcha (`CAPTCHA_ADAPTER=smartcaptcha`) — backend готов, **frontend не отправляет captcha_token**

---

## 3. Модели данных и БД

### 3.1. Миграции (16 файлов, `backend/migrations/`)

Основные таблицы с PII:

| Таблица | PII-поля | Назначение |
|---------|----------|------------|
| `users` | email, phone, name, password_hash | Аккаунты |
| `user_identities` | provider, subject | OAuth-привязки |
| `bookings` | name, phone, email, comment | Заявки на туры |
| `passengers` | name, phone, birth_date, passport | Данные пассажиров |
| `reviews` | client_name, text | Отзывы |
| `support_threads` / `support_messages` | body (с user_id) | Поддержка |
| `telegram_recipients`, `telegram_chat_map` | chat_id, username | Маршрутизация уведомлений |
| `site_settings` | contact phone/email, mail_forward_to | Контакты сайта |

### 3.2. Отсутствующие сущности

- Таблицы `consents`, `legal_documents`, `cookie_preferences`
- Модели версионирования юридических документов
- Журнал обращений субъектов ПДн
- Журнал согласий (consent_type, consent_version, accepted_at)

---

## 4. Существующие юридические документы

| Маршрут | Файл | Статус |
|---------|------|--------|
| `/privacy` | `frontend/src/app/(site)/privacy/page.tsx` | Единственная юридическая страница. Статический MVP-текст без версии, без даты, без PDF. Явно помечен как «базовая версия для MVP» |
| `/legal`, `/terms`, `/offer`, cookie policy | — | **Не существуют** |
| CMS-страницы | `/pages/[slug]` через `cms_pages` | Могут хранить произвольный контент, **без связи с версиями согласий** |
| Footer | `frontend/src/lib/site-nav.ts` | Только ссылка «Политика конфиденциальности» → `/privacy` |

**Механизм версионирования юридических документов: отсутствует полностью.**

---

## 5. Где хранятся персональные данные

### 5.1. Первичное хранение (РФ)

| Данные | Место |
|--------|-------|
| Структурированные PII | **PostgreSQL** на VM Yandex Cloud (`docker-compose.prod.yml`, `docs/DEPLOY.md`) |
| Сессии / rate limits | **Redis** на той же VM |
| Загруженные изображения туров | Локальный volume `uploads_data` → `api.tikhvin-palomnik.ru/uploads/` |
| Резервные копии БД | Опционально **Yandex Object Storage** (`S3_*`, `ru-central1`) |

Хостинг: Yandex Cloud Compute, Ubuntu 24.04, IP `93.77.165.81`, каталог `/opt/palomnik`.

### 5.2. Cookies (frontend)

| Cookie | Назначение | Параметры |
|--------|------------|-----------|
| `palomnik_token` | JWT-сессия пользователя | httpOnly, sameSite=lax, secure (prod), 7 дней |
| `palomnik_admin_session` | Сессия админки | SHA256 от `ADMIN_TOKEN` или management JWT |

**Cookie banner и политика cookie: отсутствуют.**

---

## 6. Внешние сервисы, получающие PII

| Сервис | Статус | Какие данные | Файлы / env |
|--------|--------|--------------|-------------|
| **Telegram Bot API** | Prod (настроен) | Полные PII заявок: имя, телефон, email, комментарий | `notification/telegram/messages.go`, `NOTIFICATION_ADAPTER=telegram` |
| **Max messenger** | noop до ключей | Уведомления поддержки/заявок | `notification/max/`, `notification_routing` |
| **Bitrix24 CRM** | noop (код готов) | Контакты + сделки: name, phone, email, comment | `integration/bitrix/crm.go`, `CRM_ADAPTER=bitrix` |
| **1C** | noop (код готов) | Экспорт заявок/контрагентов | `integration/onec/accounting.go`, `ACCOUNTING_ADAPTER=onec` |
| **sms.ru** | noop / опционально | Номера телефонов (callcheck) | `PHONE_ADAPTER=smsru` |
| **SMTP** | noop / опционально | Email (регистрация, сброс пароля) | `MAIL_ADAPTER=smtp` |
| **OAuth: Яндекс** | опционально | Профиль (имя, email, subject id) | `frontend/src/app/api/auth/social/yandex/` |
| **OAuth: VK** | опционально | Профиль | `.../social/vk/` |
| **OAuth: Max** | опционально | Профиль | `.../social/max/` |
| **OAuth: Telegram Login** | опционально | Telegram user id, имя | `.../social/telegram/` |
| **Yandex Metrika** | опционально (`NEXT_PUBLIC_YM_ID`) | Поведенческие данные, cookies | `frontend/src/components/analytics.tsx` |
| **Google Analytics** | опционально (`NEXT_PUBLIC_GA_ID`) | Поведенческие данные, cookies | `frontend/src/components/analytics.tsx` |
| **Cloudflare Worker** | Prod (обход блокировки Telegram на VM) | Прокси Telegram API | `TELEGRAM_API_BASE` в `.env.production.example` |
| **Unsplash** | Статические изображения | PII не передаётся | `frontend/src/lib/tour-cover.ts` |

### 6.1. Платежи

- **Не реализованы.** Только `payment/noop/gateway.go`, порт `PaymentPort`.
- UI: «Без оплаты на сайте».

### 6.2. Туроператоры / перевозчики

- **Прямых интеграций нет.** Данные передаются оператором вручную (телефон, email, мессенджеры) — это нужно отразить в политике, но конкретные контрагенты **неизвестны** (placeholder).

---

## 7. Аналитика и cookie

| Компонент | Файл | Поведение |
|-----------|------|-----------|
| Yandex Metrika | `frontend/src/components/analytics.tsx` | Загружается при `NEXT_PUBLIC_YM_ID`; webvisor **выключен по умолчанию** (`NEXT_PUBLIC_YM_WEBVISOR=1` для включения) |
| Google Analytics | там же | Загружается при `NEXT_PUBLIC_GA_ID` |
| События | `frontend/src/lib/analytics.ts` | tour_view, begin_checkout, booking_submit, support_contact |

**Проблемы:**
- Аналитика загружается **без согласия пользователя** (нет cookie banner)
- Документация (`docs/V2_OWNER_SETUP.md`) рекомендует webvisor «только после согласия», но **код не блокирует** загрузку
- Нет cookie policy

---

## 8. API (релевантные маршруты)

Источник: `backend/internal/adapters/http/fiber/router.go`

### Публичные
```
POST /api/v1/auth/register, /login, /forgot-password, /reset-password
POST /api/v1/auth/phone/start|complete
POST /api/v1/auth/oauth (internal)
POST /api/v1/bookings
GET  /api/v1/pages, /pages/:slug
GET  /api/v1/site-settings
```

### Личный кабинет (`/api/v1/me/*`, JWT)
```
GET/PATCH /me
GET/POST/PATCH/DELETE /me/passengers
GET /me/bookings, /me/favorites, /me/support
POST /me/support/messages
```

### Админка (`/api/v1/management/*`)
Заявки (включая CSV), отзывы, поддержка, настройки, RBAC, загрузки.

**API согласий: отсутствует.**

---

## 9. Админка

Базовый URL: `/management`

| Раздел | PII-доступ |
|--------|------------|
| Заявки | Полный PII, CSV-экспорт |
| Поддержка | Сообщения пользователей |
| Отзывы | Создание/модерация с client_name |
| Настройки | Контакты, Telegram-получатели, RBAC |

**Управление юридическими документами и просмотр согласий: отсутствует.**

---

## 10. Оператор (текущее состояние)

| Поле | Значение в коде |
|------|-----------------|
| Название (публичное) | Под Покровом Божией Матери "Тихвинская" (`NEXT_PUBLIC_SITE_NAME`) |
| Полное название | «РПЦ Тихвинская Епархия. Паломническая служба «Под покровом Божией Матери «Тихвинская»»» |
| Родительская организация | Тихвинская епархия |
| Контакты | `+7 966 933-43-21`, `info@tikhvin-palomnik.ru` |
| Регион | Санкт-Петербург (departure city) |
| ИНН, ОГРН, юр. адрес | **Не указаны в проекте** |

Единого конфигурационного источника реквизитов оператора **нет** — данные разбросаны по `site-config.ts` и env.

---

## 11. Трансграничная передача — предварительная оценка

| Сервис | Риск трансграничной обработки | Требует проверки |
|--------|-------------------------------|------------------|
| PostgreSQL / Redis на Yandex Cloud | Низкий (РФ) | — |
| Yandex Object Storage (ru-central1) | Низкий (РФ) | — |
| Yandex Metrika | Требует уточнения | Да |
| Google Analytics | **Высокий** (если включён) | Да |
| Telegram Bot API | **Высокий** (международная инфраструктура) | Да |
| Cloudflare Worker (прокси Telegram) | **Высокий** | Да |
| Bitrix24 | Зависит от размещения облака клиента | Да (при подключении) |
| 1C | Обычно on-prem РФ | Да (при подключении) |
| sms.ru | Низкий (РФ) | — |
| OAuth Яндекс / VK | Низкий–средний (РФ) | Да |
| OAuth Max | Неизвестно | Да |
| SMTP | Зависит от провайдера | Да |

**Вывод:** при текущей prod-конфигурации (Telegram + опционально GA) возможна трансграничная обработка. Не скрывать в документах; зафиксировать в `LEGAL-REVIEW.md`.

---

## 12. Существующая инфраструктура — что переиспользовать

| Механизм | Применимость для системы согласий |
|----------|-----------------------------------|
| CMS (`cms_pages`) | Может хранить HTML документов, но **нет версионирования и связи с согласиями** — недостаточно |
| `site_settings` | Контакты оператора — частично |
| `site-config.ts` | Публичные названия — частично |
| Goose migrations | Стиль миграций — использовать |
| Hexagonal architecture | Новые domain/application/adapters — следовать |
| RBAC админки | Расширить для legal admin |

**Рекомендация:** создать отдельные таблицы `legal_documents` и `consents` (не дублировать CMS).

---

## 13. Gaps — что нужно реализовать

### Критичные пробелы

1. **5 юридических документов** (политика, согласие на обработку, согласие на распространение, рекламное согласие, cookie policy)
2. **Версионирование документов** с хранением исторических редакций
3. **Checkbox согласия** в регистрации, заявке (обязательные); отзывы, реклама, распространение (добровольные)
4. **Фиксация согласий** в БД (consent_type, document_version, accepted_at, user_id/request_id)
5. **Backend API** для документов и согласий
6. **Cookie banner** с блокировкой аналитики до согласия
7. **Раздел `/legal`** с документами, версиями, скачиванием
8. **Единый конфиг оператора** (реквизиты-placeholder)
9. **Админка:** просмотр версий документов, публикация новых, просмотр согласий
10. **Тесты** всех сценариев

### Формы без согласия (требуют доработки)

- Регистрация
- Заявка (заменить текст «нажимая кнопку» на checkbox)
- Пассажиры (паспортные данные)
- Поддержка
- OAuth-вход (согласие при первой регистрации через соцсеть)

### Не требуют доработки сейчас

- Оплата (не реализована)
- Публичная форма отзывов (не реализована)
- Загрузка фото пользователем (не реализована) — механизм подготовить на будущее

---

## 14. Следующие этапы

| Этап | Содержание |
|------|------------|
| 1 ✓ | Аудит проекта (этот документ) |
| 2 | Аудит законодательства → `LEGAL-REVIEW.md` |
| 3 | Юридическая модель и список документов |
| 4 | Тексты документов |
| 5 | Версионирование и хранение |
| 6 | Миграции БД (`consents`, `legal_documents`) |
| 7 | Backend API |
| 8–14 | Frontend: формы, cookie, ЛК, админка |
| 15 | Тестирование |
| 16 | Финальный аудит |

---

*Документ подготовлен автоматически на основе анализа кодовой базы. Перед публикацией юридических документов требуется проверка оператором/юристом.*
