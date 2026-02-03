package pr

import (
	"github.com/spf13/cobra"

	"gitflow/internal/config"
	"gitflow/internal/workflow"
)

func listCmd() *cobra.Command {
	var state string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List pull requests",
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := config.Load()
			if err != nil {
				return err
			}

			out, err := workflow.ListPRs(res.Config, workflow.PRListOptions{
				State: state,
			})
			if err != nil {
				return err
			}

			if len(out.PRs) == 0 {
				cmd.Println("No pull requests found")
				return nil
			}

			for _, pr := range out.PRs {
				draftText := ""
				if pr.Draft {
					draftText = " draft"
				}
				cmd.Printf("#%d %s %s %s%s\n", pr.Number, pr.State, pr.Author, pr.Title, draftText)
				cmd.Printf("  head %s  base %s\n", pr.HeadBranch, pr.BaseBranch)
				cmd.Printf("  %s\n", pr.URL)
				cmd.Println()
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&state, "state", "open", "State: open, closed, all")
	return cmd
}
