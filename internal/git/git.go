package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Client struct {
	repoPath string
}

func NewClient(repoPath string) (*Client, error) {
	c := &Client{repoPath: repoPath}

	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = repoPath
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return nil, fmt.Errorf("not a git repository at %s: %s", repoPath, msg)
	}

	return c, nil
}

func (c *Client) Run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = c.repoPath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stdout.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("git %s failed: %s", strings.Join(args, " "), msg)
	}

	return strings.TrimSpace(stdout.String()), nil

}

func (c *Client) CurrentBranch() (string, error) {
	out, err := c.Run("symbolic-ref", "--short", "HEAD")
	if err == nil && out != "" {
		return out, nil
	}

	out, err = c.Run("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}
	return out, nil
}

func (c *Client) IsDirty() (bool, error) {
	out, err := c.Run("status", "--porcelain")
	if err != nil {
		return false, err
	}
	return len(out) > 0, nil
}
