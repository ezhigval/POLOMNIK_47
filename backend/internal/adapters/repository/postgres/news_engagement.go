package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

func (s *Store) CountNewsLikes(ctx context.Context, newsID uuid.UUID) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM news_likes WHERE news_id = $1`, newsID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count news likes: %w", err)
	}
	return count, nil
}

func (s *Store) HasNewsLike(ctx context.Context, newsID uuid.UUID, visitorID string) (bool, error) {
	var exists bool
	err := s.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM news_likes WHERE news_id = $1 AND visitor_id = $2)`,
		newsID, visitorID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("has news like: %w", err)
	}
	return exists, nil
}

func (s *Store) AddNewsLike(ctx context.Context, newsID uuid.UUID, visitorID string) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO news_likes (news_id, visitor_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		newsID, visitorID,
	)
	if err != nil {
		return fmt.Errorf("add news like: %w", err)
	}
	return nil
}

func (s *Store) RemoveNewsLike(ctx context.Context, newsID uuid.UUID, visitorID string) error {
	_, err := s.db.ExecContext(ctx,
		`DELETE FROM news_likes WHERE news_id = $1 AND visitor_id = $2`,
		newsID, visitorID,
	)
	if err != nil {
		return fmt.Errorf("remove news like: %w", err)
	}
	return nil
}

func (s *Store) CountNewsComments(ctx context.Context, newsID uuid.UUID) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM news_comments WHERE news_id = $1`, newsID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count news comments: %w", err)
	}
	return count, nil
}

func (s *Store) ListNewsComments(ctx context.Context, newsID uuid.UUID, pagination ports.Pagination) ([]domain.NewsComment, error) {
	pagination = ports.NormalizePagination(pagination.Page, pagination.Limit)
	offset := (pagination.Page - 1) * pagination.Limit

	rows, err := s.db.QueryContext(ctx, `
SELECT c.id, c.news_id, c.user_id, COALESCE(NULLIF(TRIM(u.name), ''), 'Паломник'), c.body, c.created_at
FROM news_comments c
JOIN users u ON u.id = c.user_id
WHERE c.news_id = $1
ORDER BY c.created_at ASC
LIMIT $2 OFFSET $3`, newsID, pagination.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list news comments: %w", err)
	}
	defer rows.Close()

	var comments []domain.NewsComment
	for rows.Next() {
		var comment domain.NewsComment
		if err := rows.Scan(&comment.ID, &comment.NewsID, &comment.UserID, &comment.Author, &comment.Body, &comment.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan news comment: %w", err)
		}
		comments = append(comments, comment)
	}
	return comments, rows.Err()
}

func (s *Store) AddNewsComment(ctx context.Context, comment domain.NewsComment) (domain.NewsComment, error) {
	_, err := s.db.ExecContext(ctx, `
INSERT INTO news_comments (id, news_id, user_id, body, created_at)
VALUES ($1, $2, $3, $4, $5)`,
		comment.ID, comment.NewsID, comment.UserID, comment.Body, comment.CreatedAt,
	)
	if err != nil {
		return domain.NewsComment{}, fmt.Errorf("add news comment: %w", err)
	}
	return comment, nil
}

var _ ports.NewsEngagementRepository = (*Store)(nil)
