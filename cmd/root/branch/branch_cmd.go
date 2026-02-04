// Package branch defines branch-related CLI commands.
package branch

import "github.com/spf13/cobra"

// Cmd builds the branch command tree.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "branch",
		Short: "Branch utilities",
	}
	cmd.AddCommand(listCmd())
	return cmd
}
