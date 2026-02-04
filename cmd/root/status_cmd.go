package root

import (
	"fmt"
	"gitflow/internal/cli"
	"gitflow/internal/workflow"
	"os"

	"github.com/spf13/cobra"
)

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show repository status summary",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := cli.CommonFromCmd(cmd)
			if err != nil {
				return err
			}

			repoPath, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			s, err := workflow.GetStatus(repoPath)
			if err != nil {
				return err
			}

			c.UI.Header("Repository status")
			cli.PrintConfigSource(c.UI, c.ConfigResult.Path)

			c.UI.Line("Branch: %s", s.Branch)
			if s.Dirty {
				c.UI.Warn("Working tree: dirty")
			} else {
				c.UI.Success("Working tree: clean")
			}

			return nil
		},
	}
}
