// Package version exposes build-time version metadata.
package version

import "runtime"

var (
	// Version is the gitflow build version.
	Version = "dev"
	// Commit is the git commit hash for the build.
	Commit = "none"
	// Date is the build timestamp.
	Date = "unkown"
)

// String formats version information for display.
func String() string {
	return Version + "(" + runtime.Version() + ")"
}
