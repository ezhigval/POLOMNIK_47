package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

func (s *Store) CreateSMMPost(ctx context.Context, post domain.SMMPost) (domain.SMMPost, error) {
	results, err := json.Marshal(post.Results)
	if err != nil {
		return domain.SMMPost{}, err
	}
	if results == nil {
		results = []byte("[]")
	}
	row := s.db.QueryRowContext(ctx, `
INSERT INTO smm_posts (id, title, body, url, publish_at, channels, published_at, results, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id, title, body, url, publish_at, channels, published_at, results, created_at, updated_at
`, post.ID, post.Title, post.Body, post.URL, post.PublishAt, pq.Array(post.Channels), post.PublishedAt, results, post.CreatedAt, post.UpdatedAt)
	return scanSMMPost(row)
}

func (s *Store) GetSMMPost(ctx context.Context, id uuid.UUID) (domain.SMMPost, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT id, title, body, url, publish_at, channels, published_at, results, created_at, updated_at
FROM smm_posts
WHERE id = $1
`, id)
	return scanSMMPost(row)
}

func (s *Store) ListSMMPosts(ctx context.Context, pagination ports.Pagination) (ports.SMMPostList, error) {
	pagination = ports.NormalizePagination(pagination.Page, pagination.Limit)
	total, err := count(ctx, s.db, `SELECT COUNT(*) FROM smm_posts`)
	if err != nil {
		return ports.SMMPostList{}, err
	}
	rows, err := s.db.QueryContext(ctx, `
SELECT id, title, body, url, publish_at, channels, published_at, results, created_at, updated_at
FROM smm_posts
ORDER BY publish_at DESC, created_at DESC
LIMIT $1 OFFSET $2
`, pagination.Limit, offset(pagination))
	if err != nil {
		return ports.SMMPostList{}, fmt.Errorf("list smm posts: %w", err)
	}
	defer rows.Close()
	items, err := scanSMMPosts(rows)
	if err != nil {
		return ports.SMMPostList{}, err
	}
	return ports.SMMPostList{Items: items, Meta: pageMeta(pagination, total)}, nil
}

func (s *Store) ListDueSMMPosts(ctx context.Context, now time.Time) ([]domain.SMMPost, error) {
	rows, err := s.db.QueryContext(ctx, `
SELECT id, title, body, url, publish_at, channels, published_at, results, created_at, updated_at
FROM smm_posts
WHERE published_at IS NULL AND publish_at <= $1
ORDER BY publish_at ASC
`, now.UTC())
	if err != nil {
		return nil, fmt.Errorf("list due smm posts: %w", err)
	}
	defer rows.Close()
	return scanSMMPosts(rows)
}

func (s *Store) SaveSMMPost(ctx context.Context, post domain.SMMPost) (domain.SMMPost, error) {
	results, err := json.Marshal(post.Results)
	if err != nil {
		return domain.SMMPost{}, err
	}
	if results == nil {
		results = []byte("[]")
	}
	row := s.db.QueryRowContext(ctx, `
UPDATE smm_posts
SET title = $2, body = $3, url = $4, publish_at = $5, channels = $6, published_at = $7, results = $8, updated_at = $9
WHERE id = $1
RETURNING id, title, body, url, publish_at, channels, published_at, results, created_at, updated_at
`, post.ID, post.Title, post.Body, post.URL, post.PublishAt, pq.Array(post.Channels), post.PublishedAt, results, post.UpdatedAt)
	return scanSMMPost(row)
}

func (s *Store) DeleteSMMPost(ctx context.Context, id uuid.UUID) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM smm_posts WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete smm post: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func scanSMMPosts(rows *sql.Rows) ([]domain.SMMPost, error) {
	var items []domain.SMMPost
	for rows.Next() {
		post, err := scanSMMPost(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, post)
	}
	return items, rows.Err()
}

func scanSMMPost(row scanner) (domain.SMMPost, error) {
	var post domain.SMMPost
	var channels []string
	var publishedAt sql.NullTime
	var results []byte
	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Body,
		&post.URL,
		&post.PublishAt,
		pq.Array(&channels),
		&publishedAt,
		&results,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.SMMPost{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.SMMPost{}, fmt.Errorf("scan smm post: %w", err)
	}
	post.Channels = append([]string(nil), channels...)
	if publishedAt.Valid {
		ts := publishedAt.Time.UTC()
		post.PublishedAt = &ts
	}
	if len(results) > 0 {
		if err := json.Unmarshal(results, &post.Results); err != nil {
			return domain.SMMPost{}, fmt.Errorf("scan smm results: %w", err)
		}
	}
	return post, nil
}
