package types

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

type Branch struct {
	Name          string
	Current       bool
	LastCommitMsg string
	Author        string

	AgeDays int

	Ahead  int
	Behind int
}

type Release struct {
	Tag  string
	Name string
	URL  string
}
