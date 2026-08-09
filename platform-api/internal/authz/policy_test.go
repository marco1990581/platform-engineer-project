package authz

import "testing"

func TestEnforce(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		resource string
		action   string
		want     bool
	}{
		{name: "admin can read system", role: "admin", resource: "system", action: "read", want: true},
		{name: "admin can write users", role: "admin", resource: "users", action: "write", want: true},
		{name: "viewer can read system", role: "viewer", resource: "system", action: "read", want: true},
		{name: "viewer can read network", role: "viewer", resource: "network", action: "read", want: true},
		{name: "viewer can read filesystem", role: "viewer", resource: "filesystem", action: "read", want: true},
		{name: "viewer cannot read users", role: "viewer", resource: "users", action: "read", want: false},
		{name: "viewer cannot write system", role: "viewer", resource: "system", action: "write", want: false},
		{name: "unknown role denied", role: "unknown", resource: "system", action: "read", want: false},
		{name: "empty role denied", role: "", resource: "system", action: "read", want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Enforce(test.role, test.resource, test.action)
			if got != test.want {
				t.Errorf("Enforce(%q, %q, %q) = %v, want %v", test.role, test.resource, test.action, got, test.want)
			}
		})
	}
}
