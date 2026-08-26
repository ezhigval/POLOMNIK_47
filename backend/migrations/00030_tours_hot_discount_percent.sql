-- +goose Up
ALTER TABLE tours
    ADD COLUMN IF NOT EXISTS hot_discount_percent INTEGER NOT NULL DEFAULT 0;

ALTER TABLE tours
    DROP CONSTRAINT IF EXISTS tours_hot_discount_percent_check;

ALTER TABLE tours
    ADD CONSTRAINT tours_hot_discount_percent_check CHECK (
        hot_discount_percent >= 0 AND hot_discount_percent <= 100
    );

-- +goose Down
ALTER TABLE tours
    DROP CONSTRAINT IF EXISTS tours_hot_discount_percent_check;

ALTER TABLE tours
    DROP COLUMN IF EXISTS hot_discount_percent;
