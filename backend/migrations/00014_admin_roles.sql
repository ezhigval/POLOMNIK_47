-- +goose Up
CREATE TABLE IF NOT EXISTS admin_roles (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS admin_role_permissions (
    role_id UUID NOT NULL REFERENCES admin_roles(id) ON DELETE CASCADE,
    permission TEXT NOT NULL,
    PRIMARY KEY (role_id, permission)
);

CREATE TABLE IF NOT EXISTS admin_role_user_assignments (
    role_id UUID NOT NULL REFERENCES admin_roles(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (role_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_admin_role_user_assignments_user
    ON admin_role_user_assignments (user_id);

-- +goose Down
DROP TABLE IF EXISTS admin_role_user_assignments;
DROP TABLE IF EXISTS admin_role_permissions;
DROP TABLE IF EXISTS admin_roles;
