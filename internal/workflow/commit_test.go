package workflow

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"gitflow/internal/config"
)

func runGitCommitTest(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v output: %s", args, err, string(out))
	}
}

func setupCommitRepo(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	runGitCommitTest(t, dir, "init")
	runGitCommitTest(t, dir, "config", "user.email", "test@example.com")
	runGitCommitTest(t, dir, "config", "user.name", "Test User")

	f := filepath.Join(dir, "a.txt")
	if err := os.WriteFile(f, []byte("a"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	runGitCommitTest(t, dir, "add", "-A")
	runGitCommitTest(t, dir, "commit", "-m", "initial")

	return dir
}

func TestCommitConventional(t *testing.T) {
	repo := setupCommitRepo(t)

	f := filepath.Join(repo, "a.txt")
	if err := os.WriteFile(f, []byte("b"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	runGitCommitTest(t, repo, "add", "-A")

	cfg := config.Default()
	cfg.Commits.Conventional = true
	cfg.Commits.Types = []string{"feat", "fix"}

	out, err := Commit(cfg, CommitOptions{
		RepoPath: repo,
		Message:  "add login",
		Type:     "feat",
		Scope:    "auth",
		Breaking: false,
	})
	if err != nil {
		t.Fatalf("Commit: %v", err)
	}

	if out.Message == "" {
		t.Fatalf("expected message")
	}
}

func TestCommitRequiresStagedChanges(t *testing.T) {
	repo := setupCommitRepo(t)

	cfg := config.Default()
	_, err := Commit(cfg, CommitOptions{
		RepoPath: repo,
		Message:  "x",
	})
	if err == nil {
		t.Fatalf("expected error")
	}
}
