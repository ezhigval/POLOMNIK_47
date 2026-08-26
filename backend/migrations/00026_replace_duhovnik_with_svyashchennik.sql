-- +goose Up
-- Prod already applied goose 00023–00025. Replace the role word in live rows
-- rather than rewriting those migrations. Inflections follow the stem:
-- духовник/а/у/ом/е/и/ов/ами/ах → священник… ; same for Духовник / ДУХОВНИК.
-- Do not touch «духовный» / «духовность» (different stem).

UPDATE cms_blocks
SET
    content = replace(replace(replace(content::text, 'ДУХОВНИК', 'СВЯЩЕННИК'), 'Духовник', 'Священник'), 'духовник', 'священник')::jsonb,
    updated_at = NOW()
WHERE content::text ILIKE '%духовник%';

UPDATE cms_pages
SET
    title = replace(replace(replace(title, 'ДУХОВНИК', 'СВЯЩЕННИК'), 'Духовник', 'Священник'), 'духовник', 'священник'),
    meta_title = replace(replace(replace(meta_title, 'ДУХОВНИК', 'СВЯЩЕННИК'), 'Духовник', 'Священник'), 'духовник', 'священник'),
    meta_description = replace(replace(replace(meta_description, 'ДУХОВНИК', 'СВЯЩЕННИК'), 'Духовник', 'Священник'), 'духовник', 'священник'),
    updated_at = NOW()
WHERE title ILIKE '%духовник%'
   OR meta_title ILIKE '%духовник%'
   OR meta_description ILIKE '%духовник%';

UPDATE tours
SET
    title = replace(replace(replace(title, 'ДУХОВНИК', 'СВЯЩЕННИК'), 'Духовник', 'Священник'), 'духовник', 'священник'),
    description = replace(replace(replace(description, 'ДУХОВНИК', 'СВЯЩЕННИК'), 'Духовник', 'Священник'), 'духовник', 'священник'),
    updated_at = NOW()
WHERE title ILIKE '%духовник%'
   OR description ILIKE '%духовник%';

UPDATE reviews
SET
    text = replace(replace(replace(text, 'ДУХОВНИК', 'СВЯЩЕННИК'), 'Духовник', 'Священник'), 'духовник', 'священник'),
    company_reply = replace(replace(replace(company_reply, 'ДУХОВНИК', 'СВЯЩЕННИК'), 'Духовник', 'Священник'), 'духовник', 'священник'),
    updated_at = NOW()
WHERE text ILIKE '%духовник%'
   OR company_reply ILIKE '%духовник%';

UPDATE news_articles
SET
    title = replace(replace(replace(title, 'ДУХОВНИК', 'СВЯЩЕННИК'), 'Духовник', 'Священник'), 'духовник', 'священник'),
    excerpt = replace(replace(replace(excerpt, 'ДУХОВНИК', 'СВЯЩЕННИК'), 'Духовник', 'Священник'), 'духовник', 'священник'),
    body = replace(replace(replace(body, 'ДУХОВНИК', 'СВЯЩЕННИК'), 'Духовник', 'Священник'), 'духовник', 'священник'),
    updated_at = NOW()
WHERE title ILIKE '%духовник%'
   OR excerpt ILIKE '%духовник%'
   OR body ILIKE '%духовник%';

UPDATE site_settings
SET
    description = replace(replace(replace(description, 'ДУХОВНИК', 'СВЯЩЕННИК'), 'Духовник', 'Священник'), 'духовник', 'священник'),
    tagline = replace(replace(replace(tagline, 'ДУХОВНИК', 'СВЯЩЕННИК'), 'Духовник', 'Священник'), 'духовник', 'священник'),
    updated_at = NOW()
WHERE description ILIKE '%духовник%'
   OR tagline ILIKE '%духовник%';

UPDATE legal_documents
SET
    title = replace(replace(replace(title, 'ДУХОВНИК', 'СВЯЩЕННИК'), 'Духовник', 'Священник'), 'духовник', 'священник'),
    content = replace(replace(replace(content, 'ДУХОВНИК', 'СВЯЩЕННИК'), 'Духовник', 'Священник'), 'духовник', 'священник'),
    updated_at = NOW()
WHERE title ILIKE '%духовник%'
   OR content ILIKE '%духовник%';

UPDATE smm_posts
SET
    title = replace(replace(replace(title, 'ДУХОВНИК', 'СВЯЩЕННИК'), 'Духовник', 'Священник'), 'духовник', 'священник'),
    body = replace(replace(replace(body, 'ДУХОВНИК', 'СВЯЩЕННИК'), 'Духовник', 'Священник'), 'духовник', 'священник'),
    updated_at = NOW()
WHERE title ILIKE '%духовник%'
   OR body ILIKE '%духовник%';

-- +goose Down
UPDATE cms_blocks
SET
    content = replace(replace(replace(content::text, 'СВЯЩЕННИК', 'ДУХОВНИК'), 'Священник', 'Духовник'), 'священник', 'духовник')::jsonb,
    updated_at = NOW()
WHERE content::text ILIKE '%священник%';

UPDATE cms_pages
SET
    title = replace(replace(replace(title, 'СВЯЩЕННИК', 'ДУХОВНИК'), 'Священник', 'Духовник'), 'священник', 'духовник'),
    meta_title = replace(replace(replace(meta_title, 'СВЯЩЕННИК', 'ДУХОВНИК'), 'Священник', 'Духовник'), 'священник', 'духовник'),
    meta_description = replace(replace(replace(meta_description, 'СВЯЩЕННИК', 'ДУХОВНИК'), 'Священник', 'Духовник'), 'священник', 'духовник'),
    updated_at = NOW()
WHERE title ILIKE '%священник%'
   OR meta_title ILIKE '%священник%'
   OR meta_description ILIKE '%священник%';

UPDATE tours
SET
    title = replace(replace(replace(title, 'СВЯЩЕННИК', 'ДУХОВНИК'), 'Священник', 'Духовник'), 'священник', 'духовник'),
    description = replace(replace(replace(description, 'СВЯЩЕННИК', 'ДУХОВНИК'), 'Священник', 'Духовник'), 'священник', 'духовник'),
    updated_at = NOW()
WHERE title ILIKE '%священник%'
   OR description ILIKE '%священник%';

UPDATE reviews
SET
    text = replace(replace(replace(text, 'СВЯЩЕННИК', 'ДУХОВНИК'), 'Священник', 'Духовник'), 'священник', 'духовник'),
    company_reply = replace(replace(replace(company_reply, 'СВЯЩЕННИК', 'ДУХОВНИК'), 'Священник', 'Духовник'), 'священник', 'духовник'),
    updated_at = NOW()
WHERE text ILIKE '%священник%'
   OR company_reply ILIKE '%священник%';

UPDATE news_articles
SET
    title = replace(replace(replace(title, 'СВЯЩЕННИК', 'ДУХОВНИК'), 'Священник', 'Духовник'), 'священник', 'духовник'),
    excerpt = replace(replace(replace(excerpt, 'СВЯЩЕННИК', 'ДУХОВНИК'), 'Священник', 'Духовник'), 'священник', 'духовник'),
    body = replace(replace(replace(body, 'СВЯЩЕННИК', 'ДУХОВНИК'), 'Священник', 'Духовник'), 'священник', 'духовник'),
    updated_at = NOW()
WHERE title ILIKE '%священник%'
   OR excerpt ILIKE '%священник%'
   OR body ILIKE '%священник%';

UPDATE site_settings
SET
    description = replace(replace(replace(description, 'СВЯЩЕННИК', 'ДУХОВНИК'), 'Священник', 'Духовник'), 'священник', 'духовник'),
    tagline = replace(replace(replace(tagline, 'СВЯЩЕННИК', 'ДУХОВНИК'), 'Священник', 'Духовник'), 'священник', 'духовник'),
    updated_at = NOW()
WHERE description ILIKE '%священник%'
   OR tagline ILIKE '%священник%';

UPDATE legal_documents
SET
    title = replace(replace(replace(title, 'СВЯЩЕННИК', 'ДУХОВНИК'), 'Священник', 'Духовник'), 'священник', 'духовник'),
    content = replace(replace(replace(content, 'СВЯЩЕННИК', 'ДУХОВНИК'), 'Священник', 'Духовник'), 'священник', 'духовник'),
    updated_at = NOW()
WHERE title ILIKE '%священник%'
   OR content ILIKE '%священник%';

UPDATE smm_posts
SET
    title = replace(replace(replace(title, 'СВЯЩЕННИК', 'ДУХОВНИК'), 'Священник', 'Духовник'), 'священник', 'духовник'),
    body = replace(replace(replace(body, 'СВЯЩЕННИК', 'ДУХОВНИК'), 'Священник', 'Духовник'), 'священник', 'духовник'),
    updated_at = NOW()
WHERE title ILIKE '%священник%'
   OR body ILIKE '%священник%';
