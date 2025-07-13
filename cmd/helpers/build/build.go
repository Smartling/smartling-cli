package build

import (
	"fmt"
)

// These variables will be set via ldflags at build time.
var (
	// CliVersion is version information
	CliVersion = "2.1"
	// ReleaseTag is the git tag
	ReleaseTag = ""
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
ReleaseTag:    %s
GitCommit:     %s
BuildDate:     %s
BuiltBy:       %s
GoVersion:     %s
Platform:      %s
`,
		CliVersion,
		ReleaseTag,
		CommitID,
		Date,
		BuiltBy,
		GoVersion,
		Platform,
	)
}
