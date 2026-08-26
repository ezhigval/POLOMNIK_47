# Правила для агентов

Канон для Cursor, Claude, Codex и любых других агентов. Релиз: [docs/RELEASE.md](docs/RELEASE.md). Архитектура: [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md). Кодстайл: [docs/CONVENTIONS.md](docs/CONVENTIONS.md). План: [docs/V4_PLAN.md](docs/V4_PLAN.md).

## Прод

Единственный публичный сайт: **https://tikhvin-palomnik.ru**  
API: **https://api.tikhvin-palomnik.ru**  
Отдельного preview-стенда нет.

## Конвейер

```text
правки в отдельной ветке → проверки → PR в main → мерж → деплой на palomnik (пока v4 этап 12 не готов — только если владелец попросил)
```

`main` соответствует продовому деплою. В `main` напрямую не коммитить. Автодеплой с `main` — [V4_PLAN.md](docs/V4_PLAN.md) этап 12 (не включать, пока нет workflow и секретов Actions).

Каждое **существенное** изменение — своя ветка и свой PR (hotfix, багфикс, этап v4 не смешивать). Так проще читать git и откатывать.

```bash
make deploy
```

Не деплоить на сервер без просьбы владельца (до этапа 12). Не два `make deploy` параллельно. Не `compose down -v`. Не коммитить секреты. Не force-push в `main`. Goose на проде только `up`.

## Код

- Не выдумывай бизнес-правила. Неясно — спроси.
- Backend: гексагон, `adapters → application → domain`. Логика не в Fiber-handlers и не во frontend.
- Bitrix24 / 1C live — только по явной просьбе. Telegram на проде уже включён: один бот, получатели в админке «Настройки». Не выключать без просьбы.
- Юридические тексты и контент туров не выдумывать.

## Проверки перед деплоем

```bash
cd backend && go test ./... && go vet ./...
cd frontend && npm run lint && npm run build
```
