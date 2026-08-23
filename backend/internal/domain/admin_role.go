package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Permission string

const (
	PermManageTours        Permission = "manage_tours"
	PermManageBookings     Permission = "manage_bookings"
	PermManageContent      Permission = "manage_content"
	PermManageSupport      Permission = "manage_support"
	PermManageSettingsSite Permission = "manage_settings_site"
	PermManageRecipients   Permission = "manage_recipients"
	PermManageRoles        Permission = "manage_roles"
	PermViewStats          Permission = "view_stats"
	PermManageIntegrations Permission = "manage_integrations"
)

func AllPermissions() []Permission {
	return []Permission{
		PermManageTours,
		PermManageBookings,
		PermManageContent,
		PermManageSupport,
		PermManageSettingsSite,
		PermManageRecipients,
		PermManageRoles,
		PermViewStats,
		PermManageIntegrations,
	}
}

func ValidPermission(p Permission) bool {
	for _, item := range AllPermissions() {
		if item == p {
			return true
		}
	}
	return false
}

// AssignablePermissions are rights a DB role may hold. manage_roles stays on ADMIN_TOKEN only.
func AssignablePermissions() []Permission {
	out := make([]Permission, 0, len(AllPermissions())-1)
	for _, p := range AllPermissions() {
		if p == PermManageRoles {
			continue
		}
		out = append(out, p)
	}
	return out
}

// RoleTemplate is a named preset for the create-role form. It does not insert a row.
type RoleTemplate struct {
	ID          string
	Label       string
	RoleName    string
	Permissions []Permission
}

// RoleTemplates are the stage-2 names from V3_PLAN. Permission sets are the
// existing rights those jobs already use in the admin UI (not new business rules).
func RoleTemplates() []RoleTemplate {
	return []RoleTemplate{
		{
			ID:          "booking_manager",
			Label:       "Менеджер заявок",
			RoleName:    "booking_manager",
			Permissions: []Permission{PermManageBookings, PermManageSupport},
		},
		{
			ID:          "advertiser",
			Label:       "Рекламщик",
			RoleName:    "advertiser",
			Permissions: []Permission{PermViewStats},
		},
		{
			ID:          "smm",
			Label:       "Сммщик",
			RoleName:    "smm",
			Permissions: []Permission{PermManageContent},
		},
		{
			ID:          "director",
			Label:       "Директор",
			RoleName:    "director",
			Permissions: AssignablePermissions(),
		},
		{
			ID:          "developer",
			Label:       "Разработчик",
			RoleName:    "developer",
			Permissions: []Permission{PermManageIntegrations, PermViewStats},
		},
	}
}

type AdminRole struct {
	ID           uuid.UUID
	Name         string
	PasswordHash string
	Permissions  []Permission
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type AdminRoleAssignment struct {
	RoleID    uuid.UUID
	UserID    uuid.UUID
	CreatedAt time.Time
}

const (
	minAdminRoleNameLen = 2
	maxAdminRoleNameLen = 64
	minAdminRolePassLen = 8
)

func NormalizeAdminRoleName(raw string) (string, error) {
	name := strings.TrimSpace(strings.ToLower(raw))
	if len(name) < minAdminRoleNameLen || len(name) > maxAdminRoleNameLen {
		return "", ErrInvalidAdminRoleName
	}
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			continue
		}
		return "", ErrInvalidAdminRoleName
	}
	return name, nil
}

func HashAdminRolePassword(password string) (string, error) {
	if len(strings.TrimSpace(password)) < minAdminRolePassLen {
		return "", ErrInvalidPassword
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func VerifyAdminRolePassword(hash, password string) bool {
	if strings.TrimSpace(hash) == "" || password == "" {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func NormalizePermissions(raw []Permission) ([]Permission, error) {
	seen := make(map[Permission]struct{})
	out := make([]Permission, 0, len(raw))
	for _, item := range raw {
		p := Permission(strings.TrimSpace(string(item)))
		if !ValidPermission(p) {
			return nil, ErrInvalidPermission
		}
		if _, ok := seen[p]; ok {
			continue
		}
		// Roles must not receive manage_roles — only full env admin mutates RBAC.
		if p == PermManageRoles {
			continue
		}
		seen[p] = struct{}{}
		out = append(out, p)
	}
	return out, nil
}

func RoleHasPermission(role AdminRole, permission Permission) bool {
	for _, p := range role.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

func NewAdminRole(id uuid.UUID, name, password string, permissions []Permission, now time.Time) (AdminRole, error) {
	if id == uuid.Nil {
		return AdminRole{}, ErrInvalidID
	}
	normalizedName, err := NormalizeAdminRoleName(name)
	if err != nil {
		return AdminRole{}, err
	}
	hash, err := HashAdminRolePassword(password)
	if err != nil {
		return AdminRole{}, err
	}
	perms, err := NormalizePermissions(permissions)
	if err != nil {
		return AdminRole{}, err
	}
	if now.IsZero() {
		now = time.Now().UTC()
	}
	return AdminRole{
		ID:           id,
		Name:         normalizedName,
		PasswordHash: hash,
		Permissions:  perms,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func PublicAdminRole(role AdminRole) AdminRole {
	role.PasswordHash = ""
	role.Permissions = append([]Permission(nil), role.Permissions...)
	return role
}
