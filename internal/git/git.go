// Package git provides a thin wrapper over git commands.
package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Client runs git commands scoped to a repository path.
type Client struct {
	repoPath string
}

// NewClient creates a client after verifying the repo exists.
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

// Run executes a git command in the repository and returns stdout.
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

// CurrentBranch returns the checked-out branch name.
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

// IsDirty reports whether the working tree has changes.
func (c *Client) IsDirty() (bool, error) {
	out, err := c.Run("status", "--porcelain", "--untracked-files=all")
	if err != nil {
		return false, err
	}
	return len(out) > 0, nil
}

// Fetch updates refs from the named remote.
func (c *Client) Fetch(remote string) error {
	_, err := c.Run("fetch", remote)
	return err
}

// Checkout switches to the named branch.
func (c *Client) Checkout(branch string) error {
	_, err := c.Run("checkout", branch)
	return err
}

// CheckoutNew creates and switches to a new branch.
func (c *Client) CheckoutNew(branch string) error {
	_, err := c.Run("checkout", "-b", branch)
	return err
}

// Pull updates the branch from the remote.
func (c *Client) Pull(remote, branch string) error {
	_, err := c.Run("pull", remote, branch)
	return err
}

// PushSetUpstream pushes the branch and sets upstream tracking.
func (c *Client) PushSetUpstream(remote, branch string) error {
	_, err := c.Run("push", "-u", remote, branch)
	return err
}

// HasRemote reports whether a named remote is configured.
func (c *Client) HasRemote(remote string) (bool, error) {
	out, err := c.Run("remote")
	if err != nil {
		return false, err
	}
	for _, line := range strings.Split(out, "\n") {
		if strings.TrimSpace(line) == remote {
			return true, nil
		}
	}
	return false, nil
}
