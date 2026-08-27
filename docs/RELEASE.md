# Релиз: local → GitHub → palomnik

Публичный сайт: **https://tikhvin-palomnik.ru**  
API: **https://api.tikhvin-palomnik.ru**

Отдельного preview-окружения нет.

Актуальный статус и теги: [STATUS.md](STATUS.md).

## Шаги

`main` соответствует тому, что на проде. Все правки — в отдельной ветке, в `main` только через PR. Одно существенное изменение — один PR (удобный git и откат).

1. Ветка от `main`: правка, `docker compose up --build -d`, тесты `go test` / `npm run lint`.
2. GitHub: https://github.com/ezhigval/POLOMNIK_47 — PR в `main`, мерж после проверки.
3. Прод:

Пока GitHub Actions Deploy красный (`DEPLOY_SSH_KEY`):

```bash
make deploy
```

После исправления секрета: зелёный CI на push в `main` запускает [`.github/workflows/deploy.yml`](../.github/workflows/deploy.yml). Вручную: Actions → **Deploy production**. Нужен **приватный** OpenSSH-ключ; опционально **`DEPLOY_SSH_HOST`**. Миграции: [MIGRATIONS.md](MIGRATIONS.md). Не `compose down -v`.

SSH: `ssh smailikin70@93.77.165.81`, каталог **`/opt/polomnik`**.

## Теги

Промежуточные: `v1.n.m`. Линейка продукта: `v2.0.0`, `v2.0.1`, `v2.1.0`, `v3.0.0`, …  
Откат: checkout тега → `make deploy` (миграции только аддитивные). Не двигать уже выпущенные теги.

## После деплоя

- [x] https://tikhvin-palomnik.ru — HTTPS 200 (сверка 2026-08-27, после #46)
- [x] `/robots.txt` и `sitemap.xml`
- [x] `https://api.tikhvin-palomnik.ru/health/ready`
- [ ] Автодеплой Actions — чинить `DEPLOY_SSH_KEY`
- [ ] Секреты из [V4_OWNER_SETUP.md](V4_OWNER_SETUP.md) — по мере готовности кабинетов

## Cron на ВМ

- `0 3 * * *` — `scripts/backup-postgres.sh` (контейнер Postgres из `docker compose ps`)
- `*/10 * * * *` — health API через `docker exec`
