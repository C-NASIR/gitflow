// Package provider exposes integrations for hosted git providers.
package provider

import (
	"context"
	"fmt"
)

// Provider defines the hosting provider behaviors needed by the app.
type Provider interface {
	ValidateAuth(ctx context.Context) error
	GetDefaultBranch(ctx context.Context) (string, error)
}
type ProviderConfig struct {
	Type    string
	BaseURL string

	Token string
	Owner string
	Repo  string
}

// New constructs a provider implementation from config.
func New(cfg ProviderConfig) (Provider, error) {
	switch cfg.Type {
	case "github":
		return NewGitHub(cfg)
	case "gitlab":
		return nil, fmt.Errorf("gitlab provider not implemented yet")
	case "":
		return nil, fmt.Errorf("provider type is empty")
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", cfg.Type)
	}
}
