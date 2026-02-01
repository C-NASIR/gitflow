package workflow

import (
	"fmt"
	"gitflow/internal/git"
)

// Status captures current repository state.
type Status struct {
	RepoPath string
	Branch   string
	Dirty    bool
}

// GetStatus inspects a repo and returns its status.
func GetStatus(repoPath string) (*Status, error) {
	client, err := git.NewClient(repoPath)
	if err != nil {
		return nil, err
	}

	branch, err := client.CurrentBranch()
	if err != nil {
		return nil, fmt.Errorf("failed to determine current branch: %w", err)
	}

	dirty, err := client.IsDirty()
	if err != nil {
		return nil, fmt.Errorf("failed to determine working tree state: %w", err)
	}

	return &Status{
		RepoPath: repoPath,
		Branch:   branch,
		Dirty:    dirty,
	}, nil

}
