package cli

import (
	"fmt"
	"strings"

	"github.com/JamesIOmete/platformctl/internal/plugin"
	"github.com/JamesIOmete/platformctl/internal/version"
)

// PrintHelp renders usage with discovered plugins.
func PrintHelp() {
	usage := fmt.Sprintf(`platformctl %s — internal platform helper CLI (POC)

Usage:
	platformctl <command> [args]

Built-in commands:
	help                     Show this help
	version                  Show CLI version
	doctor                   Environment checks
	auth status              Show mock auth status
	fleet [cmd]              Manage robots (ls, status, logs, ssh)
	secrets [cmd]            Manage secrets (ls, get, set)
	init                     Run interactive setup
	env bootstrap <env>      Dry-run infra bootstrap plan (requires scope infra:write)

Plugins:
	Git-style: executables named platformctl-<command> on PATH are invoked as subcommands.
	Example: platformctl-hello => platformctl hello
`, version.Version)

	fmt.Print(usage)
	plugins := plugin.ListPluginsOnPath()
	if len(plugins) == 0 {
		fmt.Println("Discovered plugins: (none found on PATH)")
		fmt.Println("Tip: add an executable named platformctl-<name> to PATH to extend the CLI.")
		return
	}
	fmt.Println("Discovered plugins:")
	fmt.Println("  " + strings.Join(plugins, ", "))
}
