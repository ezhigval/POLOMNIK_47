package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

func (s *Store) ListPages(ctx context.Context, filters ports.CMSPageFilters) ([]domain.Page, error) {
	query := `
SELECT id, slug, title, path, meta_title, meta_description, is_published, created_at, updated_at
FROM cms_pages`
	args := []any{}
	if filters.PublishedOnly {
		query += ` WHERE is_published = TRUE`
	}
	query += ` ORDER BY title ASC, created_at ASC`

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list cms pages: %w", err)
	}
	defer rows.Close()

	var pages []domain.Page
	for rows.Next() {
		page, err := scanCMSPage(rows)
		if err != nil {
			return nil, err
		}
		pages = append(pages, page)
	}
	return pages, rows.Err()
}

func (s *Store) GetPage(ctx context.Context, id uuid.UUID) (domain.Page, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT id, slug, title, path, meta_title, meta_description, is_published, created_at, updated_at
FROM cms_pages WHERE id = $1
`, id)
	return scanCMSPage(row)
}

func (s *Store) GetPageBySlug(ctx context.Context, slug string) (domain.Page, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT id, slug, title, path, meta_title, meta_description, is_published, created_at, updated_at
FROM cms_pages WHERE slug = $1
`, slug)
	return scanCMSPage(row)
}

func (s *Store) GetPageByPath(ctx context.Context, path string) (domain.Page, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT id, slug, title, path, meta_title, meta_description, is_published, created_at, updated_at
FROM cms_pages WHERE path = $1
`, path)
	return scanCMSPage(row)
}

func (s *Store) CreatePage(ctx context.Context, page domain.Page) (domain.Page, error) {
	_, err := s.db.ExecContext(ctx, `
INSERT INTO cms_pages (id, slug, title, path, meta_title, meta_description, is_published, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
`, page.ID, page.Slug, page.Title, page.Path, page.MetaTitle, page.MetaDescription, page.IsPublished, page.CreatedAt, page.UpdatedAt)
	if err != nil {
		return domain.Page{}, mapCMSWriteError(err)
	}
	return page, nil
}

func (s *Store) UpdatePage(ctx context.Context, page domain.Page) (domain.Page, error) {
	result, err := s.db.ExecContext(ctx, `
UPDATE cms_pages
SET title = $2, path = $3, meta_title = $4, meta_description = $5, is_published = $6, updated_at = $7
WHERE id = $1
`, page.ID, page.Title, page.Path, page.MetaTitle, page.MetaDescription, page.IsPublished, page.UpdatedAt)
	if err != nil {
		return domain.Page{}, mapCMSWriteError(err)
	}
	if err := requireAffected(result); err != nil {
		return domain.Page{}, err
	}
	return page, nil
}

func (s *Store) DeletePage(ctx context.Context, id uuid.UUID) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM cms_pages WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete cms page: %w", err)
	}
	return requireAffected(result)
}

func (s *Store) ListBlocks(ctx context.Context, pageID uuid.UUID) ([]domain.Block, error) {
	rows, err := s.db.QueryContext(ctx, `
SELECT id, page_id, type, sort_order, content, is_visible, created_at, updated_at
FROM cms_blocks
WHERE page_id = $1
ORDER BY sort_order ASC, created_at ASC
`, pageID)
	if err != nil {
		return nil, fmt.Errorf("list cms blocks: %w", err)
	}
	defer rows.Close()

	var blocks []domain.Block
	for rows.Next() {
		block, err := scanCMSBlock(rows)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}
	return blocks, rows.Err()
}

func (s *Store) GetBlock(ctx context.Context, id uuid.UUID) (domain.Block, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT id, page_id, type, sort_order, content, is_visible, created_at, updated_at
FROM cms_blocks WHERE id = $1
`, id)
	return scanCMSBlock(row)
}

func (s *Store) CreateBlock(ctx context.Context, block domain.Block) (domain.Block, error) {
	_, err := s.db.ExecContext(ctx, `
INSERT INTO cms_blocks (id, page_id, type, sort_order, content, is_visible, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`, block.ID, block.PageID, block.Type, block.SortOrder, []byte(block.Content), block.IsVisible, block.CreatedAt, block.UpdatedAt)
	if err != nil {
		return domain.Block{}, fmt.Errorf("create cms block: %w", err)
	}
	return block, nil
}

func (s *Store) UpdateBlock(ctx context.Context, block domain.Block) (domain.Block, error) {
	result, err := s.db.ExecContext(ctx, `
UPDATE cms_blocks
SET sort_order = $2, content = $3, is_visible = $4, updated_at = $5
WHERE id = $1
`, block.ID, block.SortOrder, []byte(block.Content), block.IsVisible, block.UpdatedAt)
	if err != nil {
		return domain.Block{}, fmt.Errorf("update cms block: %w", err)
	}
	if err := requireAffected(result); err != nil {
		return domain.Block{}, err
	}
	return block, nil
}

func (s *Store) DeleteBlock(ctx context.Context, id uuid.UUID) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM cms_blocks WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete cms block: %w", err)
	}
	return requireAffected(result)
}

func (s *Store) ReorderBlocks(ctx context.Context, pageID uuid.UUID, orderedIDs []uuid.UUID) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin reorder blocks: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	for index, id := range orderedIDs {
		result, err := tx.ExecContext(ctx, `
UPDATE cms_blocks SET sort_order = $3, updated_at = NOW()
WHERE id = $1 AND page_id = $2
`, id, pageID, index)
		if err != nil {
			return fmt.Errorf("reorder cms block: %w", err)
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if affected == 0 {
			return domain.ErrNotFound
		}
	}
	return tx.Commit()
}

func scanCMSPage(row scanner) (domain.Page, error) {
	var page domain.Page
	err := row.Scan(
		&page.ID,
		&page.Slug,
		&page.Title,
		&page.Path,
		&page.MetaTitle,
		&page.MetaDescription,
		&page.IsPublished,
		&page.CreatedAt,
		&page.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Page{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Page{}, fmt.Errorf("scan cms page: %w", err)
	}
	return page, nil
}

func scanCMSBlock(row scanner) (domain.Block, error) {
	var block domain.Block
	var content []byte
	err := row.Scan(&block.ID, &block.PageID, &block.Type, &block.SortOrder, &content, &block.IsVisible, &block.CreatedAt, &block.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Block{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Block{}, fmt.Errorf("scan cms block: %w", err)
	}
	if len(content) == 0 {
		block.Content = json.RawMessage(`{}`)
	} else {
		block.Content = json.RawMessage(content)
	}
	return block, nil
}

func mapCMSWriteError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23505" {
		switch pqErr.Constraint {
		case "cms_pages_slug_unique":
			return domain.ErrDuplicateSlug
		case "cms_pages_path_unique":
			return domain.ErrDuplicatePath
		}
		return domain.ErrDuplicateSlug
	}
	return fmt.Errorf("write cms page: %w", err)
}
