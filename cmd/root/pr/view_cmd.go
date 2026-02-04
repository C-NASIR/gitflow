package pr

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"gitflow/internal/cli"
	"gitflow/internal/config"
	"gitflow/internal/workflow"
)

func viewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view <number>",
		Short: "View a pull request",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("number is required")
			}
			_, err := strconv.Atoi(args[0])
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			n, _ := strconv.Atoi(args[0])

			res, err := config.Load()
			if err != nil {
				return err
			}

			out, err := workflow.ViewPR(res.Config, n)
			if err != nil {
				return err
			}

			c, err := cli.CommonFromCmd(cmd)
			if err != nil {
				return err
			}

			c.UI.Header("Pull request")
			cli.PrintConfigSource(c.UI, c.ConfigResult.Path)

			pr := out.PR
			c.UI.Line("Number: %d", pr.Number)
			c.UI.Line("Title: %s", pr.Title)
			c.UI.Line("State: %s", pr.State)
			c.UI.Line("Author: %s", pr.Author)
			c.UI.Line("Head: %s", pr.HeadBranch)
			c.UI.Line("Base: %s", pr.BaseBranch)
			c.UI.Line("Draft: %v", pr.Draft)
			c.UI.Line("URL: %s", pr.URL)

			if pr.Description != "" {
				c.UI.Line("")
				c.UI.Line(pr.Description)
			}

			return nil
		},
	}

	return cmd
}
