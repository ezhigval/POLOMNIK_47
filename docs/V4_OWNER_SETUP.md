# Чеклист владельца к v4

Сайт: **https://tikhvin-palomnik.ru**  
API: **https://api.tikhvin-palomnik.ru**  
План: [V4_PLAN.md](V4_PLAN.md). Секреты только в gitignored `.env.production` на ВМ. После правок env — `make deploy` (без `compose down -v`).

v2-секреты: [V2_OWNER_SETUP.md](V2_OWNER_SETUP.md). v3-адаптеры: [V3_OWNER_SETUP.md](V3_OWNER_SETUP.md).  
OAuth: [OAUTH_SETUP.md](OAUTH_SETUP.md) · Telegram: [TELEGRAM_SETUP.md](TELEGRAM_SETUP.md) · SEO: [SEO_ADS.md](SEO_ADS.md).

Код не заменяет клики в кабинетах поисковиков и не выдумывает ИНН, цены и статусы оплаты.

---

## Где что живёт

| Что | Куда |
|-----|------|
| Секреты (токены, пароли, ключи API) | `.env.production` (имена ниже, **не** значения в git) |
| Имя сайта, телефон, почта, пересылка, получатели уведомлений, роли | Админка **Настройки** |
| Счётчик Метрики | только env `NEXT_PUBLIC_YM_ID` (на проде **111985266**). Второго поля в админке нет |
| Реквизиты оператора ПДн | env `OPERATOR_*` / `NEXT_PUBLIC_OPERATOR_*` (пустые = плейсхолдеры `название` / «—»). **Не выдумывать** ИНН/ОГРН |
| Тексты туров, новости, юридичка | Админка / материалы владельца |

Compose с v4 этапа 1 прокидывает `OPERATOR_*` в API и `NEXT_PUBLIC_OPERATOR_*` в сборку фронта. Пустая переменная = как раньше, плейсхолдер из кода.

---

## 1. Секреты — имена переменных (без значений)

### Уже нужны / уже на проде

| Имя | Зачем |
|-----|--------|
| `ADMIN_TOKEN` | Полный админ |
| `JWT_SECRET` | Сессии кабинета |
| `INTERNAL_API_SECRET` | Next → API (OAuth) |
| `POSTGRES_PASSWORD` | БД |
| `TELEGRAM_BOT_TOKEN` | Один бот: уведомления, webhook, login |
| `TELEGRAM_API_BASE` | Cloudflare Worker (с ВМ нельзя `api.telegram.org`) |
| `NOTIFICATION_ADAPTER=telegram` | Не выключать |
| `NEXT_PUBLIC_YM_ID` | Метрика; прод 111985266 (публичный ID) |
| `NEXT_PUBLIC_YM_WEBVISOR` | `1` на проде |

### Вписать, когда кабинеты готовы (если ещё пусто)

| Имя | Документ |
|-----|----------|
| `YANDEX_OAUTH_CLIENT_ID` / `YANDEX_OAUTH_CLIENT_SECRET` | [OAUTH_SETUP.md](OAUTH_SETUP.md) |
| `VK_OAUTH_CLIENT_ID` / `VK_OAUTH_CLIENT_SECRET` | то же |
| `MAX_OAUTH_*` | то же; URL authorize/token/userinfo когда выдадут |
| `MAIL_ADAPTER=smtp`, `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `SMTP_FROM` | [V2_OWNER_SETUP.md](V2_OWNER_SETUP.md) |
| `PHONE_ADAPTER=smsru`, `SMSRU_API_ID` | callcheck, не SMS |
| `NEXT_PUBLIC_GA_ID` | опционально |
| `NEXT_PUBLIC_YANDEX_VERIFICATION` / `NEXT_PUBLIC_GOOGLE_SITE_VERIFICATION` | сменить коды Вебмастера/GSC |

Несекретные адреса пересылки: админка «Настройки» (`mail_forward_to`) или запасной `MAIL_FORWARD_TO` в env.

### GitHub Actions — только этап 12 (автодеплой), не в git

Пока workflow нет — деплой `make deploy`. Когда этап 12 включат, в **Settings → Secrets** репозитория (значения не сюда):

| Имя (предложение) | Зачем |
|-------------------|--------|
| `DEPLOY_SSH_KEY` | Приватный ключ на ВМ `smailikin70@93.77.165.81` |
| `DEPLOY_SSH_HOST` | если не зашивать IP в workflow |

Не класть в Actions весь `.env.production`. Правка только env на ВМ по-прежнему требует пересборки фронта/`make deploy` (NEXT_PUBLIC_* в образе).

### v3-адаптеры — оставлять noop, пока не попросите live

| Имя | По умолчанию |
|-----|----------------|
| `MESSENGER_ADAPTER` | `noop` |
| `PUBLISHER_ADAPTER` | `noop` |
| `AI_ADAPTER` | `noop` |
| `PAYMENT_ADAPTER` | `noop` — **не** включать до решения по `AWAITING_PAYMENT` / `PAID` |
| `CAPTCHA_ADAPTER` | `noop` |
| `BACKUP_STORAGE_ADAPTER` | `noop` |
| `CRM_ADAPTER` / `ACCOUNTING_ADAPTER` | `noop` |

Ключи к ним (имена): `WHATSAPP_TOKEN`, `WHATSAPP_PHONE_NUMBER_ID`, `TELEGRAM_CHANNEL_ID`, `VK_WALL_TOKEN`, `VK_WALL_OWNER_ID`, `MAX_BOT_TOKEN`, `MAX_FEED_CHAT_ID`, `YANDEXGPT_API_KEY`, `YANDEXGPT_FOLDER_ID`, `SBER_USERNAME`, `SBER_PASSWORD`, `YOOKASSA_SHOP_ID`, `YOOKASSA_SECRET_KEY`, `SMARTCAPTCHA_*`, `S3_*`. Значения не печатать и не коммитить.

Cloudflare token воркера — только на ВМ / в кабинете Cloudflare, не в git.

---

## 2. Админка «Настройки» (уже в коде)

- Короткое / полное имя, слоган, описание, регион, город выезда, родительская орг., телефон, email, пересылка почты.
- Получатели: канал + адрес (telegram / max / …) на события заявок и поддержки.
- Роли: полный админ = пустая роль + `ADMIN_TOKEN`.

Метрику и OAuth-секреты здесь **не** дублируем.

---

## 3. Только клики владельца (не админ-экраны сайта)

- [ ] Яндекс.Вебмастер: «Проверить», sitemap `https://tikhvin-palomnik.ru/sitemap.xml`
- [ ] Яндекс.Wordstat: частоты по кластерам из [SEO_WORDSTAT.md](SEO_WORDSTAT.md) (СПб + ЛО); без выдуманных цифр
- [ ] Google Search Console: то же
- [ ] Карточка организации в Яндекс Бизнесе / Картах и 2ГИС
- [ ] Ссылка с сайта епархии [tikhvin-eparhia.ru](https://www.tikhvin-eparhia.ru/)
- [ ] MX / SPF / DKIM для `info@` (Яндекс 360 / DNS, не эта ВМ)
- [ ] Получатели Telegram написали боту `/start`
- [ ] Цены и даты туров из листовок (пока скрыты)
- [ ] Вопрос: сумма заявки на регулярный тур сейчас **0** — какая внутренняя цена?
- [ ] `AWAITING_PAYMENT` / `PAID` — подтвердить, если нужен эквайринг
- [ ] Юрист: тексты согласий; реквизиты в `OPERATOR_*` когда будут
- [ ] Перед живым ботом/GPT: ВМ **4 vCPU / 8 ГБ**, диск ≥80 ГБ

Подробности SEO: [SEO_ADS.md](SEO_ADS.md). Проход по мета-тегам и Wordstat (v4 этап 10) — [SEO_WORDSTAT.md](SEO_WORDSTAT.md); кнопки «Проверить» в кабинетах поисковиков по-прежнему только здесь.

---

## 4. После заполнения секретов

1. Обновить `.env.production` на ВМ (не коммитить).
2. Пока нет этапа 12: `make deploy`. После этапа 12: мерж в `main` выкладывает код; **только env** — всё равно пересборка на ВМ (секреты не в git).
3. Hard-refresh сайта.
4. Проверить: главная 200, `/news`, `/management`, заявка → Telegram.
