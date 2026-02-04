package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ApplyEnvOverrides(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	if value, ok := os.LookupEnv("GITFLOW_RELEASE_TAG_PREFIX"); ok {
		cfg.Release.TagPrefix = strings.TrimSpace(value)
	}
	if value, ok := os.LookupEnv("GITFLOW_RELEASE_DEFAULT_BUMP"); ok {
		cfg.Release.DefaultBump = strings.TrimSpace(strings.ToLower(value))
	}
	if value, ok := os.LookupEnv("GITFLOW_UI_NO_COLOR"); ok {
		enabled, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid GITFLOW_UI_NO_COLOR: %w", err)
		}
		cfg.UI.Color = !enabled
	}
	if value, ok := os.LookupEnv("GITFLOW_UI_EMOJI"); ok {
		enabled, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid GITFLOW_UI_EMOJI: %w", err)
		}
		cfg.UI.Emoji = enabled
	}
	if value, ok := os.LookupEnv("GITFLOW_UI_VERBOSE"); ok {
		enabled, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid GITFLOW_UI_VERBOSE: %w", err)
		}
		cfg.UI.Verbose = enabled
	}

	switch cfg.Release.DefaultBump {
	case "major", "minor", "patch":
	case "":
		cfg.Release.DefaultBump = "patch"
	default:
		return fmt.Errorf("unsupported release default bump: %s", cfg.Release.DefaultBump)
	}

	return nil
}
