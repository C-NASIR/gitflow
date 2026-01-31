package config

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func findConfig(startDir string) (string, error) {
	if p := filepath.Join(startDir, ".gitflow.yaml"); fileExists(p) {
		return p, nil
	}

	gitRoot, err := gitTopLevel(startDir)
	if err == nil {
		if p := filepath.Join(gitRoot, ".gitflow.yaml"); fileExists(p) {
			return p, nil
		}
	}

	home, err := os.UserHomeDir()
	if err == nil {
		if p := filepath.Join(home, ".gitflow.yaml"); fileExists(p) {
			return p, nil
		}
	}

	return "", fmt.Errorf("no config file found")
}

func fileExists(path string) bool {
	st, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !st.IsDir()
}

func gitTopLevel(dir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = dir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s", stderr.String())
	}

	out := strings.TrimSpace(stdout.String())
	if out == "" {
		return "", fmt.Errorf("empty git root")
	}

	return out, nil
}
