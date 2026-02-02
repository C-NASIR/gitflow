package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGitHubCreatePR(t *testing.T) {
	var sawCreate bool
	var sawReviewers bool
	var sawLabels bool

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer testtoken" {
			t.Fatalf("missing auth header")
		}

		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/repos/acme/repo/pulls":
			sawCreate = true
			var body map[string]any
			_ = json.NewDecoder(r.Body).Decode(&body)

			if body["title"] != "Hello" {
				t.Fatalf("unexpected title")
			}
			if body["head"] != "feature/x" {
				t.Fatalf("unexpected head")
			}
			if body["base"] != "main" {
				t.Fatalf("unexpected base")
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{
				"number": 12,
				"title": "Hello",
				"body": "Body",
				"state": "open",
				"html_url": "https://example/pr/12",
				"draft": false,
				"head": {"ref":"feature/x"},
				"base": {"ref":"main"}
			}`))
			return

		case r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/requested_reviewers"):
			sawReviewers = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{}`))
			return

		case r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/labels"):
			sawLabels = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{}`))
			return
		}

		t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
	}))
	defer server.Close()

	g, err := NewGitHub(ProviderConfig{
		Type:    "github",
		BaseURL: server.URL,
		Token:   "testtoken",
		Owner:   "acme",
		Repo:    "repo",
	})
	if err != nil {
		t.Fatalf("NewGitHub: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	pr, err := g.CreatePR(ctx, CreatePROptions{
		Title:       "Hello",
		Description: "Body",
		HeadBranch:  "feature/x",
		BaseBranch:  "main",
		Draft:       false,
		Reviewers:   []string{"alice"},
		Labels:      []string{"needs-review"},
	})
	if err != nil {
		t.Fatalf("CreatePR: %v", err)
	}

	if !sawCreate || !sawReviewers || !sawLabels {
		t.Fatalf("expected create reviewers labels calls, got create=%v reviewers=%v labels=%v", sawCreate, sawReviewers, sawLabels)
	}

	if pr.Number != 12 {
		t.Fatalf("expected 12 got %d", pr.Number)
	}
	if pr.URL == "" {
		t.Fatalf("expected url")
	}
}
