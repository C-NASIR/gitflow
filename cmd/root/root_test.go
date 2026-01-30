package root

import (
	"bytes"
	"testing"
)

func TestVersionTemplateAndCommand(t *testing.T) {
	buf := new(bytes.Buffer)

	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs([]string{"version"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected no error, go %v", err)
	}

	out := buf.String()
	if out == "" {
		t.Fatalf("expected version output, got empty string")
	}
}