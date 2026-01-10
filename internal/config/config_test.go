package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultPaths(t *testing.T) {
	paths := DefaultPaths()
	if len(paths) == 0 {
		t.Fatal("Expected default paths, got none")
	}

	home, _ := os.UserHomeDir()
	expectedBase := filepath.Join(home, ".config", "platformctl")
	
	found := false
	for _, p := range paths {
		if filepath.Dir(p) == expectedBase {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected paths to be under %s, got %v", expectedBase, paths)
	}
}

func TestLoadFile_Missing(t *testing.T) {
	cfg, found, err := loadFile("/path/to/non/existent/file.yaml")
	if err != nil {
		t.Fatalf("Expected no error for missing file, got %v", err)
	}
	if found {
		t.Error("Expected found=false for missing file")
	}
	if cfg.Principal != "" {
		t.Error("Expected empty config")
	}
}
