package cli

import (
	"fmt"
	"gitflow/internal/ui"
)

func PrintConfigSource(u *ui.UI, path string) {
	u.Line("Config source: %s", ConfigSource(path))
}

func ConfigSource(path string) string {
	if path == "" {
		return "defaults"
	}
	return path
}

func EnsureOneKind(bugfix bool, hotfix bool) error {
	if bugfix && hotfix {
		return fmt.Errorf("choose only one of bugfix or hotfix")
	}
	return nil
}
