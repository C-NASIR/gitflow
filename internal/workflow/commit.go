package workflow

import (
	"fmt"
	"strings"

	"gitflow/internal/config"
	"gitflow/internal/git"
)

// CommitOptions defines inputs for creating a commit.
type CommitOptions struct {
	RepoPath string

	All bool

	Message string
	Body    string

	Type     string
	Scope    string
	Breaking bool

	Interactive bool
}

// CommitResult reports the created commit message.
type CommitResult struct {
	Message string
}

// Commit creates a commit using workflow settings.
func Commit(cfg *config.Config, opts CommitOptions) (*CommitResult, error) {
	if opts.RepoPath == "" {
		return nil, fmt.Errorf("repo path is required")
	}

	client, err := git.NewClient(opts.RepoPath)
	if err != nil {
		return nil, err
	}

	if opts.All {
		if err := client.AddAll(); err != nil {
			return nil, err
		}
	}

	hasStaged, err := client.HasStagedChanges()
	if err != nil {
		return nil, err
	}
	if !hasStaged {
		return nil, fmt.Errorf("no staged changes to commit")
	}

	msg, err := buildCommitMessage(cfg, opts)
	if err != nil {
		return nil, err
	}

	if err := client.CommitMessage(msg); err != nil {
		return nil, err
	}

	return &CommitResult{Message: msg}, nil
}

func buildCommitMessage(cfg *config.Config, opts CommitOptions) (string, error) {
	if cfg.Commits.Conventional {
		t := strings.TrimSpace(opts.Type)
		if t == "" {
			return "", fmt.Errorf("type is required for conventional commits")
		}

		summary := strings.TrimSpace(opts.Message)
		if summary == "" {
			return "", fmt.Errorf("summary is required for conventional commits")
		}

		scope := strings.TrimSpace(opts.Scope)
		header := ""
		if scope != "" {
			header = fmt.Sprintf("%s(%s)", t, scope)
		} else {
			header = t
		}

		if opts.Breaking {
			header += "!"
		}

		header = header + ": " + summary

		body := strings.TrimSpace(opts.Body)
		if body == "" {
			return header, nil
		}

		return header + "\n\n" + body, nil
	}

	summary := strings.TrimSpace(opts.Message)
	if summary == "" {
		return "", fmt.Errorf("message is required")
	}

	body := strings.TrimSpace(opts.Body)
	if body == "" {
		return summary, nil
	}
	return summary + "\n\n" + body, nil
}
