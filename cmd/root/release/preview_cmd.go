package release

import (
	"fmt"
	"os"
	"strings"

	"gitflow/internal/cli"
	"gitflow/internal/ui"
	"gitflow/internal/workflow"

	"github.com/spf13/cobra"
)

func previewCmd() *cobra.Command {
	var versionOverride string

	cmd := &cobra.Command{
		Use:   "preview",
		Short: "Preview the next release",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := cli.CommonFromCmd(cmd)
			if err != nil {
				return err
			}

			repoPath, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			opts := workflow.ReleaseOptions{
				RepoPath: repoPath,
				DryRun:   true,
			}
			if versionOverride != "" {
				version, ok := parseVersion(versionOverride)
				if !ok {
					return fmt.Errorf("invalid version override: %s", versionOverride)
				}
				opts.VersionOverride = &version
			}

			out, err := workflow.Release(opts)
			if err != nil {
				return err
			}

			c.UI.Header("Release preview")
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

			return nil
		},
	}

	cmd.Flags().StringVar(&versionOverride, "version", "", "Override computed version")
	return cmd
}
