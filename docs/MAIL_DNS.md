# Почта `info@tikhvin-palomnik.ru`: MX, SPF, DKIM

Почту принимает **Яндекс 360**, не эта ВМ. DNS зоны — **REG.RU** (`ns1.reg.ru` / `ns2.reg.ru`). Записи **A** для `@`, `www` и `api` (сайт и API) **не трогать**.

Официально у Яндекса для почты нужны три записи: [MX](https://yandex.ru/support/yandex-360/business/admin/ru/domains/dns/mx), [SPF](https://yandex.ru/support/yandex-360/business/admin/ru/domains/dns/spf), [DKIM](https://yandex.ru/support/yandex-360/business/admin/ru/domains/dns/dkim). Сводка: [Первоначальная настройка почты](https://yandex.ru/support/yandex-360/business/admin/ru/mail/start).

Секреты SMTP после DNS — в [SECRETS.md](SECRETS.md). Сайт сам почту не принимает: без MX на `mx.yandex.net` письма на `info@` никуда не доходят.

---

## Как сейчас (проверка 2026-08-27)

| Запись | Факт | Что сделать |
|--------|------|-------------|
| NS | `ns1.reg.ru.`, `ns2.reg.ru.` | **Оставить.** Не делегировать зону на `dns1.yandex.net` — иначе сломаются A сайта, если их не скопировать. |
| A `@` / `www` / `api` | `93.77.165.81` | **Не менять.** |
| MX `@` | `0 mail.tikhvin-palomnik.ru.` | **Заменить.** У `mail.tikhvin-palomnik.ru` нет A и нет CNAME — входящая почта сейчас некуда. |
| TXT SPF | нет | **Добавить.** |
| TXT `mail._domainkey` (DKIM) | нет | **Добавить** ключ из кабинета Яндекс 360 (он уникальный, не копировать из этой инструкции). |
| TXT `_dmarc` | нет | По желанию, после того как MX/SPF/DKIM зелёные. |

TTL у Яндекса в примерах — `21600`. У REG.RU может быть своё поле TTL; если просят — `21600`. Распространение DNS — до **72 часов**, часто меньше.

---

## 1. Подключить домен в Яндекс 360

1. Войдите в [https://admin.yandex.ru/](https://admin.yandex.ru/) (организация Яндекс 360 для бизнеса).
2. **Общие настройки → Домены** ([https://admin.yandex.ru/domains](https://admin.yandex.ru/domains)).
3. Добавьте домен **`tikhvin-palomnik.ru`**, если его ещё нет.
4. **Подтвердите владение**, как покажет кабинет (обычно отдельная **TXT** на `@`). Это **ещё одна** TXT рядом с уже существующей Google-верификацией — старую не удалять, A-записи не трогать.
5. После появления TXT в REG.RU в кабинете Яндекса нажмите **Проверить**.

Не меняйте NS-серверы домена на Яндекс, пока сайт живёт на REG.RU DNS.

---

## 2. Записи в REG.RU

1. [https://www.reg.ru](https://www.reg.ru) → логин → **Домены и услуги** → **tikhvin-palomnik.ru** → **DNS-серверы и управление зоной**.
2. **Удалите все текущие MX** (сейчас это `mail.tikhvin-palomnik.ru`). Два MX сразу (старый + Яндекс) Яндекс не рекомендует.

### MX

| Поле в REG.RU | Значение |
|---------------|----------|
| Тип | MX |
| Subdomain / хост | `@` |
| Mail Server | `mx.yandex.net.` (точка в конце, если панель её не дописывает сама) |
| Priority | `10` |
| TTL | `21600`, если поле есть |

### SPF

Одна TXT на `@` (не вместо Google-верификации, **рядом**):

| Поле | Значение |
|------|----------|
| Тип | TXT |
| Subdomain / хост | `@` |
| Text | `v=spf1 redirect=_spf.yandex.net` |

Если панель уже содержит другую TXT на `@` — оставьте её и добавьте эту второй строкой. Не склеивайте SPF и `google-site-verification` в одну запись.

Письма сайт шлёт через SMTP Яндекса (`smtp.yandex.ru`), IP ВМ в SPF **не** добавлять.

### DKIM

Ключ **только** из кабинета, он у каждого домена свой.

1. [https://admin.yandex.ru/domains](https://admin.yandex.ru/domains) → у `tikhvin-palomnik.ru` → **Настроить DKIM**.
2. Скопируйте значение публичного ключа целиком (вид `v=DKIM1; k=rsa; … p=…`).
3. В REG.RU новая TXT:

| Поле | Значение |
|------|----------|
| Тип | TXT |
| Subdomain | `mail._domainkey` (без имени домена; панель сама допишет `.tikhvin-palomnik.ru`) |
| Text | **весь** скопированный ключ, без кавычек и без обрезки |

Официальная шпаргалка REG.RU у Яндекса: [DKIM → reg.ru](https://yandex.ru/support/yandex-360/business/admin/ru/domains/dns/dkim).

---

## 3. Проверка

Подождите (минуты… часы, редко до 72 ч). Затем в Яндекс 360: **Домены → Проверить**.

Снаружи (или [digwebinterface.com](https://www.digwebinterface.com/)):

```text
dig MX tikhvin-palomnik.ru +short
# ожидается: 10 mx.yandex.net.

dig TXT tikhvin-palomnik.ru +short
# среди ответов: "v=spf1 redirect=_spf.yandex.net"

dig TXT mail._domainkey.tikhvin-palomnik.ru +short
# строка, начинающаяся с v=DKIM1
```

Старого MX `mail.tikhvin-palomnik.ru` в ответе быть не должно.

---

## 4. Ящик `info@`

1. В Яндекс 360: сотрудники / почтовые ящики ([https://admin.yandex.ru/users](https://admin.yandex.ru/users)) → создайте пользователя с адресом **`info@tikhvin-palomnik.ru`**.
2. Зайдите в этот ящик через [https://mail.yandex.ru](https://mail.yandex.ru) (логин — **полный** адрес `info@tikhvin-palomnik.ru`).
3. Настройки Почты → **Почтовые программы**: включите IMAP и **Пароли приложений и OAuth-токены**.
4. Пересылка входящих на личные ящики — **в Яндекс 360 / настройках ящика**, не в коде сайта. Несекретный список «куда ещё слать копии с сайта» — админка **Настройки** (поле пересылки) или запасной `MAIL_FORWARD_TO` в env.

Пока MX не указывает на Яндекс, ящик можно создать, но входящие с интернета не придут.

---

## 5. SMTP сайта (после зелёных MX/SPF/DKIM)

Сайт шлёт письма (подтверждение почты, сброс пароля) через `smtp.yandex.ru`, порт **587** (так в коде). Обычный пароль Яндекс ID для SMTP **не** подходит — нужен **пароль приложения**.

1. Под учёткой ящика `info@` откройте [Пароли приложений](https://id.yandex.ru/security/app-passwords) (Яндекс ID → Безопасность → Пароли приложений).
2. Создать → тип **Почта** → имя вроде `palomnik-smtp`.
3. Скопируйте пароль **сразу** (показывается один раз). Начинает действовать через 2–3 часа.
4. На ВМ в `/opt/polomnik/.env.production` (значения не в git):

```bash
MAIL_ADAPTER=smtp
SMTP_HOST=smtp.yandex.ru
SMTP_PORT=587
SMTP_USERNAME=info@tikhvin-palomnik.ru
SMTP_PASSWORD=          # пароль приложения, не пароль входа на Яндекс
SMTP_FROM=info@tikhvin-palomnik.ru
# опционально, если в админке пересылка пустая:
# MAIL_FORWARD_TO=smailikin70@yandex.ru
```

5. Скажите агенту сделать `make deploy` (без `compose down -v`). Пока `DEPLOY_SSH_KEY` в GitHub неверный — только ручной деплой.
6. Проверка: письмо на `info@` с внешнего ящика; на сайте «Забыли пароль?» (после деплоя SMTP) — ссылка должна уйти.

Без `MAIL_ADAPTER=smtp` регистрация жива, письмо не уходит, «Забыли пароль?» пишет «пока недоступно» — это норма.

---

## По желанию после рабочей почты

**DMARC** (TXT, хост `_dmarc`):

```text
v=DMARC1; p=none; rua=mailto:info@tikhvin-palomnik.ru
```

`p=none` только наблюдает, письма не режет. Ужесточать (`quarantine` / `reject`) — после того как увидите, что SPF/DKIM проходят.

**CNAME `mail`** — не обязателен для SMTP и для MX. Если нужен ярлык веб-почты на поддомене, смотрите актуальную подсказку в кабинете Яндекс 360; не ставьте A `mail` на IP сайта.

---

## Чего не делать

- Не менять A `@` / `www` / `api` и не переносить NS на Яндекс «для почты».
- Не оставлять старый MX `mail.tikhvin-palomnik.ru` рядом с `mx.yandex.net`.
- Не коммитить `SMTP_PASSWORD` и не вставлять его в чат / PR.
- Не включать Bitrix/1С/эквайринг «заодно».
