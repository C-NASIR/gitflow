package release

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"gitflow/internal/ui"
	"gitflow/internal/workflow"
)

const (
	outputText = "text"
	outputJSON = "json"
	outputEnv  = "env"
)

type previewOutput struct {
	CurrentVersion string `json:"current_version"`
	NextVersion    string `json:"next_version"`
	CommitCount    int    `json:"commit_count"`
	Changelog      string `json:"changelog"`
}

type publishOutput struct {
	Provider string `json:"provider"`
	Version  string `json:"version"`
	URL      string `json:"url"`
	DryRun   bool   `json:"dry_run"`
}

type versionOutput struct {
	Version string `json:"version"`
}

func parseOutputFormat(jsonFlag bool, envFlag bool) (string, error) {
	if jsonFlag && envFlag {
		return "", fmt.Errorf("choose only one of --json or --env")
	}
	if jsonFlag {
		return outputJSON, nil
	}
	if envFlag {
		return outputEnv, nil
	}
	return outputText, nil
}

func outputReleasePreview(u *ui.UI, out io.Writer, format string, result *workflow.ReleaseResult) error {
	switch format {
	case outputJSON:
		payload := previewOutput{
			CurrentVersion: result.BaseVersion.String(),
			NextVersion:    result.NextVersion.String(),
			CommitCount:    result.CommitCount,
			Changelog:      result.Changelog,
		}
		return writeJSON(u, payload)
	case outputEnv:
		changelog := escapeEnvValue(result.Changelog)
		lines := []string{
			fmt.Sprintf("GITFLOW_RELEASE_CURRENT_VERSION=%s", result.BaseVersion.String()),
			fmt.Sprintf("GITFLOW_RELEASE_NEXT_VERSION=%s", result.NextVersion.String()),
			fmt.Sprintf("GITFLOW_RELEASE_COMMIT_COUNT=%d", result.CommitCount),
			fmt.Sprintf("GITFLOW_RELEASE_CHANGELOG=%s", changelog),
		}
		for _, line := range lines {
			u.Line("%s", line)
		}
		return nil
	default:
		u.Header("Release preview")
		t := ui.NewTable(out)
		t.Header("KEY", "VALUE")
		t.KeyValue("Current version", result.BaseVersion.String())
		t.KeyValue("Next version", result.NextVersion.String())
		t.KeyValue("Commit count", result.CommitCount)
		t.Flush()
		u.Line("")
		for _, line := range strings.Split(result.Changelog, "\n") {
			u.Line("%s", line)
		}
		return nil
	}
}

func outputReleaseVersion(outWriter *ui.UI, format string, version string) error {
	switch format {
	case outputJSON:
		return writeJSON(outWriter, versionOutput{Version: version})
	case outputEnv:
		outWriter.Line("GITFLOW_RELEASE_VERSION=%s", version)
		return nil
	default:
		outWriter.Line("%s", version)
		return nil
	}
}

func outputReleasePublish(u *ui.UI, out io.Writer, format string, result *workflow.ReleaseResult, publishResult *workflow.ReleasePublishResult) error {
	switch format {
	case outputJSON:
		payload := publishOutput{
			Provider: publishResult.Provider,
			Version:  result.NextVersion.String(),
			URL:      publishResult.URL,
			DryRun:   publishResult.DryRun,
		}
		return writeJSON(u, payload)
	case outputEnv:
		lines := []string{
			fmt.Sprintf("GITFLOW_RELEASE_PROVIDER=%s", publishResult.Provider),
			fmt.Sprintf("GITFLOW_RELEASE_VERSION=%s", result.NextVersion.String()),
			fmt.Sprintf("GITFLOW_RELEASE_URL=%s", publishResult.URL),
			fmt.Sprintf("GITFLOW_RELEASE_DRY_RUN=%t", publishResult.DryRun),
		}
		for _, line := range lines {
			u.Line("%s", line)
		}
		return nil
	default:
		u.Header("Release published")
		t := ui.NewTable(out)
		t.Header("KEY", "VALUE")
		t.KeyValue("Provider", publishResult.Provider)
		t.KeyValue("Version", result.NextVersion.String())
		t.KeyValue("URL", publishResult.URL)
		t.Flush()
		if publishResult.DryRun {
			u.Warn("Dry run: no release published")
		}
		return nil
	}
}

func writeJSON(u *ui.UI, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	u.Line("%s", string(data))
	return nil
}

func escapeEnvValue(value string) string {
	return strings.ReplaceAll(value, "\n", "\\n")
}
