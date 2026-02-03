package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGitHubListPRsAndGetPR(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/repos/acme/repo/pulls":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{
					"number": 1,
					"title": "First",
					"state": "open",
					"html_url": "https://example/pr/1",
					"draft": false,
					"user": {"login":"alice"},
					"head": {"ref":"feature/a"},
					"base": {"ref":"main"}
				}
			]`))
			return

		case "/repos/acme/repo/pulls/1":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"number": 1,
				"title": "First",
				"body": "Details",
				"state": "open",
				"html_url": "https://example/pr/1",
				"draft": false,
				"user": {"login":"alice"},
				"head": {"ref":"feature/a"},
				"base": {"ref":"main"}
			}`))
			return
		}

		t.Fatalf("unexpected path %s", r.URL.Path)
	}))
	defer server.Close()

	g, err := NewGitHub(ProviderConfig{
		Type:    "github",
		BaseURL: server.URL,
		Token:   "t",
		Owner:   "acme",
		Repo:    "repo",
	})
	if err != nil {
		t.Fatalf("NewGitHub: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	prs, err := g.ListPRs(ctx, "open")
	if err != nil {
		t.Fatalf("ListPRs: %v", err)
	}
	if len(prs) != 1 {
		t.Fatalf("expected 1 pr")
	}
	if prs[0].Number != 1 || prs[0].Author != "alice" {
		t.Fatalf("unexpected pr")
	}

	pr, err := g.GetPR(ctx, 1)
	if err != nil {
		t.Fatalf("GetPR: %v", err)
	}
	if pr.Description != "Details" {
		t.Fatalf("expected details")
	}
}
