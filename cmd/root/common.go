package root

import (
	"fmt"
	"gitflow/internal/config"
	"gitflow/internal/ui"

	"github.com/spf13/cobra"
)

type Common struct {
	ConfigResult *config.LoadResult
	UI           *ui.UI
}

func commonFromCmd(cmd *cobra.Command) (*Common, error) {
	res, err := config.Load()
	if err != nil {
		return nil, err
	}

	u := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), res.Config.UI.Color, res.Config.UI.Emoji, res.Config.UI.Verbose)
	return &Common{
		ConfigResult: res,
		UI:           u,
	}, nil
}

func printConfigSource(u *ui.UI, path string) {
	src := path
	if src == "" {
		src = "defaults"
	}
	u.Line("Config source: %s", src)
}

func ensureOneKind(bugfix bool, hotfix bool) error {
	if bugfix && hotfix {
		return fmt.Errorf("choose only one of bugfix or hotfix")
	}
	return nil
}
