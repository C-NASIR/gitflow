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
