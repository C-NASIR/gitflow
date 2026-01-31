package root

import (
	"gitflow/internal/version"

	"github.com/spf13/cobra"
)

func VersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print gitflow version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Println(version.String())
			return nil
		},
	}
}
