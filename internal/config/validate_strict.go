package config

import (
	"errors"
	"strings"
)

// ValidateStrict validates config without applying defaults.
func ValidateStrict(cfg *Config) error {
	var errs []string

	if cfg.Branches.MainBranch == "" {
		errs = append(errs, "branches.main_branch is required")
	}
	if cfg.Branches.FeaturePrefix == "" {
		errs = append(errs, "branches.feature_prefix is required")
	}
	if cfg.Workflows.Start.BaseBranch == "" {
		errs = append(errs, "workflows.start.base_branch is required")
	}
	if cfg.Workflows.Sync.Strategy != "" && cfg.Workflows.Sync.Strategy != "rebase" && cfg.Workflows.Sync.Strategy != "merge" {
		errs = append(errs, "workflows.sync.strategy must be rebase or merge")
	}
	if cfg.Workflows.Cleanup.AgeThresholdDays < 0 {
		errs = append(errs, "workflows.cleanup.age_threshold_days must be >= 0")
	}

	if cfg.Provider.Type != "" {
		if cfg.Provider.Type != "github" && cfg.Provider.Type != "gitlab" {
			errs = append(errs, "provider.type must be github or gitlab")
		}
		if cfg.Provider.TokenEnv == "" {
			errs = append(errs, "provider.token_env is required when provider is enabled")
		}
		if cfg.Provider.Owner == "" {
			errs = append(errs, "provider.owner is required when provider is enabled")
		}
		if cfg.Provider.Repo == "" {
			errs = append(errs, "provider.repo is required when provider is enabled")
		}
	}

	if cfg.Commits.Conventional {
		if len(cfg.Commits.Types) == 0 {
			errs = append(errs, "commits.types must be set when commits.conventional is true")
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}
