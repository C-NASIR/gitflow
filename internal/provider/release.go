package provider

import "errors"

var ErrReleaseExists = errors.New("release already exists")

func IsReleaseExists(err error) bool {
	return errors.Is(err, ErrReleaseExists)
}
