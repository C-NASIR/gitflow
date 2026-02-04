package release

import (
	"fmt"
	"os"
	"strings"

	"gitflow/internal/cli"
	"gitflow/internal/workflow"

	"github.com/spf13/cobra"
)

func changelogCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "changelog",
		Short: "Generate changelog since last release",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := cli.CommonFromCmd(cmd)
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

			for _, line := range strings.Split(out.Changelog, "\n") {
				c.UI.Line("%s", line)
			}
			return nil
		},
	}

	return cmd
}
