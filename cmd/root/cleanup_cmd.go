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
			c, err := commonFromCmd(cmd)
			if err != nil {
				return err
			}

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

			c.UI.Header("Cleanup branches")
			printConfigSource(c.UI, c.ConfigResult.Path)

			c.UI.Line("Base branch: %s", out.BaseBranch)
			c.UI.Line("Current branch: %s", out.Current)

			if len(out.Deleted) == 0 {
				c.UI.Success("Deleted: none")
				return nil
			}

			c.UI.Line("Deleted: %d", len(out.Deleted))
			for _, b := range out.Deleted {
				c.UI.Line("  %s", b)
			}

			if len(out.RemoteDeleted) > 0 {
				c.UI.Line("Remote deleted: %d", len(out.RemoteDeleted))
				for _, b := range out.RemoteDeleted {
					c.UI.Line("  %s", b)
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
