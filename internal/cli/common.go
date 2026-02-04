// Package cli provides shared helpers for command wiring.
package cli

import (
	"io"

	"github.com/spf13/cobra"

	"gitflow/internal/config"
	"gitflow/internal/ui"
)

// Common bundles config and UI context for commands.
type Common struct {
	ConfigResult *config.LoadResult
	UI           *ui.UI
}

// UIOverrides applies runtime UI overrides to config.
type UIOverrides struct {
	Color   *bool
	Emoji   *bool
	Verbose *bool
}

var uiOverrides UIOverrides

// SetUIOverrides stores overrides applied during CLI execution.
func SetUIOverrides(overrides UIOverrides) {
	uiOverrides = overrides
}

// CommonFromCmd constructs Common from command output and config.
func CommonFromCmd(cmd *cobra.Command) (*Common, error) {
	res, err := config.Load()
	if err != nil {
		return nil, err
	}

	colorEnabled := res.Config.UI.Color
	emojiEnabled := res.Config.UI.Emoji
	verboseEnabled := res.Config.UI.Verbose
	if uiOverrides.Color != nil {
		colorEnabled = *uiOverrides.Color
	}
	if uiOverrides.Emoji != nil {
		emojiEnabled = *uiOverrides.Emoji
	}
	if uiOverrides.Verbose != nil {
		verboseEnabled = *uiOverrides.Verbose
	}

	out := cmd.OutOrStdout()
	u := ui.New(ui.Options{
		Out:     out,
		Color:   colorEnabled,
		Emoji:   emojiEnabled,
		Verbose: verboseEnabled,
	})

	return &Common{
		ConfigResult: res,
		UI:           u,
	}, nil
}

func outWriter(cmd *cobra.Command) io.Writer {
	return cmd.OutOrStdout()
}
