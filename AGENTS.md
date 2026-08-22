# Правила для агентов

Канон для Cursor, Claude, Codex и любых других агентов. Релиз: [docs/RELEASE.md](docs/RELEASE.md). Архитектура: [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md). Кодстайл: [docs/CONVENTIONS.md](docs/CONVENTIONS.md).

## Прод

Единственный публичный сайт: **https://tikhvin-palomnik.ru**  
API: **https://api.tikhvin-palomnik.ru**  
Старый домен `tikhvin-polomnik.ru` только редиректит на palomnik (301). Отдельного preview-стенда нет.

## Конвейер

```text
правки локально → проверки → GitHub (commit/push только если владелец попросил) → деплой на palomnik
```

```bash
make deploy
```

Не деплоить на сервер без просьбы владельца. Не коммитить секреты. Не force-push в `main`.

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
