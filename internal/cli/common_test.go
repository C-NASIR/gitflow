package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"gitflow/internal/config"

	"github.com/spf13/cobra"
)

func TestNoColorOverridesConfig(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Default()
	cfg.UI.Color = true
	if err := config.WriteFile(filepath.Join(dir, ".gitflow.yml"), cfg); err != nil {
		t.Fatalf("write config: %v", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer func() {
		_ = os.Chdir(wd)
	}()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	color := false
	SetUIOverrides(UIOverrides{Color: &color})
	defer SetUIOverrides(UIOverrides{})

	var buf bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&buf)

	common, err := CommonFromCmd(cmd)
	if err != nil {
		t.Fatalf("CommonFromCmd: %v", err)
	}
	if common.UI.ColorEnabled() {
		t.Fatalf("expected color override to disable output")
	}
}

func TestEmojiOverridesConfig(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Default()
	cfg.UI.Emoji = false
	if err := config.WriteFile(filepath.Join(dir, ".gitflow.yml"), cfg); err != nil {
		t.Fatalf("write config: %v", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer func() {
		_ = os.Chdir(wd)
	}()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	emoji := true
	SetUIOverrides(UIOverrides{Emoji: &emoji})
	defer SetUIOverrides(UIOverrides{})

	var buf bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&buf)

	common, err := CommonFromCmd(cmd)
	if err != nil {
		t.Fatalf("CommonFromCmd: %v", err)
	}
	if !common.UI.EmojiEnabled() {
		t.Fatalf("expected emoji override to enable output")
	}
}
