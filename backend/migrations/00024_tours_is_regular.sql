-- +goose Up
ALTER TABLE tours
    ADD COLUMN IF NOT EXISTS is_regular BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE tours
    DROP CONSTRAINT IF EXISTS tours_date_range_check;

ALTER TABLE tours
    ALTER COLUMN date_start DROP NOT NULL,
    ALTER COLUMN date_end DROP NOT NULL;

ALTER TABLE tours
    ADD CONSTRAINT tours_date_range_check CHECK (
        (date_start IS NULL AND date_end IS NULL)
        OR (date_start IS NOT NULL AND date_end IS NOT NULL AND date_start <= date_end)
    );

ALTER TABLE tours
    DROP CONSTRAINT IF EXISTS tours_regular_schedule_check;

ALTER TABLE tours
    ADD CONSTRAINT tours_regular_schedule_check CHECK (
        (is_regular AND date_start IS NULL AND date_end IS NULL)
        OR (NOT is_regular AND date_start IS NOT NULL AND date_end IS NOT NULL)
    );

-- +goose Down
ALTER TABLE tours
    DROP CONSTRAINT IF EXISTS tours_regular_schedule_check;

ALTER TABLE tours
    DROP CONSTRAINT IF EXISTS tours_date_range_check;

UPDATE tours
SET date_start = CURRENT_DATE,
    date_end = CURRENT_DATE
WHERE date_start IS NULL OR date_end IS NULL;

ALTER TABLE tours
    ALTER COLUMN date_start SET NOT NULL,
    ALTER COLUMN date_end SET NOT NULL;

ALTER TABLE tours
    ADD CONSTRAINT tours_date_range_check CHECK (date_start <= date_end);

ALTER TABLE tours
    DROP COLUMN IF EXISTS is_regular;
