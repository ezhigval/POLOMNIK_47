# Релиз: local → GitHub → palomnik

Публичный сайт: **https://tikhvin-palomnik.ru**  
API: **https://api.tikhvin-palomnik.ru**

Отдельного preview-окружения нет.

Актуальный статус и теги: [STATUS.md](STATUS.md).

## Шаги

`main` соответствует тому, что на проде. Все правки — в отдельной ветке, в `main` только через PR. Одно существенное изменение — один PR (удобный git и откат).

1. Ветка от `main`: правка, `docker compose up --build -d`, тесты `go test` / `npm run lint`.
2. GitHub: https://github.com/ezhigval/POLOMNIK_47 — PR в `main`, мерж после проверки.
3. Прод (пока нет v4 этапа 12 — автодеплой с `main`):

```bash
make deploy
```

После этапа 12: зелёный CI на `main` сам выкладывает; два деплоя сразу не гонять. Миграции только `goose up`, без `compose down -v`.

SSH: `ssh smailikin70@93.77.165.81` (каталог `/opt/palomnik`, либо `DEPLOY_DIR`).  
Не `compose down -v` — не сбрасывать Postgres.

## Теги

Промежуточные: `v1.n.m`. Линейка продукта: `v2.0.0`, `v2.0.1`, `v2.1.0`, `v3.0.0`, …  
Откат: checkout тега → `make deploy` (миграции только аддитивные). Не двигать уже выпущенные теги.

## После деплоя

- [x] https://tikhvin-palomnik.ru — HTTPS 200 (сверка freeze v3.0.0, 2026-08-26)
- [x] `/robots.txt` и `sitemap.xml`
- [x] `https://api.tikhvin-palomnik.ru/health/ready`
- [ ] `/management/login` — проверить после hard-refresh
- [ ] Секреты из [V4_OWNER_SETUP.md](V4_OWNER_SETUP.md) — по мере готовности кабинетов

## Cron на ВМ

- `0 3 * * *` — `scripts/backup-postgres.sh` (контейнер Postgres из `docker compose ps`)
- `*/10 * * * *` — health API через `docker exec`
