package domain

import (
	"slices"
	"testing"
)

func TestRoleTemplatesMatchV3PlanNames(t *testing.T) {
	templates := RoleTemplates()
	if len(templates) != 5 {
		t.Fatalf("expected 5 templates from V3_PLAN, got %d", len(templates))
	}

	want := map[string]struct {
		label string
		perms []Permission
	}{
		"booking_manager": {label: "Менеджер заявок", perms: []Permission{PermManageBookings, PermManageSupport}},
		"advertiser":      {label: "Рекламщик", perms: []Permission{PermViewStats}},
		"smm":             {label: "Сммщик", perms: []Permission{PermManageContent}},
		"developer":       {label: "Разработчик", perms: []Permission{PermManageIntegrations, PermViewStats}},
	}

	seen := map[string]bool{}
	for _, tmpl := range templates {
		seen[tmpl.ID] = true
		if _, err := NormalizeAdminRoleName(tmpl.RoleName); err != nil {
			t.Fatalf("template %s role name %q: %v", tmpl.ID, tmpl.RoleName, err)
		}
		if slices.Contains(tmpl.Permissions, PermManageRoles) {
			t.Fatalf("template %s must not include manage_roles", tmpl.ID)
		}
		normalized, err := NormalizePermissions(tmpl.Permissions)
		if err != nil {
			t.Fatalf("template %s permissions: %v", tmpl.ID, err)
		}
		if len(normalized) != len(tmpl.Permissions) {
			t.Fatalf("template %s dropped permissions: %v", tmpl.ID, tmpl.Permissions)
		}
		if tmpl.ID == "director" {
			if tmpl.Label != "Директор" {
				t.Fatalf("director label: %q", tmpl.Label)
			}
			wantPerms := AssignablePermissions()
			if len(tmpl.Permissions) != len(wantPerms) {
				t.Fatalf("director perms %v, want %v", tmpl.Permissions, wantPerms)
			}
			for i, p := range wantPerms {
				if tmpl.Permissions[i] != p {
					t.Fatalf("director perms %v, want %v", tmpl.Permissions, wantPerms)
				}
			}
			continue
		}
		spec, ok := want[tmpl.ID]
		if !ok {
			t.Fatalf("unexpected template id %q", tmpl.ID)
		}
		if tmpl.Label != spec.label {
			t.Fatalf("template %s label %q, want %q", tmpl.ID, tmpl.Label, spec.label)
		}
		if len(tmpl.Permissions) != len(spec.perms) {
			t.Fatalf("template %s perms %v, want %v", tmpl.ID, tmpl.Permissions, spec.perms)
		}
		for i, p := range spec.perms {
			if tmpl.Permissions[i] != p {
				t.Fatalf("template %s perms %v, want %v", tmpl.ID, tmpl.Permissions, spec.perms)
			}
		}
	}
	if !seen["director"] {
		t.Fatal("missing director template")
	}
}

func TestAssignablePermissionsOmitManageRoles(t *testing.T) {
	for _, p := range AssignablePermissions() {
		if p == PermManageRoles {
			t.Fatal("assignable permissions must omit manage_roles")
		}
		if !ValidPermission(p) {
			t.Fatalf("unexpected permission %q", p)
		}
	}
	if len(AssignablePermissions()) != len(AllPermissions())-1 {
		t.Fatalf("assignable count %d, all %d", len(AssignablePermissions()), len(AllPermissions()))
	}
}
