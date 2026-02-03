package workflow

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gitflow/internal/config"
	"gitflow/internal/provider"
	"gitflow/pkg/types"
)

type PRListOptions struct {
	State string
}

type PRListResult struct {
	PRs []*types.PullRequest
}

func ListPRs(cfg *config.Config, opts PRListOptions) (*PRListResult, error) {
	if !provider.Enabled(cfg) {
		return nil, fmt.Errorf("provider is not configured in .gitflow.yml")
	}

	pcfg, err := provider.FromAppConfig(cfg)
	if err != nil {
		return nil, err
	}

	p, err := provider.New(pcfg)
	if err != nil {
		return nil, err
	}

	state := strings.TrimSpace(opts.State)
	if state == "" {
		state = "open"
	}
	if state == "all" {
		state = "all"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	prs, err := p.ListPRs(ctx, state)
	if err != nil {
		return nil, err
	}

	return &PRListResult{PRs: prs}, nil
}
