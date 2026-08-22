# Релиз: local → GitHub → palomnik

Публичный сайт: **https://tikhvin-palomnik.ru**  
API: **https://api.tikhvin-palomnik.ru**

`tikhvin-polomnik.ru` и `www` редиректят 301 на palomnik. Отдельного preview-окружения нет.

## Шаги

1. Локально: `docker compose up --build -d`, тесты `go test` / `npm run lint`.
2. GitHub: https://github.com/ezhigval/POLOMNIK_47 — commit/push только по просьбе владельца.
3. Прод:

```bash
make deploy
```

SSH: `ssh polomnik-yc` (`smailikin70@93.77.165.81`, `/opt/polomnik`).

## После деплоя

- [ ] https://tikhvin-palomnik.ru — HTTPS 200
- [ ] https://www.tikhvin-palomnik.ru → 301 на apex
- [ ] `/robots.txt` указывает на sitemap palomnik
- [ ] `/management/login` работает
- [ ] `https://api.tikhvin-palomnik.ru/health/ready` (нужна A-запись `api` на REG.RU)

## Cron на ВМ

- `0 3 * * *` — `scripts/backup-postgres.sh` (контейнер `polomnik-postgres-1`)
- `*/10 * * * *` — health API через `docker exec`
