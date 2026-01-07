package env

import "testing"

func TestBootstrapPlanValidation(t *testing.T) {
	cases := []struct {
		name    string
		envName string
		wantErr bool
	}{
		{"empty", "", true},
		{"spaces", "   ", true},
		{"invalid chars", "dev!", true},
		{"valid", "dev", false},
		{"valid dash", "prod-east-1", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := BootstrapPlan(tc.envName)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error for %q", tc.envName)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestFormatPlan(t *testing.T) {
	plan := Plan{EnvName: "dev", Steps: []string{"a", "b"}, Notes: []string{"note"}}
	out := FormatPlan(plan)
	if out == "" {
		t.Fatalf("expected output")
	}
	if want := "Bootstrap plan for dev"; out[:len(want)] != want {
		t.Fatalf("output prefix mismatch: %q", out)
	}
}
