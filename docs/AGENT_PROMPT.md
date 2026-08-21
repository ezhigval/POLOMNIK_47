# Agent Prompt

Ты агент разработки проекта POLOMNIK_47.

Работай строго по документации из `docs/`.

## 1. Обязательные решения

- Backend: Go + Fiber.
- Go module: `polomnik`.
- Repo layout: `docs/`, `backend/`, `frontend/`.
- Backend architecture: hexagonal architecture.
- MVP 1: backend logic without active Bitrix24 and 1C.
- MVP 1: no user auth, no личный кабинет, no online payments.
- Integrations with Bitrix24 and 1C are Stage C.

## 2. Перед работой прочитай

1. `docs/DECISIONS.md`.
2. `docs/PROJECT_SPEC.md`.
3. `docs/ARCHITECTURE.md`.
4. `docs/API.md`.
5. `docs/DATA_MODEL.md`.
6. `docs/ROADMAP.md`.
7. `docs/INTEGRATIONS.md`.
8. `docs/INTEGRATION_DISCOVERY.md` — перед Stage C.

Если документ противоречит другому документу в `docs/`, остановись и спроси.

## 3. Главные правила

1. Не придумывай бизнес-логику.
2. Не подключай Bitrix24 в MVP 1.
3. Не подключай 1C в MVP 1.
4. Не добавляй auth/JWT в MVP 1.
5. Не добавляй payments в MVP 1.
6. Не размещай business logic во frontend.
7. Не размещай business logic в HTTP handlers.
8. Не размещай business logic в PostgreSQL adapter.
9. Не размещай business logic в Bitrix24/1C adapters.
10. Domain and application должны работать через ports.

## 4. Как действовать

Иди по `docs/ROADMAP.md`.

Порядок:

1. Stage A - Backend Logic.
2. Stage B - Frontend Logic.
3. Stage C - Integrations.

Не переходи к Stage C, пока владелец проекта явно не попросит.

## 5. Если неясно

Спроси владельца проекта, если неизвестно:

- business behavior;
- status transition;
- public API contract;
- management API protection;
- frontend design;
- Bitrix24 fields;
- Bitrix24 entity model;
- 1C exchange format;
- payment rules;
- auth flow;
- source of truth после интеграций.

Не угадывай.

## 6. Backend rules

Backend должен быть расположен в `backend/`.

Используй структуру из `docs/ARCHITECTURE.md`.

Минимальные проверки:

```bash
cd backend
go test ./...
go vet ./...
```

Fiber является HTTP adapter, а не центром приложения.

## 7. Frontend rules

Frontend должен быть расположен в `frontend/`.

Frontend:

- вызывает только backend API;
- не обращается к Bitrix24;
- не обращается к 1C;
- делает только UX validation;
- не содержит booking business rules.

## 8. Integration rules

В MVP 1 реализуй только noop/fake adapters for:

- `CRMPort`;
- `AccountingPort`.

Real Bitrix24 and 1C adapters делать только на Stage C после уточнения внешних моделей.

## 9. Отчет после задачи

После задачи сообщи:

- что сделано;
- какие файлы изменены;
- какие проверки запущены;
- какие проверки прошли или упали;
- какие вопросы остались;
- следующий логичный шаг.

