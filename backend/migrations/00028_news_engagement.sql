-- +goose Up
CREATE TABLE IF NOT EXISTS news_likes (
    news_id UUID NOT NULL REFERENCES news_articles (id) ON DELETE CASCADE,
    visitor_id TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    PRIMARY KEY (news_id, visitor_id)
);

CREATE INDEX IF NOT EXISTS news_likes_news_id_idx ON news_likes (news_id);

CREATE TABLE IF NOT EXISTS news_comments (
    id UUID PRIMARY KEY,
    news_id UUID NOT NULL REFERENCES news_articles (id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    body TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc')
);

CREATE INDEX IF NOT EXISTS news_comments_news_created_idx ON news_comments (news_id, created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS news_comments;
DROP TABLE IF EXISTS news_likes;
