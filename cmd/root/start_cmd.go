package root

import (
	"fmt"
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
			kind := "feature"
			if bugfix {
				kind = "bugfix"
			}
			if hotfix {
				kind = "hotfix"
			}
			if bugfix && hotfix {
				return fmt.Errorf("choose only one of bugfix or hotfix")
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

			cmd.Printf("Base branch: %s\n", out.BaseBranch)
			cmd.Printf("New branch: %s\n", out.NewBranch)
			if out.Pushed {
				cmd.Printf("Remote: pushed\n")
			} else {
				cmd.Printf("Remote: not pushed\n")
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&bugfix, "bugfix", false, "Use bugfix prefix")
	cmd.Flags().BoolVar(&hotfix, "hotfix", false, "Use hotfix prefix")
	cmd.Flags().StringVar(&remote, "remote", "origin", "Remote name")
	return cmd
}
