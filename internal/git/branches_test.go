package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v output: %s", args, err, string(out))
	}
}

func setupBranchRepo(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.email", "test@example.com")
	runGit(t, dir, "config", "user.name", "Test User")

	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	runGit(t, dir, "add", "-A")
	runGit(t, dir, "commit", "-m", "initial")

	return dir
}

func TestListLocalBranches(t *testing.T) {
	repo := setupBranchRepo(t)

	runGit(t, repo, "checkout", "-b", "feature/x")
	if err := os.WriteFile(filepath.Join(repo, "b.txt"), []byte("b"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	runGit(t, repo, "add", "-A")
	runGit(t, repo, "commit", "-m", "add b")

	c, err := NewClient(repo)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	branches, err := c.ListLocalBranches("master")
	if err != nil {
		branches, err = c.ListLocalBranches("main")
	}
	if err != nil {
		t.Fatalf("ListLocalBranches: %v", err)
	}

	if len(branches) < 1 {
		t.Fatalf("expected branches")
	}

	foundCurrent := false
	for _, b := range branches {
		if b.Current {
			foundCurrent = true
		}
	}
	if !foundCurrent {
		t.Fatalf("expected a current branch")
	}
}
