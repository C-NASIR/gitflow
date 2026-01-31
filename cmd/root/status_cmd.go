package root

import (
	"fmt"
	"gitflow/internal/workflow"
	"os"

	"github.com/spf13/cobra"
)

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show repository status summary",
		RunE: func(cmd *cobra.Command, args []string) error {
			repoPath, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			s, err := workflow.GetStatus(repoPath)
			if err != nil {
				return err
			}

			dirtyText := "clean"
			if s.Dirty {
				dirtyText = "dirty"
			}

			cmd.Printf("Branch: %s\n", s.Branch)
			cmd.Printf("Working tree: %s\n", dirtyText)
			return nil
		},
	}
}
