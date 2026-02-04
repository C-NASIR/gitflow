package cli

import (
	"fmt"
	"gitflow/internal/ui"
)

// PrintConfigSource renders the config source summary line.
func PrintConfigSource(u *ui.UI, path string) {
	u.Line("Config source: %s", ConfigSource(path))
}

// ConfigSource returns a printable config source string.
func ConfigSource(path string) string {
	if path == "" {
		return "defaults"
	}
	return path
}

// EnsureOneKind ensures only one of bugfix or hotfix is selected.
func EnsureOneKind(bugfix bool, hotfix bool) error {
	if bugfix && hotfix {
		return fmt.Errorf("choose only one of bugfix or hotfix")
	}
	return nil
}
