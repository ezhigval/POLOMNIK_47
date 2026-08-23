-- +goose Up
CREATE TABLE smm_posts (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    url TEXT NOT NULL DEFAULT '',
    publish_at TIMESTAMPTZ NOT NULL,
    channels TEXT[] NOT NULL,
    published_at TIMESTAMPTZ,
    results JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_smm_posts_publish_at ON smm_posts (publish_at);

-- +goose Down
DROP INDEX IF EXISTS idx_smm_posts_publish_at;
DROP TABLE IF EXISTS smm_posts;
