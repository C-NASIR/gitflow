package root

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestStatusCommandOutputHeaders(t *testing.T) {
	repo := setupStatusRepo(t)
	defer os.RemoveAll(repo)

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer func() {
		_ = os.Chdir(wd)
	}()
	if err := os.Chdir(repo); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	var buf bytes.Buffer
	cmd := statusCmd()
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("status command: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "KEY") || !strings.Contains(output, "VALUE") {
		t.Fatalf("expected table headers in output")
	}
}

func setupStatusRepo(t *testing.T) string {
	t.Helper()

	dir, err := os.MkdirTemp("", "gitflow-status-repo-*")
	if err != nil {
		t.Fatalf("temp dir: %v", err)
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		_ = os.RemoveAll(dir)
		t.Fatalf("git init failed: %v output: %s", err, string(out))
	}

	return dir
}
