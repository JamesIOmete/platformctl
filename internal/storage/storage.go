package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/JamesIOmete/platformctl/internal/fleet"
	"github.com/JamesIOmete/platformctl/internal/secrets"
	"github.com/JamesIOmete/platformctl/internal/sim"
)

// State represents the persisted application state.
type State struct {
	Devices     []fleet.Device            `json:"devices"`
	Simulations []sim.Job                 `json:"simulations"`
	Secrets     map[string]secrets.Secret `json:"secrets"`
}

// DefaultState returns the initial MVP data.
func DefaultState() *State {
	return &State{
		Devices: []fleet.Device{
			{ID: "robot-001", Model: "digit", State: "online", IP: "10.0.0.45", Battery: 88, Firmware: "v2.1.0", LastSeen: time.Now().Add(-5 * time.Minute).Format(time.RFC3339)},
			{ID: "robot-002", Model: "digit", State: "maintenance", IP: "10.0.0.46", Battery: 12, Firmware: "v2.0.0", LastSeen: time.Now().Add(-48 * time.Hour).Format(time.RFC3339)},
			{ID: "edge-101", Model: "edge-node", State: "offline", IP: "192.168.1.10", Battery: 0, Firmware: "v1.5.4", LastSeen: time.Now().Add(-720 * time.Hour).Format(time.RFC3339)},
		},
		Simulations: []sim.Job{},
		Secrets:     make(map[string]secrets.Secret),
	}
}

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".config", "platformctl")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "mock-state.json"), nil
}

// Load reads the state from disk or returns default if missing.
func Load() (*State, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}
	
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		// First run: save and return default
		s := DefaultState()
		if err := Save(s); err != nil {
			return nil, fmt.Errorf("initializing state: %w", err)
		}
		return s, nil
	}
	if err != nil {
		return nil, err
	}

	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parsing mock-state.json: %w", err)
	}
	return &s, nil
}

// Save writes the state to disk.
func Save(s *State) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}
	
	bytes, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(path, bytes, 0644)
}
