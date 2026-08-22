# SEO и реклама (паломник)

Канон: **https://tikhvin-palomnik.ru** (старый `tikhvin-polomnik.ru` только 301).  
API: `https://api.tikhvin-palomnik.ru`.

## 1. Счётчики на сайте

В `.env.production` (и пересборка фронта):

| Переменная | Зачем |
|------------|--------|
| `NEXT_PUBLIC_YM_ID` | ID счётчика Яндекс.Метрики |
| `NEXT_PUBLIC_GA_ID` | ID Google Analytics (G-…) |
| `NEXT_PUBLIC_YM_WEBVISOR` | `1` — включить вебвизор (по умолчанию выкл.) |
| `NEXT_PUBLIC_YM_CLICKMAP` | `0` — выключить карту кликов (по умолчанию вкл. при Метрике) |
| `NEXT_PUBLIC_SITE_URL` | `https://tikhvin-palomnik.ru` |

Пустые ID — скрипты не подключаются, сайт не ломается.

Цели (reachGoal / GA events): `tour_view`, `begin_checkout`, `booking_submit`, `support_contact`.

## 2. Поисковики

1. [Яндекс.Вебмастер](https://webmaster.yandex.ru/) — добавить `tikhvin-palomnik.ru`, подтвердить, указать sitemap: `https://tikhvin-palomnik.ru/sitemap.xml`.
2. [Google Search Console](https://search.google.com/search-console) — то же.
3. Проверить `https://tikhvin-palomnik.ru/robots.txt`.

## 3. Яндекс.Директ / VK Ads

В объявлениях:

- **URL сайта:** `https://tikhvin-palomnik.ru`
- **Целевые страницы:** `/` · `/search` · `/tours/{id}` · `/support`
- **Счётчик Метрики:** тот же `NEXT_PUBLIC_YM_ID`
- UTM-пример:  
  `https://tikhvin-palomnik.ru/search?utm_source=yandex&utm_medium=cpc&utm_campaign=pilgrimage`

Не указывайте в рекламе выдуманные цены и отзывы — только то, что на сайте.

## 4. После деплоя (чеклист)

- [ ] `NEXT_PUBLIC_YM_ID` / при необходимости `NEXT_PUBLIC_GA_ID` в prod
- [ ] В Метрике видны `tour_view` / `booking_submit`
- [ ] Sitemap в Вебмастере и Search Console
- [ ] OG-превью: ссылка на главную в мессенджере
