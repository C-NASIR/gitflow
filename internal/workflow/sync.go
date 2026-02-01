package workflow

import (
	"fmt"
	"gitflow/internal/config"
	"gitflow/internal/git"
	"strings"
)

// SyncOptions defines inputs for syncing the current branch.
type SyncOptions struct {
	RepoPath          string
	Remote            string
	StrategyOverride  string
	AutoPushOverride  *bool
	ForcePushOverride *bool
}

// SyncResult reports the outcome of a sync operation.
type SyncResult struct {
	BaseBranch    string
	CurrentBranch string
	Strategy      string
	Pushed        bool
	ForcePushed   bool
}

// Sync updates the current branch from the base branch.
func Sync(cfg *config.Config, opts SyncOptions) (*SyncResult, error) {
	if opts.RepoPath == "" {
		return nil, fmt.Errorf("repo path is required")
	}

	if opts.Remote == "" {
		opts.Remote = "origin"
	}

	client, err := git.NewClient(opts.RepoPath)
	if err != nil {
		return nil, err
	}

	dirty, err := client.IsDirty()
	if err != nil {
		return nil, err
	}

	if dirty {
		return nil, fmt.Errorf("working tree is not clean")
	}

	current, err := client.CurrentBranch()
	if err != nil {
		return nil, err
	}

	base := cfg.Workflows.Start.BaseBranch
	if base == "" {
		base = cfg.Branches.MainBranch
	}

	if base == "" {
		base = "main"
	}

	if strings.TrimSpace(current) == strings.TrimSpace(base) {
		return nil, fmt.Errorf("already on base branch %s", base)
	}

	remoteExists, err := client.HasRemote(opts.Remote)
	if err != nil {
		return nil, err
	}
	if !remoteExists {
		return nil, fmt.Errorf("remote %s not found", opts.Remote)
	}

	strategy := cfg.Workflows.Sync.Strategy
	if strategy == "" {
		strategy = "rebase"
	}
	if opts.StrategyOverride != "" {
		strategy = opts.StrategyOverride
	}

	autoPush := cfg.Workflows.Sync.AutoPush
	if opts.AutoPushOverride != nil {
		autoPush = *opts.AutoPushOverride
	}

	forcePush := cfg.Workflows.Sync.ForcePush
	if opts.ForcePushOverride != nil {
		forcePush = *opts.ForcePushOverride
	}

	if err := client.Fetch(opts.Remote); err != nil {
		return nil, err
	}
	if err := client.Checkout(base); err != nil {
		return nil, err
	}
	if err := client.Pull(opts.Remote, base); err != nil {
		return nil, err
	}
	if err := client.Checkout(current); err != nil {
		return nil, err
	}

	switch strategy {
	case "rebase":
		if err := client.Rebase(base); err != nil {
			return nil, err
		}
	case "merge":
		if err := client.Merge(base); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported sync strategy: %s", strategy)
	}

	pushed := false
	forcePushed := false
	if autoPush {
		useForce := false
		if strategy == "rebase" && forcePush {
			useForce = true
		}

		if err := client.Push(opts.Remote, current, useForce); err != nil {
			return nil, err
		}

		pushed = true
		forcePushed = useForce
	}

	return &SyncResult{
		BaseBranch:    base,
		CurrentBranch: current,
		Strategy:      strategy,
		Pushed:        pushed,
		ForcePushed:   forcePushed,
	}, nil
}
