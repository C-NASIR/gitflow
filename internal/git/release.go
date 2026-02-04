package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Commit represents git commit metadata used in releases.
type Commit struct {
	Hash    string
	Subject string
	Body    string
	Date    string
}

// ListTags returns all tags in the repository.
func (c *Client) ListTags() ([]string, error) {
	out, err := c.Run("tag", "--list")
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}
	lines := strings.Split(out, "\n")
	var tags []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		tags = append(tags, line)
	}
	return tags, nil
}

// CommitsBetween returns commits between two refs.
func (c *Client) CommitsBetween(fromRef, toRef string) ([]Commit, error) {
	format := "%H%x1f%s%x1f%b%x1f%cs%x1e"
	args := []string{"log", "--pretty=format:" + format}
	if toRef == "" {
		toRef = "HEAD"
	}
	if fromRef != "" {
		args = append(args, fmt.Sprintf("%s..%s", fromRef, toRef))
	} else {
		args = append(args, toRef)
	}
	out, err := c.Run(args...)
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}
	entries := strings.Split(out, "\x1e")
	var commits []Commit
	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		parts := strings.Split(entry, "\x1f")
		if len(parts) < 4 {
			continue
		}
		commits = append(commits, Commit{
			Hash:    parts[0],
			Subject: strings.TrimSpace(parts[1]),
			Body:    strings.TrimSpace(parts[2]),
			Date:    strings.TrimSpace(parts[3]),
		})
	}
	return commits, nil
}

// TagExists reports whether a tag is present.
func (c *Client) TagExists(tag string) (bool, error) {
	cmd := exec.Command("git", "show-ref", "--tags", "--verify", "--quiet", "refs/tags/"+tag)
	cmd.Dir = c.repoPath
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() != 0 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CreateAnnotatedTag creates an annotated tag with a message.
func (c *Client) CreateAnnotatedTag(tag, message string) error {
	exists, err := c.TagExists(tag)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("tag %s already exists", tag)
	}
	cmd := exec.Command("git", "tag", "-a", tag, "-m", message)
	cmd.Dir = c.repoPath
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return fmt.Errorf("git tag failed: %s", msg)
	}
	return nil
}
