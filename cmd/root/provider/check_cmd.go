package provider

import (
	"context"
	"time"

	"github.com/spf13/cobra"

	"gitflow/internal/config"
	"gitflow/internal/provider"
)

func checkCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "Validate provider credentials and repo access",
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := config.Load()
			if err != nil {
				return err
			}

			if !provider.Enabled(res.Config) {
				cmd.Println("Provider not configured")
				cmd.Println("Add provider fields to .gitflow.yml to enable")
				return nil
			}

			pcfg, err := provider.FromAppConfig(res.Config)
			if err != nil {
				return err
			}

			p, err := provider.New(pcfg)
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			if err := p.ValidateAuth(ctx); err != nil {
				return err
			}

			branch, err := p.GetDefaultBranch(ctx)
			if err != nil {
				return err
			}

			cmd.Println("Provider auth ok")
			cmd.Printf("Default branch: %s\n", branch)
			return nil
		},
	}
}
