// Package version exposes build metadata (version, commit, date) injected at compile time via ldflags.
package version

import "fmt"

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func String() string {
	return fmt.Sprintf("%s (commit=%s date=%s)", Version, Commit, Date)
}
