-- +goose Up
ALTER TABLE site_settings ADD COLUMN logo_path TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE site_settings DROP COLUMN logo_path;
