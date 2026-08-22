# Project Specification

Актуальное состояние продукта: [DECISIONS.md](DECISIONS.md) §10 и [ROADMAP.md](ROADMAP.md). Ниже — исходный набросок мая 2026. Кабинет пользователя и Telegram уже есть; Bitrix24 и 1С по-прежнему не live.

## 1. Цель продукта

PALOMNIK 47 — платформа паломнической службы для туроператора.

Целевой процесс:

```text
Пользователь смотрит туры
  -> выбирает тур
  -> отправляет заявку
  -> менеджер обрабатывает заявку
  -> в следующих версиях данные синхронизируются с Bitrix24 и 1C
```

Изначальная документация рассматривала Bitrix24 как источник правды. После уточнения стратегии первый MVP строится без активного Bitrix24. Backend должен иметь собственную доменную логику и адаптеры хранения, а Bitrix24 и 1C подключаются позже через порты.

## 2. Глобальные этапы

### Этап A - Backend Logic

Создать backend, который умеет работать самостоятельно:

- доменные сущности;
- бизнес-правила;
- API;
- PostgreSQL storage;
- Redis cache, если нужен;
- management API;
- тесты;
- порты под будущие интеграции.

### Этап B - Frontend Logic

Создать сайт поверх backend API:

- список туров;
- фильтры;
- страница тура;
- отзывы;
- форма заявки;
- базовые management-интерфейсы, если они нужны для MVP;
- обработка loading/empty/error/success states.

### Этап C - Integrations

Подключить внешние системы:

- Bitrix24 для CRM-сценариев;
- 1C для учетных/операционных сценариев;
- синхронизацию, webhooks, import/export;
- переключение источников данных через adapters.

## 3. Пользователи и роли

MVP 1:

- Public visitor - смотрит туры и отзывы, отправляет заявку.
- Operator/manager - обрабатывает заявки через management API или внутренний интерфейс.

Не входит в MVP 1:

- client account;
- user login;
- JWT;
- полноценные роли;
- история поездок;
- payments.

## 4. Основные сущности

- Tour - тур.
- Tour image - изображения тура.
- Booking - заявка/бронирование.
- Review - отзыв.
- Client contact - контактные данные из заявки.
- Integration reference - связь локальной записи с внешней системой, появится при подключении Bitrix24/1C.

`User` как аккаунт клиента не является обязательной сущностью MVP 1. Контактные данные заявки хранятся в booking/contact модели.

## 5. MVP 1 scope

MVP 1 включает:

- backend skeleton на Go + Fiber;
- гексагональную архитектуру;
- domain layer без зависимостей от Fiber/PostgreSQL/Bitrix24/1C;
- PostgreSQL adapter;
- fake/in-memory adapter для тестов;
- optional Redis adapter;
- public tours API;
- public reviews API;
- booking API;
- management API для наполнения и обработки данных;
- единый error format;
- валидацию;
- request id;
- structured logging;
- Docker Compose;
- тесты business logic.

## 6. MVP 1 не включает

- Bitrix24 как активный источник данных;
- 1C как активный источник данных;
- online payments;
- user auth;
- личный кабинет;
- мобильное приложение;
- email/SMS/Telegram notifications;
- production-grade admin panel.

## 7. Booking lifecycle в MVP 1

Минимальный жизненный цикл заявки:

```text
NEW
  -> CONTACTED
  -> CONFIRMED
  -> COMPLETED
```

Отмена возможна из любого активного состояния:

```text
CANCELLED
```

Дополнительные будущие статусы:

```text
AWAITING_PAYMENT
PAID
IN_TRIP
SYNC_PENDING
SYNC_FAILED
```

Будущие статусы не используются в MVP 1 без отдельной задачи.

## 8. Booking rules

При создании заявки backend должен:

1. Проверить request.
2. Проверить существование активного тура.
3. Проверить даты тура.
4. Проверить количество мест.
5. Рассчитать total price.
6. Создать booking в локальном хранилище.
7. Уменьшить доступные места или зарезервировать их по выбранной модели.
8. Вернуть booking id и текущий статус.

Overbooking:

- управляется настройкой тура `overbooking_enabled`;
- не передается пользователем в booking request;
- если мест достаточно, booking разрешен;
- если мест недостаточно и `overbooking_enabled = true`, booking переводится в `NEW` с признаком `overbooked = true` или в отдельный future status после уточнения;
- если мест недостаточно и overbooking запрещен, API возвращает ошибку.

Если точное поведение overbooking нужно поменять, агент должен спросить.

## 9. Future integration principle

Backend не должен быть написан "под PostgreSQL" или "под Bitrix24". Он должен быть написан под domain и ports.

Пример:

```text
BookingService
  -> BookingRepository port
  -> TourRepository port
  -> ExternalCRMPort
  -> AccountingPort
```

В MVP 1:

```text
BookingRepository -> PostgreSQL
ExternalCRMPort -> Noop/Fake
AccountingPort -> Noop/Fake
```

В integration stage:

```text
ExternalCRMPort -> Bitrix24
AccountingPort -> 1C
```

## 10. Принцип "не угадывать"

Агент обязан спрашивать, если неизвестны:

- точные поля Bitrix24;
- 1C exchange format;
- платежные правила;
- auth flow;
- учетные статусы;
- production admin security;
- дизайн frontend;
- любые бизнес-правила, которые меняют поведение заявки.

