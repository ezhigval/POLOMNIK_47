package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

func (s *Store) ListNews(ctx context.Context, filters ports.NewsFilters, pagination ports.Pagination) (ports.NewsList, error) {
	pagination = ports.NormalizePagination(pagination.Page, pagination.Limit)
	publishedOnly := false
	if filters.PublishedOnly {
		publishedOnly = true
	}

	total, err := count(ctx, s.db, `
SELECT COUNT(*) FROM news_articles
WHERE ($1::boolean = FALSE OR is_published = TRUE)
`, publishedOnly)
	if err != nil {
		return ports.NewsList{}, err
	}

	rows, err := s.db.QueryContext(ctx, `
SELECT id, slug, title, excerpt, body, image_url, published_at, is_published, sort_order, created_at, updated_at
FROM news_articles
WHERE ($1::boolean = FALSE OR is_published = TRUE)
ORDER BY published_at DESC, sort_order ASC, created_at DESC
LIMIT $2 OFFSET $3
`, publishedOnly, pagination.Limit, offset(pagination))
	if err != nil {
		return ports.NewsList{}, fmt.Errorf("list news: %w", err)
	}
	defer rows.Close()

	items, err := scanNewsArticles(rows)
	if err != nil {
		return ports.NewsList{}, err
	}

	return ports.NewsList{Items: items, Meta: pageMeta(pagination, total)}, nil
}

func (s *Store) GetNewsBySlug(ctx context.Context, slug string) (domain.NewsArticle, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT id, slug, title, excerpt, body, image_url, published_at, is_published, sort_order, created_at, updated_at
FROM news_articles
WHERE slug = $1
`, strings.TrimSpace(slug))
	return scanNewsArticle(row)
}

func (s *Store) GetNews(ctx context.Context, id uuid.UUID) (domain.NewsArticle, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT id, slug, title, excerpt, body, image_url, published_at, is_published, sort_order, created_at, updated_at
FROM news_articles
WHERE id = $1
`, id)
	return scanNewsArticle(row)
}

func (s *Store) CreateNews(ctx context.Context, article domain.NewsArticle) (domain.NewsArticle, error) {
	row := s.db.QueryRowContext(ctx, `
INSERT INTO news_articles (
    id, slug, title, excerpt, body, image_url, published_at, is_published, sort_order, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING id, slug, title, excerpt, body, image_url, published_at, is_published, sort_order, created_at, updated_at
`, article.ID, article.Slug, article.Title, article.Excerpt, article.Body, article.ImageURL,
		article.PublishedAt, article.IsPublished, article.SortOrder, article.CreatedAt, article.UpdatedAt)
	article, err := scanNewsArticle(row)
	if err != nil {
		return domain.NewsArticle{}, mapNewsWriteError(err)
	}
	return article, nil
}

func (s *Store) UpdateNews(ctx context.Context, article domain.NewsArticle) (domain.NewsArticle, error) {
	row := s.db.QueryRowContext(ctx, `
UPDATE news_articles
SET slug = $2,
    title = $3,
    excerpt = $4,
    body = $5,
    image_url = $6,
    published_at = $7,
    is_published = $8,
    sort_order = $9,
    updated_at = $10
WHERE id = $1
RETURNING id, slug, title, excerpt, body, image_url, published_at, is_published, sort_order, created_at, updated_at
`, article.ID, article.Slug, article.Title, article.Excerpt, article.Body, article.ImageURL,
		article.PublishedAt, article.IsPublished, article.SortOrder, article.UpdatedAt)
	article, err := scanNewsArticle(row)
	if err != nil {
		return domain.NewsArticle{}, mapNewsWriteError(err)
	}
	return article, nil
}

func (s *Store) DeleteNews(ctx context.Context, id uuid.UUID) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM news_articles WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete news: %w", err)
	}
	return requireAffected(result)
}

func scanNewsArticle(row scanner) (domain.NewsArticle, error) {
	var article domain.NewsArticle
	err := row.Scan(
		&article.ID,
		&article.Slug,
		&article.Title,
		&article.Excerpt,
		&article.Body,
		&article.ImageURL,
		&article.PublishedAt,
		&article.IsPublished,
		&article.SortOrder,
		&article.CreatedAt,
		&article.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.NewsArticle{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.NewsArticle{}, fmt.Errorf("scan news: %w", err)
	}
	return article, nil
}

func scanNewsArticles(rows *sql.Rows) ([]domain.NewsArticle, error) {
	var items []domain.NewsArticle
	for rows.Next() {
		article, err := scanNewsArticle(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, article)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate news: %w", err)
	}
	return items, nil
}

func mapNewsWriteError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23505" {
		return domain.ErrDuplicateSlug
	}
	if errors.Is(err, domain.ErrNotFound) {
		return err
	}
	return fmt.Errorf("write news: %w", err)
}
