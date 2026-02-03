// Package config loads and validates gitflow configuration.
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the top-level gitflow configuration.
type Config struct {
	Provider  ProviderConfig `yaml:"provider"`
	Branches  BranchConfig   `yaml:"branches"`
	Workflows WorkflowConfig `yaml:"workflows"`
	Commits   CommitConfig   `yaml:"commits"`
	UI        UIConfig       `yaml:"ui"`
}

// ProviderConfig controls optional hosting provider integration.
type ProviderConfig struct {
	Type     string `yaml:"type"`
	BaseURL  string `yaml:"base_url"`
	TokenEnv string `yaml:"token_env"`
	Owner    string `yaml:"owner"`
	Repo     string `yaml:"repo"`
}

// BranchConfig contains naming conventions for branches.
type BranchConfig struct {
	FeaturePrefix string `yaml:"feature_prefix"`
	BugfixPrefix  string `yaml:"bugfix_prefix"`
	HotfixPrefix  string `yaml:"hotfix_prefix"`
	MainBranch    string `yaml:"main_branch"`
	DevelopBranch string `yaml:"develop_branch"`
}

// WorkflowConfig groups workflow-specific settings.
type WorkflowConfig struct {
	Start   StartConfig   `yaml:"start"`
	PR      PRConfig      `yaml:"pr"`
	Sync    SyncConfig    `yaml:"sync"`
	Cleanup CleanupConfig `yaml:"cleanup"`
}

// StartConfig governs the start workflow behavior.
type StartConfig struct {
	BaseBranch string `yaml:"base_branch"`
	AutoPush   bool   `yaml:"auto_push"`
	FetchFirst bool   `yaml:"fetch_first"`
}

type PRConfig struct {
	Draft            bool     `yaml:"draft"`
	DefaultReviewers []string `yaml:"default_reviewers"`
	Labels           []string `yaml:"labels"`
}

// SyncConfig governs syncing behavior.
type SyncConfig struct {
	Strategy  string `yaml:"strategy"`
	AutoPush  bool   `yaml:"auto_push"`
	ForcePush bool   `yaml:"force_push"`
}

// CleanupConfig governs cleanup workflow behavior.
type CleanupConfig struct {
	MergedOnly        bool     `yaml:"merged_only"`
	AgeThresholdDays  int      `yaml:"age_threshold_days"`
	ProtectedBranches []string `yaml:"protected_branches"`
}

// UIConfig controls CLI output styling.
type UIConfig struct {
	Color   bool `yaml:"color"`
	Emoji   bool `yaml:"emoji"`
	Verbose bool `yaml:"verbose"`
}

type CommitConfig struct {
	Conventional bool     `yaml:"conventional"`
	Types        []string `yaml:"types"`
	Scopes       []string `yaml:"scopes"`
	RequireScope bool     `yaml:"require_scope"`
}

// LoadResult captures the config and its source path.
type LoadResult struct {
	Path   string
	Config *Config
}

// Load searches for a config starting from the working directory.
func Load() (*LoadResult, error) {
	startDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}

	return LoadFromDir(startDir)
}

// LoadFromDir searches for a config starting from startDir.
func LoadFromDir(startDir string) (*LoadResult, error) {
	path, err := findConfig(startDir)
	if err != nil {
		cfg := Default()
		return &LoadResult{Path: "", Config: cfg}, nil
	}

	cfg, err := loadFromPath(path)
	if err != nil {
		return nil, err
	}

	return &LoadResult{Path: path, Config: cfg}, nil
}

func loadFromPath(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config at %s: %w", path, err)
	}

	cfg := Default()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse yaml at %s: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate fills defaults and checks for invalid settings.
func (c *Config) Validate() error {
	if c.Branches.MainBranch == "" {
		c.Branches.MainBranch = "main"
	}
	if c.Branches.FeaturePrefix == "" {
		c.Branches.FeaturePrefix = "feature/"
	}
	if c.Branches.BugfixPrefix == "" {
		c.Branches.BugfixPrefix = "bugfix/"
	}
	if c.Branches.HotfixPrefix == "" {
		c.Branches.HotfixPrefix = "hotfix/"
	}
	if c.Workflows.Start.BaseBranch == "" {
		c.Workflows.Start.BaseBranch = c.Branches.MainBranch
	}
	if c.Workflows.Sync.Strategy == "" {
		c.Workflows.Sync.Strategy = "rebase"
	}
	if c.Workflows.Cleanup.AgeThresholdDays == 0 {
		c.Workflows.Cleanup.AgeThresholdDays = 30
	}
	if len(c.Workflows.Cleanup.ProtectedBranches) == 0 {
		c.Workflows.Cleanup.ProtectedBranches = []string{"main", "master", "develop"}
	}

	if c.Provider.Type != "" && c.Provider.Type != "github" && c.Provider.Type != "gitlab" {
		return fmt.Errorf("unsupported provider type: %s", c.Provider.Type)
	}

	if c.Provider.Type != "" && c.Provider.TokenEnv != "" {
		if os.Getenv(c.Provider.TokenEnv) == "" {
			return fmt.Errorf("provider token not found in env var: %s", c.Provider.TokenEnv)
		}
	}

	if c.Workflows.Sync.Strategy != "rebase" && c.Workflows.Sync.Strategy != "merge" {
		return fmt.Errorf("unsupported sync strategy: %s", c.Workflows.Sync.Strategy)
	}
	return nil
}
