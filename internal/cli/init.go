package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/JamesIOmete/platformctl/internal/config"
	"github.com/JamesIOmete/platformctl/internal/storage"
	"gopkg.in/yaml.v3"
)

func handleInit(args []string) int {
	fmt.Println("Welcome to platformctl setup!")
	fmt.Println("This wizard will generate your configuration file.")
	fmt.Println("")

	reader := bufio.NewReader(os.Stdin)

	// Principal
	defaultUser := os.Getenv("USER")
	if defaultUser == "" {
		defaultUser = "dev-user"
	}
	fmt.Printf("Enter your Principal identity [%s]: ", defaultUser)
	principal, _ := reader.ReadString('\n')
	principal = strings.TrimSpace(principal)
	if principal == "" {
		principal = defaultUser
	}

	// Scopes (Simulated selection)
	fmt.Println("Select scopes to assign:")
	fmt.Println("1. fleet:read (Read-only access to robots)")
	fmt.Println("2. infra:write (Manage environments)")
	fmt.Println("3. Full Access (All of the above)")
	fmt.Print("Enter choice [3]: ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var scopes []string
	switch choice {
	case "1":
		scopes = []string{"fleet:read"}
	case "2":
		scopes = []string{"infra:write"}
	default:
		scopes = []string{"fleet:read", "infra:write"}
	}

	// Generate Config
	cfg := config.Config{
		Principal: principal,
		Scopes:    scopes,
	}

	// Save Config
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error finding home dir: %v\n", err)
		return 1
	}
	configDir := filepath.Join(home, ".config", "platformctl")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Printf("Error creating config dir: %v\n", err)
		return 1
	}
	configPath := filepath.Join(configDir, "config.yaml")

	data, err := yaml.Marshal(cfg)
	if err != nil {
		fmt.Printf("Error marshalling config: %v\n", err)
		return 1
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		fmt.Printf("Error writing config file: %v\n", err)
		return 1
	}

	fmt.Printf("\nConfiguration saved to %s\n", configPath)

	// Initialize Storage (Mock State) if missing
	fmt.Println("Initializing local state...")
	_, err = storage.Load()
	if err != nil {
		fmt.Printf("Warning: failed to initialize storage: %v\n", err)
	} else {
		fmt.Println("Local state initialized.")
	}

	fmt.Println("\nSetup complete! Run 'platformctl doctor' to verify.")
	return 0
}
