-- +goose Up
ALTER TABLE cms_blocks DROP CONSTRAINT cms_blocks_type_check;
ALTER TABLE cms_blocks ADD CONSTRAINT cms_blocks_type_check CHECK (
    type IN (
        'hero',
        'about',
        'why_us',
        'how_it_works',
        'faq',
        'cta',
        'rich_text',
        'popular_destinations',
        'testimonials',
        'featured_route'
    )
);

-- +goose Down
ALTER TABLE cms_blocks DROP CONSTRAINT cms_blocks_type_check;
ALTER TABLE cms_blocks ADD CONSTRAINT cms_blocks_type_check CHECK (
    type IN (
        'hero',
        'about',
        'why_us',
        'how_it_works',
        'faq',
        'cta',
        'rich_text',
        'popular_destinations',
        'testimonials'
    )
);
