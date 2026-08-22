package application

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

const managementSessionTTL = 8 * time.Hour

type AdminRoleService struct {
	roles     ports.AdminRoleRepository
	users     ports.UserRepository
	adminTok  string
	jwtSecret []byte
}

func NewAdminRoleService(
	roles ports.AdminRoleRepository,
	users ports.UserRepository,
	adminToken string,
	jwtSecret string,
) *AdminRoleService {
	return &AdminRoleService{
		roles:     roles,
		users:     users,
		adminTok:  strings.TrimSpace(adminToken),
		jwtSecret: []byte(jwtSecret),
	}
}

type ManagementPrincipal struct {
	FullAdmin   bool
	RoleID      uuid.UUID
	RoleName    string
	Permissions []domain.Permission
}

type ManagementLoginResult struct {
	Token       string
	FullAdmin   bool
	RoleID      uuid.UUID
	RoleName    string
	Permissions []domain.Permission
}

type managementClaims struct {
	FullAdmin bool     `json:"full"`
	RoleID    string   `json:"role_id,omitempty"`
	RoleName  string   `json:"role_name,omitempty"`
	Perms     []string `json:"perms,omitempty"`
	jwt.RegisteredClaims
}

func (s *AdminRoleService) Login(ctx context.Context, roleName, password string) (ManagementLoginResult, error) {
	roleName = strings.TrimSpace(roleName)
	password = strings.TrimSpace(password)
	if password == "" {
		return ManagementLoginResult{}, domain.ErrInvalidCredentials
	}

	// Full admin: empty role + ADMIN_TOKEN (env only, never stored in DB).
	if roleName == "" {
		if s.adminTok == "" || password != s.adminTok {
			return ManagementLoginResult{}, domain.ErrInvalidCredentials
		}
		token, err := s.issueToken(ManagementPrincipal{
			FullAdmin:   true,
			Permissions: domain.AllPermissions(),
		})
		if err != nil {
			return ManagementLoginResult{}, err
		}
		return ManagementLoginResult{
			Token:       token,
			FullAdmin:   true,
			Permissions: domain.AllPermissions(),
		}, nil
	}

	if s.roles == nil {
		return ManagementLoginResult{}, domain.ErrInvalidCredentials
	}
	normalized, err := domain.NormalizeAdminRoleName(roleName)
	if err != nil {
		return ManagementLoginResult{}, domain.ErrInvalidCredentials
	}
	role, err := s.roles.GetAdminRoleByName(ctx, normalized)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return ManagementLoginResult{}, domain.ErrInvalidCredentials
		}
		return ManagementLoginResult{}, err
	}
	if !domain.VerifyAdminRolePassword(role.PasswordHash, password) {
		return ManagementLoginResult{}, domain.ErrInvalidCredentials
	}
	principal := ManagementPrincipal{
		FullAdmin:   false,
		RoleID:      role.ID,
		RoleName:    role.Name,
		Permissions: append([]domain.Permission(nil), role.Permissions...),
	}
	token, err := s.issueToken(principal)
	if err != nil {
		return ManagementLoginResult{}, err
	}
	return ManagementLoginResult{
		Token:       token,
		FullAdmin:   false,
		RoleID:      role.ID,
		RoleName:    role.Name,
		Permissions: principal.Permissions,
	}, nil
}

func (s *AdminRoleService) ParseSession(token string) (ManagementPrincipal, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return ManagementPrincipal{}, domain.ErrInvalidCredentials
	}
	parsed, err := jwt.ParseWithClaims(token, &managementClaims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, domain.ErrInvalidCredentials
		}
		return s.jwtSecret, nil
	})
	if err != nil || !parsed.Valid {
		return ManagementPrincipal{}, domain.ErrInvalidCredentials
	}
	claims, ok := parsed.Claims.(*managementClaims)
	if !ok {
		return ManagementPrincipal{}, domain.ErrInvalidCredentials
	}
	if claims.FullAdmin {
		return ManagementPrincipal{
			FullAdmin:   true,
			Permissions: domain.AllPermissions(),
		}, nil
	}
	roleID, err := uuid.Parse(claims.RoleID)
	if err != nil || roleID == uuid.Nil {
		return ManagementPrincipal{}, domain.ErrInvalidCredentials
	}
	perms := make([]domain.Permission, 0, len(claims.Perms))
	for _, p := range claims.Perms {
		perm := domain.Permission(p)
		if domain.ValidPermission(perm) && perm != domain.PermManageRoles {
			perms = append(perms, perm)
		}
	}
	return ManagementPrincipal{
		FullAdmin:   false,
		RoleID:      roleID,
		RoleName:    claims.RoleName,
		Permissions: perms,
	}, nil
}

func (s *AdminRoleService) AuthenticateHeader(adminTokenHeader, sessionHeader string) (ManagementPrincipal, error) {
	adminTokenHeader = strings.TrimSpace(adminTokenHeader)
	sessionHeader = strings.TrimSpace(sessionHeader)
	if s.adminTok != "" && adminTokenHeader != "" && adminTokenHeader == s.adminTok {
		return ManagementPrincipal{
			FullAdmin:   true,
			Permissions: domain.AllPermissions(),
		}, nil
	}
	if sessionHeader != "" {
		return s.ParseSession(sessionHeader)
	}
	return ManagementPrincipal{}, domain.ErrInvalidCredentials
}

func (s *AdminRoleService) HasPermission(p ManagementPrincipal, perm domain.Permission) bool {
	if p.FullAdmin {
		return true
	}
	for _, item := range p.Permissions {
		if item == perm {
			return true
		}
	}
	return false
}

func (s *AdminRoleService) ListRoles(ctx context.Context) ([]domain.AdminRole, error) {
	if s.roles == nil {
		return nil, nil
	}
	roles, err := s.roles.ListAdminRoles(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]domain.AdminRole, 0, len(roles))
	for _, role := range roles {
		out = append(out, domain.PublicAdminRole(role))
	}
	return out, nil
}

func (s *AdminRoleService) CreateRole(ctx context.Context, name, password string, permissions []domain.Permission) (domain.AdminRole, error) {
	if s.roles == nil {
		return domain.AdminRole{}, domain.ErrForbidden
	}
	role, err := domain.NewAdminRole(uuid.New(), name, password, permissions, time.Time{})
	if err != nil {
		return domain.AdminRole{}, err
	}
	created, err := s.roles.CreateAdminRole(ctx, role)
	if err != nil {
		return domain.AdminRole{}, err
	}
	return domain.PublicAdminRole(created), nil
}

func (s *AdminRoleService) UpdateRole(
	ctx context.Context,
	id uuid.UUID,
	permissions []domain.Permission,
	newPassword string,
) (domain.AdminRole, error) {
	if s.roles == nil {
		return domain.AdminRole{}, domain.ErrForbidden
	}
	role, err := s.roles.GetAdminRole(ctx, id)
	if err != nil {
		return domain.AdminRole{}, err
	}
	if permissions != nil {
		perms, err := domain.NormalizePermissions(permissions)
		if err != nil {
			return domain.AdminRole{}, err
		}
		role.Permissions = perms
	}
	if strings.TrimSpace(newPassword) != "" {
		hash, err := domain.HashAdminRolePassword(newPassword)
		if err != nil {
			return domain.AdminRole{}, err
		}
		role.PasswordHash = hash
	}
	role.UpdatedAt = time.Now().UTC()
	updated, err := s.roles.UpdateAdminRole(ctx, role)
	if err != nil {
		return domain.AdminRole{}, err
	}
	return domain.PublicAdminRole(updated), nil
}

func (s *AdminRoleService) DeleteRole(ctx context.Context, id uuid.UUID) error {
	if s.roles == nil {
		return domain.ErrForbidden
	}
	return s.roles.DeleteAdminRole(ctx, id)
}

func (s *AdminRoleService) AssignUser(ctx context.Context, roleID, userID uuid.UUID) error {
	if s.roles == nil {
		return domain.ErrForbidden
	}
	if roleID == uuid.Nil || userID == uuid.Nil {
		return domain.ErrInvalidID
	}
	if s.users != nil {
		if _, err := s.users.GetUserByID(ctx, userID); err != nil {
			return err
		}
	}
	if _, err := s.roles.GetAdminRole(ctx, roleID); err != nil {
		return err
	}
	return s.roles.AssignUserToRole(ctx, domain.AdminRoleAssignment{
		RoleID:    roleID,
		UserID:    userID,
		CreatedAt: time.Now().UTC(),
	})
}

func (s *AdminRoleService) UnassignUser(ctx context.Context, roleID, userID uuid.UUID) error {
	if s.roles == nil {
		return domain.ErrForbidden
	}
	return s.roles.UnassignUserFromRole(ctx, roleID, userID)
}

func (s *AdminRoleService) ListAssignments(ctx context.Context, roleID uuid.UUID) ([]domain.AdminRoleAssignment, error) {
	if s.roles == nil {
		return nil, nil
	}
	return s.roles.ListRoleAssignments(ctx, roleID)
}

func (s *AdminRoleService) issueToken(p ManagementPrincipal) (string, error) {
	now := time.Now().UTC()
	claims := managementClaims{
		FullAdmin: p.FullAdmin,
		RoleName:  p.RoleName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(managementSessionTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   "management",
		},
	}
	if !p.FullAdmin {
		claims.RoleID = p.RoleID.String()
		claims.Perms = make([]string, 0, len(p.Permissions))
		for _, perm := range p.Permissions {
			claims.Perms = append(claims.Perms, string(perm))
		}
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
