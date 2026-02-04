package workflow

import (
	"fmt"
	"os"
	"path/filepath"

	"gitflow/internal/config"
	"gitflow/internal/git"
)

const (
	// DoctorOK indicates a successful check.
	DoctorOK    = "OK"
	// DoctorWarn indicates a warning check.
	DoctorWarn  = "WARN"
	// DoctorError indicates a failed check.
	DoctorError = "ERROR"
)

// DoctorCheck describes a single doctor check result.
type DoctorCheck struct {
	Name    string
	Level   string
	Message string
}

// DoctorResult aggregates doctor checks.
type DoctorResult struct {
	Checks []DoctorCheck
}

// Doctor inspects repository health and config readiness.
func Doctor(repoPath string) (*DoctorResult, error) {
	if repoPath == "" {
		return nil, fmt.Errorf("repo path is required")
	}

	result := &DoctorResult{}
	client, err := git.NewClient(repoPath)
	if err != nil {
		result.Checks = append(result.Checks, DoctorCheck{
			Name:    "Git repository",
			Level:   DoctorError,
			Message: err.Error(),
		})
		return result, nil
	}

	result.Checks = append(result.Checks, DoctorCheck{
		Name:    "Git repository",
		Level:   DoctorOK,
		Message: "Repository detected",
	})

	dirty, err := client.IsDirty()
	if err != nil {
		result.Checks = append(result.Checks, DoctorCheck{
			Name:    "Working tree",
			Level:   DoctorError,
			Message: err.Error(),
		})
	} else if dirty {
		result.Checks = append(result.Checks, DoctorCheck{
			Name:    "Working tree",
			Level:   DoctorWarn,
			Message: "Working tree has uncommitted changes",
		})
	} else {
		result.Checks = append(result.Checks, DoctorCheck{
			Name:    "Working tree",
			Level:   DoctorOK,
			Message: "Working tree is clean",
		})
	}

	root, err := client.Run("rev-parse", "--show-toplevel")
	if err != nil {
		root = repoPath
	}

	configPath, configExists := repoConfigPath(root)
	if configExists {
		result.Checks = append(result.Checks, DoctorCheck{
			Name:    "Config presence",
			Level:   DoctorOK,
			Message: fmt.Sprintf("Found %s", configPath),
		})
	} else {
		result.Checks = append(result.Checks, DoctorCheck{
			Name:    "Config presence",
			Level:   DoctorWarn,
			Message: "Config file not found; run gitflow init",
		})
	}

	res, loadErr := config.LoadFromDir(root)
	if loadErr != nil {
		result.Checks = append(result.Checks, DoctorCheck{
			Name:    "Config validity",
			Level:   DoctorError,
			Message: loadErr.Error(),
		})
	} else if configExists {
		if err := config.ValidateStrict(res.Config); err != nil {
			result.Checks = append(result.Checks, DoctorCheck{
				Name:    "Config validity",
				Level:   DoctorError,
				Message: err.Error(),
			})
		} else {
			result.Checks = append(result.Checks, DoctorCheck{
				Name:    "Config validity",
				Level:   DoctorOK,
				Message: "Config valid",
			})
		}
	} else {
		result.Checks = append(result.Checks, DoctorCheck{
			Name:    "Config validity",
			Level:   DoctorWarn,
			Message: "Config missing; skipping validation",
		})
	}

	if loadErr != nil {
		result.Checks = append(result.Checks, DoctorCheck{
			Name:    "Provider token",
			Level:   DoctorWarn,
			Message: "Skipping provider token check due to invalid config",
		})
		return result, nil
	}

	provider := res.Config.Provider
	if provider.Type == "" {
		result.Checks = append(result.Checks, DoctorCheck{
			Name:    "Provider token",
			Level:   DoctorOK,
			Message: "Provider not configured",
		})
		return result, nil
	}

	if provider.TokenEnv == "" {
		result.Checks = append(result.Checks, DoctorCheck{
			Name:    "Provider token",
			Level:   DoctorWarn,
			Message: "provider.token_env is not set",
		})
		return result, nil
	}

	if os.Getenv(provider.TokenEnv) == "" {
		result.Checks = append(result.Checks, DoctorCheck{
			Name:    "Provider token",
			Level:   DoctorWarn,
			Message: fmt.Sprintf("Missing %s environment variable", provider.TokenEnv),
		})
		return result, nil
	}

	result.Checks = append(result.Checks, DoctorCheck{
		Name:    "Provider token",
		Level:   DoctorOK,
		Message: fmt.Sprintf("%s is set", provider.TokenEnv),
	})

	return result, nil
}

func repoConfigPath(root string) (string, bool) {
	if root == "" {
		return "", false
	}
	if p := filepath.Join(root, ".gitflow.yml"); fileExists(p) {
		return p, true
	}
	if p := filepath.Join(root, ".gitflow.yaml"); fileExists(p) {
		return p, true
	}
	return "", false
}

func fileExists(path string) bool {
	st, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !st.IsDir()
}
