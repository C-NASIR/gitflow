package workflow

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gitflow/internal/config"
	"gitflow/internal/git"
	"gitflow/internal/provider"
	"gitflow/pkg/types"
)

type PRCreateOptions struct {
	RepoPath string
	Remote   string

	Title       string
	Description string
	BaseBranch  string
	Draft       *bool

	Reviewers []string
	Labels    []string
}

type PRCreateResult struct {
	PR *types.PullRequest
}

func CreatePR(cfg *config.Config, opts PRCreateOptions) (*PRCreateResult, error) {
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

	currentBranch, err := client.CurrentBranch()
	if err != nil {
		return nil, err
	}

	base := strings.TrimSpace(opts.BaseBranch)
	if base == "" {
		base = strings.TrimSpace(cfg.Workflows.Start.BaseBranch)
	}
	if base == "" {
		base = strings.TrimSpace(cfg.Branches.MainBranch)
	}
	if base == "" {
		base = "main"
	}

	if currentBranch == base {
		return nil, fmt.Errorf("current branch is base branch %s", base)
	}

	hasRemote, err := client.HasRemote(opts.Remote)
	if err != nil {
		return nil, err
	}
	if !hasRemote {
		return nil, fmt.Errorf("remote %s not found", opts.Remote)
	}

	hasUpstream, err := client.HasUpstream()
	if err != nil {
		return nil, err
	}
	if !hasUpstream {
		if err := client.PushSetUpstream(opts.Remote, currentBranch); err != nil {
			return nil, err
		}
	}

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

	title := strings.TrimSpace(opts.Title)
	if title == "" {
		title = defaultTitleFromBranch(currentBranch)
	}

	draft := cfg.Workflows.PR.Draft
	if opts.Draft != nil {
		draft = *opts.Draft
	}

	reviewers := opts.Reviewers
	if len(reviewers) == 0 {
		reviewers = cfg.Workflows.PR.DefaultReviewers
	}

	labels := opts.Labels
	if len(labels) == 0 {
		labels = cfg.Workflows.PR.Labels
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	pr, err := p.CreatePR(ctx, provider.CreatePROptions{
		Title:       title,
		Description: opts.Description,
		HeadBranch:  currentBranch,
		BaseBranch:  base,
		Draft:       draft,
		Reviewers:   reviewers,
		Labels:      labels,
	})
	if err != nil {
		return nil, err
	}

	return &PRCreateResult{PR: pr}, nil
}

func defaultTitleFromBranch(branch string) string {
	b := branch
	b = strings.TrimPrefix(b, "feature/")
	b = strings.TrimPrefix(b, "bugfix/")
	b = strings.TrimPrefix(b, "hotfix/")
	b = strings.ReplaceAll(b, "-", " ")
	b = strings.ReplaceAll(b, "_", " ")
	b = strings.TrimSpace(b)
	if b == "" {
		return "Update"
	}
	return strings.Title(b)
}
