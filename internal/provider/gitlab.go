package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gitflow/pkg/types"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// GitLab implements provider access for the GitLab API.
type GitLab struct {
	baseURL string
	token   string
	project string
	client  *http.Client
}

// NewGitLab builds a GitLab provider with the supplied configuration.
func NewGitLab(cfg ProviderConfig) (*GitLab, error) {
	baseURL := cfg.BaseURL
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "https://gitlab.com/api/v4"
	}

	project := strings.TrimSpace(cfg.Owner)
	if cfg.Repo != "" {
		project = fmt.Sprintf("%s/%s", strings.TrimSpace(cfg.Owner), strings.TrimSpace(cfg.Repo))
	}
	if project == "" {
		return nil, fmt.Errorf("project is required")
	}

	return &GitLab{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   cfg.Token,
		project: url.PathEscape(project),
		client: &http.Client{
			Timeout: 20 * time.Second,
		},
	}, nil
}

// ValidateAuth verifies the token with a basic API request.
func (g *GitLab) ValidateAuth(ctx context.Context) error {
	_, err := g.do(ctx, http.MethodGet, "", nil, nil)
	return err
}

// GetDefaultBranch fetches the repo's default branch.
func (g *GitLab) GetDefaultBranch(ctx context.Context) (string, error) {
	var resp struct {
		DefaultBranch string `json:"default_branch"`
	}
	_, err := g.do(ctx, http.MethodGet, "", nil, &resp)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(resp.DefaultBranch) == "" {
		return "", fmt.Errorf("default branch missing in response")
	}
	return resp.DefaultBranch, nil
}

// CreatePR returns an error because GitLab PRs are not implemented.
func (g *GitLab) CreatePR(ctx context.Context, opts CreatePROptions) (*types.PullRequest, error) {
	return nil, fmt.Errorf("gitlab pull requests not implemented")
}

// GetPR returns an error because GitLab PRs are not implemented.
func (g *GitLab) GetPR(ctx context.Context, number int) (*types.PullRequest, error) {
	return nil, fmt.Errorf("gitlab pull requests not implemented")
}

// ListPRs returns an error because GitLab PRs are not implemented.
func (g *GitLab) ListPRs(ctx context.Context, state string) ([]*types.PullRequest, error) {
	return nil, fmt.Errorf("gitlab pull requests not implemented")
}

// CreateRelease creates a GitLab release for a tag.
func (g *GitLab) CreateRelease(tag string, name string, body string) (*types.Release, error) {
	reqBody := map[string]any{
		"tag_name":    tag,
		"name":        name,
		"description": body,
	}

	var respBody struct {
		TagName string `json:"tag_name"`
		Name    string `json:"name"`
		Links   struct {
			Self string `json:"self"`
		} `json:"_links"`
	}

	_, err := g.do(context.Background(), http.MethodPost, "/releases", reqBody, &respBody)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return nil, ErrReleaseExists
		}
		return nil, err
	}

	return &types.Release{
		Tag:  respBody.TagName,
		Name: respBody.Name,
		URL:  respBody.Links.Self,
	}, nil
}

// UpdateRelease updates an existing GitLab release.
func (g *GitLab) UpdateRelease(tag string, name string, body string) (*types.Release, error) {
	reqBody := map[string]any{
		"name":        name,
		"description": body,
	}

	var respBody struct {
		TagName string `json:"tag_name"`
		Name    string `json:"name"`
		Links   struct {
			Self string `json:"self"`
		} `json:"_links"`
	}

	_, err := g.do(context.Background(), http.MethodPut, fmt.Sprintf("/releases/%s", url.PathEscape(tag)), reqBody, &respBody)
	if err != nil {
		return nil, err
	}

	return &types.Release{
		Tag:  respBody.TagName,
		Name: respBody.Name,
		URL:  respBody.Links.Self,
	}, nil
}

func (g *GitLab) do(ctx context.Context, method string, path string, body any, out any) (*http.Response, error) {
	url := fmt.Sprintf("%s/projects/%s%s", g.baseURL, g.project, path)

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

	req.Header.Set("PRIVATE-TOKEN", g.token)
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
			return nil, fmt.Errorf("gitlab api error: %s (failed to read body)", resp.Status)
		}
		msg := strings.TrimSpace(string(b))
		if msg == "" {
			msg = resp.Status
		}
		return nil, fmt.Errorf("gitlab api error: %s", msg)
	}

	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return nil, fmt.Errorf("json decode: %w", err)
		}
	}

	return resp, nil
}
