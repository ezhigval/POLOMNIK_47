-- +goose Up
CREATE TABLE tours (
    id UUID PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    price INTEGER NOT NULL CHECK (price >= 0),
    currency TEXT NOT NULL DEFAULT 'RUB',
    date_start DATE NOT NULL,
    date_end DATE NOT NULL,
    slots_total INTEGER NOT NULL CHECK (slots_total >= 0),
    slots_left INTEGER NOT NULL CHECK (slots_left >= 0),
    location TEXT NOT NULL DEFAULT '',
    images TEXT[] NOT NULL DEFAULT '{}',
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    is_hot BOOLEAN NOT NULL DEFAULT FALSE,
    overbooking_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT tours_date_range_check CHECK (date_start <= date_end),
    CONSTRAINT tours_slots_left_lte_total_check CHECK (slots_left <= slots_total)
);

CREATE TABLE bookings (
    id UUID PRIMARY KEY,
    tour_id UUID NOT NULL REFERENCES tours(id) ON DELETE RESTRICT,
    name TEXT NOT NULL,
    phone TEXT NOT NULL,
    email TEXT NOT NULL DEFAULT '',
    people_count INTEGER NOT NULL CHECK (people_count > 0),
    status TEXT NOT NULL,
    total_price INTEGER NOT NULL CHECK (total_price >= 0),
    comment TEXT NOT NULL DEFAULT '',
    overbooked BOOLEAN NOT NULL DEFAULT FALSE,
    source TEXT NOT NULL DEFAULT 'web',
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT bookings_status_check CHECK (
        status IN ('NEW', 'CONTACTED', 'CONFIRMED', 'COMPLETED', 'CANCELLED')
    )
);

CREATE TABLE reviews (
    id UUID PRIMARY KEY,
    tour_id UUID NOT NULL REFERENCES tours(id) ON DELETE CASCADE,
    client_name TEXT NOT NULL,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    text TEXT NOT NULL,
    is_approved BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE integration_references (
    id UUID PRIMARY KEY,
    local_entity_type TEXT NOT NULL,
    local_entity_id UUID NOT NULL,
    external_system TEXT NOT NULL,
    external_entity_type TEXT NOT NULL,
    external_entity_id TEXT NOT NULL,
    sync_status TEXT NOT NULL,
    last_sync_at TIMESTAMPTZ,
    last_error TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT integration_references_sync_status_check CHECK (
        sync_status IN ('not_configured', 'synced', 'pending', 'failed')
    ),
    CONSTRAINT integration_references_local_unique UNIQUE (
        local_entity_type,
        local_entity_id,
        external_system,
        external_entity_type
    ),
    CONSTRAINT integration_references_external_unique UNIQUE (
        external_system,
        external_entity_type,
        external_entity_id
    )
);

CREATE INDEX idx_tours_active ON tours(is_active);
CREATE INDEX idx_tours_hot ON tours(is_hot);
CREATE INDEX idx_tours_dates ON tours(date_start, date_end);
CREATE INDEX idx_bookings_tour ON bookings(tour_id);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_bookings_created_at ON bookings(created_at);
CREATE INDEX idx_reviews_tour ON reviews(tour_id);
CREATE INDEX idx_reviews_approved ON reviews(is_approved);
CREATE INDEX idx_integration_refs_local ON integration_references(local_entity_type, local_entity_id);
CREATE INDEX idx_integration_refs_external ON integration_references(external_system, external_entity_id);

-- +goose Down
DROP TABLE IF EXISTS integration_references;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS tours;

