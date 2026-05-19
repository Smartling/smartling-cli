package build

import (
	"fmt"
)

// These variables will be set via ldflags at build time.
var (
	// CliVersion is version information
	CliVersion = "3.2"
	// CommitID is the git commit hash
	CommitID = ""
	// Date is date of binary built
	Date = ""
	// BuiltBy is into about builder
	BuiltBy = "GoReleaser"
	// GoVersion is Go version
	GoVersion = ""
	// Platform is platform
	Platform = ""
)

// String returns the version information as a formatted string.
func String() string {
	return fmt.Sprintf(`
Smartling-cli is a library and CLI tool for managing Smartling projects.

Version:       %s
GitCommit:     %s
BuildDate:     %s
BuiltBy:       %s
GoVersion:     %s
Platform:      %s
`,
		CliVersion,
		CommitID,
		Date,
		BuiltBy,
		GoVersion,
		Platform,
	)
}
