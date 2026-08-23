package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"palomnik/internal/domain"
)

type notificationRoutesDTO map[string][]notificationRecipientDTO

type notificationRecipientDTO struct {
	Channel string `json:"channel"`
	Address string `json:"address"`
}

func (s *Store) GetNotificationRouting(ctx context.Context) (domain.NotificationRouting, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT routes, updated_at
FROM notification_routing
WHERE id = $1
`, domain.NotificationRoutingID())

	var raw []byte
	var updatedAt time.Time
	err := row.Scan(&raw, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.NotificationRouting{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.NotificationRouting{}, fmt.Errorf("get notification routing: %w", err)
	}
	return decodeNotificationRouting(raw, updatedAt)
}

func (s *Store) UpsertNotificationRouting(ctx context.Context, routing domain.NotificationRouting) (domain.NotificationRouting, error) {
	payload, err := encodeNotificationRouting(routing)
	if err != nil {
		return domain.NotificationRouting{}, err
	}
	row := s.db.QueryRowContext(ctx, `
INSERT INTO notification_routing (id, routes, updated_at)
VALUES ($1, $2::jsonb, $3)
ON CONFLICT (id) DO UPDATE SET
    routes = EXCLUDED.routes,
    updated_at = EXCLUDED.updated_at
RETURNING routes, updated_at
`, domain.NotificationRoutingID(), payload, routing.UpdatedAt)

	var raw []byte
	var updatedAt time.Time
	if err := row.Scan(&raw, &updatedAt); err != nil {
		return domain.NotificationRouting{}, fmt.Errorf("upsert notification routing: %w", err)
	}
	return decodeNotificationRouting(raw, updatedAt)
}

func encodeNotificationRouting(routing domain.NotificationRouting) ([]byte, error) {
	dto := notificationRoutesDTO{}
	for _, kind := range domain.AllNotificationEventKinds() {
		list := domain.RecipientsForEvent(routing, kind)
		items := make([]notificationRecipientDTO, 0, len(list))
		for _, recipient := range list {
			items = append(items, notificationRecipientDTO{
				Channel: string(recipient.Channel),
				Address: recipient.Address,
			})
		}
		dto[string(kind)] = items
	}
	return json.Marshal(dto)
}

func decodeNotificationRouting(raw []byte, updatedAt time.Time) (domain.NotificationRouting, error) {
	dto := notificationRoutesDTO{}
	if len(raw) > 0 && string(raw) != "null" {
		if err := json.Unmarshal(raw, &dto); err != nil {
			return domain.NotificationRouting{}, fmt.Errorf("decode notification routing: %w", err)
		}
	}
	events := make(map[domain.NotificationEventKind][]domain.NotificationRecipient, len(dto))
	for kind, items := range dto {
		list := make([]domain.NotificationRecipient, 0, len(items))
		for _, item := range items {
			list = append(list, domain.NotificationRecipient{
				Channel: domain.NotificationChannel(item.Channel),
				Address: item.Address,
			})
		}
		events[domain.NotificationEventKind(kind)] = list
	}
	return domain.NewNotificationRouting(events, updatedAt)
}
