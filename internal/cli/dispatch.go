package cli

import (
	"fmt"

	"github.com/JamesIOmete/platformctl/internal/auth"
	"github.com/JamesIOmete/platformctl/internal/doctor"
	"github.com/JamesIOmete/platformctl/internal/env"
	"github.com/JamesIOmete/platformctl/internal/fleet"
	"github.com/JamesIOmete/platformctl/internal/output"
	"github.com/JamesIOmete/platformctl/internal/plugin"
	"github.com/JamesIOmete/platformctl/internal/version"
)

// Run dispatches CLI args to built-ins or plugins. Returns exit code.
func Run(args []string) int {
	if len(args) == 0 {
		PrintHelp()
		return 0
	}

	switch args[0] {
	case "-h", "--help", "help":
		PrintHelp()
		return 0
	case "version":
		fmt.Println(version.Version)
		return 0
	case "doctor":
		report := doctor.Run()
		fmt.Print(doctor.Format(report))
		return 0
	case "auth":
		return handleAuth(args[1:])
	case "fleet":
		return handleFleet(args[1:])
	case "env":
		return handleEnv(args[1:])
	default:
		if err := plugin.ExecPlugin(args[0], args[1:]); err != nil {
			output.PrintError(err.Error())
			output.PrintError("Tip: run platformctl help to see built-ins and installed plugins.")
			return 1
		}
		return 0
	}
}

func handleAuth(args []string) int {
	if len(args) == 0 || args[0] == "help" {
		fmt.Println("Usage: platformctl auth status")
		return 0
	}
	if args[0] != "status" {
		output.PrintError("Unknown auth subcommand. Try: platformctl auth status")
		return 1
	}
	status := auth.LoadStatus()
	fmt.Print(auth.FormatStatus(status))
	return 0
}

func handleFleet(args []string) int {
	if len(args) == 0 || args[0] == "help" {
		fmt.Println("Usage: platformctl fleet ls")
		return 0
	}
	if args[0] != "ls" {
		output.PrintError("Unknown fleet subcommand. Try: platformctl fleet ls")
		return 1
	}
	status := auth.LoadStatus()
	if !auth.HasScope(status, "fleet:read") {
		output.PrintError("Access denied: missing scope fleet:read. See platformctl auth status.")
		return 1
	}
	devices := fleet.ListDevices()
	fmt.Print(fleet.FormatDevices(devices))
	return 0
}

func handleEnv(args []string) int {
	if len(args) == 0 || args[0] == "help" {
		fmt.Println("Usage: platformctl env bootstrap <env>")
		return 0
	}
	if args[0] != "bootstrap" {
		output.PrintError("Unknown env subcommand. Try: platformctl env bootstrap <env>")
		return 1
	}
	if len(args) < 2 {
		output.PrintError("Environment name required. Example: platformctl env bootstrap dev")
		return 1
	}
	envName := args[1]
	status := auth.LoadStatus()
	if !auth.HasScope(status, "infra:write") {
		output.PrintError("Access denied: missing scope infra:write. See platformctl auth status.")
		return 1
	}
	plan, err := env.BootstrapPlan(envName)
	if err != nil {
		output.PrintError(err.Error())
		return 1
	}
	fmt.Print(env.FormatPlan(plan))
	return 0
}
