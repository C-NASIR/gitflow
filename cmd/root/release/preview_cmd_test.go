package release

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestReleasePreviewJSONOutput(t *testing.T) {
	repo := setupReleasePreviewRepo(t)
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
	cmd := previewCmd()
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"--json"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("preview command: %v", err)
	}

	var payload struct {
		CurrentVersion string `json:"current_version"`
		NextVersion    string `json:"next_version"`
		CommitCount    int    `json:"commit_count"`
		Changelog      string `json:"changelog"`
	}
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatalf("json decode: %v", err)
	}
	if payload.CurrentVersion == "" || payload.NextVersion == "" {
		t.Fatalf("expected version values")
	}
	if payload.CommitCount == 0 {
		t.Fatalf("expected commit count")
	}
	if payload.Changelog == "" {
		t.Fatalf("expected changelog")
	}
}

func setupReleasePreviewRepo(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	runGitPreviewCmd(t, dir, "init")
	runGitPreviewCmd(t, dir, "config", "user.email", "test@example.com")
	runGitPreviewCmd(t, dir, "config", "user.name", "Test User")
	writePreviewFile(t, dir, "readme.md", "init")
	runGitPreviewCmd(t, dir, "add", "-A")
	runGitPreviewCmd(t, dir, "commit", "-m", "chore: init")
	runGitPreviewCmd(t, dir, "tag", "v0.1.0")

	writePreviewFile(t, dir, "readme.md", "feat")
	runGitPreviewCmd(t, dir, "add", "-A")
	runGitPreviewCmd(t, dir, "commit", "-m", "feat: add preview")

	return dir
}

func runGitPreviewCmd(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v output: %s", args, err, string(out))
	}
}

func writePreviewFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}
