package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
