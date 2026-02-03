package root

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"gitflow/internal/workflow"
)

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

			out, err := workflow.Cleanup(c.ConfigResult.Config, workflow.CleanupOptions{
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

			if len(out.Candidates) == 0 {
				c.UI.Success("Nothing to cleanup")
				return nil
			}

			if !yes && len(out.Deleted) == 0 {
				choices := make([]huh.Option[string], 0, len(out.Candidates))
				for _, b := range out.Candidates {
					label := fmt.Sprintf("%s  reason %s  age %d  ahead %d  behind %d", b.Name, b.Reason, b.AgeDays, b.Ahead, b.Behind)
					choices = append(choices, huh.NewOption(label, b.Name))
				}

				var selected []string
				form := huh.NewForm(
					huh.NewGroup(
						huh.NewMultiSelect[string]().
							Title("Select branches to delete").
							Options(choices...).
							Value(&selected),
					),
				)

				if err := form.Run(); err != nil {
					return err
				}

				if len(selected) == 0 {
					c.UI.Warn("No branches selected")
					return nil
				}

				out2, err := workflow.Cleanup(c.ConfigResult.Config, workflow.CleanupOptions{
					RepoPath:     repoPath,
					Remote:       remoteName,
					Yes:          true,
					All:          all,
					AgeThreshold: age,
					DeleteRemote: remote,
					Selected:     selected,
				})
				if err != nil {
					return err
				}

				c.UI.Line("Deleted: %d", len(out2.Deleted))
				for _, b := range out2.Deleted {
					c.UI.Line("  %s", b)
				}

				if remote && len(out2.RemoteDeleted) > 0 {
					c.UI.Line("Remote deleted: %d", len(out2.RemoteDeleted))
					for _, b := range out2.RemoteDeleted {
						c.UI.Line("  %s", b)
					}
				}

				return nil

			}

			if len(out.Deleted) == 0 {
				c.UI.Success("Deleted: none")
				return nil
			}

			c.UI.Line("Deleted: %d", len(out.Deleted))
			for _, b := range out.Deleted {
				c.UI.Line("  %s", b)
			}

			if remote && len(out.RemoteDeleted) > 0 {
				c.UI.Line("Remote deleted: %d", len(out.RemoteDeleted))
				for _, b := range out.RemoteDeleted {
					c.UI.Line("  %s", b)
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&yes, "yes", false, "Skip confirmation and delete all candidates")
	cmd.Flags().BoolVar(&all, "all", false, "Include stale unmerged branches by age")
	cmd.Flags().IntVar(&age, "age", 0, "Age threshold in days")
	cmd.Flags().BoolVar(&remote, "remote", false, "Also delete remote branches")
	cmd.Flags().StringVar(&remoteName, "remote-name", "origin", "Remote name")

	return cmd
}

func contains(list []string, v string) bool {
	for _, x := range list {
		if strings.TrimSpace(x) == v {
			return true
		}
	}
	return false
}
