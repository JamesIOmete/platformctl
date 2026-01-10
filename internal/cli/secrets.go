package cli

import (
	"fmt"
	"time"

	"github.com/JamesIOmete/platformctl/internal/auth"
	"github.com/JamesIOmete/platformctl/internal/output"
	"github.com/JamesIOmete/platformctl/internal/secrets"
	"github.com/JamesIOmete/platformctl/internal/storage"
)

func handleSecrets(args []string) int {
	if len(args) == 0 || args[0] == "help" {
		fmt.Println("Usage: platformctl secrets [ls | get <key> | set <key> <value>]")
		return 0
	}

	status := auth.LoadStatus()
	// In real life, secrets require high privilege.
	if !status.Authenticated {
		output.PrintError("Access denied: You must be authenticated to manage secrets.")
		return 1
	}

	switch args[0] {
	case "ls":
		return runSecretsList()
	case "get":
		if len(args) < 2 {
			output.PrintError("Missing key. Usage: platformctl secrets get <key>")
			return 1
		}
		return runSecretsGet(args[1])
	case "set":
		if len(args) < 3 {
			output.PrintError("Missing arguments. Usage: platformctl secrets set <key> <value>")
			return 1
		}
		return runSecretsSet(args[1], args[2], status.Principal)
	default:
		output.PrintError("Unknown secrets subcommand. Try: platformctl secrets help")
		return 1
	}
}

func runSecretsList() int {
	s, err := storage.Load()
	if err != nil {
		output.PrintError(fmt.Sprintf("Failed to load state: %v", err))
		return 1
	}

	if len(s.Secrets) == 0 {
		fmt.Println("No secrets found.")
		return 0
	}

	fmt.Println("KEY                 \tCREATED BY          \tCREATED AT")
	fmt.Println("--------------------\t--------------------\t--------------------")
	for _, sec := range s.Secrets {
		fmt.Printf("%-20s\t%-20s\t%s\n", sec.Key, sec.CreatedBy, sec.CreatedAt.Format(time.RFC3339))
	}
	return 0
}

func runSecretsGet(key string) int {
	s, err := storage.Load()
	if err != nil {
		output.PrintError(fmt.Sprintf("Failed to load state: %v", err))
		return 1
	}

	sec, ok := s.Secrets[key]
	if !ok {
		output.PrintError(fmt.Sprintf("Secret '%s' not found.", key))
		return 1
	}

	// In a real CLI, we might mask this or require a --show flag.
	// For POC, just print it.
	fmt.Println(sec.Value)
	return 0
}

func runSecretsSet(key, value, user string) int {
	s, err := storage.Load()
	if err != nil {
		output.PrintError(fmt.Sprintf("Failed to load state: %v", err))
		return 1
	}

	// Initialize map if nil (safety check, though DefaultState handles it)
	if s.Secrets == nil {
		s.Secrets = make(map[string]secrets.Secret)
	}

	s.Secrets[key] = secrets.Secret{
		Key:       key,
		Value:     value,
		CreatedBy: user,
		CreatedAt: time.Now(),
	}

	if err := storage.Save(s); err != nil {
		output.PrintError(fmt.Sprintf("Failed to save secret: %v", err))
		return 1
	}

	fmt.Printf("Secret '%s' set successfully.\n", key)
	return 0
}
