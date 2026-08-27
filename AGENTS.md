# Правила для агентов

Канон для Cursor, Claude, Codex и любых других агентов. Релиз: [docs/RELEASE.md](docs/RELEASE.md). Архитектура: [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md). Кодстайл: [docs/CONVENTIONS.md](docs/CONVENTIONS.md). План: [docs/V4_PLAN.md](docs/V4_PLAN.md). Статус: [docs/STATUS.md](docs/STATUS.md).

## Прод

Единственный публичный сайт: **https://tikhvin-palomnik.ru**  
API: **https://api.tikhvin-palomnik.ru**  
Отдельного preview-стенда нет.

## Конвейер

```text
ветка cursor/<имя>-… → проверки → PR в main → мерж
  → GitHub Actions Deploy (сейчас падает: SSH Permission denied)
  → запасной путь: make deploy
```

`main` = то, что должно быть на проде. В `main` напрямую не коммитить. Одно существенное изменение — свой PR (hotfix и фичи не смешивать). Не мержить черновики **#25** и **#27**.

```bash
make deploy
```

Не два `make deploy` параллельно. Не `compose down -v`. Не коммитить секреты. Не force-push в `main`. Goose на проде только `up`. Пока `DEPLOY_SSH_KEY` не совпадает с ключом на ВМ, автодеплой из Actions не считать релизом.

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
