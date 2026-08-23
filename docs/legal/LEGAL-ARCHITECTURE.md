# LEGAL-ARCHITECTURE: система согласий

**Дата:** 2026-08-23  
**Статус:** реализовано согласно подготовленной модели; перед публикацией — проверка юристом

---

## 1. Документы

| № | type | Название | Назначение |
|---|------|----------|------------|
| 1 | `privacy_policy` | Политика обработки ПД | Информационный документ для всех субъектов |
| 2 | `personal_data` | Согласие на обработку ПД | Обязательное для регистрации, заявок, ЛК |
| 3 | `distribution` | Согласие на распространение ПД | Добровольное: публикация отзывов, фото |
| 4 | `marketing` | Согласие на рекламу | Добровольное: рассылки, акции |
| 5 | `cookie` | Cookie Policy | Категории cookie, аналитика |

---

## 2. Где используется каждый документ

| Документ | Форма / механизм | Checkbox |
|----------|------------------|----------|
| `personal_data` | Регистрация, заявка, отзыв (обработка) | Обязательный, unchecked |
| `marketing` | Регистрация, заявка | Необязательный, unchecked |
| `distribution` | Отзыв (публикация), фото с туров | Необязательный, unchecked |
| `cookie` | Cookie banner | Три варианта выбора |
| `privacy_policy` | Footer, ссылки из других документов | — |

---

## 3. Типы согласий в БД (`consents.consent_type`)

| consent_type | Связанный документ | Обязательность |
|--------------|-------------------|----------------|
| `personal_data` | `personal_data` | Да (регистрация, заявка) |
| `marketing` | `marketing` | Нет |
| `distribution` | `distribution` | Нет |
| `cookie_all` | `cookie` | Нет (выбор пользователя) |
| `cookie_essential` | `cookie` | Нет |
| `cookie_reject` | `cookie` | Нет |

---

## 4. Что сохраняется в БД

### `legal_documents`

```
id, type, version, title, content, published_at, updated_at, is_active
```

- Одна активная версия на тип (`is_active = true`)
- Исторические версии сохраняются (`is_active = false`)
- При публикации новой версии старая деактивируется

### `consents`

```
id, user_id, request_id, consent_type, document_id, document_version, accepted_at, ip, user_agent
```

- `document_version` и `accepted_at` определяются **сервером**
- Клиент не может подменить версию или дату

---

## 5. API

### Публичные

```
GET  /api/v1/legal/documents
GET  /api/v1/legal/documents/:type
GET  /api/v1/legal/documents/:type/versions/:version
POST /api/v1/consents
```

### Личный кабинет

```
GET /api/v1/me/consents
```

### Админка

```
GET  /api/v1/management/legal/documents
POST /api/v1/management/legal/documents
GET  /api/v1/management/consents
```

---

## 6. Версионирование

1. Тексты хранятся в `backend/internal/legal/content/`
2. Bootstrap v1.0 при первом запуске API (`BootstrapInitialDocuments`)
3. Новая версия публикуется через админ-API (`PublishNewVersion`)
4. Изменение текста → новая версия (1.0 → 1.1 → 2.0)
5. Согласия привязаны к `document_id` + `document_version`

---

## 7. Конфиг оператора

Единый источник:

- `backend/internal/legal/operator/operator.go`
- `frontend/src/lib/operator-config.ts`

---

## 8. Разделение сущностей

```
personal_data ≠ marketing ≠ distribution ≠ cookie ≠ privacy_policy ≠ оферта ≠ пользовательское соглашение
```

Каждый тип — отдельный документ, отдельный checkbox (где применимо), отдельный `consent_type`.
