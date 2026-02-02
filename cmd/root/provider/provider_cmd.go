package provider

import "github.com/spf13/cobra"

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "provider",
		Short: "Provider integration utilities",
	}
	cmd.AddCommand(checkCmd())
	return cmd
}
