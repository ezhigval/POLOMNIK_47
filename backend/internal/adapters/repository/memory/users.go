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

func (s *Store) GetUserByOAuth(ctx context.Context, provider, subject string) (domain.User, error) {
	identity, err := s.GetIdentity(ctx, provider, subject)
	if err != nil {
		return domain.User{}, err
	}
	return s.GetUserByID(ctx, identity.UserID)
}

func (s *Store) GetIdentity(_ context.Context, provider, subject string) (domain.UserIdentity, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	identity, ok := s.identities[identityKey(provider, subject)]
	if !ok {
		return domain.UserIdentity{}, domain.ErrNotFound
	}
	return identity, nil
}

func (s *Store) ListIdentities(_ context.Context, userID uuid.UUID) ([]domain.UserIdentity, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]domain.UserIdentity, 0)
	for _, identity := range s.identities {
		if identity.UserID == userID {
			out = append(out, identity)
		}
	}
	return out, nil
}

func (s *Store) AddIdentity(_ context.Context, identity domain.UserIdentity) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[identity.UserID]; !ok {
		return domain.ErrNotFound
	}
	key := identityKey(identity.Provider, identity.Subject)
	if existing, ok := s.identities[key]; ok {
		if existing.UserID == identity.UserID {
			return nil
		}
		return domain.ErrDuplicateIdentity
	}
	s.identities[key] = identity
	return nil
}

func (s *Store) UpdateUserProfile(_ context.Context, user domain.User) (domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.users[user.ID]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	for _, other := range s.users {
		if other.ID == user.ID {
			continue
		}
		if user.Phone != "" && other.Phone == user.Phone {
			return domain.User{}, domain.ErrDuplicatePhone
		}
		if user.Email != "" && other.Email == user.Email {
			return domain.User{}, domain.ErrDuplicateEmail
		}
	}
	existing.Name = user.Name
	existing.Email = user.Email
	existing.Phone = user.Phone
	existing.UpdatedAt = time.Now().UTC()
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

func (s *Store) MergeAccountInto(_ context.Context, targetID, sourceID uuid.UUID) error {
	if targetID == sourceID {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[targetID]; !ok {
		return domain.ErrNotFound
	}
	if _, ok := s.users[sourceID]; !ok {
		return domain.ErrNotFound
	}

	for id, booking := range s.bookings {
		if booking.UserID != nil && *booking.UserID == sourceID {
			uid := targetID
			booking.UserID = &uid
			s.bookings[id] = booking
		}
	}

	for key, favorite := range s.favorites {
		if favorite.UserID != sourceID {
			continue
		}
		delete(s.favorites, key)
		targetKey := favoriteKey(targetID, favorite.TourID)
		if _, exists := s.favorites[targetKey]; exists {
			continue
		}
		favorite.UserID = targetID
		s.favorites[targetKey] = favorite
	}

	var targetOpen, sourceOpen *domain.SupportThread
	for _, thread := range s.supportThreads {
		thread := thread
		if thread.Status != "open" {
			continue
		}
		switch thread.UserID {
		case targetID:
			targetOpen = &thread
		case sourceID:
			sourceOpen = &thread
		}
	}
	if targetOpen != nil && sourceOpen != nil && targetOpen.ID != sourceOpen.ID {
		for id, message := range s.supportMessages {
			if message.ThreadID == sourceOpen.ID {
				message.ThreadID = targetOpen.ID
				s.supportMessages[id] = message
			}
		}
		delete(s.supportThreads, sourceOpen.ID)
	}
	for id, thread := range s.supportThreads {
		if thread.UserID == sourceID {
			thread.UserID = targetID
			s.supportThreads[id] = thread
		}
	}

	keptRoles := make(map[uuid.UUID]struct{})
	for _, item := range s.adminAssignments {
		if item.UserID == targetID {
			keptRoles[item.RoleID] = struct{}{}
		}
	}
	filtered := make([]domain.AdminRoleAssignment, 0, len(s.adminAssignments))
	for _, item := range s.adminAssignments {
		if item.UserID == sourceID {
			if _, exists := keptRoles[item.RoleID]; exists {
				continue
			}
			item.UserID = targetID
			keptRoles[item.RoleID] = struct{}{}
			filtered = append(filtered, item)
			continue
		}
		filtered = append(filtered, item)
	}
	s.adminAssignments = filtered

	for key, identity := range s.identities {
		if identity.UserID == sourceID {
			identity.UserID = targetID
			s.identities[key] = identity
		}
	}

	mergePassengersLocked(s, targetID, sourceID)

	delete(s.users, sourceID)
	return nil
}

func identityKey(provider, subject string) string {
	return domain.NormalizeOAuthProvider(provider) + "\x00" + strings.TrimSpace(subject)
}
