package provider

import "errors"

// ErrReleaseExists indicates a release already exists for a tag.
var ErrReleaseExists = errors.New("release already exists")

// IsReleaseExists reports whether the error indicates an existing release.
func IsReleaseExists(err error) bool {
	return errors.Is(err, ErrReleaseExists)
}
