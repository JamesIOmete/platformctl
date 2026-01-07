package doctor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/JamesIOmete/platformctl/internal/config"
)

// Report holds doctor check results.
type Report struct {
	Checks []CheckResult
}

// CheckResult represents a single check.
type CheckResult struct {
	Name   string
	Status string // "OK" | "WARN" | "FAIL"
	Detail string
}

// Run executes checks.
func Run() Report {
	checks := []CheckResult{
		checkCommand("Go", "go", "version"),
		checkCommand("Git", "git", "--version"),
		checkConfigPath(),
		checkShell(),
	}
	return Report{Checks: checks}
}

// Format renders the report.
func Format(r Report) string {
	var b strings.Builder
	fmt.Fprintln(&b, "platformctl doctor")
	fmt.Fprintln(&b, strings.Repeat("-", 20))
	for _, c := range r.Checks {
		fmt.Fprintf(&b, "[%s] %s: %s\n", c.Status, c.Name, c.Detail)
	}
	return b.String()
}

func checkCommand(name string, bin string, arg string) CheckResult {
	cmd := exec.Command(bin, arg)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return CheckResult{Name: name, Status: "WARN", Detail: fmt.Sprintf("%s not available", bin)}
	}
	return CheckResult{Name: name, Status: "OK", Detail: strings.TrimSpace(string(out))}
}

func checkConfigPath() CheckResult {
	paths := config.DefaultPaths()
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return CheckResult{Name: "Config", Status: "OK", Detail: fmt.Sprintf("found %s", p)}
		}
	}
	if len(paths) == 0 {
		return CheckResult{Name: "Config", Status: "WARN", Detail: "could not determine config path"}
	}
	return CheckResult{Name: "Config", Status: "WARN", Detail: fmt.Sprintf("not found (%s)", strings.Join(paths, ", "))}
}

func checkShell() CheckResult {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return CheckResult{Name: "Shell", Status: "WARN", Detail: "SHELL not set"}
	}
	return CheckResult{Name: "Shell", Status: "OK", Detail: fmt.Sprintf("%s (GOOS=%s GOARCH=%s)", shell, runtime.GOOS, runtime.GOARCH)}
}
