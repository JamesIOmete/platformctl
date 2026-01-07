package env

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Plan describes a dry-run bootstrap plan.
type Plan struct {
	EnvName string
	Steps   []string
	Notes   []string
}

var envPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// BootstrapPlan validates envName and returns a mocked plan.
func BootstrapPlan(envName string) (Plan, error) {
	if strings.TrimSpace(envName) == "" {
		return Plan{}, errors.New("environment name is required")
	}
	if !envPattern.MatchString(envName) {
		return Plan{}, fmt.Errorf("invalid environment name: %s", envName)
	}

	steps := []string{
		fmt.Sprintf("Validate credentials for env %s", envName),
		"Check Terraform state access (dry-run)",
		fmt.Sprintf("Render Terraform plan for %s (no apply)", envName),
		"Validate network prerequisites",
		"Generate kubeconfig and store in ~/.kube (skipped in POC)",
	}
	notes := []string{
		"This is a dry-run. No changes are applied.",
		"Requires scope infra:write (mocked).",
	}
	return Plan{EnvName: envName, Steps: steps, Notes: notes}, nil
}

// FormatPlan renders the plan text.
func FormatPlan(p Plan) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Bootstrap plan for %s (dry-run)\n", p.EnvName)
	fmt.Fprintln(&b, strings.Repeat("-", 40))
	for i, step := range p.Steps {
		fmt.Fprintf(&b, "%d. %s\n", i+1, step)
	}
	if len(p.Notes) > 0 {
		fmt.Fprintln(&b, "Notes:")
		for _, note := range p.Notes {
			fmt.Fprintf(&b, "- %s\n", note)
		}
	}
	return b.String()
}
