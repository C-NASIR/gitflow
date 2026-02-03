package root

import (
	"github.com/spf13/cobra"

	"gitflow/internal/config"
)

func configValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate configuration and show actionable errors",
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := config.Load()
			if err != nil {
				return err
			}

			if err := config.ValidateStrict(res.Config); err != nil {
				cmd.Println("Config invalid")
				cmd.Println(err.Error())
				return nil
			}

			cmd.Println("Config valid")
			if res.Path == "" {
				cmd.Println("Config source: defaults")
			} else {
				cmd.Printf("Config source: %s\n", res.Path)
			}
			return nil
		},
	}
}
