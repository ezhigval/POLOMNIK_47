package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"palomnik/internal/domain"
)

func (s *Store) ListAdminRoles(ctx context.Context) ([]domain.AdminRole, error) {
	rows, err := s.db.QueryContext(ctx, `
SELECT id, name, password_hash, created_at, updated_at
FROM admin_roles
ORDER BY name
`)
	if err != nil {
		return nil, fmt.Errorf("list admin roles: %w", err)
	}
	defer rows.Close()

	var roles []domain.AdminRole
	for rows.Next() {
		role, err := scanAdminRole(rows)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for i := range roles {
		perms, err := s.listRolePermissions(ctx, roles[i].ID)
		if err != nil {
			return nil, err
		}
		roles[i].Permissions = perms
	}
	return roles, nil
}

func (s *Store) GetAdminRole(ctx context.Context, id uuid.UUID) (domain.AdminRole, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT id, name, password_hash, created_at, updated_at
FROM admin_roles
WHERE id = $1
`, id)
	role, err := scanAdminRole(row)
	if err != nil {
		return domain.AdminRole{}, err
	}
	perms, err := s.listRolePermissions(ctx, role.ID)
	if err != nil {
		return domain.AdminRole{}, err
	}
	role.Permissions = perms
	return role, nil
}

func (s *Store) GetAdminRoleByName(ctx context.Context, name string) (domain.AdminRole, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT id, name, password_hash, created_at, updated_at
FROM admin_roles
WHERE name = $1
`, name)
	role, err := scanAdminRole(row)
	if err != nil {
		return domain.AdminRole{}, err
	}
	perms, err := s.listRolePermissions(ctx, role.ID)
	if err != nil {
		return domain.AdminRole{}, err
	}
	role.Permissions = perms
	return role, nil
}

func (s *Store) CreateAdminRole(ctx context.Context, role domain.AdminRole) (domain.AdminRole, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.AdminRole{}, err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.ExecContext(ctx, `
INSERT INTO admin_roles (id, name, password_hash, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
`, role.ID, role.Name, role.PasswordHash, role.CreatedAt, role.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.AdminRole{}, domain.ErrDuplicateAdminRoleName
		}
		return domain.AdminRole{}, fmt.Errorf("create admin role: %w", err)
	}
	if err := replaceRolePermissionsTx(ctx, tx, role.ID, role.Permissions); err != nil {
		return domain.AdminRole{}, err
	}
	if err := tx.Commit(); err != nil {
		return domain.AdminRole{}, err
	}
	return s.GetAdminRole(ctx, role.ID)
}

func (s *Store) UpdateAdminRole(ctx context.Context, role domain.AdminRole) (domain.AdminRole, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.AdminRole{}, err
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx, `
UPDATE admin_roles
SET password_hash = $2, updated_at = $3
WHERE id = $1
`, role.ID, role.PasswordHash, role.UpdatedAt)
	if err != nil {
		return domain.AdminRole{}, fmt.Errorf("update admin role: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.AdminRole{}, domain.ErrNotFound
	}
	if err := replaceRolePermissionsTx(ctx, tx, role.ID, role.Permissions); err != nil {
		return domain.AdminRole{}, err
	}
	if err := tx.Commit(); err != nil {
		return domain.AdminRole{}, err
	}
	return s.GetAdminRole(ctx, role.ID)
}

func (s *Store) DeleteAdminRole(ctx context.Context, id uuid.UUID) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM admin_roles WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete admin role: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (s *Store) AssignUserToRole(ctx context.Context, assignment domain.AdminRoleAssignment) error {
	_, err := s.db.ExecContext(ctx, `
INSERT INTO admin_role_user_assignments (role_id, user_id, created_at)
VALUES ($1, $2, $3)
ON CONFLICT (role_id, user_id) DO NOTHING
`, assignment.RoleID, assignment.UserID, assignment.CreatedAt)
	if err != nil {
		return fmt.Errorf("assign user to role: %w", err)
	}
	return nil
}

func (s *Store) UnassignUserFromRole(ctx context.Context, roleID, userID uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, `
DELETE FROM admin_role_user_assignments
WHERE role_id = $1 AND user_id = $2
`, roleID, userID)
	if err != nil {
		return fmt.Errorf("unassign user from role: %w", err)
	}
	return nil
}

func (s *Store) ListRoleAssignments(ctx context.Context, roleID uuid.UUID) ([]domain.AdminRoleAssignment, error) {
	rows, err := s.db.QueryContext(ctx, `
SELECT role_id, user_id, created_at
FROM admin_role_user_assignments
WHERE role_id = $1
ORDER BY created_at
`, roleID)
	if err != nil {
		return nil, fmt.Errorf("list role assignments: %w", err)
	}
	defer rows.Close()

	var out []domain.AdminRoleAssignment
	for rows.Next() {
		var item domain.AdminRoleAssignment
		if err := rows.Scan(&item.RoleID, &item.UserID, &item.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (s *Store) listRolePermissions(ctx context.Context, roleID uuid.UUID) ([]domain.Permission, error) {
	rows, err := s.db.QueryContext(ctx, `
SELECT permission FROM admin_role_permissions WHERE role_id = $1 ORDER BY permission
`, roleID)
	if err != nil {
		return nil, fmt.Errorf("list role permissions: %w", err)
	}
	defer rows.Close()
	var out []domain.Permission
	for rows.Next() {
		var raw string
		if err := rows.Scan(&raw); err != nil {
			return nil, err
		}
		out = append(out, domain.Permission(raw))
	}
	return out, rows.Err()
}

func replaceRolePermissionsTx(ctx context.Context, tx *sql.Tx, roleID uuid.UUID, perms []domain.Permission) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM admin_role_permissions WHERE role_id = $1`, roleID); err != nil {
		return err
	}
	for _, perm := range perms {
		if _, err := tx.ExecContext(ctx, `
INSERT INTO admin_role_permissions (role_id, permission) VALUES ($1, $2)
`, roleID, string(perm)); err != nil {
			return err
		}
	}
	return nil
}

func scanAdminRole(row scanner) (domain.AdminRole, error) {
	var role domain.AdminRole
	err := row.Scan(&role.ID, &role.Name, &role.PasswordHash, &role.CreatedAt, &role.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.AdminRole{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.AdminRole{}, fmt.Errorf("scan admin role: %w", err)
	}
	return role, nil
}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}
	return false
}
