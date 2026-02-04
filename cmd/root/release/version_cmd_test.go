package release

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestReleaseVersionCommandOutput(t *testing.T) {
	repo := setupReleaseVersionRepo(t)
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
	cmd := versionCmd()
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("version command: %v", err)
	}

	if got := buf.String(); got != "0.1.1\n" {
		t.Fatalf("expected version output, got %q", got)
	}
}

func setupReleaseVersionRepo(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	runGitReleaseCmd(t, dir, "init")
	runGitReleaseCmd(t, dir, "config", "user.email", "test@example.com")
	runGitReleaseCmd(t, dir, "config", "user.name", "Test User")
	writeReleaseFile(t, dir, "readme.md", "init")
	runGitReleaseCmd(t, dir, "add", "-A")
	runGitReleaseCmd(t, dir, "commit", "-m", "chore: init")
	runGitReleaseCmd(t, dir, "tag", "v0.1.0")

	writeReleaseFile(t, dir, "readme.md", "fix")
	runGitReleaseCmd(t, dir, "add", "-A")
	runGitReleaseCmd(t, dir, "commit", "-m", "fix: patch bug")

	return dir
}

func runGitReleaseCmd(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v output: %s", args, err, string(out))
	}
}

func writeReleaseFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}
