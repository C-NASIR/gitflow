package workflow

import (
	"testing"

	"gitflow/internal/config"
)

func TestWorkflowListBranches(t *testing.T) {
	repo := setupCommitRepo(t)

	cfg := config.Default()
	out, err := ListBranches(cfg, BranchListOptions{
		RepoPath: repo,
		Base:     "",
	})
	if err != nil {
		t.Fatalf("ListBranches: %v", err)
	}
	if out.Base == "" {
		t.Fatalf("expected base")
	}
	if len(out.Branches) == 0 {
		t.Fatalf("expected branches")
	}
}
