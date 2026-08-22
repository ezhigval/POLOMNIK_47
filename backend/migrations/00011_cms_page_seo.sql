-- +goose Up
ALTER TABLE cms_pages
    ADD COLUMN IF NOT EXISTS meta_title TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS meta_description TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE cms_pages
    DROP COLUMN IF EXISTS meta_description,
    DROP COLUMN IF EXISTS meta_title;
