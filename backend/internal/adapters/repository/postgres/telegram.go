package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"

	"polomnik/internal/domain"
)

func (s *Store) GetTelegramRecipients(ctx context.Context) (domain.TelegramRecipients, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT booking_usernames, support_usernames, updated_at
FROM telegram_recipients
WHERE id = $1
`, domain.TelegramRecipientsID())
	return scanTelegramRecipients(row)
}

func (s *Store) UpsertTelegramRecipients(ctx context.Context, settings domain.TelegramRecipients) (domain.TelegramRecipients, error) {
	row := s.db.QueryRowContext(ctx, `
INSERT INTO telegram_recipients (id, booking_usernames, support_usernames, updated_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (id) DO UPDATE SET
    booking_usernames = EXCLUDED.booking_usernames,
    support_usernames = EXCLUDED.support_usernames,
    updated_at = EXCLUDED.updated_at
RETURNING booking_usernames, support_usernames, updated_at
`, domain.TelegramRecipientsID(),
		domain.FormatTelegramUsernameList(settings.BookingUsernames),
		domain.FormatTelegramUsernameList(settings.SupportUsernames),
		settings.UpdatedAt)
	return scanTelegramRecipients(row)
}

func (s *Store) UpsertTelegramChatBinding(ctx context.Context, binding domain.TelegramChatBinding) error {
	_, err := s.db.ExecContext(ctx, `
INSERT INTO telegram_chat_map (username, chat_id, updated_at)
VALUES ($1, $2, $3)
ON CONFLICT (username) DO UPDATE SET
    chat_id = EXCLUDED.chat_id,
    updated_at = EXCLUDED.updated_at
`, binding.Username, binding.ChatID, binding.UpdatedAt)
	if err != nil {
		return fmt.Errorf("upsert telegram chat map: %w", err)
	}
	return nil
}

func (s *Store) ListTelegramChatBindings(ctx context.Context, usernames []string) (map[string]string, error) {
	out := make(map[string]string, len(usernames))
	if len(usernames) == 0 {
		return out, nil
	}

	rows, err := s.db.QueryContext(ctx, `
SELECT username, chat_id
FROM telegram_chat_map
WHERE username = ANY($1)
`, pq.Array(usernames))
	if err != nil {
		return nil, fmt.Errorf("list telegram chat map: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var username, chatID string
		if err := rows.Scan(&username, &chatID); err != nil {
			return nil, fmt.Errorf("scan telegram chat map: %w", err)
		}
		out[username] = chatID
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate telegram chat map: %w", err)
	}
	return out, nil
}

func scanTelegramRecipients(row scanner) (domain.TelegramRecipients, error) {
	var bookingRaw, supportRaw string
	var settings domain.TelegramRecipients
	err := row.Scan(&bookingRaw, &supportRaw, &settings.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.TelegramRecipients{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.TelegramRecipients{}, fmt.Errorf("scan telegram recipients: %w", err)
	}
	booking, err := domain.ParseTelegramUsernameList(bookingRaw)
	if err != nil {
		return domain.TelegramRecipients{}, err
	}
	support, err := domain.ParseTelegramUsernameList(supportRaw)
	if err != nil {
		return domain.TelegramRecipients{}, err
	}
	settings.BookingUsernames = booking
	settings.SupportUsernames = support
	return settings, nil
}
