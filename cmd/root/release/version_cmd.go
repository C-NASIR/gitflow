package release

import (
	"fmt"
	"os"

	"gitflow/internal/cli"
	"gitflow/internal/ui"
	"gitflow/internal/workflow"

	"github.com/spf13/cobra"
)

func versionCmd() *cobra.Command {
	var jsonOutput bool
	var envOutput bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the next release version",
		RunE: func(cmd *cobra.Command, args []string) error {
			format, err := parseOutputFormat(jsonOutput, envOutput)
			if err != nil {
				return cli.ExitError{Err: err, Code: exitCodeConfig}
			}

			repoPath, err := os.Getwd()
			if err != nil {
				return cli.ExitError{Err: fmt.Errorf("failed to get current directory: %w", err), Code: exitCodeComputation}
			}

			out, err := workflow.Release(workflow.ReleaseOptions{
				RepoPath: repoPath,
				DryRun:   true,
			})
			if err != nil {
				return releaseExitError(err)
			}

			plainUI := ui.New(ui.Options{
				Out:     cmd.OutOrStdout(),
				Color:   false,
				Emoji:   false,
				Verbose: false,
			})
			if err := outputReleaseVersion(plainUI, format, out.NextVersion.String()); err != nil {
				return cli.ExitError{Err: err, Code: exitCodeComputation}
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output machine readable JSON")
	cmd.Flags().BoolVar(&envOutput, "env", false, "Output KEY=VALUE lines")
	return cmd
}
