package version

import "runtime"

var (
	Version = "dev"
	Commit  = "none"
	Date 	= "unkown"
)

func String() string {
	return Version + "(" + runtime.Version() + ")"
}