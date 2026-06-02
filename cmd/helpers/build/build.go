package build

import (
	"fmt"
	"runtime"
)

// These variables will be set via ldflags at build time.
var (
	// CliVersion is version information
	CliVersion = "3.4"
	// CommitID is the git commit hash
	CommitID = ""
	// Date is date of binary built
	Date = ""
	// BuiltBy identifies the build pipeline (e.g. "GoReleaser", "make").
	// Defaults to "unknown" so binaries built without ldflag injection
	// don't mislabel themselves.
	BuiltBy = "unknown"
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
		runtime.Version(),
		runtime.GOOS+"/"+runtime.GOARCH,
	)
}
