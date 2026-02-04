package branch

import (
	"fmt"
	"os"

	"gitflow/internal/cli"
	"gitflow/internal/config"
	"gitflow/internal/ui"
	"gitflow/internal/workflow"

	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	var base string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List local branches with age and ahead behind counts",
		RunE: func(cmd *cobra.Command, args []string) error {
			repoPath, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			res, err := config.Load()
			if err != nil {
				return err
			}

			out, err := workflow.ListBranches(res.Config, workflow.BranchListOptions{
				RepoPath: repoPath,
				Base:     base,
			})
			if err != nil {
				return err
			}

			if len(out.Branches) == 0 {
				c, _ := cli.CommonFromCmd(cmd)
				c.UI.Success("No branches found")
				return nil
			}

			c, err := cli.CommonFromCmd(cmd)
			if err != nil {
				return err
			}

			c.UI.Header("Branches")
			cli.PrintConfigSource(c.UI, c.ConfigResult.Path)
			c.UI.Line("Base: %s", out.Base)
			c.UI.Line("")

			t := ui.NewTable(cmd.OutOrStdout())
			t.Header("CUR", "NAME", "AGE", "AHEAD", "BEHIND", "AUTHOR", "LAST")
			for _, b := range out.Branches {
				cur := ""
				if b.Current {
					cur = "*"
				}
				t.Row(cur, b.Name, b.AgeDays, b.Ahead, b.Behind, b.Author, b.LastCommitMsg)
			}
			t.Flush()

			return nil
		},
	}

	cmd.Flags().StringVar(&base, "base", "", "Base branch to compare against")
	return cmd
}
