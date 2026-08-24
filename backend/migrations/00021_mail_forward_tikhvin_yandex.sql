-- +goose Up
-- Add the palomnik Yandex inbox to mailbox-forward recipients without wiping
-- addresses already saved in admin Settings.
UPDATE site_settings
SET
    mail_forward_to = CASE
        WHEN trim(both FROM coalesce(mail_forward_to, '')) = '' THEN 'tikhvin-palomnik@yandex.ru'
        WHEN ',' || regexp_replace(lower(mail_forward_to), '[\s;]+', ',', 'g') || ','
            LIKE '%,tikhvin-palomnik@yandex.ru,%' THEN mail_forward_to
        ELSE trim(both FROM mail_forward_to) || E'\ntikhvin-palomnik@yandex.ru'
    END,
    updated_at = NOW()
WHERE id = 1;

INSERT INTO site_settings (id, mail_forward_to, updated_at)
SELECT 1, 'tikhvin-palomnik@yandex.ru', NOW()
WHERE NOT EXISTS (SELECT 1 FROM site_settings WHERE id = 1);

-- +goose Down
UPDATE site_settings
SET
    mail_forward_to = trim(both FROM regexp_replace(
        regexp_replace(
            coalesce(mail_forward_to, ''),
            '(^|[\s,;]+)tikhvin-palomnik@yandex\.ru($|[\s,;]+)',
            '\1',
            'gi'
        ),
        '[\s,;]+$',
        ''
    )),
    updated_at = NOW()
WHERE id = 1;
