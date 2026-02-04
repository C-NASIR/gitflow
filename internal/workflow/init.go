package workflow

import (
	"fmt"
	"os"
	"path/filepath"

	"gitflow/internal/config"
)

// InitOptions defines inputs for gitflow init.
type InitOptions struct {
	RepoPath string
	Force    bool
}

// InitResult reports the config path written.
type InitResult struct {
	Path string
}

// Init writes a gitflow config file to the repository.
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
