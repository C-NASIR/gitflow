package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func setupRepo(t *testing.T) string {
	t.Helper()

	dir, err := os.MkdirTemp("", "gitflow-repo-*")
	if err != nil {
		t.Fatalf("temp dir: %v", err)
	}

	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v failed: %v output: %s", args, err, string(out))
		}
	}

	run("init")
	run("config", "user.email", "test@example.com")
	run("config", "user.name", "Test User")

	return dir
}

func TestNewClientRejectsNonRepo(t *testing.T) {
	dir, err := os.MkdirTemp("", "gitflow-nonrepo-*")
	if err != nil {
		t.Fatalf("temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	_, err = NewClient(dir)
	if err == nil {
		t.Fatalf("expected error for non repo directory")
	}
}

func TestCurrentBranchAndDirty(t *testing.T) {
	dir := setupRepo(t)
	defer os.RemoveAll(dir)

	c, err := NewClient(dir)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	branch, err := c.CurrentBranch()
	if err != nil {
		t.Fatalf("CurrentBranch: %v", err)
	}
	if branch == "" {
		t.Fatalf("expected non empty branch")
	}

	dirty, err := c.IsDirty()
	if err != nil {
		t.Fatalf("IsDirty: %v", err)
	}
	if dirty {
		t.Fatalf("expected clean repo right after init")
	}

	file := filepath.Join(dir, "hello.txt")
	if err := os.WriteFile(file, []byte("hi"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	dirty, err = c.IsDirty()
	if err != nil {
		t.Fatalf("IsDirty after change: %v", err)
	}
	if !dirty {
		t.Fatalf("expected dirty repo after creating a file")
	}
}
