package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"palomnik/internal/domain"
)

func (s *Store) ListPassengers(ctx context.Context, userID uuid.UUID) ([]domain.Passenger, error) {
	rows, err := s.conn(ctx).QueryContext(ctx, `
SELECT id, user_id, name, phone, birth_date, passport, created_at, updated_at
FROM passengers
WHERE user_id = $1
ORDER BY created_at ASC, id ASC
`, userID)
	if err != nil {
		return nil, fmt.Errorf("list passengers: %w", err)
	}
	defer rows.Close()

	var items []domain.Passenger
	for rows.Next() {
		passenger, err := scanPassenger(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, passenger)
	}
	if items == nil {
		items = []domain.Passenger{}
	}
	return items, rows.Err()
}

func (s *Store) GetPassenger(ctx context.Context, userID, id uuid.UUID) (domain.Passenger, error) {
	row := s.conn(ctx).QueryRowContext(ctx, `
SELECT id, user_id, name, phone, birth_date, passport, created_at, updated_at
FROM passengers
WHERE id = $1 AND user_id = $2
`, id, userID)
	return scanPassenger(row)
}

func (s *Store) CreatePassenger(ctx context.Context, passenger domain.Passenger) (domain.Passenger, error) {
	_, err := s.conn(ctx).ExecContext(ctx, `
INSERT INTO passengers (id, user_id, name, phone, birth_date, passport, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`, passenger.ID, passenger.UserID, passenger.Name, passenger.Phone, passenger.BirthDate, passenger.Passport, passenger.CreatedAt, passenger.UpdatedAt)
	if err != nil {
		return domain.Passenger{}, fmt.Errorf("create passenger: %w", err)
	}
	return passenger, nil
}

func (s *Store) UpdatePassenger(ctx context.Context, passenger domain.Passenger) (domain.Passenger, error) {
	res, err := s.conn(ctx).ExecContext(ctx, `
UPDATE passengers
SET name = $3, phone = $4, birth_date = $5, passport = $6, updated_at = $7
WHERE id = $1 AND user_id = $2
`, passenger.ID, passenger.UserID, passenger.Name, passenger.Phone, passenger.BirthDate, passenger.Passport, passenger.UpdatedAt)
	if err != nil {
		return domain.Passenger{}, fmt.Errorf("update passenger: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return domain.Passenger{}, err
	}
	if n == 0 {
		return domain.Passenger{}, domain.ErrNotFound
	}
	return passenger, nil
}

func (s *Store) DeletePassenger(ctx context.Context, userID, id uuid.UUID) error {
	res, err := s.conn(ctx).ExecContext(ctx, `DELETE FROM passengers WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("delete passenger: %w", err)
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

func scanPassenger(row scanner) (domain.Passenger, error) {
	var passenger domain.Passenger
	err := row.Scan(
		&passenger.ID,
		&passenger.UserID,
		&passenger.Name,
		&passenger.Phone,
		&passenger.BirthDate,
		&passenger.Passport,
		&passenger.CreatedAt,
		&passenger.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Passenger{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Passenger{}, fmt.Errorf("scan passenger: %w", err)
	}
	return passenger, nil
}
