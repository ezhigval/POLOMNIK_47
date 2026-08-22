-- +goose Up
-- Prod already applied goose 10 and 11 (site_settings historically, then CMS SEO).
-- Keep 00010/00011 filenames stable; this version is idempotent.
CREATE TABLE IF NOT EXISTS telegram_recipients (
    id SMALLINT PRIMARY KEY CHECK (id = 1),
    booking_usernames TEXT NOT NULL DEFAULT '',
    support_usernames TEXT NOT NULL DEFAULT '',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO telegram_recipients (id, booking_usernames, support_usernames)
VALUES (1, '', '')
ON CONFLICT (id) DO NOTHING;

CREATE TABLE IF NOT EXISTS telegram_chat_map (
    username TEXT PRIMARY KEY,
    chat_id TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS telegram_chat_map;
DROP TABLE IF EXISTS telegram_recipients;
