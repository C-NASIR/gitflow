// Package provider defines provider CLI commands.
package provider

import "github.com/spf13/cobra"

// Cmd builds the provider command tree.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "provider",
		Short: "Provider integration utilities",
	}
	cmd.AddCommand(checkCmd())
	return cmd
}
