package workflow

import (
	"fmt"
	"strings"

	"gitflow/internal/config"
	"gitflow/internal/git"
	"gitflow/pkg/types"
)

// BranchListOptions defines inputs for listing branches.
type BranchListOptions struct {
	RepoPath string
	Base     string
}

// BranchListResult contains branch list results.
type BranchListResult struct {
	Base     string
	Branches []*types.Branch
}

// ListBranches lists local branches and metadata.
func ListBranches(cfg *config.Config, opts BranchListOptions) (*BranchListResult, error) {
	if strings.TrimSpace(opts.RepoPath) == "" {
		return nil, fmt.Errorf("repo path is required")
	}

	base := strings.TrimSpace(opts.Base)
	if base == "" {
		base = strings.TrimSpace(cfg.Workflows.Start.BaseBranch)
	}
	if base == "" {
		base = strings.TrimSpace(cfg.Branches.MainBranch)
	}
	if base == "" {
		base = "main"
	}

	client, err := git.NewClient(opts.RepoPath)
	if err != nil {
		return nil, err
	}

	branches, err := client.ListLocalBranches(base)
	if err != nil {
		return nil, err
	}

	return &BranchListResult{
		Base:     base,
		Branches: branches,
	}, nil
}
