// Package release defines release CLI commands.
package release

import "github.com/spf13/cobra"

// Cmd builds the release command tree.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release",
		Short: "Release automation",
	}
	cmd.AddCommand(previewCmd())
	cmd.AddCommand(createCmd())
	cmd.AddCommand(versionCmd())
	cmd.AddCommand(publishCmd())
	cmd.AddCommand(changelogCmd())
	return cmd
}
