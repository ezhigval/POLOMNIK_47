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
id, email, phone, name, password_hash, oauth_provider, oauth_subject, created_at, updated_at
`

func (s *Store) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	row := s.db.QueryRowContext(ctx, `
INSERT INTO users (id, email, phone, name, password_hash, oauth_provider, oauth_subject, created_at, updated_at)
VALUES ($1, $2, NULLIF($3, ''), $4, NULLIF($5, ''), $6, $7, $8, $9)
RETURNING `+userSelectColumns+`
`, user.ID, user.Email, user.Phone, user.Name, user.PasswordHash, user.OAuthProvider, user.OAuthSubject, user.CreatedAt, user.UpdatedAt)
	created, err := scanUser(row)
	if err != nil {
		return domain.User{}, mapUserCreateError(err)
	}
	return created, nil
}

func (s *Store) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	row := s.db.QueryRowContext(ctx, `SELECT `+userSelectColumns+` FROM users WHERE id = $1`, id)
	return scanUser(row)
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	normalized := strings.TrimSpace(strings.ToLower(email))
	if normalized == "" {
		return domain.User{}, domain.ErrNotFound
	}
	row := s.db.QueryRowContext(ctx, `SELECT `+userSelectColumns+` FROM users WHERE email = $1`, normalized)
	return scanUser(row)
}

func (s *Store) GetUserByPhone(ctx context.Context, phone string) (domain.User, error) {
	normalized := domain.NormalizePhone(phone)
	if normalized == "" {
		return domain.User{}, domain.ErrNotFound
	}
	row := s.db.QueryRowContext(ctx, `SELECT `+userSelectColumns+` FROM users WHERE phone = $1`, normalized)
	return scanUser(row)
}

func (s *Store) GetUserByOAuth(ctx context.Context, provider, subject string) (domain.User, error) {
	provider = strings.TrimSpace(strings.ToLower(provider))
	subject = strings.TrimSpace(subject)
	if provider == "" || subject == "" {
		return domain.User{}, domain.ErrNotFound
	}
	row := s.db.QueryRowContext(ctx, `
SELECT `+userSelectColumns+`
FROM users
WHERE oauth_provider = $1 AND oauth_subject = $2
`, provider, subject)
	return scanUser(row)
}

func (s *Store) UpdateUserProfile(ctx context.Context, user domain.User) (domain.User, error) {
	row := s.db.QueryRowContext(ctx, `
UPDATE users
SET email = $2,
    phone = NULLIF($3, ''),
    name = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING `+userSelectColumns+`
`, user.ID, user.Email, user.Phone, user.Name)
	return scanUser(row)
}

func (s *Store) UpdateUserPassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	res, err := s.db.ExecContext(ctx, `
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
		&user.OAuthProvider,
		&user.OAuthSubject,
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

func mapUserCreateError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23505" {
		switch pqErr.Constraint {
		case "users_phone_unique", "idx_users_phone_unique":
			return domain.ErrDuplicatePhone
		case "idx_users_email_unique":
			return domain.ErrDuplicateEmail
		case "idx_users_oauth_unique":
			return domain.ErrDuplicateEmail
		}
	}
	return err
}
