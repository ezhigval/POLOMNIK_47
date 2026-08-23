-- +goose Up
ALTER TABLE reviews
    ADD COLUMN IF NOT EXISTS allow_distribution BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX IF NOT EXISTS idx_reviews_allow_distribution ON reviews (allow_distribution)
    WHERE allow_distribution = TRUE;

ALTER TABLE consents DROP CONSTRAINT IF EXISTS consents_type_check;
ALTER TABLE consents ADD CONSTRAINT consents_type_check CHECK (
    consent_type IN (
        'personal_data',
        'marketing',
        'marketing_revoked',
        'distribution',
        'distribution_revoked',
        'cookie_all',
        'cookie_essential',
        'cookie_reject'
    )
);

-- +goose Down
ALTER TABLE consents DROP CONSTRAINT IF EXISTS consents_type_check;
ALTER TABLE consents ADD CONSTRAINT consents_type_check CHECK (
    consent_type IN (
        'personal_data',
        'marketing',
        'distribution',
        'cookie_all',
        'cookie_essential',
        'cookie_reject'
    )
);
DROP INDEX IF EXISTS idx_reviews_allow_distribution;
ALTER TABLE reviews DROP COLUMN IF EXISTS allow_distribution;
