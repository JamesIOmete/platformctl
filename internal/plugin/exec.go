package plugin

import (
	"fmt"
	"os"
	"os/exec"
)

// ExecPlugin invokes a Git-style plugin executable named "platformctl-<name>".
func ExecPlugin(name string, args []string) error {
	bin := "platformctl-" + name
	path, err := exec.LookPath(bin)
	if err != nil {
		return fmt.Errorf("unknown command %q (no plugin %q found on PATH)", name, bin)
	}

	cmd := exec.Command(path, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
