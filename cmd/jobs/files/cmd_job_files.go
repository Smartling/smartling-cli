package jobfiles

import (
	"github.com/spf13/cobra"
)

// NewJobFilesCmd returns new job files command
func NewJobFilesCmd() *cobra.Command {
	jobFilesCmd := &cobra.Command{
		Use:   "files",
		Short: "Manage the source files attached to a translation job.",
		Long: `Add, remove, or list the source files attached to a translation job.

Files are identified by their Smartling fileUri. The add and remove commands
accept repeatable --file glob patterns, matched against the project's files.`,
		Example: `
# List a job's files

  smartling-cli jobs files list <translationJobUid|translationJobName>

# Add files to a job

  smartling-cli jobs files add <translationJobUid|translationJobName> --file "**/*.json"

# Remove files from a job

  smartling-cli jobs files remove <translationJobUid|translationJobName> --file "old/*.xml"

`,
	}

	return jobFilesCmd
}
