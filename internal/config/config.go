package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Provider  ProviderConfig `yaml:"provider"`
	Branches  BranchConfig   `yaml:"branches"`
	Workflows WorkflowConfig `yaml:"workflows"`
	UI        UIConfig       `yaml:"ui"`
}

type ProviderConfig struct {
	Type     string `yaml:"type"`
	BaseURL  string `yaml:"base_url"`
	TokenEnv string `yaml:"token_env"`
	Owner    string `yaml:"owner"`
	Repo     string `yaml:"repo"`
}

type BranchConfig struct {
	FeaturePrefix string `yaml:"feature_prefix"`
	BugfixPrefix  string `yaml:"bugfix_prefix"`
	HotfixPrefix  string `yaml:"hotfix_prefix"`
	MainBranch    string `yaml:"main_branch"`
	DevelopBranch string `yaml:"develop_branch"`
}

type WorkflowConfig struct {
	Start   StartConfig   `yaml:"start"`
	Sync    SyncConfig    `yaml:"sync"`
	Cleanup CleanupConfig `yaml:"cleanup"`
}

type StartConfig struct {
	BaseBranch string `yaml:"base_branch"`
	AutoPush   bool   `yaml:"auto_push"`
	FetchFirst bool   `yaml:"fetch_first"`
}

type SyncConfig struct {
	Strategy  string `yaml:"strategy"`
	AutoPush  bool   `yaml:"auto_push"`
	ForcePush bool   `yaml:"force_push"`
}

type CleanupConfig struct {
	MergedOnly        bool     `yaml:"merged_only"`
	AgeThresholdDays  int      `yaml:"age_threshold_days"`
	ProtectedBranches []string `yaml:"protected_branches"`
}

type UIConfig struct {
	Color   bool `yaml:"color"`
	Emoji   bool `yaml:"emoji"`
	Verbose bool `yaml:"verbose"`
}

type LoadResult struct {
	Path   string
	Config *Config
}

func Load() (*LoadResult, error) {
	startDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}

	return LoadFromDir(startDir)
}

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
