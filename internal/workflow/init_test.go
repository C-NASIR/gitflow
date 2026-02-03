package workflow

import (
	"os"
	"path/filepath"
	"testing"

	"gitflow/internal/config"
)

func TestInitWritesFile(t *testing.T) {
	dir := t.TempDir()

	cfg := config.Default()
	out, err := Init(cfg, InitOptions{
		RepoPath: dir,
		Force:    false,
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	if out.Path == "" {
		t.Fatalf("expected path")
	}

	_, err = os.Stat(filepath.Join(dir, ".gitflow.yml"))
	if err != nil {
		t.Fatalf("expected file to exist")
	}
}
