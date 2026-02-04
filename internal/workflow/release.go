package workflow

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"gitflow/internal/config"
	"gitflow/internal/git"
)

type SemanticVersion struct {
	Major int
	Minor int
	Patch int
}

func (v SemanticVersion) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

type ReleaseOptions struct {
	RepoPath        string
	DryRun          bool
	VersionOverride *SemanticVersion
}

type ReleaseResult struct {
	BaseVersion SemanticVersion
	NextVersion SemanticVersion
	CommitCount int
	Changelog   string
	Tag         string
}

type CommitGroups struct {
	Breaking []string
	Features []string
	Fixes    []string
	Other    []string
}

func Release(opts ReleaseOptions) (*ReleaseResult, error) {
	if opts.RepoPath == "" {
		return nil, fmt.Errorf("repo path is required")
	}

	client, err := git.NewClient(opts.RepoPath)
	if err != nil {
		return nil, err
	}

	if !opts.DryRun {
		dirty, err := client.IsDirty()
		if err != nil {
			return nil, err
		}
		if dirty {
			return nil, fmt.Errorf("working tree is dirty")
		}
	}

	cfgResult, err := config.LoadFromDir(opts.RepoPath)
	if err != nil {
		return nil, ConfigError{Err: err}
	}

	baseVersion, baseTag, err := latestVersion(client, cfgResult.Config.Release.TagPrefix)
	if err != nil {
		return nil, err
	}

	commits, err := client.CommitsBetween(baseTag, "HEAD")
	if err != nil {
		return nil, err
	}

	groups := classifyCommits(commits)
	nextVersion := baseVersion
	if opts.VersionOverride != nil {
		nextVersion = *opts.VersionOverride
	} else {
		nextVersion = bumpVersion(baseVersion, groups, cfgResult.Config.Release.DefaultBump)
	}

	releaseDate := latestCommitDate(commits)
	if releaseDate == "" {
		releaseDate = time.Now().Format("2006-01-02")
	}

	changelog := renderChangelog(nextVersion, cfgResult.Config.Release.TagPrefix, releaseDate, groups, cfgResult.Config.Release.ChangelogSections)
	return &ReleaseResult{
		BaseVersion: baseVersion,
		NextVersion: nextVersion,
		CommitCount: len(commits),
		Changelog:   changelog,
		Tag:         cfgResult.Config.Release.TagPrefix + nextVersion.String(),
	}, nil
}

func latestVersion(client *git.Client, prefix string) (SemanticVersion, string, error) {
	if prefix == "" {
		prefix = "v"
	}
	tags, err := client.ListTags()
	if err != nil {
		return SemanticVersion{}, "", err
	}
	var versions []SemanticVersion
	for _, tag := range tags {
		version, ok := parseVersionTag(tag, prefix)
		if !ok {
			continue
		}
		versions = append(versions, version)
	}
	if len(versions) == 0 {
		return SemanticVersion{}, "", nil
	}
	sort.Slice(versions, func(i, j int) bool {
		return compareVersion(versions[i], versions[j]) < 0
	})
	latest := versions[len(versions)-1]
	return latest, prefix + latest.String(), nil
}

func parseVersionTag(tag string, prefix string) (SemanticVersion, bool) {
	if !strings.HasPrefix(tag, prefix) {
		return SemanticVersion{}, false
	}
	ver := strings.TrimPrefix(tag, prefix)
	parts := strings.Split(ver, ".")
	if len(parts) != 3 {
		return SemanticVersion{}, false
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return SemanticVersion{}, false
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return SemanticVersion{}, false
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return SemanticVersion{}, false
	}
	return SemanticVersion{Major: major, Minor: minor, Patch: patch}, true
}

func compareVersion(a, b SemanticVersion) int {
	if a.Major != b.Major {
		return a.Major - b.Major
	}
	if a.Minor != b.Minor {
		return a.Minor - b.Minor
	}
	return a.Patch - b.Patch
}

func classifyCommits(commits []git.Commit) CommitGroups {
	groups := CommitGroups{}
	for _, commit := range commits {
		commitType, breaking := parseCommitType(commit.Subject)
		if strings.Contains(commit.Body, "BREAKING CHANGE") {
			breaking = true
		}

		if commitType == "" && !breaking {
			continue
		}

		if breaking {
			groups.Breaking = append(groups.Breaking, commit.Subject)
			continue
		}

		switch commitType {
		case "feat":
			groups.Features = append(groups.Features, commit.Subject)
		case "fix", "perf":
			groups.Fixes = append(groups.Fixes, commit.Subject)
		case "refactor", "docs", "test", "chore":
			groups.Other = append(groups.Other, commit.Subject)
		}
	}
	return groups
}

func parseCommitType(subject string) (string, bool) {
	parts := strings.SplitN(subject, ":", 2)
	if len(parts) == 0 {
		return "", false
	}
	prefix := strings.TrimSpace(parts[0])
	if prefix == "" {
		return "", false
	}
	breaking := strings.Contains(prefix, "!")
	prefix = strings.TrimSuffix(prefix, "!")
	if idx := strings.Index(prefix, "("); idx >= 0 {
		prefix = prefix[:idx]
	}
	prefix = strings.TrimSpace(prefix)
	switch prefix {
	case "feat", "fix", "perf", "refactor", "docs", "test", "chore":
		return prefix, breaking
	default:
		return "", breaking
	}
}

func bumpVersion(base SemanticVersion, groups CommitGroups, defaultBump string) SemanticVersion {
	if len(groups.Breaking) > 0 {
		return SemanticVersion{Major: base.Major + 1, Minor: 0, Patch: 0}
	}
	if len(groups.Features) > 0 {
		return SemanticVersion{Major: base.Major, Minor: base.Minor + 1, Patch: 0}
	}
	if len(groups.Fixes) > 0 {
		return SemanticVersion{Major: base.Major, Minor: base.Minor, Patch: base.Patch + 1}
	}
	return applyDefaultBump(base, defaultBump)
}

func applyDefaultBump(base SemanticVersion, defaultBump string) SemanticVersion {
	switch defaultBump {
	case "major":
		return SemanticVersion{Major: base.Major + 1, Minor: 0, Patch: 0}
	case "minor":
		return SemanticVersion{Major: base.Major, Minor: base.Minor + 1, Patch: 0}
	default:
		return SemanticVersion{Major: base.Major, Minor: base.Minor, Patch: base.Patch + 1}
	}
}

func renderChangelog(version SemanticVersion, prefix, date string, groups CommitGroups, sectionOrder []string) string {
	if prefix == "" {
		prefix = "v"
	}
	if len(sectionOrder) == 0 {
		sectionOrder = []string{"breaking", "features", "fixes", "other"}
	}

	sections := map[string]struct {
		title   string
		entries []string
	}{
		"breaking": {title: "Breaking Changes", entries: groups.Breaking},
		"features": {title: "Features", entries: groups.Features},
		"fixes":    {title: "Fixes", entries: groups.Fixes},
		"other":    {title: "Other", entries: groups.Other},
	}

	var b strings.Builder
	b.WriteString("## ")
	b.WriteString(prefix)
	b.WriteString(version.String())
	if date != "" {
		b.WriteString(" - ")
		b.WriteString(date)
	}
	b.WriteString("\n")

	for _, key := range sectionOrder {
		section, ok := sections[key]
		if !ok || len(section.entries) == 0 {
			continue
		}
		b.WriteString("\n### ")
		b.WriteString(section.title)
		b.WriteString("\n")
		for _, entry := range section.entries {
			b.WriteString("- ")
			b.WriteString(entry)
			b.WriteString("\n")
		}
	}

	return strings.TrimSpace(b.String())
}

func latestCommitDate(commits []git.Commit) string {
	if len(commits) == 0 {
		return ""
	}
	return commits[0].Date
}
