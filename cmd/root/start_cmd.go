package root

import (
	"fmt"
	"gitflow/internal/cli"
	"gitflow/internal/config"
	"gitflow/internal/workflow"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func startCmd() *cobra.Command {
	var bugfix bool
	var hotfix bool
	var remote string

	cmd := &cobra.Command{
		Use:   "start <name>",
		Short: "Start a new branch using conventions and optional push",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("name is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cli.EnsureOneKind(bugfix, hotfix); err != nil {
				return err
			}

			c, err := cli.CommonFromCmd(cmd)
			if err != nil {
				return err
			}

			kind := "feature"
			if bugfix {
				kind = "bugfix"
			}
			if hotfix {
				kind = "hotfix"
			}

			repoPath, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			res, err := config.Load()
			if err != nil {
				return err
			}

			name := strings.Join(args, " ")
			out, err := workflow.Start(res.Config, workflow.StartOptions{
				Kind:     kind,
				RepoPath: repoPath,
				Remote:   remote,
				Name:     name,
			})

			if err != nil {
				return err
			}

			c.UI.Header("Start branch")
			cli.PrintConfigSource(c.UI, c.ConfigResult.Path)

			c.UI.Line("Base branch: %s", out.BaseBranch)
			c.UI.Line("New branch: %s", out.NewBranch)
			if out.Pushed {
				c.UI.Success("Remote: pushed")
			} else {
				c.UI.Warn("Remote: not pushed")
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&bugfix, "bugfix", false, "Use bugfix prefix")
	cmd.Flags().BoolVar(&hotfix, "hotfix", false, "Use hotfix prefix")
	cmd.Flags().StringVar(&remote, "remote", "origin", "Remote name")
	return cmd
}
