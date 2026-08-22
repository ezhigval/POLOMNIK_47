package memory

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"polomnik/internal/domain"
)

func (s *Store) ListAdminRoles(_ context.Context) ([]domain.AdminRole, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]domain.AdminRole, 0, len(s.adminRoles))
	for _, role := range s.adminRoles {
		out = append(out, cloneAdminRole(role))
	}
	return out, nil
}

func (s *Store) GetAdminRole(_ context.Context, id uuid.UUID) (domain.AdminRole, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	role, ok := s.adminRoles[id]
	if !ok {
		return domain.AdminRole{}, domain.ErrNotFound
	}
	return cloneAdminRole(role), nil
}

func (s *Store) GetAdminRoleByName(_ context.Context, name string) (domain.AdminRole, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := strings.ToLower(strings.TrimSpace(name))
	for _, role := range s.adminRoles {
		if role.Name == key {
			return cloneAdminRole(role), nil
		}
	}
	return domain.AdminRole{}, domain.ErrNotFound
}

func (s *Store) CreateAdminRole(_ context.Context, role domain.AdminRole) (domain.AdminRole, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.adminRoles == nil {
		s.adminRoles = make(map[uuid.UUID]domain.AdminRole)
	}
	for _, existing := range s.adminRoles {
		if existing.Name == role.Name {
			return domain.AdminRole{}, domain.ErrDuplicateAdminRoleName
		}
	}
	s.adminRoles[role.ID] = cloneAdminRole(role)
	return cloneAdminRole(role), nil
}

func (s *Store) UpdateAdminRole(_ context.Context, role domain.AdminRole) (domain.AdminRole, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.adminRoles[role.ID]; !ok {
		return domain.AdminRole{}, domain.ErrNotFound
	}
	s.adminRoles[role.ID] = cloneAdminRole(role)
	return cloneAdminRole(role), nil
}

func (s *Store) DeleteAdminRole(_ context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.adminRoles[id]; !ok {
		return domain.ErrNotFound
	}
	delete(s.adminRoles, id)
	if s.adminAssignments != nil {
		filtered := s.adminAssignments[:0]
		for _, item := range s.adminAssignments {
			if item.RoleID != id {
				filtered = append(filtered, item)
			}
		}
		s.adminAssignments = filtered
	}
	return nil
}

func (s *Store) AssignUserToRole(_ context.Context, assignment domain.AdminRoleAssignment) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.adminRoles[assignment.RoleID]; !ok {
		return domain.ErrNotFound
	}
	for _, item := range s.adminAssignments {
		if item.RoleID == assignment.RoleID && item.UserID == assignment.UserID {
			return nil
		}
	}
	s.adminAssignments = append(s.adminAssignments, assignment)
	return nil
}

func (s *Store) UnassignUserFromRole(_ context.Context, roleID, userID uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filtered := s.adminAssignments[:0]
	for _, item := range s.adminAssignments {
		if item.RoleID == roleID && item.UserID == userID {
			continue
		}
		filtered = append(filtered, item)
	}
	s.adminAssignments = filtered
	return nil
}

func (s *Store) ListRoleAssignments(_ context.Context, roleID uuid.UUID) ([]domain.AdminRoleAssignment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var out []domain.AdminRoleAssignment
	for _, item := range s.adminAssignments {
		if item.RoleID == roleID {
			out = append(out, item)
		}
	}
	return out, nil
}

func cloneAdminRole(role domain.AdminRole) domain.AdminRole {
	role.Permissions = append([]domain.Permission(nil), role.Permissions...)
	return role
}
