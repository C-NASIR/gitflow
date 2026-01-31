package workflow

import (
	"fmt"
	"gitflow/internal/git"
)

type Status struct {
	RepoPath string
	Branch   string
	Dirty    bool
}

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
