-- +goose Up
ALTER TABLE legal_documents DROP CONSTRAINT IF EXISTS legal_documents_type_check;
ALTER TABLE legal_documents ADD CONSTRAINT legal_documents_type_check CHECK (
    type IN (
        'privacy_policy',
        'personal_data',
        'distribution',
        'marketing',
        'cookie',
        'terms',
        'offer'
    )
);

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
        'cookie_reject',
        'terms'
    )
);

CREATE TABLE IF NOT EXISTS user_photos (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    caption TEXT NOT NULL DEFAULT '',
    allow_distribution BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_user_photos_user_id ON user_photos (user_id);
CREATE INDEX IF NOT EXISTS idx_user_photos_distribution ON user_photos (allow_distribution)
    WHERE allow_distribution = TRUE;

-- +goose Down
DROP TABLE IF EXISTS user_photos;
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
ALTER TABLE legal_documents DROP CONSTRAINT IF EXISTS legal_documents_type_check;
ALTER TABLE legal_documents ADD CONSTRAINT legal_documents_type_check CHECK (
    type IN (
        'privacy_policy',
        'personal_data',
        'distribution',
        'marketing',
        'cookie'
    )
);
