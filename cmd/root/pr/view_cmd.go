package pr

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

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

			pr := out.PR
			cmd.Printf("PR #%d\n", pr.Number)
			cmd.Printf("Title: %s\n", pr.Title)
			cmd.Printf("State: %s\n", pr.State)
			cmd.Printf("Author: %s\n", pr.Author)
			cmd.Printf("Head: %s\n", pr.HeadBranch)
			cmd.Printf("Base: %s\n", pr.BaseBranch)
			cmd.Printf("Draft: %v\n", pr.Draft)
			cmd.Printf("URL: %s\n", pr.URL)

			if pr.Description != "" {
				cmd.Println()
				cmd.Println(pr.Description)
			}

			return nil
		},
	}

	return cmd
}
