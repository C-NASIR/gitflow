package pr

import "github.com/spf13/cobra"

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pr",
		Short: "Pull request workflows",
	}
	cmd.AddCommand(createCmd())
	return cmd
}
