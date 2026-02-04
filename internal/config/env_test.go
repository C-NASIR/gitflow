package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnvOverridesConfig(t *testing.T) {
	if err := os.Setenv("GITFLOW_RELEASE_TAG_PREFIX", "r"); err != nil {
		t.Fatalf("set env: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Unsetenv("GITFLOW_RELEASE_TAG_PREFIX")
	})

	dir := t.TempDir()
	data := []byte("release:\n  tag_prefix: v\n")
	if err := os.WriteFile(filepath.Join(dir, ".gitflow.yml"), data, 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	res, err := LoadFromDir(dir)
	if err != nil {
		t.Fatalf("LoadFromDir: %v", err)
	}
	if res.Config.Release.TagPrefix != "r" {
		t.Fatalf("expected env override tag prefix, got %s", res.Config.Release.TagPrefix)
	}
}
