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
			res, err := config.Load()
			if err != nil {
				return err
			}

			path := res.Path
			if path == "" {
				path = "defaults"
			}

			cmd.Printf("source: %s\n", path)

			out, err := yaml.Marshal(res.Config)
			if err != nil {
				return fmt.Errorf("failed to render yaml: %w", err)
			}

			cmd.Println(string(out))
			return nil
		},
	}
}
