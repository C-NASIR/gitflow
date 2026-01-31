package workflow

import (
	"os"
	"os/exec"
	"testing"
)

func setupRepo(t *testing.T) string {
	t.Helper()

	dir, err := os.MkdirTemp("", "gitflow-workflow-repo-*")
	if err != nil {
		t.Fatalf("temp dir: %v", err)
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git init failed: %v output: %s", err, string(out))
	}

	return dir
}

func TestGetStatus(t *testing.T) {
	dir := setupRepo(t)
	defer os.RemoveAll(dir)

	s, err := GetStatus(dir)
	if err != nil {
		t.Fatalf("GetStatus: %v", err)
	}
	if s.Branch == "" {
		t.Fatalf("expected non empty branch")
	}
	if s.Dirty {
		t.Fatalf("expected clean repo")
	}
}
