package workflow

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"gitflow/internal/config"
	"gitflow/internal/git"
)

func runGitSync(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v output: %s", args, err, string(out))
	}
	return string(out)
}

func setupTwoClones(t *testing.T) (origin string, a string, b string) {
	t.Helper()

	base := t.TempDir()
	origin = filepath.Join(base, "origin.git")
	a = filepath.Join(base, "a")
	b = filepath.Join(base, "b")

	runGitSync(t, base, "init", "--bare", origin)
	runGitSync(t, base, "clone", origin, a)
	runGitSync(t, base, "clone", origin, b)

	runGitSync(t, a, "config", "user.email", "test@example.com")
	runGitSync(t, a, "config", "user.name", "Test User")
	runGitSync(t, b, "config", "user.email", "test@example.com")
	runGitSync(t, b, "config", "user.name", "Test User")

	readme := filepath.Join(a, "README.md")
	if err := os.WriteFile(readme, []byte("hello"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	runGitSync(t, a, "add", ".")
	runGitSync(t, a, "commit", "-m", "initial")
	runGitSync(t, a, "branch", "-M", "main")
	runGitSync(t, a, "push", "-u", "origin", "main")

	runGitSync(t, b, "fetch", "origin")
	runGitSync(t, b, "checkout", "main")

	return origin, a, b
}

func TestSyncRebaseUpdatesBranchAndPushes(t *testing.T) {
	_, a, b := setupTwoClones(t)

	runGitSync(t, a, "checkout", "-b", "feature/user-auth")
	f := filepath.Join(a, "feature.txt")
	if err := os.WriteFile(f, []byte("feature"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	runGitSync(t, a, "add", ".")
	runGitSync(t, a, "commit", "-m", "feature commit")
	runGitSync(t, a, "push", "-u", "origin", "feature/user-auth")

	runGitSync(t, b, "checkout", "main")
	m := filepath.Join(b, "main.txt")
	if err := os.WriteFile(m, []byte("main change"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	runGitSync(t, b, "add", ".")
	runGitSync(t, b, "commit", "-m", "main advance")
	runGitSync(t, b, "push", "origin", "main")

	runGitSync(t, a, "checkout", "feature/user-auth")

	cfg := config.Default()
	cfg.Branches.MainBranch = "main"
	cfg.Workflows.Start.BaseBranch = "main"
	cfg.Workflows.Sync.Strategy = "rebase"
	cfg.Workflows.Sync.AutoPush = true
	cfg.Workflows.Sync.ForcePush = true

	out, err := Sync(cfg, SyncOptions{
		RepoPath: a,
		Remote:   "origin",
	})
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}

	if out.Strategy != "rebase" {
		t.Fatalf("expected rebase got %s", out.Strategy)
	}

	client, err := git.NewClient(a)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	logOut, err := client.Run("log", "--oneline", "--decorate", "-n", "50")
	if err != nil {
		t.Fatalf("log: %v", err)
	}
	if logOut == "" {
		t.Fatalf("expected non empty log")
	}

	remoteRef, err := client.Run("rev-parse", "origin/feature/user-auth")
	if err != nil {
		t.Fatalf("rev-parse remote: %v", err)
	}
	headRef, err := client.Run("rev-parse", "HEAD")
	if err != nil {
		t.Fatalf("rev-parse head: %v", err)
	}
	if remoteRef != headRef {
		t.Fatalf("expected remote branch to match head after push")
	}
}

func TestSyncRejectsDirtyRepo(t *testing.T) {
	_, a, _ := setupTwoClones(t)

	runGitSync(t, a, "checkout", "-b", "feature/x")
	if err := os.WriteFile(filepath.Join(a, "x.txt"), []byte("x"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	cfg := config.Default()
	cfg.Workflows.Sync.AutoPush = false

	_, err := Sync(cfg, SyncOptions{
		RepoPath: a,
		Remote:   "origin",
	})
	if err == nil {
		t.Fatalf("expected error for dirty repo")
	}
}
