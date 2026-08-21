-- +goose Up
CREATE TABLE cms_pages (
    id UUID PRIMARY KEY,
    slug TEXT NOT NULL,
    title TEXT NOT NULL,
    path TEXT NOT NULL,
    is_published BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT cms_pages_slug_unique UNIQUE (slug),
    CONSTRAINT cms_pages_path_unique UNIQUE (path)
);

CREATE INDEX idx_cms_pages_published ON cms_pages (is_published);

CREATE TABLE cms_blocks (
    id UUID PRIMARY KEY,
    page_id UUID NOT NULL REFERENCES cms_pages(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    content JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_visible BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT cms_blocks_type_check CHECK (
        type IN ('hero', 'about', 'why_us', 'how_it_works', 'faq', 'cta', 'rich_text', 'popular_destinations', 'testimonials')
    )
);

CREATE INDEX idx_cms_blocks_page_order ON cms_blocks (page_id, sort_order);

-- +goose Down
DROP TABLE IF EXISTS cms_blocks;
DROP TABLE IF EXISTS cms_pages;
