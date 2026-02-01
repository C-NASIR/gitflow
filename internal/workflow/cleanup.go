package workflow

import (
	"bufio"
	"fmt"
	"gitflow/internal/config"
	"gitflow/internal/git"
	"os"
	"sort"
	"strings"
)

// CleanupOptions defines inputs for pruning local and remote branches.
type CleanupOptions struct {
	RepoPath string
	Remote   string

	Yes bool

	All            bool
	AgeThreshold   int
	DeleteRemote   bool
	MergedOnlyHint *bool
}

// CleanupCandidate captures branch metadata used for cleanup decisions.
type CleanupCandidate struct {
	Name       string
	AgeDays    int
	MergedInto bool
	WillDelete bool
	RemoteAlso bool
	Protected  bool
	Reason     string
}

// CleanupResult summarizes cleanup decisions and actions.
type CleanupResult struct {
	BaseBranch    string
	Current       string
	Candidates    []CleanupCandidate
	Deleted       []string
	RemoteDeleted []string
}

// Cleanup determines stale branches and deletes them if confirmed.
func Cleanup(cfg *config.Config, opts CleanupOptions) (*CleanupResult, error) {
	if opts.RepoPath == "" {
		return nil, fmt.Errorf("repo path is required")
	}
	if opts.Remote == "" {
		opts.Remote = "origin"
	}

	client, err := git.NewClient(opts.RepoPath)
	if err != nil {
		return nil, err
	}

	dirty, err := client.IsDirty()
	if err != nil {
		return nil, err
	}
	if dirty {
		return nil, fmt.Errorf("working tree is not clean")
	}

	current, err := client.CurrentBranch()
	if err != nil {
		return nil, err
	}

	base := cfg.Workflows.Start.BaseBranch
	if base == "" {
		base = cfg.Branches.MainBranch
	}
	if base == "" {
		base = "main"
	}

	remoteExists, err := client.HasRemote(opts.Remote)
	if err != nil {
		return nil, err
	}
	if !remoteExists {
		return nil, fmt.Errorf("remote %s not found", opts.Remote)
	}

	ageThreshold := opts.AgeThreshold
	if ageThreshold == 0 {
		ageThreshold = cfg.Workflows.Cleanup.AgeThresholdDays
	}
	if ageThreshold == 0 {
		ageThreshold = 30
	}

	mergedOnly := cfg.Workflows.Cleanup.MergedOnly
	if opts.MergedOnlyHint != nil {
		mergedOnly = *opts.MergedOnlyHint
	}
	if opts.All {
		mergedOnly = false
	}

	protectedSet := map[string]bool{}
	for _, b := range cfg.Workflows.Cleanup.ProtectedBranches {
		protectedSet[b] = true
	}
	protectedSet[base] = true

	mergedSet := map[string]bool{}
	mergedBranches, err := client.MergedBranches(base)
	if err != nil {
		return nil, err
	}
	for _, b := range mergedBranches {
		mergedSet[b] = true
	}

	localBranches, err := client.ListLocalBranches()
	if err != nil {
		return nil, err
	}

	var candidates []CleanupCandidate
	for _, b := range localBranches {
		c := CleanupCandidate{Name: b}

		if protectedSet[b] {
			c.Protected = true
			c.Reason = "protected"
			candidates = append(candidates, c)
			continue
		}
		if b == current {
			c.Protected = true
			c.Reason = "current"
			candidates = append(candidates, c)
			continue
		}

		ageDays, err := client.BranchAgeDays(b)
		if err != nil {
			c.Reason = "age check failed"
			candidates = append(candidates, c)
			continue
		}

		c.AgeDays = ageDays
		c.MergedInto = mergedSet[b]

		if mergedOnly {
			if c.MergedInto {
				c.WillDelete = true
				c.Reason = "merged"
			} else {
				c.Reason = "not merged"
			}
		} else {
			if c.MergedInto {
				c.WillDelete = true
				c.Reason = "merged"
			} else if c.AgeDays >= ageThreshold {
				c.WillDelete = true
				c.Reason = "stale"
			} else {
				c.Reason = "recent"
			}
		}

		if c.WillDelete && opts.DeleteRemote {
			c.RemoteAlso = true
		}

		candidates = append(candidates, c)
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Name < candidates[j].Name
	})

	var plan []CleanupCandidate
	for _, c := range candidates {
		if c.WillDelete {
			plan = append(plan, c)
		}
	}

	if len(plan) == 0 {
		return &CleanupResult{
			BaseBranch: base,
			Current:    current,
			Candidates: candidates,
		}, nil
	}

	if !opts.Yes {
		ok, err := confirmCleanup(plan, opts.DeleteRemote)
		if err != nil {
			return nil, err
		}
		if !ok {
			return &CleanupResult{
				BaseBranch: base,
				Current:    current,
				Candidates: candidates,
			}, nil
		}
	}

	var deleted []string
	var remoteDeleted []string

	for _, c := range plan {
		if err := client.DeleteBranch(c.Name, false); err != nil {
			return nil, err
		}
		deleted = append(deleted, c.Name)

		if c.RemoteAlso {
			if err := client.DeleteRemoteBranch(opts.Remote, c.Name); err != nil {
				return nil, err
			}
			remoteDeleted = append(remoteDeleted, c.Name)
		}
	}

	return &CleanupResult{
		BaseBranch:    base,
		Current:       current,
		Candidates:    candidates,
		Deleted:       deleted,
		RemoteDeleted: remoteDeleted,
	}, nil

}

func confirmCleanup(plan []CleanupCandidate, deleteRemote bool) (bool, error) {
	fmt.Println()
	fmt.Println("Branches to delete")
	for _, c := range plan {
		line := fmt.Sprintf(" %s reason=%s age=%dd", c.Name, c.Reason, c.AgeDays)
		if deleteRemote {
			line += " remote=yes"
		}
		fmt.Println(line)
	}
	fmt.Println()
	fmt.Print("Proceed, type yes to confirm: ")

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	text = strings.ToLower(strings.TrimSpace(text))
	return text == "yes", nil

}
