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
)

const userSelectColumns = `
id, email, phone, name, password_hash, created_at, updated_at
`

func (s *Store) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	row := s.conn(ctx).QueryRowContext(ctx, `
INSERT INTO users (id, email, phone, name, password_hash, created_at, updated_at)
VALUES ($1, $2, NULLIF($3, ''), $4, NULLIF($5, ''), $6, $7)
RETURNING `+userSelectColumns+`
`, user.ID, user.Email, user.Phone, user.Name, user.PasswordHash, user.CreatedAt, user.UpdatedAt)
	created, err := scanUser(row)
	if err != nil {
		return domain.User{}, mapUserWriteError(err)
	}
	return created, nil
}

func (s *Store) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	row := s.conn(ctx).QueryRowContext(ctx, `SELECT `+userSelectColumns+` FROM users WHERE id = $1`, id)
	return scanUser(row)
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	normalized := strings.TrimSpace(strings.ToLower(email))
	if normalized == "" {
		return domain.User{}, domain.ErrNotFound
	}
	row := s.conn(ctx).QueryRowContext(ctx, `SELECT `+userSelectColumns+` FROM users WHERE email = $1`, normalized)
	return scanUser(row)
}

func (s *Store) GetUserByPhone(ctx context.Context, phone string) (domain.User, error) {
	normalized := domain.NormalizePhone(phone)
	if normalized == "" {
		return domain.User{}, domain.ErrNotFound
	}
	row := s.conn(ctx).QueryRowContext(ctx, `SELECT `+userSelectColumns+` FROM users WHERE phone = $1`, normalized)
	return scanUser(row)
}

func (s *Store) GetUserByOAuth(ctx context.Context, provider, subject string) (domain.User, error) {
	provider = domain.NormalizeOAuthProvider(provider)
	subject = strings.TrimSpace(subject)
	if provider == "" || subject == "" {
		return domain.User{}, domain.ErrNotFound
	}
	row := s.conn(ctx).QueryRowContext(ctx, `
SELECT u.id, u.email, u.phone, u.name, u.password_hash, u.created_at, u.updated_at
FROM users u
JOIN user_identities i ON i.user_id = u.id
WHERE i.provider = $1 AND i.subject = $2
`, provider, subject)
	return scanUser(row)
}

func (s *Store) GetIdentity(ctx context.Context, provider, subject string) (domain.UserIdentity, error) {
	provider = domain.NormalizeOAuthProvider(provider)
	subject = strings.TrimSpace(subject)
	if provider == "" || subject == "" {
		return domain.UserIdentity{}, domain.ErrNotFound
	}
	row := s.conn(ctx).QueryRowContext(ctx, `
SELECT user_id, provider, subject, created_at
FROM user_identities
WHERE provider = $1 AND subject = $2
`, provider, subject)
	return scanIdentity(row)
}

func (s *Store) ListIdentities(ctx context.Context, userID uuid.UUID) ([]domain.UserIdentity, error) {
	rows, err := s.conn(ctx).QueryContext(ctx, `
SELECT user_id, provider, subject, created_at
FROM user_identities
WHERE user_id = $1
ORDER BY created_at ASC, provider ASC
`, userID)
	if err != nil {
		return nil, fmt.Errorf("list identities: %w", err)
	}
	defer rows.Close()

	out := make([]domain.UserIdentity, 0)
	for rows.Next() {
		identity, err := scanIdentity(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, identity)
	}
	return out, rows.Err()
}

func (s *Store) AddIdentity(ctx context.Context, identity domain.UserIdentity) error {
	_, err := s.conn(ctx).ExecContext(ctx, `
INSERT INTO user_identities (user_id, provider, subject, created_at)
VALUES ($1, $2, $3, $4)
`, identity.UserID, identity.Provider, identity.Subject, identity.CreatedAt)
	if err != nil {
		return mapUserWriteError(err)
	}
	return nil
}

func (s *Store) UpdateUserProfile(ctx context.Context, user domain.User) (domain.User, error) {
	row := s.conn(ctx).QueryRowContext(ctx, `
UPDATE users
SET email = $2,
    phone = NULLIF($3, ''),
    name = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING `+userSelectColumns+`
`, user.ID, user.Email, user.Phone, user.Name)
	updated, err := scanUser(row)
	if err != nil {
		return domain.User{}, mapUserWriteError(err)
	}
	return updated, nil
}

func (s *Store) UpdateUserPassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	res, err := s.conn(ctx).ExecContext(ctx, `
UPDATE users
SET password_hash = $2,
    updated_at = NOW()
WHERE id = $1
`, userID, passwordHash)
	if err != nil {
		return fmt.Errorf("update user password: %w", err)
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

func (s *Store) MergeAccountInto(ctx context.Context, targetID, sourceID uuid.UUID) error {
	if targetID == sourceID {
		return nil
	}
	if _, err := s.GetUserByID(ctx, targetID); err != nil {
		return err
	}
	if _, err := s.GetUserByID(ctx, sourceID); err != nil {
		return err
	}

	conn := s.conn(ctx)
	if _, err := conn.ExecContext(ctx, `UPDATE bookings SET user_id = $1 WHERE user_id = $2`, targetID, sourceID); err != nil {
		return fmt.Errorf("merge bookings: %w", err)
	}
	if _, err := conn.ExecContext(ctx, `
INSERT INTO favorites (user_id, tour_id, created_at)
SELECT $1, tour_id, created_at
FROM favorites
WHERE user_id = $2
ON CONFLICT (user_id, tour_id) DO NOTHING
`, targetID, sourceID); err != nil {
		return fmt.Errorf("merge favorites: %w", err)
	}
	if _, err := conn.ExecContext(ctx, `DELETE FROM favorites WHERE user_id = $1`, sourceID); err != nil {
		return fmt.Errorf("merge favorites cleanup: %w", err)
	}

	var targetOpen, sourceOpen uuid.UUID
	_ = conn.QueryRowContext(ctx, `
SELECT id FROM support_threads WHERE user_id = $1 AND status = 'open' LIMIT 1
`, targetID).Scan(&targetOpen)
	_ = conn.QueryRowContext(ctx, `
SELECT id FROM support_threads WHERE user_id = $1 AND status = 'open' LIMIT 1
`, sourceID).Scan(&sourceOpen)
	if targetOpen != uuid.Nil && sourceOpen != uuid.Nil && targetOpen != sourceOpen {
		if _, err := conn.ExecContext(ctx, `UPDATE support_messages SET thread_id = $1 WHERE thread_id = $2`, targetOpen, sourceOpen); err != nil {
			return fmt.Errorf("merge support messages: %w", err)
		}
		if _, err := conn.ExecContext(ctx, `DELETE FROM support_threads WHERE id = $1`, sourceOpen); err != nil {
			return fmt.Errorf("merge support thread: %w", err)
		}
	}
	if _, err := conn.ExecContext(ctx, `UPDATE support_threads SET user_id = $1 WHERE user_id = $2`, targetID, sourceID); err != nil {
		return fmt.Errorf("merge support threads: %w", err)
	}

	if _, err := conn.ExecContext(ctx, `
INSERT INTO admin_role_user_assignments (role_id, user_id, created_at)
SELECT role_id, $1, created_at
FROM admin_role_user_assignments
WHERE user_id = $2
ON CONFLICT (role_id, user_id) DO NOTHING
`, targetID, sourceID); err != nil {
		return fmt.Errorf("merge role assignments: %w", err)
	}
	if _, err := conn.ExecContext(ctx, `DELETE FROM admin_role_user_assignments WHERE user_id = $1`, sourceID); err != nil {
		return fmt.Errorf("merge role assignments cleanup: %w", err)
	}

	if _, err := conn.ExecContext(ctx, `UPDATE user_identities SET user_id = $1 WHERE user_id = $2`, targetID, sourceID); err != nil {
		return fmt.Errorf("merge identities: %w", err)
	}
	if _, err := conn.ExecContext(ctx, `UPDATE passengers SET user_id = $1, updated_at = NOW() WHERE user_id = $2`, targetID, sourceID); err != nil {
		return fmt.Errorf("merge passengers: %w", err)
	}
	if _, err := conn.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, sourceID); err != nil {
		return fmt.Errorf("merge delete source user: %w", err)
	}
	return nil
}

func scanUser(row scanner) (domain.User, error) {
	var user domain.User
	var phone sql.NullString
	var passwordHash sql.NullString
	err := row.Scan(
		&user.ID,
		&user.Email,
		&phone,
		&user.Name,
		&passwordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.User{}, fmt.Errorf("scan user: %w", err)
	}
	if phone.Valid {
		user.Phone = phone.String
	}
	if passwordHash.Valid {
		user.PasswordHash = passwordHash.String
	}
	return user, nil
}

func scanIdentity(row scanner) (domain.UserIdentity, error) {
	var identity domain.UserIdentity
	err := row.Scan(&identity.UserID, &identity.Provider, &identity.Subject, &identity.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.UserIdentity{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.UserIdentity{}, fmt.Errorf("scan identity: %w", err)
	}
	return identity, nil
}

func mapUserWriteError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23505" {
		switch pqErr.Constraint {
		case "users_phone_unique", "idx_users_phone_unique":
			return domain.ErrDuplicatePhone
		case "idx_users_email_unique":
			return domain.ErrDuplicateEmail
		case "user_identities_pkey":
			return domain.ErrDuplicateIdentity
		}
	}
	return err
}
