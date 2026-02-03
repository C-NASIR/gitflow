package workflow

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"gitflow/internal/config"
	"gitflow/internal/git"
)

func runGitCleanup(t *testing.T, dir string, env []string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), env...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v output: %s", args, err, string(out))
	}
}

func setupRepoForCleanup(t *testing.T) string {
	t.Helper()

	base := t.TempDir()
	origin := filepath.Join(base, "origin.git")
	work := filepath.Join(base, "work")

	runGitCleanup(t, base, nil, "init", "--bare", origin)
	runGitCleanup(t, base, nil, "clone", origin, work)

	runGitCleanup(t, work, nil, "config", "user.email", "test@example.com")
	runGitCleanup(t, work, nil, "config", "user.name", "Test User")

	readme := filepath.Join(work, "README.md")
	if err := os.WriteFile(readme, []byte("hello"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	runGitCleanup(t, work, nil, "add", ".")
	runGitCleanup(t, work, nil, "commit", "-m", "initial")
	runGitCleanup(t, work, nil, "branch", "-M", "main")
	runGitCleanup(t, work, nil, "push", "-u", "origin", "main")

	return work
}

func TestCleanupDeletesMergedBranch(t *testing.T) {
	repo := setupRepoForCleanup(t)

	runGitCleanup(t, repo, nil, "checkout", "-b", "feature/old")
	f := filepath.Join(repo, "f.txt")
	if err := os.WriteFile(f, []byte("x"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	old := time.Now().Add(-10 * 24 * time.Hour).Format(time.RFC3339)
	env := []string{
		"GIT_AUTHOR_DATE=" + old,
		"GIT_COMMITTER_DATE=" + old,
	}

	runGitCleanup(t, repo, env, "add", ".")
	runGitCleanup(t, repo, env, "commit", "-m", "feature")
	runGitCleanup(t, repo, nil, "checkout", "main")
	runGitCleanup(t, repo, nil, "merge", "--no-ff", "feature/old", "-m", "merge feature")

	cfg := config.Default()
	cfg.Branches.MainBranch = "main"
	cfg.Workflows.Start.BaseBranch = "main"
	cfg.Workflows.Cleanup.MergedOnly = true
	cfg.Workflows.Cleanup.ProtectedBranches = []string{"main"}

	out, err := Cleanup(cfg, CleanupOptions{
		RepoPath: repo,
		Remote:   "origin",
		Yes:      true,
	})
	if err != nil {
		t.Fatalf("Cleanup: %v", err)
	}

	found := false
	for _, b := range out.Deleted {
		if b == "feature/old" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected feature/old to be deleted")
	}

	client, err := git.NewClient(repo)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	base := cfg.Branches.MainBranch
	branches, err := client.ListLocalBranches(base)
	if err != nil {
		t.Fatalf("ListLocalBranches: %v", err)
	}
	for _, b := range branches {
		if b.Name == "feature/old" {
			t.Fatalf("expected feature/old to be gone")
		}
	}
}

func TestCleanupDoesNotDeleteProtectedOrCurrent(t *testing.T) {
	repo := setupRepoForCleanup(t)

	runGitCleanup(t, repo, nil, "checkout", "-b", "feature/keep")

	cfg := config.Default()
	cfg.Branches.MainBranch = "main"
	cfg.Workflows.Start.BaseBranch = "main"
	cfg.Workflows.Cleanup.MergedOnly = false
	cfg.Workflows.Cleanup.AgeThresholdDays = 0
	cfg.Workflows.Cleanup.ProtectedBranches = []string{"main", "feature/keep"}

	out, err := Cleanup(cfg, CleanupOptions{
		RepoPath:     repo,
		Remote:       "origin",
		Yes:          true,
		All:          true,
		AgeThreshold: 0,
	})
	if err != nil {
		t.Fatalf("Cleanup: %v", err)
	}

	for _, b := range out.Deleted {
		if b == "main" || b == "feature/keep" {
			t.Fatalf("deleted protected branch %s", b)
		}
	}
}
