package release

import (
	"fmt"
	"os"
	"strings"

	"gitflow/internal/cli"
	"gitflow/internal/git"
	"gitflow/internal/ui"
	"gitflow/internal/workflow"

	"github.com/spf13/cobra"
)

func createCmd() *cobra.Command {
	var (
		dryRun          bool
		versionOverride string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a release tag",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := cli.CommonFromCmd(cmd)
			if err != nil {
				return cli.ExitError{Err: err, Code: exitCodeConfig}
			}

			repoPath, err := os.Getwd()
			if err != nil {
				return cli.ExitError{Err: fmt.Errorf("failed to get current directory: %w", err), Code: exitCodeComputation}
			}

			opts := workflow.ReleaseOptions{
				RepoPath: repoPath,
				DryRun:   dryRun,
			}
			if versionOverride != "" {
				version, ok := parseVersion(versionOverride)
				if !ok {
					return cli.ExitError{Err: fmt.Errorf("invalid version override: %s", versionOverride), Code: exitCodeConfig}
				}
				opts.VersionOverride = &version
			}

			out, err := workflow.Release(opts)
			if err != nil {
				return releaseExitError(err)
			}

			c.UI.Header("Release create")
			t := ui.NewTable(cmd.OutOrStdout())
			t.Header("KEY", "VALUE")
			t.KeyValue("Current version", out.BaseVersion.String())
			t.KeyValue("Next version", out.NextVersion.String())
			t.KeyValue("Commit count", out.CommitCount)
			t.Flush()

			c.UI.Line("")
			for _, line := range strings.Split(out.Changelog, "\n") {
				c.UI.Line("%s", line)
			}

			if dryRun {
				c.UI.Warn("Dry run: no tag created")
				return nil
			}

			client, err := git.NewClient(repoPath)
			if err != nil {
				return cli.ExitError{Err: err, Code: exitCodeComputation}
			}
			if err := client.CreateAnnotatedTag(out.Tag, out.Changelog); err != nil {
				return cli.ExitError{Err: err, Code: exitCodeComputation}
			}
			c.UI.Success("Created tag %s", out.Tag)
			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Compute release without creating a tag")
	cmd.Flags().StringVar(&versionOverride, "version", "", "Override computed version")
	return cmd
}
