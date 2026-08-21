-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL DEFAULT '',
    phone TEXT NOT NULL,
    name TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT users_phone_unique UNIQUE (phone)
);

CREATE UNIQUE INDEX idx_users_email_unique ON users (email) WHERE email <> '';

CREATE INDEX idx_users_phone ON users (phone);
CREATE INDEX idx_users_email ON users (email) WHERE email <> '';

ALTER TABLE bookings
    ADD COLUMN user_id UUID REFERENCES users(id) ON DELETE SET NULL;

CREATE INDEX idx_bookings_user ON bookings (user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_bookings_user;
ALTER TABLE bookings DROP COLUMN IF EXISTS user_id;
DROP INDEX IF EXISTS idx_users_email_unique;
DROP INDEX IF EXISTS idx_users_phone;
DROP TABLE IF EXISTS users;
