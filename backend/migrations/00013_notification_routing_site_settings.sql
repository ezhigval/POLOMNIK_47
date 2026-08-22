-- +goose Up
-- Notification routing (event → channel+address) and evolve historical site_settings.
-- Prod already has site_settings from an older shape (site_full_name / site_tagline / …).

CREATE TABLE IF NOT EXISTS notification_routing (
    id SMALLINT PRIMARY KEY CHECK (id = 1),
    routes JSONB NOT NULL DEFAULT '{}'::jsonb,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO notification_routing (id, routes)
VALUES (1, '{}'::jsonb)
ON CONFLICT (id) DO NOTHING;

CREATE TABLE IF NOT EXISTS site_settings (
    id SMALLINT PRIMARY KEY CHECK (id = 1),
    site_name TEXT NOT NULL DEFAULT '',
    full_name TEXT NOT NULL DEFAULT '',
    tagline TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    region TEXT NOT NULL DEFAULT '',
    departure_city TEXT NOT NULL DEFAULT '',
    parent_org_name TEXT NOT NULL DEFAULT '',
    parent_org_url TEXT NOT NULL DEFAULT '',
    contact_phone TEXT NOT NULL DEFAULT '',
    contact_phone_display TEXT NOT NULL DEFAULT '',
    contact_email TEXT NOT NULL DEFAULT '',
    mail_forward_to TEXT NOT NULL DEFAULT '',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Rename historical columns to the canonical names used by the API.
-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'site_settings' AND column_name = 'site_full_name'
    ) AND NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'site_settings' AND column_name = 'full_name'
    ) THEN
        ALTER TABLE site_settings RENAME COLUMN site_full_name TO full_name;
    END IF;

    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'site_settings' AND column_name = 'site_tagline'
    ) AND NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'site_settings' AND column_name = 'tagline'
    ) THEN
        ALTER TABLE site_settings RENAME COLUMN site_tagline TO tagline;
    END IF;

    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'site_settings' AND column_name = 'site_description'
    ) AND NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'site_settings' AND column_name = 'description'
    ) THEN
        ALTER TABLE site_settings RENAME COLUMN site_description TO description;
    END IF;
END $$;
-- +goose StatementEnd

ALTER TABLE site_settings ADD COLUMN IF NOT EXISTS site_name TEXT NOT NULL DEFAULT '';
ALTER TABLE site_settings ADD COLUMN IF NOT EXISTS full_name TEXT NOT NULL DEFAULT '';
ALTER TABLE site_settings ADD COLUMN IF NOT EXISTS tagline TEXT NOT NULL DEFAULT '';
ALTER TABLE site_settings ADD COLUMN IF NOT EXISTS description TEXT NOT NULL DEFAULT '';
ALTER TABLE site_settings ADD COLUMN IF NOT EXISTS region TEXT NOT NULL DEFAULT '';
ALTER TABLE site_settings ADD COLUMN IF NOT EXISTS departure_city TEXT NOT NULL DEFAULT '';
ALTER TABLE site_settings ADD COLUMN IF NOT EXISTS parent_org_name TEXT NOT NULL DEFAULT '';
ALTER TABLE site_settings ADD COLUMN IF NOT EXISTS parent_org_url TEXT NOT NULL DEFAULT '';
ALTER TABLE site_settings ADD COLUMN IF NOT EXISTS contact_phone TEXT NOT NULL DEFAULT '';
ALTER TABLE site_settings ADD COLUMN IF NOT EXISTS contact_phone_display TEXT NOT NULL DEFAULT '';
ALTER TABLE site_settings ADD COLUMN IF NOT EXISTS contact_email TEXT NOT NULL DEFAULT '';
ALTER TABLE site_settings ADD COLUMN IF NOT EXISTS mail_forward_to TEXT NOT NULL DEFAULT '';
ALTER TABLE site_settings ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

-- Seed only when the table is empty (never wipe existing identity row).
INSERT INTO site_settings (
    id, site_name, full_name, tagline, description, region, departure_city,
    parent_org_name, parent_org_url, contact_phone, contact_phone_display,
    contact_email, mail_forward_to
)
SELECT
    1, '', '', '', '', '', '',
    '', '', '', '',
    '', ''
WHERE NOT EXISTS (SELECT 1 FROM site_settings WHERE id = 1);

-- One-shot migrate from legacy telegram username lists into JSON routes.
UPDATE notification_routing nr
SET
    routes = subquery.routes,
    updated_at = NOW()
FROM (
    SELECT
        jsonb_build_object(
            'booking_created', COALESCE((
                SELECT jsonb_agg(jsonb_build_object('channel', 'telegram', 'address', username) ORDER BY ord)
                FROM (
                    SELECT DISTINCT ON (lower(trim(both '@' FROM trim(token))))
                        lower(trim(both '@' FROM trim(token))) AS username,
                        ord
                    FROM (
                        SELECT trim(token) AS token, ord
                        FROM telegram_recipients,
                             LATERAL unnest(string_to_array(replace(replace(booking_usernames, E'\n', ','), ';', ','), ','))
                                 WITH ORDINALITY AS t(token, ord)
                        WHERE id = 1
                    ) raw
                    WHERE trim(token) <> ''
                    ORDER BY lower(trim(both '@' FROM trim(token))), ord
                ) usernames
            ), '[]'::jsonb),
            'booking_status_changed', COALESCE((
                SELECT jsonb_agg(jsonb_build_object('channel', 'telegram', 'address', username) ORDER BY ord)
                FROM (
                    SELECT DISTINCT ON (lower(trim(both '@' FROM trim(token))))
                        lower(trim(both '@' FROM trim(token))) AS username,
                        ord
                    FROM (
                        SELECT trim(token) AS token, ord
                        FROM telegram_recipients,
                             LATERAL unnest(string_to_array(replace(replace(booking_usernames, E'\n', ','), ';', ','), ','))
                                 WITH ORDINALITY AS t(token, ord)
                        WHERE id = 1
                    ) raw
                    WHERE trim(token) <> ''
                    ORDER BY lower(trim(both '@' FROM trim(token))), ord
                ) usernames
            ), '[]'::jsonb),
            'support_message', COALESCE((
                SELECT jsonb_agg(jsonb_build_object('channel', 'telegram', 'address', username) ORDER BY ord)
                FROM (
                    SELECT DISTINCT ON (lower(trim(both '@' FROM trim(token))))
                        lower(trim(both '@' FROM trim(token))) AS username,
                        ord
                    FROM (
                        SELECT trim(token) AS token, ord
                        FROM telegram_recipients,
                             LATERAL unnest(string_to_array(replace(replace(support_usernames, E'\n', ','), ';', ','), ','))
                                 WITH ORDINALITY AS t(token, ord)
                        WHERE id = 1
                    ) raw
                    WHERE trim(token) <> ''
                    ORDER BY lower(trim(both '@' FROM trim(token))), ord
                ) usernames
            ), '[]'::jsonb)
        ) AS routes
    FROM telegram_recipients
    WHERE id = 1
) AS subquery
WHERE nr.id = 1
  AND (
      nr.routes = '{}'::jsonb
      OR nr.routes = '{"booking_created":[],"booking_status_changed":[],"support_message":[]}'::jsonb
      OR NOT (nr.routes ? 'booking_created')
  )
  AND EXISTS (SELECT 1 FROM telegram_recipients WHERE id = 1);

-- +goose Down
DROP TABLE IF EXISTS notification_routing;
-- Do not drop site_settings on down: may pre-exist from historical prod migrations.
