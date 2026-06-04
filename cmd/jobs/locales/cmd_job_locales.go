package locales

import (
	"github.com/spf13/cobra"
)

// NewJobLocalesCmd returns new job locales command
func NewJobLocalesCmd() *cobra.Command {
	jobLocalesCmd := &cobra.Command{
		Use:   "locales",
		Short: "Manage target locales on a translation job.",
		Long: `Add or remove target locales on an existing translation job.

A job's target locales determine which languages its content is translated into.
Use these commands to attach a locale to a job or detach one from it.`,
		Example: `
# Add a target locale to a job

  smartling-cli jobs locales add <translationJobUid|translationJobName> <targetLocaleId>

# Remove a target locale from a job

  smartling-cli jobs locales remove <translationJobUid|translationJobName> <targetLocaleId>

`,
	}

	return jobLocalesCmd
}
