package auth

import (
	"fmt"
	"sort"
	"strings"

	"github.com/JamesIOmete/platformctl/internal/config"
)

// Status represents mock authentication state.
type Status struct {
	Authenticated bool
	Principal     string
	Scopes        []string
}

// LoadStatus reads mock auth config from disk. Missing config => unauthenticated.
func LoadStatus() Status {
	cfg, _, err := config.Load()
	if err != nil {
		// On parse error, surface unauthenticated with detail in principal field.
		return Status{Authenticated: false, Principal: "invalid-config", Scopes: nil}
	}
	if cfg.Principal == "" {
		return Status{Authenticated: false, Principal: "unknown", Scopes: nil}
	}
	scopes := append([]string(nil), cfg.Scopes...)
	sort.Strings(scopes)
	return Status{Authenticated: true, Principal: cfg.Principal, Scopes: scopes}
}

// HasScope checks for an exact scope.
func HasScope(s Status, scope string) bool {
	for _, sc := range s.Scopes {
		if sc == scope {
			return true
		}
	}
	return false
}

// FormatStatus renders human-friendly text.
func FormatStatus(s Status) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Authenticated: %v\n", s.Authenticated)
	fmt.Fprintf(&b, "Principal: %s\n", s.Principal)
	if len(s.Scopes) == 0 {
		fmt.Fprintf(&b, "Scopes: (none)\n")
	} else {
		fmt.Fprintf(&b, "Scopes: %s\n", strings.Join(s.Scopes, ", "))
	}
	return b.String()
}
