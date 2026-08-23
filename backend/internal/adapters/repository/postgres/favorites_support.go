package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"polomnik/internal/domain"
)

func (s *Store) AddFavorite(ctx context.Context, userID, tourID uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, `
INSERT INTO favorites (user_id, tour_id, created_at)
VALUES ($1, $2, NOW())
ON CONFLICT (user_id, tour_id) DO NOTHING
`, userID, tourID)
	if err != nil {
		return fmt.Errorf("add favorite: %w", err)
	}
	return nil
}

func (s *Store) RemoveFavorite(ctx context.Context, userID, tourID uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM favorites WHERE user_id = $1 AND tour_id = $2`, userID, tourID)
	if err != nil {
		return fmt.Errorf("remove favorite: %w", err)
	}
	return nil
}

func (s *Store) ListFavoriteTourIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := s.db.QueryContext(ctx, `
SELECT tour_id
FROM favorites
WHERE user_id = $1
ORDER BY created_at DESC
`, userID)
	if err != nil {
		return nil, fmt.Errorf("list favorites: %w", err)
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan favorite: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (s *Store) IsFavorite(ctx context.Context, userID, tourID uuid.UUID) (bool, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT EXISTS(SELECT 1 FROM favorites WHERE user_id = $1 AND tour_id = $2)
`, userID, tourID)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, fmt.Errorf("is favorite: %w", err)
	}
	return exists, nil
}

func (s *Store) GetOpenThread(ctx context.Context, userID uuid.UUID) (domain.SupportThread, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT id, user_id, subject, status, created_at, updated_at
FROM support_threads
WHERE user_id = $1 AND status = 'open'
ORDER BY updated_at DESC
LIMIT 1
`, userID)
	return scanSupportThread(row)
}

func (s *Store) GetThreadByID(ctx context.Context, threadID uuid.UUID) (domain.SupportThread, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT id, user_id, subject, status, created_at, updated_at
FROM support_threads
WHERE id = $1
`, threadID)
	return scanSupportThread(row)
}

func (s *Store) ListThreads(ctx context.Context) ([]domain.SupportThread, error) {
	rows, err := s.db.QueryContext(ctx, `
SELECT id, user_id, subject, status, created_at, updated_at
FROM support_threads
ORDER BY updated_at DESC
`)
	if err != nil {
		return nil, fmt.Errorf("list support threads: %w", err)
	}
	defer rows.Close()

	var threads []domain.SupportThread
	for rows.Next() {
		thread, err := scanSupportThread(rows)
		if err != nil {
			return nil, err
		}
		threads = append(threads, thread)
	}
	return threads, rows.Err()
}

func (s *Store) CreateThread(ctx context.Context, thread domain.SupportThread) (domain.SupportThread, error) {
	row := s.db.QueryRowContext(ctx, `
INSERT INTO support_threads (id, user_id, subject, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, user_id, subject, status, created_at, updated_at
`, thread.ID, thread.UserID, thread.Subject, thread.Status, thread.CreatedAt, thread.UpdatedAt)
	return scanSupportThread(row)
}

func (s *Store) ListMessages(ctx context.Context, threadID uuid.UUID) ([]domain.SupportMessage, error) {
	rows, err := s.db.QueryContext(ctx, `
SELECT id, thread_id, sender_type, body, created_at
FROM support_messages
WHERE thread_id = $1
ORDER BY created_at ASC
`, threadID)
	if err != nil {
		return nil, fmt.Errorf("list support messages: %w", err)
	}
	defer rows.Close()

	var messages []domain.SupportMessage
	for rows.Next() {
		message, err := scanSupportMessage(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, rows.Err()
}

func (s *Store) AddMessage(ctx context.Context, message domain.SupportMessage) (domain.SupportMessage, error) {
	row := s.db.QueryRowContext(ctx, `
INSERT INTO support_messages (id, thread_id, sender_type, body, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, thread_id, sender_type, body, created_at
`, message.ID, message.ThreadID, message.SenderType, message.Body, message.CreatedAt)
	return scanSupportMessage(row)
}

func (s *Store) TouchThread(ctx context.Context, threadID uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, `UPDATE support_threads SET updated_at = NOW() WHERE id = $1`, threadID)
	return err
}

func scanSupportThread(row scanner) (domain.SupportThread, error) {
	var thread domain.SupportThread
	err := row.Scan(&thread.ID, &thread.UserID, &thread.Subject, &thread.Status, &thread.CreatedAt, &thread.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.SupportThread{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.SupportThread{}, fmt.Errorf("scan support thread: %w", err)
	}
	return thread, nil
}

func scanSupportMessage(row scanner) (domain.SupportMessage, error) {
	var message domain.SupportMessage
	var senderType string
	err := row.Scan(&message.ID, &message.ThreadID, &senderType, &message.Body, &message.CreatedAt)
	if err != nil {
		return domain.SupportMessage{}, err
	}
	message.SenderType = domain.SupportSenderType(senderType)
	return message, nil
}
