-- +goose Up
ALTER TABLE news_articles
    ADD COLUMN IF NOT EXISTS is_pinned BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX IF NOT EXISTS idx_news_articles_pinned_published
    ON news_articles (is_pinned DESC, published_at DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_news_articles_pinned_published;
ALTER TABLE news_articles DROP COLUMN IF EXISTS is_pinned;
