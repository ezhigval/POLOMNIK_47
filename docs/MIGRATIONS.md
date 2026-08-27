# Goose migrations (Postgres)

Канон: [V4_PLAN.md](V4_PLAN.md) этап 12 · [RELEASE.md](RELEASE.md).

На проде только **`goose up`**. Не `down`, не `reset`, не ручный `DROP`. Named volume Postgres не трогать.

## Нумерация

- Каталог: `backend/migrations/`
- Имя: `NNNNN_short_name.sql` — следующий свободный номер (на проде **00032**; следующий — **00033**).
- **Не править** уже применённые файлы на проде. Правка данных — новая миграция (как 00026 «духовник»→«священник»).

CI: `scripts/check-migrations.sh` — без пропусков в нумерации; INSERT/UPDATE с `;` в теле должны быть в `StatementBegin`/`StatementEnd` (урок 00025).

## StatementBegin

Если в INSERT/UPDATE внутри строк есть точка с запятой (HTML, `$article$`, списки), оборачивайте SQL:

```sql
-- +goose StatementBegin
INSERT INTO news_articles (...) VALUES (...);
-- +goose StatementEnd
```

Пример: `00025_news_ikona_v_moskvu.sql`.

## После миграций контента

Если меняются публичные туры/новости в БД — сбросить Redis `tours:namespace` / `news:namespace` (иначе главная до ~5 минут отдаёт кэш).

## Локально

```bash
docker compose up migrate
```
