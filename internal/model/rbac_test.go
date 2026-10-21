package model

import "testing"

func TestRoleHierarchy(t *testing.T) {
	cases := []struct {
		user, required string
		want           bool
	}{
		{RoleViewer, RoleViewer, true},
		{RoleViewer, RoleOperator, false},
		{RoleViewer, RoleAdministrator, false},
		{RoleOperator, RoleViewer, true},
		{RoleOperator, RoleOperator, true},
		{RoleOperator, RoleAdministrator, false},
		{RoleAdministrator, RoleViewer, true},
		{RoleAdministrator, RoleOperator, true},
		{RoleAdministrator, RoleAdministrator, true},
	}
	for _, c := range cases {
		if got := RequireRole(c.user, c.required); got != c.want {
			t.Errorf("RequireRole(%q, %q) = %v, want %v", c.user, c.required, got, c.want)
		}
	}
}

func TestRoleLevelOrdering(t *testing.T) {
	if RoleLevel[RoleViewer] >= RoleLevel[RoleOperator] {
		t.Error("Viewer must be below Operator")
	}
	if RoleLevel[RoleOperator] >= RoleLevel[RoleAdministrator] {
		t.Error("Operator must be below Administrator")
	}
}

func TestRequireRoleUnknownRoleDenied(t *testing.T) {
	if RequireRole("Bogus", RoleViewer) {
		t.Error("an unknown role must not satisfy any requirement")
	}
}
