// Package provider exposes integrations for hosted git providers.
package provider

import (
	"context"
	"fmt"
	"gitflow/pkg/types"
)

// Provider defines the hosting provider behaviors needed by the app.
type Provider interface {
	ValidateAuth(ctx context.Context) error
	GetDefaultBranch(ctx context.Context) (string, error)
	CreatePR(ctx context.Context, opts CreatePROptions) (*types.PullRequest, error)
	GetPR(ctx context.Context, number int) (*types.PullRequest, error)
	ListPRs(ctx context.Context, state string) ([]*types.PullRequest, error)
	CreateRelease(tag string, name string, body string) (*types.Release, error)
	UpdateRelease(tag string, name string, body string) (*types.Release, error)
}

// CreatePROptions defines pull request creation inputs.
type CreatePROptions struct {
	Title       string
	Description string
	HeadBranch  string
	BaseBranch  string
	Draft       bool
	Reviewers   []string
	Labels      []string
}

// ProviderConfig contains provider connection settings.
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
		return NewGitLab(cfg)
	case "":
		return nil, fmt.Errorf("provider type is empty")
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", cfg.Type)
	}
}
