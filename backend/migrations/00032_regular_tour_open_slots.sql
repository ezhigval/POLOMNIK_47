-- +goose Up
-- Регулярный тур: в админке часто заполняют только «всего мест», slots_left остаётся 0,
-- и витрина показывает «Мест нет». Открываем группу: свободно = всего.

-- +goose StatementBegin
UPDATE tours
SET slots_left = slots_total,
    updated_at = NOW()
WHERE is_regular
  AND slots_total > 0
  AND slots_left = 0;
-- +goose StatementEnd

-- +goose Down
SELECT 1;
