package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gitflow/pkg/types"
	"io"
	"net/http"
	"strings"
	"time"
)

// GitHub implements provider access for the GitHub API.
type GitHub struct {
	baseURL string
	token   string
	owner   string
	repo    string
	client  *http.Client
}

// NewGitHub builds a GitHub provider with the supplied configuration.
func NewGitHub(cfg ProviderConfig) (*GitHub, error) {
	baseURL := cfg.BaseURL
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "https://api.github.com"
	}

	return &GitHub{
		baseURL: baseURL,
		token:   cfg.Token,
		owner:   cfg.Owner,
		repo:    cfg.Repo,
		client: &http.Client{
			Timeout: 20 * time.Second,
		},
	}, nil
}

// ValidateAuth verifies the token with a basic API request.
func (g *GitHub) ValidateAuth(ctx context.Context) error {
	_, err := g.do(ctx, http.MethodGet, "", nil, nil)
	return err
}

// GetDefaultBranch fetches the repo's default branch.
func (g *GitHub) GetDefaultBranch(ctx context.Context) (string, error) {
	var repo struct {
		DefaultBranch string `json:"default_branch"`
	}

	_, err := g.do(ctx, http.MethodGet, "", nil, &repo)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(repo.DefaultBranch) == "" {
		return "", fmt.Errorf("default branch missing in response")
	}
	return repo.DefaultBranch, nil
}

// do executes a GitHub API request and optionally decodes JSON.
func (g *GitHub) do(ctx context.Context, method string, path string, body any, out any) (*http.Response, error) {
	url := fmt.Sprintf("%s/repos/%s/%s%s", strings.TrimRight(g.baseURL, "/"), g.owner, g.repo, path)

	var r io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("json encode: %w", err)
		}
		r = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, r)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+g.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("github api error: %s (failed to read body)", resp.Status)
		}
		msg := strings.TrimSpace(string(b))
		if msg == "" {
			msg = resp.Status
		}
		return nil, fmt.Errorf("github api error: %s", msg)
	}

	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return nil, fmt.Errorf("json decode: %w", err)
		}
	}

	return resp, nil
}

func (g *GitHub) CreatePR(ctx context.Context, opts CreatePROptions) (*types.PullRequest, error) {
	reqBody := map[string]any{
		"title": opts.Title,
		"head":  opts.HeadBranch,
		"base":  opts.BaseBranch,
		"body":  opts.Description,
		"draft": opts.Draft,
	}

	var respBody struct {
		Number  int    `json:"number"`
		Title   string `json:"title"`
		Body    string `json:"body"`
		State   string `json:"state"`
		HTMLURL string `json:"html_url"`
		Draft   bool   `json:"draft"`
		Head    struct {
			Ref string `json:"ref"`
		} `json:"head"`
		Base struct {
			Ref string `json:"ref"`
		} `json:"base"`
	}

	_, err := g.do(ctx, http.MethodPost, "/pulls", reqBody, &respBody)
	if err != nil {
		return nil, err
	}

	pr := &types.PullRequest{
		Number:      respBody.Number,
		Title:       respBody.Title,
		Description: respBody.Body,
		State:       respBody.State,
		HeadBranch:  respBody.Head.Ref,
		BaseBranch:  respBody.Base.Ref,
		URL:         respBody.HTMLURL,
		Draft:       respBody.Draft,
		Reviewers:   opts.Reviewers,
		Labels:      opts.Labels,
	}

	if len(opts.Reviewers) > 0 {
		reviewReq := map[string]any{
			"reviewers": opts.Reviewers,
		}
		_, err := g.do(ctx, http.MethodPost, fmt.Sprintf("/pulls/%d/requested_reviewers", pr.Number), reviewReq, nil)
		if err != nil {
			return nil, err
		}
	}

	if len(opts.Labels) > 0 {
		labelReq := map[string]any{
			"labels": opts.Labels,
		}
		_, err := g.do(ctx, http.MethodPost, fmt.Sprintf("/issues/%d/labels", pr.Number), labelReq, nil)
		if err != nil {
			return nil, err
		}
	}

	return pr, nil
}

func (g *GitHub) GetPR(ctx context.Context, number int) (*types.PullRequest, error) {
	var gh struct {
		Number  int    `json:"number"`
		Title   string `json:"title"`
		Body    string `json:"body"`
		State   string `json:"state"`
		HTMLURL string `json:"html_url"`
		Draft   bool   `json:"draft"`
		User    struct {
			Login string `json:"login"`
		} `json:"user"`
		Head struct {
			Ref string `json:"ref"`
		} `json:"head"`
		Base struct {
			Ref string `json:"ref"`
		} `json:"base"`
	}

	_, err := g.do(ctx, http.MethodGet, fmt.Sprintf("/pulls/%d", number), nil, &gh)
	if err != nil {
		return nil, err
	}

	return &types.PullRequest{
		Number:      gh.Number,
		Title:       gh.Title,
		Description: gh.Body,
		State:       gh.State,
		Author:      gh.User.Login,
		HeadBranch:  gh.Head.Ref,
		BaseBranch:  gh.Base.Ref,
		URL:         gh.HTMLURL,
		Draft:       gh.Draft,
	}, nil
}

func (g *GitHub) ListPRs(ctx context.Context, state string) ([]*types.PullRequest, error) {
	if state == "" {
		state = "open"
	}

	path := fmt.Sprintf("/pulls?state=%s", state)

	var gh []struct {
		Number  int    `json:"number"`
		Title   string `json:"title"`
		State   string `json:"state"`
		HTMLURL string `json:"html_url"`
		Draft   bool   `json:"draft"`
		User    struct {
			Login string `json:"login"`
		} `json:"user"`
		Head struct {
			Ref string `json:"ref"`
		} `json:"head"`
		Base struct {
			Ref string `json:"ref"`
		} `json:"base"`
	}

	_, err := g.do(ctx, http.MethodGet, path, nil, &gh)
	if err != nil {
		return nil, err
	}

	out := make([]*types.PullRequest, 0, len(gh))
	for _, pr := range gh {
		out = append(out, &types.PullRequest{
			Number:     pr.Number,
			Title:      pr.Title,
			State:      pr.State,
			Author:     pr.User.Login,
			HeadBranch: pr.Head.Ref,
			BaseBranch: pr.Base.Ref,
			URL:        pr.HTMLURL,
			Draft:      pr.Draft,
		})
	}

	return out, nil
}

func (g *GitHub) CreateRelease(tag string, name string, body string) (*types.Release, error) {
	reqBody := map[string]any{
		"tag_name": tag,
		"name":     name,
		"body":     body,
	}

	var respBody struct {
		TagName string `json:"tag_name"`
		Name    string `json:"name"`
		HTMLURL string `json:"html_url"`
	}

	_, err := g.do(context.Background(), http.MethodPost, "/releases", reqBody, &respBody)
	if err != nil {
		if strings.Contains(err.Error(), "already_exists") || strings.Contains(err.Error(), "already exists") {
			return nil, ErrReleaseExists
		}
		return nil, err
	}

	return &types.Release{
		Tag:  respBody.TagName,
		Name: respBody.Name,
		URL:  respBody.HTMLURL,
	}, nil
}

func (g *GitHub) UpdateRelease(tag string, name string, body string) (*types.Release, error) {
	var existing struct {
		ID      int    `json:"id"`
		Tag     string `json:"tag_name"`
		Name    string `json:"name"`
		HTMLURL string `json:"html_url"`
	}

	_, err := g.do(context.Background(), http.MethodGet, fmt.Sprintf("/releases/tags/%s", tag), nil, &existing)
	if err != nil {
		return nil, err
	}
	if existing.ID == 0 {
		return nil, fmt.Errorf("release not found for tag %s", tag)
	}

	reqBody := map[string]any{
		"tag_name": tag,
		"name":     name,
		"body":     body,
	}

	var respBody struct {
		TagName string `json:"tag_name"`
		Name    string `json:"name"`
		HTMLURL string `json:"html_url"`
	}

	_, err = g.do(context.Background(), http.MethodPatch, fmt.Sprintf("/releases/%d", existing.ID), reqBody, &respBody)
	if err != nil {
		return nil, err
	}

	return &types.Release{
		Tag:  respBody.TagName,
		Name: respBody.Name,
		URL:  respBody.HTMLURL,
	}, nil
}
