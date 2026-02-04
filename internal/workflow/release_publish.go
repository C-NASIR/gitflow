package workflow

import (
	"fmt"

	"gitflow/internal/config"
	"gitflow/internal/provider"
)

// ReleasePublishOptions defines inputs for publishing a release.
type ReleasePublishOptions struct {
	RepoPath     string
	DryRun       bool
	ProviderType string
	Result       *ReleaseResult
}

// ReleasePublishResult reports provider release output.
type ReleasePublishResult struct {
	Provider string
	URL      string
	DryRun   bool
}

// ReleasePublish creates or updates a release in the provider.
func ReleasePublish(opts ReleasePublishOptions) (*ReleasePublishResult, error) {
	if opts.RepoPath == "" {
		return nil, fmt.Errorf("repo path is required")
	}
	if opts.Result == nil {
		return nil, fmt.Errorf("release result is required")
	}

	cfgResult, err := config.LoadFromDir(opts.RepoPath)
	if err != nil {
		return nil, ConfigError{Err: err}
	}

	if opts.ProviderType != "" {
		cfgResult.Config.Provider.Type = opts.ProviderType
	}

	if !provider.Enabled(cfgResult.Config) {
		return nil, ConfigError{Err: fmt.Errorf("provider is not configured")}
	}

	if opts.DryRun {
		return &ReleasePublishResult{
			Provider: cfgResult.Config.Provider.Type,
			DryRun:   true,
		}, nil
	}

	pcfg, err := provider.FromAppConfig(cfgResult.Config)
	if err != nil {
		return nil, ConfigError{Err: err}
	}

	p, err := provider.New(pcfg)
	if err != nil {
		return nil, ConfigError{Err: err}
	}

	rel, err := p.CreateRelease(opts.Result.Tag, opts.Result.Tag, opts.Result.Changelog)
	if err != nil {
		if provider.IsReleaseExists(err) {
			rel, err = p.UpdateRelease(opts.Result.Tag, opts.Result.Tag, opts.Result.Changelog)
		}
		if err != nil {
			return nil, ProviderError{Err: err}
		}
	}

	return &ReleasePublishResult{
		Provider: cfgResult.Config.Provider.Type,
		URL:      rel.URL,
		DryRun:   false,
	}, nil
}
