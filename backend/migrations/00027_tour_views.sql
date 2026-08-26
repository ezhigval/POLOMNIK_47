-- +goose Up
CREATE TABLE IF NOT EXISTS tour_views (
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    tour_id UUID NOT NULL REFERENCES tours (id) ON DELETE CASCADE,
    viewed_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    PRIMARY KEY (user_id, tour_id)
);

CREATE INDEX IF NOT EXISTS tour_views_user_viewed_at_idx ON tour_views (user_id, viewed_at DESC);

-- +goose Down
DROP TABLE IF EXISTS tour_views;
