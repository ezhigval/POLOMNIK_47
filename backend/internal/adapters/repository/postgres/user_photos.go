package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"palomnik/internal/domain"
)

func (s *Store) ListUserPhotos(ctx context.Context, userID uuid.UUID) ([]domain.UserPhoto, error) {
	rows, err := s.conn(ctx).QueryContext(ctx, `
SELECT id, user_id, url, caption, allow_distribution, created_at, updated_at
FROM user_photos
WHERE user_id = $1
ORDER BY created_at DESC, id ASC
`, userID)
	if err != nil {
		return nil, fmt.Errorf("list user photos: %w", err)
	}
	defer rows.Close()

	var items []domain.UserPhoto
	for rows.Next() {
		photo, err := scanUserPhoto(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, photo)
	}
	if items == nil {
		items = []domain.UserPhoto{}
	}
	return items, rows.Err()
}

func (s *Store) CreateUserPhoto(ctx context.Context, photo domain.UserPhoto) (domain.UserPhoto, error) {
	_, err := s.conn(ctx).ExecContext(ctx, `
INSERT INTO user_photos (id, user_id, url, caption, allow_distribution, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
`, photo.ID, photo.UserID, photo.URL, photo.Caption, photo.AllowDistribution, photo.CreatedAt, photo.UpdatedAt)
	if err != nil {
		return domain.UserPhoto{}, fmt.Errorf("create user photo: %w", err)
	}
	return photo, nil
}

func (s *Store) DeleteUserPhoto(ctx context.Context, userID, id uuid.UUID) error {
	res, err := s.conn(ctx).ExecContext(ctx, `DELETE FROM user_photos WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("delete user photo: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func scanUserPhoto(row scanner) (domain.UserPhoto, error) {
	var photo domain.UserPhoto
	err := row.Scan(
		&photo.ID,
		&photo.UserID,
		&photo.URL,
		&photo.Caption,
		&photo.AllowDistribution,
		&photo.CreatedAt,
		&photo.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.UserPhoto{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.UserPhoto{}, fmt.Errorf("scan user photo: %w", err)
	}
	return photo, nil
}
