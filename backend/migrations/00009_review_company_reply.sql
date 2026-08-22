-- +goose Up
ALTER TABLE reviews
    ADD COLUMN company_reply TEXT NOT NULL DEFAULT '',
    ADD COLUMN company_replied_at TIMESTAMPTZ;

-- +goose Down
ALTER TABLE reviews
    DROP COLUMN IF EXISTS company_replied_at,
    DROP COLUMN IF EXISTS company_reply;
