package jobstrings

import (
	"github.com/spf13/cobra"
)

// NewJobStringsCmd returns new job strings command
func NewJobStringsCmd() *cobra.Command {
	jobStringsCmd := &cobra.Command{
		Use:   "strings",
		Short: "Manage strings on a translation job.",
		Long: `Add, remove, or list the strings on an existing translation job.

Strings are identified by hashcode. Use these commands to assign strings to a
job, detach them, or inspect which strings a job currently contains.`,
		Example: `
# Add strings to a job

  smartling-cli jobs strings add <translationJobUid|translationJobName> --hashcode <hashcode>

# Remove strings from a job

  smartling-cli jobs strings remove <translationJobUid|translationJobName> --hashcode <hashcode>

# List a job's strings

  smartling-cli jobs strings list <translationJobUid|translationJobName>

`,
	}

	return jobStringsCmd
}
