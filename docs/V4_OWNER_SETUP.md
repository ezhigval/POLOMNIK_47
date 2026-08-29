# Чеклист владельца к v4

Сайт: **https://tikhvin-palomnik.ru**  
API: **https://api.tikhvin-palomnik.ru**  
План: [V4_PLAN.md](V4_PLAN.md). Секреты только в gitignored `.env.production` на ВМ. После правок env — `make deploy` (без `compose down -v`).

v2-секреты: [V2_OWNER_SETUP.md](V2_OWNER_SETUP.md). v3-адаптеры: [V3_OWNER_SETUP.md](V3_OWNER_SETUP.md).  
Секреты (откуда взять): [SECRETS.md](SECRETS.md) · почта DNS: [MAIL_DNS.md](MAIL_DNS.md) · OAuth: [OAUTH_SETUP.md](OAUTH_SETUP.md) · Telegram: [TELEGRAM_SETUP.md](TELEGRAM_SETUP.md) · SEO: [SEO_ADS.md](SEO_ADS.md).

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

### v4 этапы 9–19 (продукт + инфра)

Этапы **9–11, 13–19** — код на `main` и на проде. Goose: **00032**. Автодеплой (этап 12): workflow есть, секрет SSH **не** подходит — см. § GitHub Actions.

| Этап | Что | Секреты / env |
|------|-----|----------------|
| 9 | Cookie `palomnik_viewed_tours`, блок «Вы смотрели» | нет (существующий `JWT_SECRET`, `palomnik_token`) |
| 10 | Компакт новости на главной, SEO pass, `seo-wordstat.ts` | Wordstat / Вебмастер / GSC — кабинеты владельца |
| 11 | Лента `/news`, лайки (cookie `palomnik_visitor`), комментарии (кабинет) | нет новых env |
| 12 | GitHub Actions deploy после CI | `DEPLOY_SSH_KEY` (**приватный** ключ). Сейчас: `Permission denied (publickey)` → `make deploy` |
| 13 | Флаги туров | не назначать цену регулярному |
| 16 | Цена если > 0; `payment_status` | live Pay не включать |
| 17–19 | Тихвинский путь, места regular, мобильный UX | нет новых секретов |

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
| `YANDEX_OAUTH_CLIENT_ID` / `YANDEX_OAUTH_CLIENT_SECRET` | [OAUTH_SETUP.md](OAUTH_SETUP.md) / [SECRETS.md](SECRETS.md). Пара из кабинета oauth.yandex.ru — **только** на ВМ, не в git |
| `VK_OAUTH_CLIENT_ID` / `VK_OAUTH_CLIENT_SECRET` | то же |
| `MAX_OAUTH_*` | то же; URL authorize/token/userinfo когда выдадут |
| `MAIL_ADAPTER=smtp`, `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `SMTP_FROM` | [MAIL_DNS.md](MAIL_DNS.md) (после MX/SPF/DKIM). SMTP-пароль ≠ OAuth Client secret |
| `PHONE_ADAPTER=smsru`, `SMSRU_API_ID` | callcheck, не SMS |
| `NEXT_PUBLIC_GA_ID` | опционально |
| `NEXT_PUBLIC_YANDEX_VERIFICATION` / `NEXT_PUBLIC_GOOGLE_SITE_VERIFICATION` | сменить коды Вебмастера/GSC |

Несекретные адреса пересылки: админка «Настройки» (`mail_forward_to`) или запасной `MAIL_FORWARD_TO` в env.

### GitHub Actions — этап 12 (автодеплой)

Workflow `.github/workflows/deploy.yml`: после **успешного CI** на push в `main` → rsync + compose на ВМ (тот же `deploy/yandex/deploy.sh`, `DEPLOY_CI=1`). Один деплой за раз (`concurrency: production-deploy`). Проверка `https://api.tikhvin-palomnik.ru/health/ready`.

В **Settings → Secrets and variables → Actions** (значения не сюда):

| Имя | Зачем |
|-----|--------|
| `DEPLOY_SSH_KEY` | **Обязательно.** Приватный ключ OpenSSH (не `.pub`). Должен открывать `smailikin70@93.77.165.81`. На сервере уже есть pubkey комментария `github-actions-deploy-palomnik`: `ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIA/Pf0QxoAH4bbrcC5bRW9Y7lN0OPigypFRn+0KI0HQZ` |
| `DEPLOY_SSH_HOST` | Опционально. Target, если не `smailikin70@93.77.165.81` |

На 2026-08-27 workflow красный: в секрет, судя по логу, попал не тот ключ (часто вставляют **публичный**). Пока не исправлено — `make deploy`. Проверка: Actions → Deploy production → Run workflow.

Не класть в Actions весь `.env.production`. Правка только env на ВМ требует пересборки фронта.

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

- [ ] **GitHub `DEPLOY_SSH_KEY`**: приватный ключ к pubkey `github-actions-deploy-palomnik` на ВМ (сейчас автодеплой падает)
- [ ] Яндекс.Вебмастер: «Проверить», sitemap `https://tikhvin-palomnik.ru/sitemap.xml`
- [ ] Яндекс.Wordstat: частоты по кластерам из [SEO_WORDSTAT.md](SEO_WORDSTAT.md) (СПб + ЛО); без выдуманных цифр
- [ ] Google Search Console: то же
- [ ] Карточка организации в Яндекс Бизнесе / Картах и 2ГИС
- [ ] Ссылка с сайта епархии [tikhvin-eparhia.ru](https://www.tikhvin-eparhia.ru/)
- [ ] Почта `info@`: либо Яндекс 360 [MAIL_DNS.md](MAIL_DNS.md), либо свой MX [MAIL_SELFHOST.md](MAIL_SELFHOST.md) (на YC исходящий :25 закрыт)
- [ ] Получатели Telegram написали боту `/start`
- [ ] Листовки на проде уже **регулярные без цены**. Если нужен датированный выезд — дата и цена в админке (не выдумывать)
- [ ] Вопрос: сумма заявки на регулярный тур сейчас **0** — какая внутренняя цена?
- [ ] `AWAITING_PAYMENT` / `PAID` как **живой** эквайринг — подтвердить, если нужен договор (ручное поле статуса уже есть)
- [ ] Юрист: тексты согласий; реквизиты в `OPERATOR_*` когда будут
- [ ] Перед живым ботом/GPT: ВМ **4 vCPU / 8 ГБ**, диск ≥80 ГБ
- [ ] Модерация комментариев новостей — пока нет; включать только по просьбе

Подробности SEO: [SEO_ADS.md](SEO_ADS.md). Проход по мета-тегам и Wordstat (v4 этап 10) — [SEO_WORDSTAT.md](SEO_WORDSTAT.md); кнопки «Проверить» в кабинетах поисковиков по-прежнему только здесь.

---

## 4. После заполнения секретов

1. Обновить `.env.production` на ВМ (не коммитить).
2. Пока `DEPLOY_SSH_KEY` неверный: `make deploy`. После исправления секрета: мерж в `main` (зелёный CI) должен выкладывать сам. **Только env** на ВМ — всё равно пересборка.
3. Hard-refresh сайта.
4. Проверить: главная 200, `/news` (лента, лайк без входа, комментарий после входа), `/management`, заявка → Telegram.

---

## 5. Автоматизировано vs вручную

| Действие | Кто |
|----------|-----|
| CI (тесты, lint, build) на push/PR | GitHub Actions |
| Деплой после зелёного CI на `main` | GitHub Actions **если ключ верный**; иначе `make deploy` |
| Goose `up` при деплое | `deploy/yandex/deploy.sh` |
| Лайки/комментарии новостей | Код (без новых секретов) |
| Wordstat частоты, Вебмастер «Проверить», GSC | Владелец |
| OAuth, SMTP, sms.ru ключи | Владелец → `.env.production` |
| Цены туров, регулярный тур (сумма 0) | Владелец (бизнес-решение) |
| `AWAITING_PAYMENT` / `PAID`, эквайринг | Владелец (не в коде) |
| Юридичка, `OPERATOR_*` реквизиты | Владелец + юрист |

---

## 6. Открытые бизнес-вопросы

- Сумма заявки на **регулярный** тур сейчас **0** — какая внутренняя цена? На витрине цены нет.
- Живой эквайринг (`PAYMENT_ADAPTER`) — нужен ли, при том что ручной `payment_status` уже есть?
- Модерация комментариев новостей — пока нет.
- Этапы **2–8** v4 — не блокируют витрину 9–19.
