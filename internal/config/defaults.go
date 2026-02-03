package config

// Default returns a Config with built-in defaults.
func Default() *Config {
	return &Config{
		Branches: BranchConfig{
			FeaturePrefix: "feature/",
			BugfixPrefix:  "bugfix/",
			HotfixPrefix:  "hotfix/",
			MainBranch:    "main",
			DevelopBranch: "",
		},
		Workflows: WorkflowConfig{
			Start: StartConfig{
				BaseBranch: "main",
				AutoPush:   true,
				FetchFirst: true,
			},
			PR: PRConfig{
				Draft:            false,
				DefaultReviewers: nil,
				Labels:           nil,
			},
			Sync: SyncConfig{
				Strategy:  "rebase",
				AutoPush:  true,
				ForcePush: true,
			},
			Cleanup: CleanupConfig{
				MergedOnly:       true,
				AgeThresholdDays: 30,
				ProtectedBranches: []string{
					"main", "master", "develop",
				},
			},
		},
		Commits: CommitConfig{
			Conventional: false,
			Types:        []string{"feat", "fix", "docs", "refactor", "test", "chore"},
			Scopes:       nil,
			RequireScope: false,
		},
		UI: UIConfig{
			Color:   true,
			Emoji:   false,
			Verbose: false,
		},
	}
}
