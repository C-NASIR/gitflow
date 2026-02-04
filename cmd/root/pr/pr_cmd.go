// Package pr defines pull request CLI commands.
package pr

import "github.com/spf13/cobra"

// Cmd builds the pull request command tree.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pr",
		Short: "Pull request workflows",
	}
	cmd.AddCommand(createCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(viewCmd())
	return cmd
}
