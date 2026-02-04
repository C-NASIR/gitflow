package release

import (
	"fmt"
	"os"

	"gitflow/internal/cli"
	"gitflow/internal/workflow"

	"github.com/spf13/cobra"
)

func publishCmd() *cobra.Command {
	var (
		dryRun     bool
		jsonOutput bool
		envOutput  bool
	)

	cmd := &cobra.Command{
		Use:   "publish",
		Short: "Publish release notes to a provider",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := cli.CommonFromCmd(cmd)
			if err != nil {
				return cli.ExitError{Err: err, Code: exitCodeConfig}
			}
			format, err := parseOutputFormat(jsonOutput, envOutput)
			if err != nil {
				return cli.ExitError{Err: err, Code: exitCodeConfig}
			}

			repoPath, err := os.Getwd()
			if err != nil {
				return cli.ExitError{Err: fmt.Errorf("failed to get current directory: %w", err), Code: exitCodeComputation}
			}

			releaseResult, err := workflow.Release(workflow.ReleaseOptions{
				RepoPath: repoPath,
				DryRun:   true,
			})
			if err != nil {
				return releaseExitError(err)
			}

			publishResult, err := workflow.ReleasePublish(workflow.ReleasePublishOptions{
				RepoPath: repoPath,
				DryRun:   dryRun,
				Result:   releaseResult,
			})
			if err != nil {
				return releaseExitError(err)
			}

			if err := outputReleasePublish(c.UI, cmd.OutOrStdout(), format, releaseResult, publishResult); err != nil {
				return cli.ExitError{Err: err, Code: exitCodeComputation}
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Skip publishing release notes")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output machine readable JSON")
	cmd.Flags().BoolVar(&envOutput, "env", false, "Output KEY=VALUE lines")
	return cmd
}
