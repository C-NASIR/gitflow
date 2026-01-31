package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaultsWhenNotFile(t *testing.T) {
	dir := t.TempDir()

	res, err := LoadFromDir(dir)
	if err != nil {
		t.Fatalf("LoadFromDir: %v", err)
	}

	if res.Path != "" {
		t.Fatalf("expected empty path, got %s", res.Path)
	}

	if res.Config.Branches.MainBranch == "" {
		t.Fatalf("expected default main branch to be set")
	}
}

func TestLoadFromCurrentDirFile(t *testing.T) {
	dir := t.TempDir()

	cfgPath := filepath.Join(dir, ".gitflow.yaml")
	data := []byte("branches:\n  main_branch: trunk\n  feature_prefix: feat/\nworkflows:\n  sync:\n    strategy: merge\n")
	if err := os.WriteFile(cfgPath, data, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	res, err := LoadFromDir(dir)
	if err != nil {
		t.Fatalf("LoadFromDir: %v", err)
	}

	if res.Path != cfgPath {
		t.Fatalf("expected path %s got %s", cfgPath, res.Path)
	}

	if res.Config.Branches.MainBranch != "trunk" {
		t.Fatalf("expected trunk got %s", res.Config.Branches.MainBranch)
	}

	if res.Config.Workflows.Sync.Strategy != "merge" {
		t.Fatalf("expected merge got %s", res.Config.Workflows.Sync.Strategy)
	}
}

func TestValidateRejectBadSyncStrategy(t *testing.T) {
	cfg := Default()
	cfg.Workflows.Sync.Strategy = "something"
	if err := cfg.Validate(); err == nil {
		t.Fatalf("expected validation error")
	}
}
