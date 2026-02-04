package workflow

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"gitflow/internal/config"
)

func TestReleasePublishDryRunSkipsProvider(t *testing.T) {
	repo := setupReleaseRepo(t)
	defer os.RemoveAll(repo)

	count := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	writeProviderConfig(t, repo, server.URL)

	rel, err := Release(ReleaseOptions{RepoPath: repo, DryRun: true})
	if err != nil {
		t.Fatalf("Release: %v", err)
	}

	_, err = ReleasePublish(ReleasePublishOptions{
		RepoPath: repo,
		DryRun:   true,
		Result:   rel,
	})
	if err != nil {
		t.Fatalf("ReleasePublish: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected no provider calls, got %d", count)
	}
}

func TestReleasePublishProviderError(t *testing.T) {
	repo := setupReleaseRepo(t)
	defer os.RemoveAll(repo)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("boom"))
	}))
	defer server.Close()

	writeProviderConfig(t, repo, server.URL)

	rel, err := Release(ReleaseOptions{RepoPath: repo, DryRun: true})
	if err != nil {
		t.Fatalf("Release: %v", err)
	}

	_, err = ReleasePublish(ReleasePublishOptions{
		RepoPath: repo,
		DryRun:   false,
		Result:   rel,
	})
	if err == nil {
		t.Fatalf("expected provider error")
	}
	var providerErr ProviderError
	if !errors.As(err, &providerErr) {
		t.Fatalf("expected provider error, got %v", err)
	}
}

func writeProviderConfig(t *testing.T, repo string, baseURL string) {
	t.Helper()
	if err := os.Setenv("GITFLOW_TOKEN", "token"); err != nil {
		t.Fatalf("set env: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Unsetenv("GITFLOW_TOKEN")
	})

	cfg := config.Default()
	cfg.Provider.Type = "github"
	cfg.Provider.BaseURL = baseURL
	cfg.Provider.TokenEnv = "GITFLOW_TOKEN"
	cfg.Provider.Owner = "octo"
	cfg.Provider.Repo = "repo"

	if err := config.WriteFile(filepath.Join(repo, ".gitflow.yml"), cfg); err != nil {
		t.Fatalf("write config: %v", err)
	}
}
