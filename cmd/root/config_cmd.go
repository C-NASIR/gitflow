package root

import "github.com/spf13/cobra"

func configCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Work with gitflow configuration",
	}
	cmd.AddCommand(configShowCmd())
	cmd.AddCommand(configValidateCmd())
	return cmd
}
