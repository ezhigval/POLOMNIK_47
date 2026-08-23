-- +goose Up
CREATE TABLE legal_documents (
    id UUID PRIMARY KEY,
    type TEXT NOT NULL,
    version TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    published_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT legal_documents_type_version_unique UNIQUE (type, version),
    CONSTRAINT legal_documents_type_check CHECK (
        type IN ('privacy_policy', 'personal_data', 'distribution', 'marketing', 'cookie')
    )
);

CREATE UNIQUE INDEX idx_legal_documents_one_active_per_type
    ON legal_documents (type) WHERE is_active = TRUE;

CREATE INDEX idx_legal_documents_type ON legal_documents (type);
CREATE INDEX idx_legal_documents_published_at ON legal_documents (published_at DESC);

CREATE TABLE consents (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    request_id UUID REFERENCES bookings(id) ON DELETE SET NULL,
    consent_type TEXT NOT NULL,
    document_id UUID NOT NULL REFERENCES legal_documents(id) ON DELETE RESTRICT,
    document_version TEXT NOT NULL,
    accepted_at TIMESTAMPTZ NOT NULL,
    ip TEXT NOT NULL DEFAULT '',
    user_agent TEXT NOT NULL DEFAULT '',
    CONSTRAINT consents_type_check CHECK (
        consent_type IN (
            'personal_data',
            'marketing',
            'distribution',
            'cookie_all',
            'cookie_essential',
            'cookie_reject'
        )
    )
);

CREATE INDEX idx_consents_user_id ON consents (user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_consents_request_id ON consents (request_id) WHERE request_id IS NOT NULL;
CREATE INDEX idx_consents_type_accepted ON consents (consent_type, accepted_at DESC);
CREATE INDEX idx_consents_document_id ON consents (document_id);

-- +goose Down
DROP TABLE IF EXISTS consents;
DROP TABLE IF EXISTS legal_documents;
