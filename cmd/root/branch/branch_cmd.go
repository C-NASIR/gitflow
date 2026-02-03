package branch

import "github.com/spf13/cobra"

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "branch",
		Short: "Branch utilities",
	}
	cmd.AddCommand(listCmd())
	return cmd
}
