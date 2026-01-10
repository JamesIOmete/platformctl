package cli

import (
	"fmt"
	"time"

	"github.com/JamesIOmete/platformctl/internal/auth"
	"github.com/JamesIOmete/platformctl/internal/fleet"
	"github.com/JamesIOmete/platformctl/internal/output"
	"github.com/JamesIOmete/platformctl/internal/storage"
)

func handleFleet(args []string) int {
	if len(args) == 0 || args[0] == "help" {
		fmt.Println("Usage: platformctl fleet [ls | status <id> | logs <id> | ssh <id>]")
		return 0
	}

	status := auth.LoadStatus()
	if !auth.HasScope(status, "fleet:read") {
		output.PrintError("Access denied: missing scope fleet:read. See platformctl auth status.")
		return 1
	}

	cmd := args[0]
	switch cmd {
	case "ls":
		return runFleetList()
	case "status":
		if len(args) < 2 {
			output.PrintError("Missing device ID. Usage: platformctl fleet status <id>")
			return 1
		}
		return runFleetStatus(args[1])
	case "logs":
		if len(args) < 2 {
			output.PrintError("Missing device ID. Usage: platformctl fleet logs <id>")
			return 1
		}
		return runFleetLogs(args[1])
	case "ssh":
		if len(args) < 2 {
			output.PrintError("Missing device ID. Usage: platformctl fleet ssh <id>")
			return 1
		}
		return runFleetSSH(args[1])
	default:
		output.PrintError("Unknown fleet subcommand. Try: platformctl fleet help")
		return 1
	}
}

func runFleetList() int {
	s, err := storage.Load()
	if err != nil {
		output.PrintError(fmt.Sprintf("Failed to load state: %v", err))
		return 1
	}
	// Use existing format helper (updated in fleet package to just print table)
	// We might need to map storage devices to fleet devices if types differed, but they share the type.
	fmt.Print(fleet.FormatDevices(s.Devices))
	return 0
}

func runFleetStatus(id string) int {
	s, err := storage.Load()
	if err != nil {
		output.PrintError(fmt.Sprintf("Failed to load state: %v", err))
		return 1
	}
	
	var dev *fleet.Device
	for _, d := range s.Devices {
		if d.ID == id {
			dev = &d
			break
		}
	}
	if dev == nil {
		output.PrintError(fmt.Sprintf("Device %s not found.", id))
		return 1
	}

	fmt.Printf("Device ID:   %s\n", dev.ID)
	fmt.Printf("Model:       %s\n", dev.Model)
	fmt.Printf("State:       %s\n", dev.State)
	fmt.Printf("IP Address:  %s\n", dev.IP)
	fmt.Printf("Battery:     %d%%\n", dev.Battery)
	fmt.Printf("Firmware:    %s\n", dev.Firmware)
	fmt.Printf("Last Seen:   %s\n", dev.LastSeen)
	return 0
}

func runFleetLogs(id string) int {
	// Mock logs
	fmt.Printf("Fetching logs for %s...\n", id)
	time.Sleep(500 * time.Millisecond) // Simulate net lag
	logs := []string{
		"[INFO] System boot complete",
		"[INFO] Connected to fleet-server",
		"[WARN] High latency on motor_controller_1",
		"[INFO] Task assigned: pick_object_42",
		"[INFO] Battery level at 88%",
	}
	for _, l := range logs {
		fmt.Printf("%s %s\n", time.Now().Format("15:04:05.000"), l)
	}
	return 0
}

func runFleetSSH(id string) int {
	fmt.Printf("Requesting tunnel for %s...\n", id)
	time.Sleep(800 * time.Millisecond)
	fmt.Println("Tunnel established. Proxy listening on 127.0.0.1:22022")
	fmt.Println("Connecting...")
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("Welcome to %s (Ubuntu 24.04 LTS)\n", id)
	fmt.Printf("root@%s:~# ", id)
	// Just exit after showing the prompt to simulate "dropping into shell" visually (actual interactive shell is hard)
	// Or we could wait for input but `read_terminal` handles input poorly. 
	// MVP: just show the success message.
	fmt.Println("\n(Mock SSH session ended)")
	return 0
}
