package root

import (
	"fmt"
	"gitflow/internal/config"
	"gitflow/internal/workflow"
	"os"

	"github.com/spf13/cobra"
)

func syncCmd() *cobra.Command {
	var remote string
	var merge bool
	var rebase bool
	var noPush bool
	var force bool

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync current branch with base branch using rebase or merge",
		RunE: func(cmd *cobra.Command, args []string) error {
			if merge && rebase {
				return fmt.Errorf("choose only one of merge or rebase")
			}

			repoPath, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			res, err := config.Load()
			if err != nil {
				return err
			}

			var strategy string
			if merge {
				strategy = "merge"
			}
			if rebase {
				strategy = "rebase"
			}

			var autoPushOverride *bool
			if noPush {
				v := false
				autoPushOverride = &v
			}

			var forcePushOverride *bool
			if force {
				v := true
				forcePushOverride = &v
			}

			out, err := workflow.Sync(res.Config, workflow.SyncOptions{
				RepoPath:          repoPath,
				Remote:            remote,
				StrategyOverride:  strategy,
				AutoPushOverride:  autoPushOverride,
				ForcePushOverride: forcePushOverride,
			})

			if err != nil {
				return err
			}

			cmd.Printf("Base Branch: %s\n", out.BaseBranch)
			cmd.Printf("Current branch: %s\n", out.CurrentBranch)
			cmd.Printf("Strategy: %s\n", out.Strategy)

			if out.Pushed {
				if out.ForcePushed {
					cmd.Printf("Remote: pushed with force lease\n")
				} else {
					cmd.Printf("Remote: pushed\n")
				}
			} else {
				cmd.Printf("Remote: not pushed\n")
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&remote, "remote", "origin", "Remote name")
	cmd.Flags().BoolVar(&merge, "merge", false, "Use merge strategy")
	cmd.Flags().BoolVar(&rebase, "rebase", false, "Use rebase strategy")
	cmd.Flags().BoolVar(&noPush, "no-push", false, "Do not push after syncing")
	cmd.Flags().BoolVar(&force, "force", false, "Allow force with lease when rebasing and pushing")

	return cmd
}
