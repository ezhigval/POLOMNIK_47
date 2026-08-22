-- +goose Up
CREATE TABLE news_articles (
    id UUID PRIMARY KEY,
    slug TEXT NOT NULL,
    title TEXT NOT NULL,
    excerpt TEXT NOT NULL,
    body TEXT NOT NULL,
    image_url TEXT NOT NULL DEFAULT '',
    published_at DATE NOT NULL,
    is_published BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT news_articles_slug_unique UNIQUE (slug)
);

CREATE INDEX idx_news_articles_published ON news_articles (is_published, published_at DESC);

-- +goose Down
DROP TABLE IF EXISTS news_articles;
