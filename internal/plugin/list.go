package plugin

import (
	"os"
	"sort"
	"strings"
)

// ListPluginsOnPath returns plugin names discovered via PATH entries.
func ListPluginsOnPath() []string {
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return nil
	}
	seen := make(map[string]struct{})
	var names []string
	for _, dir := range strings.Split(pathEnv, string(os.PathListSeparator)) {
		if dir == "" {
			continue
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			name := e.Name()
			if !strings.HasPrefix(name, "platformctl-") {
				continue
			}
			info, err := e.Info()
			if err != nil {
				continue
			}
			if info.Mode()&0111 == 0 { // skip non-executable files
				continue
			}
			trimmed := strings.TrimPrefix(name, "platformctl-")
			if _, ok := seen[trimmed]; ok {
				continue
			}
			seen[trimmed] = struct{}{}
			names = append(names, trimmed)
		}
	}
	sort.Strings(names)
	return names
}
