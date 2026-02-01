package root

import (
	"fmt"
	"gitflow/internal/config"
	"gitflow/internal/workflow"
	"os"

	"github.com/spf13/cobra"
)

// cleanupCmd builds the cleanup subcommand.
func cleanupCmd() *cobra.Command {
	var yes bool
	var all bool
	var age int
	var remote bool
	var remoteName string

	cmd := &cobra.Command{
		Use:   "cleanup",
		Short: "Delete merged or stale branches safely",
		RunE: func(cmd *cobra.Command, args []string) error {
			repoPath, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			res, err := config.Load()
			if err != nil {
				return err
			}

			out, err := workflow.Cleanup(res.Config, workflow.CleanupOptions{
				RepoPath:     repoPath,
				Remote:       remoteName,
				Yes:          yes,
				All:          all,
				AgeThreshold: age,
				DeleteRemote: remote,
			})
			if err != nil {
				return err
			}

			cmd.Printf("Base branch: %s\n", out.BaseBranch)
			cmd.Printf("Current branch: %s\n", out.Current)

			if len(out.Deleted) == 0 {
				cmd.Println("Deleted: none")
				return nil
			}

			cmd.Printf("Deleted: %d\n", len(out.Deleted))
			for _, b := range out.Deleted {
				cmd.Printf("  %s\n", b)
			}

			if len(out.RemoteDeleted) > 0 {
				cmd.Printf("Remote deleted: %d\n", len(out.RemoteDeleted))
				for _, b := range out.RemoteDeleted {
					cmd.Printf("  %s\n", b)
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&yes, "yes", false, "Skip confirmation")
	cmd.Flags().BoolVar(&all, "all", false, "Include stale unmerged branches by age")
	cmd.Flags().IntVar(&age, "age", 0, "Age threshold in days")
	cmd.Flags().BoolVar(&remote, "remote", false, "Also delete remote branches")
	cmd.Flags().StringVar(&remoteName, "remote-name", "origin", "Remote name")

	return cmd
}
