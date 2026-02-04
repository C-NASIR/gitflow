package ui

import (
	"bytes"
	"strings"
	"testing"
)

func TestUIWrites(t *testing.T) {
	var buf bytes.Buffer
	u := New(Options{
		Out:     &buf,
		Color:   false,
		Emoji:   false,
		Verbose: true,
	})

	u.Header("Test")
	u.Success("ok")
	u.Warn("warn")
	u.Error("err")
	u.Info("info")
	u.Verbose("v")

	s := buf.String()
	if !strings.Contains(s, "Test") {
		t.Fatalf("expected header")
	}
	if !strings.Contains(s, "OK") {
		t.Fatalf("expected ok prefix")
	}
}
