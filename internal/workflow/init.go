package workflow

import (
	"fmt"
	"os"
	"path/filepath"

	"gitflow/internal/config"
)

type InitOptions struct {
	RepoPath string
	Force    bool
}

type InitResult struct {
	Path string
}

func Init(cfg *config.Config, opts InitOptions) (*InitResult, error) {
	if opts.RepoPath == "" {
		return nil, fmt.Errorf("repo path is required")
	}

	path := filepath.Join(opts.RepoPath, ".gitflow.yml")

	_, err := os.Stat(path)
	if err == nil && !opts.Force {
		return nil, fmt.Errorf(".gitflow.yml already exists, use --force to overwrite")
	}

	if err := config.WriteFile(path, cfg); err != nil {
		return nil, err
	}

	return &InitResult{Path: path}, nil
}
