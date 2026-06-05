package remove

import (
	"fmt"
	"os"

	filescmd "github.com/Smartling/smartling-cli/cmd/jobs/files"
	"github.com/Smartling/smartling-cli/output"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/spf13/cobra"
)

const fileFlag = "file"

// NewJobFilesRemoveCmd returns new command to job file remove
func NewJobFilesRemoveCmd(initializer filescmd.SrvInitializer) *cobra.Command {
	var filePatterns []string
	removeCmd := &cobra.Command{
		Use:   "remove <translationJobUid|translationJobName>",
		Short: "Remove files from a translation job.",
		Long:  `Detach files from an existing translation job. Each --file is a glob pattern matched against the project's files.`,
		Args:  cobra.ExactArgs(1),
		Example: `
# Remove all XML files from a job

  smartling-cli jobs files remove aabbccdd1122 --file "**/*.xml"

# Remove a specific file from a job by name

  smartling-cli jobs files remove "Website Q1 2026" --file "menu/old.json"
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			params, err := resolveParams(args[0], filePatterns)
			if err != nil {
				return fmt.Errorf("failed to resolve remove params: %w", err)
			}
			format, err := cmd.Flags().GetString("output")
			if err != nil {
				return err
			}
			return run(ctx, initializer, params, output.Params{Format: format})
		},
	}

	removeCmd.Flags().StringArrayVar(&filePatterns, fileFlag, nil, "File URI glob pattern to remove (repeatable, required).")
	if err := removeCmd.MarkFlagRequired(fileFlag); err != nil {
		rlog.Errorf("failed to mark --%s required: %s", fileFlag, err)
		os.Exit(1)
	}

	return removeCmd
}
