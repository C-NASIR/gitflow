package root

import (
	"fmt"
	"gitflow/internal/config"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func configShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Print the resolved configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := commonFromCmd(cmd)
			if err != nil {
				return err
			}

			res, err := config.Load()
			if err != nil {
				return err
			}

			printConfigSource(c.UI, c.ConfigResult.Path)
			out, err := yaml.Marshal(res.Config)
			if err != nil {
				return fmt.Errorf("failed to render yaml: %w", err)
			}

			cmd.Println(string(out))
			return nil
		},
	}
}
