package cli

import (
	"io"

	"github.com/spf13/cobra"

	"gitflow/internal/config"
	"gitflow/internal/ui"
)

type Common struct {
	ConfigResult *config.LoadResult
	UI           *ui.UI
}

func CommonFromCmd(cmd *cobra.Command) (*Common, error) {
	res, err := config.Load()
	if err != nil {
		return nil, err
	}

	out := cmd.OutOrStdout()
	u := ui.New(ui.Options{
		Out:     out,
		Color:   res.Config.UI.Color,
		Emoji:   res.Config.UI.Emoji,
		Verbose: res.Config.UI.Verbose,
	})

	return &Common{
		ConfigResult: res,
		UI:           u,
	}, nil
}

func outWriter(cmd *cobra.Command) io.Writer {
	return cmd.OutOrStdout()
}
