# Релиз: local → GitHub → palomnik

Публичный сайт: **https://tikhvin-palomnik.ru**  
API: **https://api.tikhvin-palomnik.ru**

`tikhvin-polomnik.ru` и `www` редиректят 301 на palomnik. Отдельного preview-окружения нет.

Актуальный статус и теги: [STATUS.md](STATUS.md).

## Шаги

1. Локально: `docker compose up --build -d`, тесты `go test` / `npm run lint`.
2. GitHub: https://github.com/ezhigval/POLOMNIK_47 — commit/push только по просьбе владельца.
3. Прод:

```bash
make deploy
```

SSH: `ssh polomnik-yc` (`smailikin70@93.77.165.81`, `/opt/polomnik`).  
Не `compose down -v` — не сбрасывать Postgres.

## Теги

Промежуточные: `v1.n.m`. Линейка продукта: `v2.0.0`, `v2.0.1`, …  
Откат: checkout тега → `make deploy` (миграции только аддитивные).

## После деплоя

- [x] https://tikhvin-palomnik.ru — HTTPS 200 (сверка 2026-08-23)
- [x] `/robots.txt` и `sitemap.xml`
- [x] `https://api.tikhvin-palomnik.ru/health/ready`
- [ ] `/management/login` — проверить после hard-refresh
- [ ] Секреты из [V2_OWNER_SETUP.md](V2_OWNER_SETUP.md) — по мере готовности кабинетов

## Cron на ВМ

- `0 3 * * *` — `scripts/backup-postgres.sh` (контейнер `polomnik-postgres-1`)
- `*/10 * * * *` — health API через `docker exec`
