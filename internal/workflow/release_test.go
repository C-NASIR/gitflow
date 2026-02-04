package workflow

import (
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"gitflow/internal/git"
)

func TestParseVersionTagAndSort(t *testing.T) {
	if _, ok := parseVersionTag("v1.2.3", "v"); !ok {
		t.Fatalf("expected tag to parse")
	}
	if _, ok := parseVersionTag("1.2.3", "v"); ok {
		t.Fatalf("expected tag without prefix to fail")
	}

	versions := []SemanticVersion{{Major: 1, Minor: 10, Patch: 0}, {Major: 2, Minor: 0, Patch: 0}, {Major: 1, Minor: 2, Patch: 3}}
	sort.Slice(versions, func(i, j int) bool {
		return compareVersion(versions[i], versions[j]) < 0
	})
	if versions[len(versions)-1].Major != 2 {
		t.Fatalf("expected highest version to be 2.x.x")
	}
}

func TestCommitClassification(t *testing.T) {
	commits := []git.Commit{
		{Subject: "feat: add login"},
		{Subject: "fix: patch bug"},
		{Subject: "docs: update readme"},
		{Subject: "refactor: update api", Body: "BREAKING CHANGE: api"},
		{Subject: "feat!: change auth"},
		{Subject: "merge branch"},
	}

	groups := classifyCommits(commits)
	if len(groups.Breaking) != 2 {
		t.Fatalf("expected breaking commits")
	}
	if len(groups.Features) != 1 {
		t.Fatalf("expected feature commits")
	}
	if len(groups.Fixes) != 1 {
		t.Fatalf("expected fix commits")
	}
	if len(groups.Other) != 1 {
		t.Fatalf("expected other commits")
	}
}

func TestBumpVersion(t *testing.T) {
	base := SemanticVersion{Major: 1, Minor: 2, Patch: 3}
	groups := CommitGroups{Features: []string{"feat"}}
	if next := bumpVersion(base, groups, "patch"); next.Minor != 3 || next.Patch != 0 {
		t.Fatalf("expected minor bump")
	}
	groups = CommitGroups{}
	if next := bumpVersion(base, groups, "patch"); next.Patch != 4 {
		t.Fatalf("expected patch bump")
	}
}

func TestRenderChangelog(t *testing.T) {
	groups := CommitGroups{
		Features: []string{"feat: add"},
		Fixes:    []string{"fix: bug"},
	}
	changelog := renderChangelog(SemanticVersion{Major: 1, Minor: 0, Patch: 0}, "v", "2024-01-01", groups, []string{"features", "fixes"})
	if !strings.Contains(changelog, "## v1.0.0 - 2024-01-01") {
		t.Fatalf("expected changelog header")
	}
	if !strings.Contains(changelog, "### Features") || !strings.Contains(changelog, "### Fixes") {
		t.Fatalf("expected changelog sections")
	}
}

func TestReleaseWorkflowDryRun(t *testing.T) {
	repo := setupReleaseRepo(t)
	defer os.RemoveAll(repo)

	runGitRelease(t, repo, "tag", "v0.1.0")
	writeFile(t, repo, "a.txt", "change")
	runGitRelease(t, repo, "add", "-A")
	runGitRelease(t, repo, "commit", "-m", "feat: add login")

	res, err := Release(ReleaseOptions{RepoPath: repo, DryRun: true})
	if err != nil {
		t.Fatalf("Release: %v", err)
	}
	if res.NextVersion.Minor != 2 || res.NextVersion.Patch != 0 {
		t.Fatalf("expected minor bump")
	}
	if res.CommitCount != 1 {
		t.Fatalf("expected commit count 1 got %d", res.CommitCount)
	}
}

func TestReleaseWorkflowDirtyRepo(t *testing.T) {
	repo := setupReleaseRepo(t)
	defer os.RemoveAll(repo)

	writeFile(t, repo, "dirty.txt", "dirty")

	_, err := Release(ReleaseOptions{RepoPath: repo, DryRun: false})
	if err == nil {
		t.Fatalf("expected error for dirty repo")
	}
}

func setupReleaseRepo(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	runGitRelease(t, dir, "init")
	runGitRelease(t, dir, "config", "user.email", "test@example.com")
	runGitRelease(t, dir, "config", "user.name", "Test User")
	writeFile(t, dir, "readme.md", "init")
	runGitRelease(t, dir, "add", "-A")
	runGitRelease(t, dir, "commit", "-m", "chore: init")
	return dir
}

func runGitRelease(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v output: %s", args, err, string(out))
	}
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}
