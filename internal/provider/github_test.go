package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGitHubValidateAuthSendsBearerToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/repos/acme/repo") {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		auth := r.Header.Get("Authorization")
		if auth != "Bearer testtoken" {
			t.Fatalf("expected bearer token, got %q", auth)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"default_branch":"main"}`))
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

	if err := g.ValidateAuth(ctx); err != nil {
		t.Fatalf("ValidateAuth: %v", err)
	}

	branch, err := g.GetDefaultBranch(ctx)
	if err != nil {
		t.Fatalf("GetDefaultBranch: %v", err)
	}
	if branch != "main" {
		t.Fatalf("expected main got %s", branch)
	}
}
