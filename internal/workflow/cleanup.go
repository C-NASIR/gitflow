package workflow

import (
	"fmt"
	"sort"
	"strings"

	"gitflow/internal/config"
	"gitflow/internal/git"
)

// CleanupOptions defines inputs for cleanup.
type CleanupOptions struct {
	RepoPath     string
	Remote       string
	Yes          bool
	All          bool
	AgeThreshold int
	DeleteRemote bool
	Selected     []string
}

// CleanupResult reports cleanup candidates and deletions.
type CleanupResult struct {
	BaseBranch    string
	Current       string
	Deleted       []string
	RemoteDeleted []string
	Candidates    []CandidateBranch
}

// CandidateBranch describes a branch eligible for cleanup.
type CandidateBranch struct {
	Name    string
	Reason  string
	AgeDays int
	Ahead   int
	Behind  int
}

// Cleanup identifies and optionally deletes stale branches.
func Cleanup(cfg *config.Config, opts CleanupOptions) (*CleanupResult, error) {
	if strings.TrimSpace(opts.RepoPath) == "" {
		return nil, fmt.Errorf("repo path is required")
	}
	if opts.Remote == "" {
		opts.Remote = "origin"
	}

	base := strings.TrimSpace(cfg.Workflows.Start.BaseBranch)
	if base == "" {
		base = strings.TrimSpace(cfg.Branches.MainBranch)
	}
	if base == "" {
		base = "main"
	}

	client, err := git.NewClient(opts.RepoPath)
	if err != nil {
		return nil, err
	}

	current, err := client.CurrentBranch()
	if err != nil {
		return nil, err
	}

	protected := make(map[string]bool)
	protected[base] = true
	protected[current] = true
	for _, p := range cfg.Workflows.Cleanup.ProtectedBranches {
		v := strings.TrimSpace(p)
		if v != "" {
			protected[v] = true
		}
	}

	ageThreshold := opts.AgeThreshold
	if ageThreshold == 0 {
		ageThreshold = cfg.Workflows.Cleanup.AgeThresholdDays
	}
	if ageThreshold < 0 {
		ageThreshold = 0
	}

	var candidates []CandidateBranch

	if cfg.Workflows.Cleanup.MergedOnly && !opts.All {
		merged, err := client.MergedBranches(base)
		if err != nil {
			return nil, err
		}

		for _, b := range merged {
			if protected[b] {
				continue
			}
			age, _ := clientBranchAgeDays(client, b)
			ahead, behind := clientAheadBehind(client, b, base)
			candidates = append(candidates, CandidateBranch{
				Name:    b,
				Reason:  "merged",
				AgeDays: age,
				Ahead:   ahead,
				Behind:  behind,
			})
		}
	} else {
		branches, err := client.ListLocalBranches(base)
		if err != nil {
			return nil, err
		}

		for _, b := range branches {
			if protected[b.Name] {
				continue
			}
			if ageThreshold > 0 && b.AgeDays < ageThreshold {
				continue
			}

			reason := "stale"
			if cfg.Workflows.Cleanup.MergedOnly {
				mergedList, err := client.MergedBranches(base)
				if err == nil {
					for _, mb := range mergedList {
						if mb == b.Name {
							reason = "merged and stale"
							break
						}
					}
				}
			}

			candidates = append(candidates, CandidateBranch{
				Name:    b.Name,
				Reason:  reason,
				AgeDays: b.AgeDays,
				Ahead:   b.Ahead,
				Behind:  b.Behind,
			})
		}
	}

	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Reason != candidates[j].Reason {
			return candidates[i].Reason < candidates[j].Reason
		}
		if candidates[i].AgeDays != candidates[j].AgeDays {
			return candidates[i].AgeDays > candidates[j].AgeDays
		}
		return candidates[i].Name < candidates[j].Name
	})

	if len(candidates) == 0 {
		return &CleanupResult{
			BaseBranch: base,
			Current:    current,
			Candidates: nil,
		}, nil
	}

	var toDelete []string
	if opts.Yes {
		allow := make(map[string]bool)
		if len(opts.Selected) > 0 {
			for _, s := range opts.Selected {
				s = strings.TrimSpace(s)
				if s != "" {
					allow[s] = true
				}
			}
		}

		for _, c := range candidates {
			if len(allow) > 0 && !allow[c.Name] {
				continue
			}
			toDelete = append(toDelete, c.Name)
		}
	} else {
		return &CleanupResult{
			BaseBranch: base,
			Current:    current,
			Candidates: candidates,
		}, nil
	}

	var deleted []string
	for _, b := range toDelete {
		if err := client.DeleteBranch(b, false); err != nil {
			return nil, err
		}
		deleted = append(deleted, b)
	}

	var remoteDeleted []string
	if opts.DeleteRemote {
		for _, b := range deleted {
			if err := client.DeleteRemoteBranch(opts.Remote, b); err != nil {
				return nil, err
			}
			remoteDeleted = append(remoteDeleted, b)
		}
	}

	return &CleanupResult{
		BaseBranch:    base,
		Current:       current,
		Deleted:       deleted,
		RemoteDeleted: remoteDeleted,
		Candidates:    candidates,
	}, nil
}

func clientBranchAgeDays(c *git.Client, branch string) (int, error) {
	out, err := c.Run("log", "-1", "--format=%ct", branch)
	if err != nil {
		return 0, err
	}

	sec, err := parseInt64(strings.TrimSpace(out))
	if err != nil {
		return 0, err
	}
	ageDays := int(secondsSince(sec) / 86400)
	if ageDays < 0 {
		ageDays = 0
	}
	return ageDays, nil
}

func clientAheadBehind(c *git.Client, branch string, base string) (int, int) {
	a, b, err := cAheadBehind(c, branch, base)
	if err != nil {
		return 0, 0
	}
	return a, b
}

func cAheadBehind(c *git.Client, branch string, base string) (int, int, error) {
	out, err := c.Run("rev-list", "--left-right", "--count", fmt.Sprintf("%s...%s", base, branch))
	if err != nil {
		return 0, 0, err
	}
	fields := strings.Fields(strings.TrimSpace(out))
	if len(fields) != 2 {
		return 0, 0, fmt.Errorf("unexpected rev-list output: %q", out)
	}
	behind, err := parseInt(fields[0])
	if err != nil {
		return 0, 0, err
	}
	ahead, err := parseInt(fields[1])
	if err != nil {
		return 0, 0, err
	}
	return ahead, behind, nil
}
