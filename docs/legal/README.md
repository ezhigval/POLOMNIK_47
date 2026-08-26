# Юридическая документация — README

Проект системы согласий на обработку персональных данных для [tikhvin-palomnik.ru](https://tikhvin-palomnik.ru).

## Документы в этой папке

| Файл | Назначение |
|------|------------|
| [AUDIT.md](./AUDIT.md) | Предварительный аудит кодовой базы (этап 0) |
| [LEGAL-REVIEW.md](./LEGAL-REVIEW.md) | Проверка законодательства, открытые вопросы для юриста |
| [LEGAL-ARCHITECTURE.md](./LEGAL-ARCHITECTURE.md) | Техническая архитектура системы согласий |

## Статус реализации

Реализовано согласно подготовленной юридической модели; **перед публикацией требуется финальная проверка оператором/юристом**.

## Placeholders оператора

Реквизиты оператора задаются в:

- Backend: `backend/internal/legal/operator/operator.go` (env: `OPERATOR_*`)
- Frontend: `frontend/src/lib/operator-config.ts` (env: `NEXT_PUBLIC_OPERATOR_*`)

После получения настоящих реквизитов задать `OPERATOR_*` / `NEXT_PUBLIC_OPERATOR_*` в `.env.production` (compose прокидывает с v4 этапа 1). Не править ИНН/ОГРН в коде наугад.

## Типы юридических документов

| type | Документ | URL |
|------|----------|-----|
| `privacy_policy` | Политика обработки ПД | `/legal/privacy-policy` |
| `personal_data` | Согласие на обработку ПД | `/legal/personal-data-consent` |
| `distribution` | Согласие на распространение ПД | `/legal/distribution-consent` |
| `marketing` | Согласие на рекламу | `/legal/marketing-consent` |
| `cookie` | Cookie Policy | `/legal/cookie-policy` |
| `terms` | Пользовательское соглашение | `/legal/terms` |
| `offer` | Публичная оферта | `/legal/offer` |
