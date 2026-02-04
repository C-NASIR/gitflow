// Package types defines shared API models for gitflow.
package types

// PullRequest describes a hosted pull request.
type PullRequest struct {
	Number      int
	Title       string
	Description string
	State       string
	Author      string
	HeadBranch  string
	BaseBranch  string
	URL         string
	Draft       bool

	Reviewers []string
	Labels    []string
}

// Branch represents a local branch summary.
type Branch struct {
	Name          string
	Current       bool
	LastCommitMsg string
	Author        string

	AgeDays int

	Ahead  int
	Behind int
}

// Release represents a published release.
type Release struct {
	Tag  string
	Name string
	URL  string
}
