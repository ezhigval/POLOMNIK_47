-- +goose Up
CREATE TABLE passengers (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    phone TEXT NOT NULL,
    birth_date DATE NOT NULL,
    passport TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_passengers_user ON passengers (user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_passengers_user;
DROP TABLE IF EXISTS passengers;
