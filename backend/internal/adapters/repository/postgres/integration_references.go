package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

func (s *Store) UpsertReference(ctx context.Context, ref domain.IntegrationReference) (domain.IntegrationReference, error) {
	var lastSyncAt sql.NullTime
	if ref.LastSyncAt != nil {
		lastSyncAt = sql.NullTime{Time: ref.LastSyncAt.UTC(), Valid: true}
	}

	row := s.db.QueryRowContext(ctx, `
INSERT INTO integration_references (
    id, local_entity_type, local_entity_id, external_system, external_entity_type,
    external_entity_id, sync_status, last_sync_at, last_error, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
ON CONFLICT (local_entity_type, local_entity_id, external_system, external_entity_type)
DO UPDATE SET
    external_entity_id = EXCLUDED.external_entity_id,
    sync_status = EXCLUDED.sync_status,
    last_sync_at = EXCLUDED.last_sync_at,
    last_error = EXCLUDED.last_error,
    updated_at = EXCLUDED.updated_at
RETURNING id, local_entity_type, local_entity_id, external_system, external_entity_type,
          external_entity_id, sync_status, last_sync_at, last_error, created_at, updated_at
`, ref.ID, ref.LocalEntityType, ref.LocalEntityID, ref.ExternalSystem, ref.ExternalEntityType,
		ref.ExternalEntityID, ref.SyncStatus, lastSyncAt, ref.LastError, ref.CreatedAt, ref.UpdatedAt)

	return scanIntegrationReference(row)
}

func (s *Store) ListReferences(ctx context.Context, filters ports.IntegrationReferenceFilters, pagination ports.Pagination) (ports.IntegrationReferenceList, error) {
	pagination = ports.NormalizePagination(pagination.Page, pagination.Limit)
	args := integrationReferenceFilterArgs(filters)

	total, err := count(ctx, s.db, "SELECT COUNT(*) FROM integration_references WHERE "+integrationReferenceWhereClause, args...)
	if err != nil {
		return ports.IntegrationReferenceList{}, err
	}

	rows, err := s.db.QueryContext(ctx, `
SELECT id, local_entity_type, local_entity_id, external_system, external_entity_type,
       external_entity_id, sync_status, last_sync_at, last_error, created_at, updated_at
FROM integration_references
WHERE `+integrationReferenceWhereClause+`
ORDER BY updated_at DESC, id ASC
LIMIT $4 OFFSET $5
`, append(args, pagination.Limit, offset(pagination))...)
	if err != nil {
		return ports.IntegrationReferenceList{}, fmt.Errorf("list integration references: %w", err)
	}
	defer rows.Close()

	items, err := scanIntegrationReferences(rows)
	if err != nil {
		return ports.IntegrationReferenceList{}, err
	}

	return ports.IntegrationReferenceList{Items: items, Meta: pageMeta(pagination, total)}, nil
}

const integrationReferenceWhereClause = `
($1::text = '' OR external_system = $1) AND
($2::text = '' OR local_entity_type = $2) AND
($3::text = '' OR sync_status = $3)
`

func integrationReferenceFilterArgs(filters ports.IntegrationReferenceFilters) []any {
	return []any{
		filters.ExternalSystem,
		filters.LocalEntityType,
		filters.SyncStatus,
	}
}

func scanIntegrationReferences(rows *sql.Rows) ([]domain.IntegrationReference, error) {
	var items []domain.IntegrationReference
	for rows.Next() {
		ref, err := scanIntegrationReference(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, ref)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate integration references: %w", err)
	}
	return items, nil
}

func (s *Store) GetReference(ctx context.Context, query ports.IntegrationReferenceQuery) (domain.IntegrationReference, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT id, local_entity_type, local_entity_id, external_system, external_entity_type,
       external_entity_id, sync_status, last_sync_at, last_error, created_at, updated_at
FROM integration_references
WHERE local_entity_type = $1
  AND local_entity_id = $2
  AND external_system = $3
  AND external_entity_type = $4
`, query.LocalEntityType, query.LocalEntityID, query.ExternalSystem, query.ExternalEntityType)
	return scanIntegrationReference(row)
}

type integrationReferenceScanner interface {
	Scan(dest ...any) error
}

func scanIntegrationReference(row integrationReferenceScanner) (domain.IntegrationReference, error) {
	var ref domain.IntegrationReference
	var lastSyncAt sql.NullTime
	err := row.Scan(
		&ref.ID,
		&ref.LocalEntityType,
		&ref.LocalEntityID,
		&ref.ExternalSystem,
		&ref.ExternalEntityType,
		&ref.ExternalEntityID,
		&ref.SyncStatus,
		&lastSyncAt,
		&ref.LastError,
		&ref.CreatedAt,
		&ref.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.IntegrationReference{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.IntegrationReference{}, fmt.Errorf("scan integration reference: %w", err)
	}
	if lastSyncAt.Valid {
		syncAt := lastSyncAt.Time.UTC()
		ref.LastSyncAt = &syncAt
	}
	return ref, nil
}
