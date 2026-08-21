-- +goose Up
CREATE UNIQUE INDEX idx_outbox_events_pending_dedup
    ON outbox_events (event_type, entity_id)
    WHERE status = 'pending';

-- +goose Down
DROP INDEX IF EXISTS idx_outbox_events_pending_dedup;
