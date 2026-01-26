package jobs

import (
	"fmt"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

const (
	outputFormatFlag = "output"
)

var (
	outputFormat   string
	allowedOutputs = []string{
		"json",
		"simple",
	}
	joinedAllowedOutputs = strings.Join(allowedOutputs, ", ")
)

// NewJobsCmd returns new jobs command
func NewJobsCmd() *cobra.Command {
	jobsCmd := &cobra.Command{
		Use:   "jobs",
		Short: "Manage translation jobs and monitor their progress.",
		Long: `Translation jobs are the fundamental unit of work in Smartling TMS that organize
content for translation and track it through the translation workflow.

The jobs command group provides tools to interact with translation jobs, including
monitoring translation progress, viewing job details, and managing job workflows.

Each job contains one or more files targeted for translation into specific locales,
with defined due dates and workflow steps. Jobs help coordinate translation work between
content owners, project managers, and translators.

Available options:
  --output string   Output format: ` + joinedAllowedOutputs + ` (default "simple")
                    - simple: Human-readable format optimized for terminal display
                    - json: Raw API response for programmatic processing and automation`,
		Example: `
# View job progress in human-readable format

  smartling-cli jobs progress "Website Q1 2026"

# Get detailed progress data for automation

  smartling-cli jobs progress aabbccdd --output json

`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !slices.Contains(allowedOutputs, outputFormat) {
				return fmt.Errorf("invalid output: %s (allowed: %s)", outputFormat, joinedAllowedOutputs)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 && cmd.Flags().NFlag() == 0 {
				return cmd.Help()
			}
			return nil
		},
	}

	jobsCmd.PersistentFlags().StringVar(&outputFormat, outputFormatFlag, "simple", "Output format: "+joinedAllowedOutputs)

	return jobsCmd
}
