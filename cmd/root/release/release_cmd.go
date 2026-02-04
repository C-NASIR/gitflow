package release

import "github.com/spf13/cobra"

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release",
		Short: "Release automation",
	}
	cmd.AddCommand(previewCmd())
	cmd.AddCommand(createCmd())
	cmd.AddCommand(changelogCmd())
	return cmd
}
