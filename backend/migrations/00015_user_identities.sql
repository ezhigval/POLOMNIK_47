-- +goose Up
CREATE TABLE user_identities (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider TEXT NOT NULL,
    subject TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (provider, subject)
);

CREATE INDEX idx_user_identities_user ON user_identities (user_id);

INSERT INTO user_identities (user_id, provider, subject, created_at)
SELECT id, lower(oauth_provider), oauth_subject, created_at
FROM users
WHERE oauth_provider <> '' AND oauth_subject <> '';

DROP INDEX IF EXISTS idx_users_oauth_unique;
ALTER TABLE users DROP COLUMN IF EXISTS oauth_subject;
ALTER TABLE users DROP COLUMN IF EXISTS oauth_provider;

-- +goose Down
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS oauth_provider TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS oauth_subject TEXT NOT NULL DEFAULT '';

UPDATE users u
SET oauth_provider = i.provider,
    oauth_subject = i.subject
FROM (
    SELECT DISTINCT ON (user_id) user_id, provider, subject
    FROM user_identities
    ORDER BY user_id, created_at ASC
) i
WHERE u.id = i.user_id;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_oauth_unique ON users (oauth_provider, oauth_subject)
    WHERE oauth_provider <> '' AND oauth_subject <> '';

DROP INDEX IF EXISTS idx_user_identities_user;
DROP TABLE IF EXISTS user_identities;
