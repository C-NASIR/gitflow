package workflow

import (
	"fmt"
	"gitflow/internal/config"
	"gitflow/internal/git"
	"regexp"
	"strings"
)

type StartOptions struct {
	Kind     string
	RepoPath string
	Remote   string
	Name     string
}

type StartResult struct {
	BaseBranch string
	NewBranch  string
	Pushed     bool
}

func Start(cfg *config.Config, opts StartOptions) (*StartResult, error) {
	if opts.RepoPath == "" {
		return nil, fmt.Errorf("repo path is required")
	}
	if opts.Remote == "" {
		opts.Remote = ""
	}
	if strings.TrimSpace(opts.Name) == "" {
		return nil, fmt.Errorf("branch name is required")
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

	base := cfg.Workflows.Start.BaseBranch
	if base == "" {
		base = cfg.Branches.MainBranch
	}
	if base == "" {
		base = "main"
	}

	remoteExists, err := client.HasRemote(opts.Remote)
	if err != nil {
		return nil, err
	}
	if !remoteExists {
		return nil, fmt.Errorf("remote %s not found", opts.Remote)
	}

	if cfg.Workflows.Start.FetchFirst {
		if err := client.Fetch(opts.Remote); err != nil {
			return nil, err
		}
	}

	if err := client.Pull(opts.Remote, base); err != nil {
		return nil, err
	}

	prefix := prefixForKind(cfg, opts.Kind)
	slug := slugify(opts.Name)
	newBranch := prefix + slug

	if err := client.CheckoutNew(newBranch); err != nil {
		return nil, err
	}

	pushed := false
	if cfg.Workflows.Start.AutoPush {
		if err := client.PushSetUpstream(opts.Remote, newBranch); err != nil {
			return nil, err
		}
		pushed = true
	}

	return &StartResult{
		BaseBranch: base,
		NewBranch:  newBranch,
		Pushed:     pushed,
	}, nil

}

func prefixForKind(cfg *config.Config, kind string) string {
	switch kind {
	case "bugfix":
		if cfg.Branches.BugfixPrefix != "" {
			return cfg.Branches.BugfixPrefix
		}
		return "bugfix/"
	case "hotfix":
		if cfg.Branches.HotfixPrefix != "" {
			return cfg.Branches.HotfixPrefix
		}
		return "hotfix/"
	default:
		return "feature/"

	}
}

var nonSlug = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, "_", " ")
	s = nonSlug.ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)
	s = strings.Join(strings.Fields(s), "-")
	s = strings.Trim(s, "-")
	if s == "" {
		return "work"
	}
	return s
}
