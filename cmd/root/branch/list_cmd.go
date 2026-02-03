package branch

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"gitflow/internal/config"
	"gitflow/internal/workflow"
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
				cmd.Println("No branches found")
				return nil
			}

			cmd.Printf("Base: %s\n\n", out.Base)

			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "CURRENT\tNAME\tAGE_DAYS\tAHEAD\tBEHIND\tAUTHOR\tLAST_COMMIT")
			for _, b := range out.Branches {
				cur := ""
				if b.Current {
					cur = "*"
				}
				fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%d\t%s\t%s\n",
					cur,
					b.Name,
					b.AgeDays,
					b.Ahead,
					b.Behind,
					b.Author,
					b.LastCommitMsg,
				)
			}
			_ = w.Flush()

			return nil
		},
	}

	cmd.Flags().StringVar(&base, "base", "", "Base branch to compare against")
	return cmd
}
