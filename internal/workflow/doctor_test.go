package workflow

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDoctorNonGitDirectory(t *testing.T) {
	dir := t.TempDir()

	res, err := Doctor(dir)
	if err != nil {
		t.Fatalf("Doctor: %v", err)
	}
	if len(res.Checks) == 0 {
		t.Fatalf("expected checks")
	}
	if res.Checks[0].Level != DoctorError {
		t.Fatalf("expected error for non-git directory")
	}
}

func TestDoctorCleanRepo(t *testing.T) {
	dir := setupRepo(t)
	defer os.RemoveAll(dir)

	res, err := Doctor(dir)
	if err != nil {
		t.Fatalf("Doctor: %v", err)
	}
	level := checkLevel(t, res, "Working tree")
	if level != DoctorOK {
		t.Fatalf("expected clean repo to be OK, got %s", level)
	}
}

func TestDoctorDirtyRepo(t *testing.T) {
	dir := setupRepo(t)
	defer os.RemoveAll(dir)

	if err := os.WriteFile(filepath.Join(dir, "dirty.txt"), []byte("dirty"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	res, err := Doctor(dir)
	if err != nil {
		t.Fatalf("Doctor: %v", err)
	}
	level := checkLevel(t, res, "Working tree")
	if level != DoctorWarn {
		t.Fatalf("expected dirty repo to warn, got %s", level)
	}
}

func checkLevel(t *testing.T, res *DoctorResult, name string) string {
	t.Helper()
	for _, check := range res.Checks {
		if check.Name == name {
			return check.Level
		}
	}
	t.Fatalf("missing check %s", name)
	return ""
}
