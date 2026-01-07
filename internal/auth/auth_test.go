package auth

import "testing"

func TestHasScope(t *testing.T) {
	cases := []struct {
		name   string
		status Status
		scope  string
		want   bool
	}{
		{"empty", Status{}, "fleet:read", false},
		{"present", Status{Scopes: []string{"fleet:read", "infra:write"}}, "fleet:read", true},
		{"absent", Status{Scopes: []string{"infra:write"}}, "fleet:read", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := HasScope(tc.status, tc.scope)
			if got != tc.want {
				t.Fatalf("HasScope(%v,%s)=%v want %v", tc.status.Scopes, tc.scope, got, tc.want)
			}
		})
	}
}
