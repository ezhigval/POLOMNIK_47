-- +goose Up
CREATE TABLE outbox_events (
    id UUID PRIMARY KEY,
    event_type TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id UUID NOT NULL,
    payload JSONB NOT NULL DEFAULT '{}',
    status TEXT NOT NULL,
    attempts INT NOT NULL DEFAULT 0,
    last_error TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT outbox_events_status_check CHECK (
        status IN ('pending', 'processed', 'failed')
    )
);

CREATE INDEX idx_outbox_events_status ON outbox_events(status);
CREATE INDEX idx_outbox_events_entity ON outbox_events(entity_type, entity_id);
CREATE INDEX idx_outbox_events_created_at ON outbox_events(created_at);

-- +goose Down
DROP TABLE IF EXISTS outbox_events;
