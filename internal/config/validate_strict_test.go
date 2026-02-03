package config

import "testing"

func TestValidateStrictRejectsBadStrategy(t *testing.T) {
	cfg := Default()
	cfg.Workflows.Sync.Strategy = "bad"
	if err := ValidateStrict(cfg); err == nil {
		t.Fatalf("expected error")
	}
}
