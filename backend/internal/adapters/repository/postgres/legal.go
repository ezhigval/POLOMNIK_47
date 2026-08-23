package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

func (s *Store) HasAnyDocuments(ctx context.Context) (bool, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM legal_documents`).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("count legal documents: %w", err)
	}
	return count > 0, nil
}

func (s *Store) ListDocuments(ctx context.Context, filters ports.LegalDocumentFilters) ([]domain.LegalDocument, error) {
	activeOnly := filters.ActiveOnly && !filters.IncludeAll
	typeFilter := string(filters.Type)

	rows, err := s.db.QueryContext(ctx, `
SELECT id, type, version, title, content, published_at, updated_at, is_active
FROM legal_documents
WHERE ($1::text = '' OR type = $1)
  AND ($2::boolean = FALSE OR is_active = TRUE)
ORDER BY type ASC, published_at DESC, version DESC
`, typeFilter, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("list legal documents: %w", err)
	}
	defer rows.Close()

	return scanLegalDocuments(rows)
}

func (s *Store) GetDocument(ctx context.Context, id uuid.UUID) (domain.LegalDocument, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT id, type, version, title, content, published_at, updated_at, is_active
FROM legal_documents WHERE id = $1
`, id)
	return scanLegalDocument(row)
}

func (s *Store) GetActiveByType(ctx context.Context, docType domain.LegalDocumentType) (domain.LegalDocument, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT id, type, version, title, content, published_at, updated_at, is_active
FROM legal_documents WHERE type = $1 AND is_active = TRUE
LIMIT 1
`, string(docType))
	return scanLegalDocument(row)
}

func (s *Store) GetByTypeAndVersion(ctx context.Context, docType domain.LegalDocumentType, version string) (domain.LegalDocument, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT id, type, version, title, content, published_at, updated_at, is_active
FROM legal_documents WHERE type = $1 AND version = $2
`, string(docType), version)
	return scanLegalDocument(row)
}

func (s *Store) CreateDocument(ctx context.Context, doc domain.LegalDocument) (domain.LegalDocument, error) {
	row := s.db.QueryRowContext(ctx, `
INSERT INTO legal_documents (id, type, version, title, content, published_at, updated_at, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, type, version, title, content, published_at, updated_at, is_active
`, doc.ID, string(doc.Type), doc.Version, doc.Title, doc.Content, doc.PublishedAt, doc.UpdatedAt, doc.IsActive)
	created, err := scanLegalDocument(row)
	if err != nil {
		return domain.LegalDocument{}, mapLegalDocumentWriteError(err)
	}
	return created, nil
}

func (s *Store) DeactivateType(ctx context.Context, docType domain.LegalDocumentType) error {
	_, err := s.db.ExecContext(ctx, `
UPDATE legal_documents SET is_active = FALSE, updated_at = NOW() AT TIME ZONE 'utc'
WHERE type = $1 AND is_active = TRUE
`, string(docType))
	if err != nil {
		return fmt.Errorf("deactivate legal documents: %w", err)
	}
	return nil
}

func (s *Store) CreateConsent(ctx context.Context, consent domain.Consent) (domain.Consent, error) {
	row := s.db.QueryRowContext(ctx, `
INSERT INTO consents (id, user_id, request_id, consent_type, document_id, document_version, accepted_at, ip, user_agent)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, user_id, request_id, consent_type, document_id, document_version, accepted_at, ip, user_agent
`, consent.ID, nullableUUID(consent.UserID), nullableUUID(consent.RequestID),
		string(consent.ConsentType), consent.DocumentID, consent.DocumentVersion,
		consent.AcceptedAt, consent.IP, consent.UserAgent)
	return scanConsent(row)
}

func (s *Store) ListConsents(ctx context.Context, filters ports.ConsentFilters, pagination ports.Pagination) (ports.ConsentList, error) {
	pagination = ports.NormalizePagination(pagination.Page, pagination.Limit)

	total, err := count(ctx, s.db, `
SELECT COUNT(*) FROM consents
WHERE ($1::uuid IS NULL OR user_id = $1)
  AND ($2::uuid IS NULL OR request_id = $2)
  AND ($3::text = '' OR consent_type = $3)
`, nullableUUID(filters.UserID), nullableUUID(filters.RequestID), string(filters.ConsentType))
	if err != nil {
		return ports.ConsentList{}, err
	}

	rows, err := s.db.QueryContext(ctx, `
SELECT id, user_id, request_id, consent_type, document_id, document_version, accepted_at, ip, user_agent
FROM consents
WHERE ($1::uuid IS NULL OR user_id = $1)
  AND ($2::uuid IS NULL OR request_id = $2)
  AND ($3::text = '' OR consent_type = $3)
ORDER BY accepted_at DESC
LIMIT $4 OFFSET $5
`, nullableUUID(filters.UserID), nullableUUID(filters.RequestID), string(filters.ConsentType),
		pagination.Limit, offset(pagination))
	if err != nil {
		return ports.ConsentList{}, fmt.Errorf("list consents: %w", err)
	}
	defer rows.Close()

	items, err := scanConsents(rows)
	if err != nil {
		return ports.ConsentList{}, err
	}

	return ports.ConsentList{Items: items, Meta: pageMeta(pagination, total)}, nil
}

type legalDocumentScanner interface {
	Scan(dest ...any) error
}

func scanLegalDocument(row legalDocumentScanner) (domain.LegalDocument, error) {
	var doc domain.LegalDocument
	var docType string
	err := row.Scan(&doc.ID, &docType, &doc.Version, &doc.Title, &doc.Content,
		&doc.PublishedAt, &doc.UpdatedAt, &doc.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.LegalDocument{}, domain.ErrNotFound
		}
		return domain.LegalDocument{}, err
	}
	doc.Type = domain.LegalDocumentType(docType)
	return doc, nil
}

func scanLegalDocuments(rows *sql.Rows) ([]domain.LegalDocument, error) {
	var items []domain.LegalDocument
	for rows.Next() {
		var doc domain.LegalDocument
		var docType string
		if err := rows.Scan(&doc.ID, &docType, &doc.Version, &doc.Title, &doc.Content,
			&doc.PublishedAt, &doc.UpdatedAt, &doc.IsActive); err != nil {
			return nil, err
		}
		doc.Type = domain.LegalDocumentType(docType)
		items = append(items, doc)
	}
	return items, rows.Err()
}

func scanConsent(row legalDocumentScanner) (domain.Consent, error) {
	var c domain.Consent
	var consentType string
	var userID, requestID sql.NullString
	err := row.Scan(&c.ID, &userID, &requestID, &consentType, &c.DocumentID,
		&c.DocumentVersion, &c.AcceptedAt, &c.IP, &c.UserAgent)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Consent{}, domain.ErrNotFound
		}
		return domain.Consent{}, err
	}
	c.ConsentType = domain.ConsentType(consentType)
	if userID.Valid {
		id, err := uuid.Parse(userID.String)
		if err == nil {
			c.UserID = &id
		}
	}
	if requestID.Valid {
		id, err := uuid.Parse(requestID.String)
		if err == nil {
			c.RequestID = &id
		}
	}
	return c, nil
}

func scanConsents(rows *sql.Rows) ([]domain.Consent, error) {
	var items []domain.Consent
	for rows.Next() {
		var c domain.Consent
		var consentType string
		var userID, requestID sql.NullString
		if err := rows.Scan(&c.ID, &userID, &requestID, &consentType, &c.DocumentID,
			&c.DocumentVersion, &c.AcceptedAt, &c.IP, &c.UserAgent); err != nil {
			return nil, err
		}
		c.ConsentType = domain.ConsentType(consentType)
		if userID.Valid {
			id, err := uuid.Parse(userID.String)
			if err == nil {
				c.UserID = &id
			}
		}
		if requestID.Valid {
			id, err := uuid.Parse(requestID.String)
			if err == nil {
				c.RequestID = &id
			}
		}
		items = append(items, c)
	}
	return items, rows.Err()
}

func mapLegalDocumentWriteError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23505" {
		return domain.ErrDuplicateSlug
	}
	return err
}

func nullableUUID(id *uuid.UUID) any {
	if id == nil || *id == uuid.Nil {
		return nil
	}
	return *id
}
