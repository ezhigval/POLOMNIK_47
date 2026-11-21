-- +goose Up
-- Portraits published with the service directory at tikhvin-eparhia.ru.
-- Update matching members in place, preserving editor changes, visibility and order.
-- +goose StatementBegin
WITH portraits(name, photo) AS (VALUES
    ('Протоиерей Александр Ежов', '/images/team/ezhov.jpg'),
    ('Людмила Викторовна Богуцкая', '/images/team/bogutskaya.jpg'),
    ('Светлана Юрьевна Семагина', '/images/team/semagina.jpg'),
    ('Лариса Арнольдовна Юревичиене', '/images/team/yurevichiene.jpg')
), members AS (
    SELECT b.id, jsonb_agg(
        CASE WHEN p.photo IS NOT NULL AND NOT (m.val ? 'photo')
             THEN m.val || jsonb_build_object('photo', p.photo) ELSE m.val END
        ORDER BY m.pos
    ) AS updated
    FROM cms_blocks b JOIN cms_pages h ON h.id=b.page_id AND h.slug='home'
    CROSS JOIN LATERAL jsonb_array_elements(CASE WHEN jsonb_typeof(b.content->'members')='array'
        THEN b.content->'members' ELSE '[]'::jsonb END) WITH ORDINALITY m(val,pos)
    LEFT JOIN portraits p ON p.name=m.val->>'name'
    WHERE b.type='service_team'
    GROUP BY b.id
)
UPDATE cms_blocks b SET content=jsonb_set(content,'{members}',members.updated), updated_at=now()
FROM members WHERE b.id=members.id AND content->'members' IS DISTINCT FROM members.updated;

WITH addresses(old_address,new_address) AS (VALUES
    ('Тихвин, ул. Тихвинская, 1', 'Ленинградская обл., г. Тихвин, ул. Тихвинская, д. 1'),
    ('Лодейнопольский район, деревня Старая Слобода', 'Ленинградская обл., Лодейнопольский р-н, д. Старая Слобода'),
    ('Санкт-Петербург, ул. Боровая, 52, лит. А', '192007, г. Санкт-Петербург, ул. Боровая, д. 52, лит. А')
), places AS (
    SELECT b.id, jsonb_agg(CASE WHEN a.new_address IS NULL THEN i.val
        ELSE i.val || jsonb_build_object('address', a.new_address) END ORDER BY i.pos) AS updated
    FROM cms_blocks b JOIN cms_pages h ON h.id=b.page_id AND h.slug='home'
    CROSS JOIN LATERAL jsonb_array_elements(CASE WHEN jsonb_typeof(b.content->'items')='array'
        THEN b.content->'items' ELSE '[]'::jsonb END) WITH ORDINALITY i(val,pos)
    LEFT JOIN addresses a ON a.old_address=i.val->>'address'
    WHERE b.type='holy_places'
    GROUP BY b.id
)
UPDATE cms_blocks b SET content=jsonb_set(content,'{items}',places.updated), updated_at=now()
FROM places WHERE b.id=places.id AND content->'items' IS DISTINCT FROM places.updated;

-- Official visitor directions: https://dymskij.ru/roadhelp, physical address: /adress.
WITH guide AS (
    SELECT b.id FROM cms_blocks b JOIN cms_pages h ON h.id=b.page_id
    WHERE h.slug='home' AND b.type='holy_places'
      AND jsonb_typeof(b.content->'items')='array'
    ORDER BY b.sort_order, b.id LIMIT 1
)
UPDATE cms_blocks b SET content=jsonb_set(content, '{items}', (content->'items') || $place$[{
    "name":"Антониево-Дымский мужской монастырь",
    "region":"Дымское озеро · Бокситогорский район",
    "description":"В Казанском храме хранятся мощи преподобного Антония Дымского. На берегу Дымского озера стоит часовня в память о его уединённой молитве.",
    "directions":"Из Санкт-Петербурга подойдёт автобус в сторону Бокситогорска или Пикалёва. Выходите на остановке «Галично». Оставшиеся три километра до обители проходят по асфальтированной дороге.",
    "address":"Ленинградская обл., Бокситогорский р-н, д. Красный Броневик",
    "latitude":"59.575067","longitude":"33.683767",
    "phone":"+7 (812) 603-79-21","phoneLabel":"Свечная лавка",
    "websiteUrl":"https://dymskij.ru/",
    "bookingHref":"/tours/pod-pokrovom-bogoroditsy","bookingLabel":"Поездка «Под покровом Богородицы»"
}]$place$::jsonb), updated_at=now()
FROM guide WHERE b.id=guide.id AND NOT EXISTS (
    SELECT FROM jsonb_array_elements(content->'items') i WHERE i->>'name'='Антониево-Дымский мужской монастырь'
);
-- +goose StatementEnd

-- +goose Down
-- Content can be edited after publication.
SELECT 1;
