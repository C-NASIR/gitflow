package workflow

import (
	"gitflow/internal/config"
	"gitflow/internal/git"
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

func setupOriginAndClone(t *testing.T) (origin string, clone string) {
	t.Helper()

	base := t.TempDir()
	origin = filepath.Join(base, "origin.git")
	clone = filepath.Join(base, "work")

	runGit(t, base, "init", "--bare", origin)
	runGit(t, base, "clone", origin, clone)

	runGit(t, clone, "config", "user.email", "test@example.com")
	runGit(t, clone, "config", "user.name", "Test user")

	readme := filepath.Join(clone, "README.md")
	if err := os.WriteFile(readme, []byte("hello"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	runGit(t, clone, "add", ".")
	runGit(t, clone, "commit", "-m", "initial")
	runGit(t, clone, "branch", "-M", "main")
	runGit(t, clone, "push", "origin", "main")

	return origin, clone
}

func TestStartCreatesBranchAndPushed(t *testing.T) {
	_, repo := setupOriginAndClone(t)

	cfg := config.Default()
	cfg.Workflows.Start.BaseBranch = "main"
	cfg.Workflows.Start.FetchFirst = true
	cfg.Workflows.Start.AutoPush = true
	cfg.Branches.FeaturePrefix = "feature/"

	res, err := Start(cfg, StartOptions{
		Kind:     "feature",
		RepoPath: repo,
		Remote:   "origin",
		Name:     "user auth",
	})
	if err != nil {
		t.Fatalf("Start: %v", err)
	}

	if res.NewBranch != "feature/user-auth" {
		t.Fatalf("unexpected branch: %s", res.NewBranch)
	}

	client, err := git.NewClient(repo)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	current, err := client.CurrentBranch()
	if err != nil {
		t.Fatalf("CurrentBranch: %v", err)
	}

	if current != "feature/user-auth" {
		t.Fatalf("expected current branch feature/user-auth got %s", current)
	}

	out, err := client.Run("ls-remote", "--heads", "origin", "feature/user-auth")
	if err != nil {
		t.Fatalf("ls-remote: %v", err)
	}
	if out == "" {
		t.Fatalf("expected remote branch to exist")
	}
}

func TestStartRejectsDirtyRepo(t *testing.T) {
	_, repo := setupOriginAndClone(t)

	if err := os.WriteFile(filepath.Join(repo, "x.txt"), []byte("x"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	cfg := config.Default()
	cfg.Workflows.Start.AutoPush = true

	_, err := Start(cfg, StartOptions{
		Kind:     "feature",
		RepoPath: repo,
		Remote:   "origin",
		Name:     "test",
	})

	if err == nil {
		t.Fatalf("expected error for dirty repo")
	}
}
