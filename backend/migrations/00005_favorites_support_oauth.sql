-- +goose Up
ALTER TABLE users
    ALTER COLUMN phone DROP NOT NULL,
    ALTER COLUMN password_hash DROP NOT NULL;

ALTER TABLE users DROP CONSTRAINT IF EXISTS users_phone_unique;

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS oauth_provider TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS oauth_subject TEXT NOT NULL DEFAULT '';

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_phone_unique ON users (phone) WHERE phone IS NOT NULL AND phone <> '';
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_oauth_unique ON users (oauth_provider, oauth_subject)
    WHERE oauth_provider <> '' AND oauth_subject <> '';

CREATE TABLE favorites (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tour_id UUID NOT NULL REFERENCES tours(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, tour_id)
);

CREATE INDEX idx_favorites_user ON favorites (user_id);
CREATE INDEX idx_favorites_tour ON favorites (tour_id);

CREATE TABLE support_threads (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    subject TEXT NOT NULL DEFAULT 'Обращение в поддержку',
    status TEXT NOT NULL DEFAULT 'open',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_support_threads_user_open ON support_threads (user_id)
    WHERE status = 'open';

CREATE TABLE support_messages (
    id UUID PRIMARY KEY,
    thread_id UUID NOT NULL REFERENCES support_threads(id) ON DELETE CASCADE,
    sender_type TEXT NOT NULL CHECK (sender_type IN ('user', 'staff')),
    body TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_support_messages_thread ON support_messages (thread_id, created_at ASC);

-- +goose Down
DROP INDEX IF EXISTS idx_support_messages_thread;
DROP TABLE IF EXISTS support_messages;
DROP INDEX IF EXISTS idx_support_threads_user_open;
DROP TABLE IF EXISTS support_threads;
DROP INDEX IF EXISTS idx_favorites_tour;
DROP INDEX IF EXISTS idx_favorites_user;
DROP TABLE IF EXISTS favorites;
DROP INDEX IF EXISTS idx_users_oauth_unique;
ALTER TABLE users DROP COLUMN IF EXISTS oauth_subject;
ALTER TABLE users DROP COLUMN IF EXISTS oauth_provider;
DROP INDEX IF EXISTS idx_users_phone_unique;
ALTER TABLE users ADD CONSTRAINT users_phone_unique UNIQUE (phone);
ALTER TABLE users ALTER COLUMN phone SET NOT NULL;
ALTER TABLE users ALTER COLUMN password_hash SET NOT NULL;
