package root

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"gitflow/internal/config"
	"gitflow/internal/workflow"
)

func initCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create a starter .gitflow.yml in the current repo",
		RunE: func(cmd *cobra.Command, args []string) error {
			repoPath, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			cfg := config.Default()

			out, err := workflow.Init(cfg, workflow.InitOptions{
				RepoPath: repoPath,
				Force:    force,
			})
			if err != nil {
				return err
			}

			cmd.Printf("Wrote %s\n", out.Path)
			cmd.Println("Edit provider settings to enable PR features")
			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")
	return cmd
}
