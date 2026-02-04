package root

import (
	"fmt"
	"os"

	"gitflow/internal/cli"
	"gitflow/internal/ui"
	"gitflow/internal/workflow"

	"github.com/spf13/cobra"
)

func doctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check repository health",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := cli.CommonFromCmd(cmd)
			if err != nil {
				return err
			}

			repoPath, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			out, err := workflow.Doctor(repoPath)
			if err != nil {
				return err
			}

			c.UI.Header("Doctor report")

			t := ui.NewTable(cmd.OutOrStdout())
			t.Header("STATUS", "CHECK", "MESSAGE")

			hasErrors := false
			for _, check := range out.Checks {
				if check.Level == workflow.DoctorError {
					hasErrors = true
				}
				status := c.UI.StatusLabel(check.Level)
				t.Row(status, check.Name, check.Message)
			}
			t.Flush()

			if hasErrors {
				return fmt.Errorf("doctor found errors")
			}

			return nil
		},
	}
}
