package workflow

import (
	"context"
	"fmt"
	"time"

	"gitflow/internal/config"
	"gitflow/internal/provider"
	"gitflow/pkg/types"
)

// PRViewResult contains a pull request view result.
type PRViewResult struct {
	PR *types.PullRequest
}

// ViewPR fetches a pull request by number.
func ViewPR(cfg *config.Config, number int) (*PRViewResult, error) {
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

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	pr, err := p.GetPR(ctx, number)
	if err != nil {
		return nil, err
	}

	return &PRViewResult{PR: pr}, nil
}
