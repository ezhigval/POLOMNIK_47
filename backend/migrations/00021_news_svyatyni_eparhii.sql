-- +goose Up
INSERT INTO news_articles (
    id,
    slug,
    title,
    excerpt,
    body,
    image_url,
    published_at,
    is_published,
    sort_order,
    created_at,
    updated_at
) VALUES (
    'c0a82643-5124-47b0-aeff-193035bed12b',
    'svyatyni-tikhvinskoy-eparhii',
    'Святыни Тихвинской Епархии',
    'Явление Тихвинской иконы Божией Матери и история Обители',
    E'/images/news/svyatyni-eparhii/stend-01.jpg\n/images/news/svyatyni-eparhii/stend-02.jpg\n/images/news/svyatyni-eparhii/stend-03.jpg\n/images/news/svyatyni-eparhii/stend-04.jpg\n/images/news/svyatyni-eparhii/stend-05.jpg\n/images/news/svyatyni-eparhii/stend-06.jpg\n/images/news/svyatyni-eparhii/stend-07.jpg\n/images/news/svyatyni-eparhii/stend-08.jpg\n/images/news/svyatyni-eparhii/stend-09.jpg\n/images/news/svyatyni-eparhii/stend-10.jpg\n/images/news/svyatyni-eparhii/stend-11.jpg\n/images/news/svyatyni-eparhii/stend-12.jpg',
    '/images/news/svyatyni-eparhii/stend-01.jpg',
    '2026-08-26',
    TRUE,
    0,
    NOW(),
    NOW()
)
ON CONFLICT (slug) DO NOTHING;

-- +goose Down
DELETE FROM news_articles WHERE slug = 'svyatyni-tikhvinskoy-eparhii';
