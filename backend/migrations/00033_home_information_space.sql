-- +goose Up
ALTER TABLE cms_blocks DROP CONSTRAINT cms_blocks_type_check;
ALTER TABLE cms_blocks ADD CONSTRAINT cms_blocks_type_check CHECK (
    type IN (
        'hero', 'about', 'why_us', 'how_it_works', 'faq', 'cta', 'rich_text',
        'popular_destinations', 'testimonials', 'featured_route',
        'information_hero', 'holy_places', 'service_profile', 'service_team',
        'legal_donation', 'labor_trips'
    )
);

-- +goose StatementBegin
INSERT INTO cms_pages (id, slug, title, path, is_published)
SELECT 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'home', 'Главная', '/', TRUE
WHERE NOT EXISTS (SELECT FROM cms_pages WHERE slug='home' OR path='/');

INSERT INTO cms_blocks (id, page_id, type, sort_order, content, is_visible)
SELECT gen_random_uuid(), p.id, spec.type, spec.pos, spec.content, TRUE
FROM cms_pages p CROSS JOIN (VALUES
    ('information_hero', 100, $content${
  "eyebrow": "Тихвинская епархия",
  "title": "Путеводитель паломника",
  "subtitle": "От Успенского собора в Тихвине до обителей Присвирья. Познакомьтесь со святынями и спланируйте дорогу самостоятельно или вместе с нашей службой.",
  "primaryCta": "Монастыри и святыни",
  "primaryHref": "#holy-places",
  "secondaryCta": "О нашей службе",
  "secondaryHref": "#service"
}$content$::jsonb),
    ('holy_places', 110, $content${
  "eyebrow": "Путеводитель",
  "title": "Монастыри и святыни",
  "description": "Выберите обитель, чтобы найти святыню, проложить маршрут и связаться с принимающей стороной.",
  "items": [
    {
      "name": "Тихвинский Успенский монастырь",
      "region": "Тихвин",
      "description": "В Успенском соборе хранится Тихвинская икона Божией Матери. В 2004 году образ вернулся в обитель после десятилетий за пределами России.",
      "address": "Тихвин, ул. Тихвинская, 1",
      "directions": "Для самостоятельной поездки постройте маршрут до монастыря по адресу ниже. В программе «Тихвинского пути» посещение собора входит в первый день.",
      "phone": "+7 (81367) 5-10-84",
      "email": "info@tihvinskii-monastyr.ru",
      "websiteUrl": "http://tihvinskii-monastyr.ru/",
      "mapUrl": "https://yandex.ru/maps/?text=%D0%A2%D0%B8%D1%85%D0%B2%D0%B8%D0%BD%D1%81%D0%BA%D0%B8%D0%B9%20%D0%A3%D1%81%D0%BF%D0%B5%D0%BD%D1%81%D0%BA%D0%B8%D0%B9%20%D0%BC%D0%BE%D0%BD%D0%B0%D1%81%D1%82%D1%8B%D1%80%D1%8C%20%D0%A2%D0%B8%D1%85%D0%B2%D0%B8%D0%BD%2C%20%D1%83%D0%BB.%20%D0%A2%D0%B8%D1%85%D0%B2%D0%B8%D0%BD%D1%81%D0%BA%D0%B0%D1%8F%2C%201",
      "bookingHref": "/tours/tikhvin-path",
      "bookingLabel": "Поездка «Тихвинский путь»",
      "latitude": "",
      "longitude": ""
    },
    {
      "name": "Александро-Свирский монастырь",
      "region": "Присвирье",
      "description": "В обители паломники поклоняются мощам преподобного Александра Свирского. Об экскурсии и посещении трапезной можно договориться с монастырской паломнической службой.",
      "address": "Лодейнопольский район, деревня Старая Слобода",
      "directions": "Поездом до Лодейного Поля, затем автобусом в сторону Свирского или на такси до монастыря.",
      "phone": "+7 (81364) 44-402",
      "email": "palomnik@svirskoe.org",
      "websiteUrl": "https://svirskoe.org/",
      "mapUrl": "https://yandex.ru/maps/?text=%D0%90%D0%BB%D0%B5%D0%BA%D1%81%D0%B0%D0%BD%D0%B4%D1%80%D0%BE-%D0%A1%D0%B2%D0%B8%D1%80%D1%81%D0%BA%D0%B8%D0%B9%20%D0%BC%D0%BE%D0%BD%D0%B0%D1%81%D1%82%D1%8B%D1%80%D1%8C%20%D0%9B%D0%BE%D0%B4%D0%B5%D0%B9%D0%BD%D0%BE%D0%BF%D0%BE%D0%BB%D1%8C%D1%81%D0%BA%D0%B8%D0%B9%20%D1%80%D0%B0%D0%B9%D0%BE%D0%BD%2C%20%D0%B4%D0%B5%D1%80%D0%B5%D0%B2%D0%BD%D1%8F%20%D0%A1%D1%82%D0%B0%D1%80%D0%B0%D1%8F%20%D0%A1%D0%BB%D0%BE%D0%B1%D0%BE%D0%B4%D0%B0",
      "bookingHref": "/tours/tikhvin-path",
      "bookingLabel": "Поездка «Тихвинский путь»",
      "latitude": "",
      "longitude": ""
    },
    {
      "name": "Храм Покрова на Боровой",
      "region": "Санкт-Петербург · Подворье",
      "description": "Петербургское подворье Антониево-Дымского монастыря. В нижнем храме престол посвящён мученику Вонифатию. Время ближайшей службы можно уточнить на сайте храма.",
      "address": "Санкт-Петербург, ул. Боровая, 52, лит. А",
      "directions": "От станции метро «Обводный канал» до храма можно дойти пешком.",
      "phone": "+7 (962) 701-94-07",
      "email": "hram1897pokrova@yandex.ru",
      "websiteUrl": "https://church-pokrov.ru/",
      "mapUrl": "https://yandex.ru/maps/?text=%D0%A5%D1%80%D0%B0%D0%BC%20%D0%9F%D0%BE%D0%BA%D1%80%D0%BE%D0%B2%D0%B0%20%D0%BD%D0%B0%20%D0%91%D0%BE%D1%80%D0%BE%D0%B2%D0%BE%D0%B9%20%D0%A1%D0%B0%D0%BD%D0%BA%D1%82-%D0%9F%D0%B5%D1%82%D0%B5%D1%80%D0%B1%D1%83%D1%80%D0%B3%2C%20%D1%83%D0%BB.%20%D0%91%D0%BE%D1%80%D0%BE%D0%B2%D0%B0%D1%8F%2C%2052%2C%20%D0%BB%D0%B8%D1%82.%20%D0%90",
      "bookingHref": "",
      "bookingLabel": "Поездка «Тихвинский путь»",
      "latitude": "",
      "longitude": ""
    }
  ]
}$content$::jsonb),
    ('service_profile', 120, $content${
  "eyebrow": "О службе",
  "title": "Паломническая служба Тихвинской епархии",
  "paragraphs": [
    "Наша служба носит имя «Под Покровом Божией Матери „Тихвинская“». Мы организуем паломнические поездки по монастырям и храмам, в том числе по Старой Ладоге, Тихвину и Присвирью.",
    "Если вы собираетесь в путь впервые, свяжитесь со службой. Здесь можно уточнить программу, место встречи и как подготовиться к посещению конкретной обители."
  ],
  "links": [
    {
      "label": "Связаться со службой",
      "href": "mailto:info@tikhvin-palomnik.ru"
    },
    {
      "label": "Тихвинская епархия",
      "href": "https://www.tikhvin-eparhia.ru/"
    },
    {
      "label": "Подворье на Боровой",
      "href": "https://church-pokrov.ru/"
    }
  ]
}$content$::jsonb),
    ('service_team', 130, $content${
  "eyebrow": "Люди службы",
  "title": "Люди, к которым можно обратиться",
  "description": "Общий телефон паломнической службы +7 (966) 933-43-21.",
  "members": [
    {
      "name": "Протоиерей Александр Ежов",
      "role": "Директор",
      "bio": "",
      "phone": "",
      "email": ""
    },
    {
      "name": "Людмила Викторовна Богуцкая",
      "role": "Заместитель директора",
      "bio": "",
      "phone": "",
      "email": ""
    },
    {
      "name": "Светлана Юрьевна Семагина",
      "role": "Сотрудник службы",
      "bio": "",
      "phone": "",
      "email": ""
    },
    {
      "name": "Лариса Арнольдовна Юревичиене",
      "role": "Методист",
      "bio": "",
      "phone": "",
      "email": ""
    }
  ]
}$content$::jsonb),
    ('legal_donation', 140, $content${
  "eyebrow": "Официально",
  "title": "Документы и реквизиты",
  "legalName": "",
  "legalAddress": "",
  "inn": "",
  "kpp": "",
  "ogrn": "",
  "donationRecipient": "",
  "donationBankName": "",
  "donationAccount": "",
  "donationCorrespondentAccount": "",
  "donationBik": "",
  "donationPurpose": "",
  "donationUrl": ""
}$content$::jsonb),
    ('labor_trips', 150, $content${
  "eyebrow": "Помощь обителям",
  "title": "Трудовые поездки",
  "description": "В такой поездке паломник помогает монастырю в повседневных делах. Перед выездом уточните у службы, где сейчас нужна помощь, какие послушания предстоят, как устроены проживание и питание.",
  "items": [
    {
      "name": "Как присоединиться к трудовой поездке",
      "place": "",
      "dates": "",
      "description": "Напишите, когда вы можете поехать и какую работу готовы выполнять. Служба подскажет, с кем согласовать участие.",
      "bookingHref": "mailto:info@tikhvin-palomnik.ru"
    }
  ]
}$content$::jsonb)
) AS spec(type, pos, content)
WHERE p.slug='home' AND NOT EXISTS (
    SELECT FROM cms_blocks b WHERE b.page_id=p.id AND b.type=spec.type
);
-- +goose StatementEnd

-- +goose Down
-- Preserve editorial changes when rolling back the application.
SELECT 1;
