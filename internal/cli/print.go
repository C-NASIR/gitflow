package cli

import (
	"fmt"
	"gitflow/internal/ui"
)

func PrintConfigSource(u *ui.UI, path string) {
	if path == "" {
		u.Line("Config source: defaults")
		return
	}
	u.Line("Config source: %s", path)
}

func EnsureOneKind(bugfix bool, hotfix bool) error {
	if bugfix && hotfix {
		return fmt.Errorf("choose only one of bugfix or hotfix")
	}
	return nil
}
