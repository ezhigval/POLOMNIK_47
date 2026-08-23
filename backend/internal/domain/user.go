package domain

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

const (
	minPasswordLength = 8
	maxNameLength     = 120
)

type User struct {
	ID           uuid.UUID
	Email        string
	Phone        string
	Name         string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserIdentity struct {
	UserID    uuid.UUID
	Provider  string
	Subject   string
	CreatedAt time.Time
}

type ProfileConflict struct {
	Field string
	Kept  string
	Other string
}

type RegisterUserInput struct {
	ID           uuid.UUID
	Email        string
	Phone        string
	Name         string
	PasswordHash string
	Now          time.Time
}

func NewUser(input RegisterUserInput) (User, error) {
	if input.ID == uuid.Nil {
		return User{}, ErrInvalidID
	}

	name := strings.TrimSpace(input.Name)
	if name == "" || utf8.RuneCountInString(name) > maxNameLength {
		return User{}, ErrInvalidContactName
	}

	phone := NormalizePhone(input.Phone)
	if phone == "" {
		return User{}, ErrInvalidPhone
	}

	email := strings.TrimSpace(strings.ToLower(input.Email))
	if email != "" && !strings.Contains(email, "@") {
		return User{}, ErrInvalidEmail
	}

	if strings.TrimSpace(input.PasswordHash) == "" {
		return User{}, ErrInvalidPassword
	}

	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}

	return User{
		ID:           input.ID,
		Email:        email,
		Phone:        phone,
		Name:         name,
		PasswordHash: input.PasswordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func NormalizePhone(phone string) string {
	digits := strings.Builder{}
	for _, r := range phone {
		if r >= '0' && r <= '9' {
			digits.WriteRune(r)
		}
	}
	value := digits.String()
	if len(value) == 11 && value[0] == '8' {
		return "+7" + value[1:]
	}
	if len(value) == 10 {
		return "+7" + value
	}
	if strings.HasPrefix(phone, "+") && len(value) >= 10 {
		return "+" + value
	}
	return strings.TrimSpace(phone)
}

func ValidatePassword(password string) error {
	if utf8.RuneCountInString(password) < minPasswordLength {
		return ErrInvalidPassword
	}
	return nil
}

type OAuthUserInput struct {
	ID       uuid.UUID
	Provider string
	Subject  string
	Email    string
	Name     string
	Phone    string
	Now      time.Time
}

func NewOAuthUser(input OAuthUserInput) (User, error) {
	if input.ID == uuid.Nil {
		return User{}, ErrInvalidID
	}

	provider := strings.TrimSpace(strings.ToLower(input.Provider))
	subject := strings.TrimSpace(input.Subject)
	if provider == "" || subject == "" {
		return User{}, ErrInvalidCredentials
	}

	name := strings.TrimSpace(input.Name)
	if name == "" || utf8.RuneCountInString(name) > maxNameLength {
		return User{}, ErrInvalidContactName
	}

	email := strings.TrimSpace(strings.ToLower(input.Email))
	if email != "" && !strings.Contains(email, "@") {
		return User{}, ErrInvalidEmail
	}

	phone := strings.TrimSpace(input.Phone)
	if phone != "" {
		phone = NormalizePhone(phone)
		if phone == "" {
			return User{}, ErrInvalidPhone
		}
	}

	now := input.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}

	return User{
		ID:        input.ID,
		Email:     email,
		Phone:     phone,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func NormalizeOAuthProvider(raw string) string {
	return strings.TrimSpace(strings.ToLower(raw))
}

func NewUserIdentity(userID uuid.UUID, provider, subject string, now time.Time) (UserIdentity, error) {
	if userID == uuid.Nil {
		return UserIdentity{}, ErrInvalidID
	}
	provider = NormalizeOAuthProvider(provider)
	subject = strings.TrimSpace(subject)
	if provider == "" || subject == "" {
		return UserIdentity{}, ErrInvalidCredentials
	}
	if now.IsZero() {
		now = time.Now().UTC()
	}
	return UserIdentity{
		UserID:    userID,
		Provider:  provider,
		Subject:   subject,
		CreatedAt: now,
	}, nil
}

// FillEmptyProfile copies name/email/phone from source only when the target field is empty.
// Differing non-empty values stay on the target and are returned as conflicts.
func FillEmptyProfile(into, from User) (User, []ProfileConflict) {
	var conflicts []ProfileConflict
	into.Name, conflicts = mergeProfileField(into.Name, from.Name, "name", conflicts)
	into.Email, conflicts = mergeProfileField(into.Email, from.Email, "email", conflicts)
	into.Phone, conflicts = mergeProfileField(into.Phone, from.Phone, "phone", conflicts)
	if into.PasswordHash == "" && from.PasswordHash != "" {
		into.PasswordHash = from.PasswordHash
	}
	return into, conflicts
}

func mergeProfileField(kept, other, field string, conflicts []ProfileConflict) (string, []ProfileConflict) {
	kept = strings.TrimSpace(kept)
	other = strings.TrimSpace(other)
	if kept == "" {
		return other, conflicts
	}
	if other != "" && !strings.EqualFold(kept, other) {
		conflicts = append(conflicts, ProfileConflict{Field: field, Kept: kept, Other: other})
	}
	return kept, conflicts
}
