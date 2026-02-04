package release

import (
	"strconv"
	"strings"

	"gitflow/internal/workflow"
)

func parseVersion(input string) (workflow.SemanticVersion, bool) {
	parts := strings.Split(input, ".")
	if len(parts) != 3 {
		return workflow.SemanticVersion{}, false
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return workflow.SemanticVersion{}, false
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return workflow.SemanticVersion{}, false
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return workflow.SemanticVersion{}, false
	}
	return workflow.SemanticVersion{Major: major, Minor: minor, Patch: patch}, true
}
