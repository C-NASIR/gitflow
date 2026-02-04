package pr

import (
	"github.com/spf13/cobra"

	"gitflow/internal/cli"
	"gitflow/internal/config"
	"gitflow/internal/ui"
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

			c, err := cli.CommonFromCmd(cmd)
			if err != nil {
				return err
			}

			c.UI.Header("Pull requests")
			cli.PrintConfigSource(c.UI, c.ConfigResult.Path)

			if len(out.PRs) == 0 {
				cmd.Println("No pull requests found")
				return nil
			}

			t := ui.NewTable(cmd.OutOrStdout())
			t.Header("NUM", "STATE", "AUTHOR", "HEAD", "BASE", "TITLE")
			for _, pr := range out.PRs {
				num := pr.Number
				state := pr.State
				if pr.Draft {
					state = state + " draft"
				}
				t.Row(num, state, pr.Author, pr.HeadBranch, pr.BaseBranch, pr.Title)
			}
			t.Flush()

			return nil
		},
	}

	cmd.Flags().StringVar(&state, "state", "open", "State: open, closed, all")
	return cmd
}
