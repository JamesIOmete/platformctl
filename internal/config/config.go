package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents a minimal auth/config file.
// Example YAML:
// principal: alice
// scopes:
//   - fleet:read
//   - infra:write
//
// Example JSON:
// {"principal":"alice","scopes":["fleet:read"]}
type Config struct {
	Principal string   `json:"principal" yaml:"principal"`
	Scopes    []string `json:"scopes" yaml:"scopes"`
}

// DefaultPaths returns supported config locations in order of preference.
func DefaultPaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	base := filepath.Join(home, ".config", "platformctl")
	return []string{
		filepath.Join(base, "config.yaml"),
		filepath.Join(base, "config.yml"),
		filepath.Join(base, "config.json"),
	}
}

// Load reads the first available config file. Missing files are not an error.
// It returns the config, the path used, and an error if parsing failed.
func Load() (Config, string, error) {
	paths := DefaultPaths()
	for _, path := range paths {
		cfg, used, err := loadFile(path)
		if err != nil {
			return Config{}, path, err
		}
		if used {
			return cfg, path, nil
		}
	}
	return Config{}, "", nil
}

func loadFile(path string) (Config, bool, error) {
	if path == "" {
		return Config{}, false, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Config{}, false, nil
		}
		return Config{}, false, err
	}
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return Config{}, true, nil
	}

	var cfg Config
	// Heuristically decide whether this is JSON; otherwise treat as YAML.
	if trimmed[0] == '{' || trimmed[0] == '[' {
		if err := json.Unmarshal(trimmed, &cfg); err != nil {
			return Config{}, true, err
		}
		return cfg, true, nil
	}

	if err := yaml.Unmarshal(trimmed, &cfg); err != nil {
		return Config{}, true, err
	}
	return cfg, true, nil
}
