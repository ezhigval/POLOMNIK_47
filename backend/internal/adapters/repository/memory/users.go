package memory

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/domain"
)

func (s *Store) CreateUser(_ context.Context, user domain.User) (domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, existing := range s.users {
		if user.Phone != "" && existing.Phone == user.Phone {
			return domain.User{}, domain.ErrDuplicatePhone
		}
		if user.Email != "" && existing.Email == user.Email {
			return domain.User{}, domain.ErrDuplicateEmail
		}
		if user.OAuthProvider != "" && existing.OAuthProvider == user.OAuthProvider && existing.OAuthSubject == user.OAuthSubject {
			return domain.User{}, domain.ErrDuplicateEmail
		}
	}

	s.users[user.ID] = user
	return user, nil
}

func (s *Store) GetUserByID(_ context.Context, id uuid.UUID) (domain.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[id]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return user, nil
}

func (s *Store) GetUserByEmail(_ context.Context, email string) (domain.User, error) {
	normalized := strings.TrimSpace(strings.ToLower(email))
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.users {
		if user.Email == normalized {
			return user, nil
		}
	}
	return domain.User{}, domain.ErrNotFound
}

func (s *Store) GetUserByPhone(_ context.Context, phone string) (domain.User, error) {
	normalized := domain.NormalizePhone(phone)
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.users {
		if user.Phone == normalized {
			return user, nil
		}
	}
	return domain.User{}, domain.ErrNotFound
}

func (s *Store) GetUserByOAuth(_ context.Context, provider, subject string) (domain.User, error) {
	provider = strings.TrimSpace(strings.ToLower(provider))
	subject = strings.TrimSpace(subject)
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.users {
		if user.OAuthProvider == provider && user.OAuthSubject == subject {
			return user, nil
		}
	}
	return domain.User{}, domain.ErrNotFound
}

func (s *Store) UpdateUserProfile(_ context.Context, user domain.User) (domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.users[user.ID]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	existing.Name = user.Name
	existing.Email = user.Email
	existing.Phone = user.Phone
	s.users[user.ID] = existing
	return existing, nil
}

func (s *Store) UpdateUserPassword(_ context.Context, userID uuid.UUID, passwordHash string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.users[userID]
	if !ok {
		return domain.ErrNotFound
	}
	existing.PasswordHash = passwordHash
	existing.UpdatedAt = time.Now().UTC()
	s.users[userID] = existing
	return nil
}
