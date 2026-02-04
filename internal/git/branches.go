package git

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitflow/pkg/types"
)

// ListLocalBranches returns local branches with summary metadata.
func (c *Client) ListLocalBranches(baseBranch string) ([]*types.Branch, error) {
	current, err := c.CurrentBranch()
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(baseBranch) == "" {
		baseBranch = "main"
	}

	format := "%(refname:short)|%(authorname)|%(subject)"
	out, err := c.Run("for-each-ref", "--format="+format, "refs/heads/")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	var branches []*types.Branch

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "|", 3)
		if len(parts) != 3 {
			continue
		}

		name := strings.TrimSpace(parts[0])
		author := strings.TrimSpace(parts[1])
		subject := strings.TrimSpace(parts[2])

		ageDays, _ := c.branchAgeDays(name)

		ahead, behind := 0, 0
		if name != baseBranch {
			a, b, err := c.aheadBehind(name, baseBranch)
			if err == nil {
				ahead = a
				behind = b
			}
		}

		branches = append(branches, &types.Branch{
			Name:          name,
			Current:       name == current,
			LastCommitMsg: subject,
			Author:        author,
			AgeDays:       ageDays,
			Ahead:         ahead,
			Behind:        behind,
		})
	}

	return branches, nil
}

func (c *Client) aheadBehind(branch string, base string) (ahead int, behind int, err error) {
	out, err := c.Run("rev-list", "--left-right", "--count", fmt.Sprintf("%s...%s", base, branch))
	if err != nil {
		return 0, 0, err
	}

	fields := strings.Fields(strings.TrimSpace(out))
	if len(fields) != 2 {
		return 0, 0, fmt.Errorf("unexpected rev-list output: %q", out)
	}

	behindCount, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, err
	}
	aheadCount, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, err
	}

	return aheadCount, behindCount, nil
}

func (c *Client) branchAgeDays(branch string) (int, error) {
	out, err := c.Run("log", "-1", "--format=%ct", branch)
	if err != nil {
		return 0, err
	}

	sec, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
	if err != nil {
		return 0, err
	}

	commitTime := time.Unix(sec, 0)
	age := time.Since(commitTime)
	if age < 0 {
		return 0, nil
	}
	return int(age.Hours() / 24), nil
}
