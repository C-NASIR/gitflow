package provider

import (
	"fmt"
	"gitflow/internal/config"
	"os"
	"strings"
)

func Enabled(cfg *config.Config) bool {
	return strings.TrimSpace(cfg.Provider.Type) != ""
}

// FromAppConfig extracts provider config and validates required fields.
func FromAppConfig(cfg *config.Config) (ProviderConfig, error) {
	if strings.TrimSpace(cfg.Provider.Type) == "" {
		return ProviderConfig{}, fmt.Errorf("provider is not configured")
	}

	token := ""
	if cfg.Provider.TokenEnv != "" {
		token = os.Getenv(cfg.Provider.TokenEnv)
	}

	if strings.TrimSpace(token) == "" {
		return ProviderConfig{}, fmt.Errorf("provider token missing, set env var %s", cfg.Provider.TokenEnv)
	}
	if strings.TrimSpace(cfg.Provider.Owner) == "" {
		return ProviderConfig{}, fmt.Errorf("provider owner is required")
	}
	if strings.TrimSpace(cfg.Provider.Repo) == "" {
		return ProviderConfig{}, fmt.Errorf("provider repo is required")
	}

	return ProviderConfig{
		Type:    cfg.Provider.Type,
		BaseURL: cfg.Provider.BaseURL,
		Token:   token,
		Owner:   cfg.Provider.Owner,
		Repo:    cfg.Provider.Repo,
	}, nil
}
