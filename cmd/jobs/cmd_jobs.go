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
		Short: "Handles job subcommands.",
		Long:  `Handles job subcommands. Subcommands are a high-level abstraction layer over the underlying Job APIs.`,
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
