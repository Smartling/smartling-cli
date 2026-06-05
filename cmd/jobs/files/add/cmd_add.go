package add

import (
	"fmt"
	"os"

	filescmd "github.com/Smartling/smartling-cli/cmd/jobs/files"
	"github.com/Smartling/smartling-cli/output"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/spf13/cobra"
)

const (
	fileFlag         = "file"
	targetLocaleFlag = "target-locale"
)

// NewJobFilesAddCmd returns new command to job file add
func NewJobFilesAddCmd(initializer filescmd.SrvInitializer) *cobra.Command {
	var (
		filePatterns  []string
		targetLocales []string
	)
	addCmd := &cobra.Command{
		Use:   "add <translationJobUid|translationJobName>",
		Short: "Add files to a translation job.",
		Long:  `Attach files to an existing translation job. Each --file is a glob pattern matched against the project's files.`,
		Args:  cobra.ExactArgs(1),
		Example: `
# Add all JSON files to a job

  smartling-cli jobs files add aabbccdd1122 --file "**/*.json"

# Add files for specific locales

  smartling-cli jobs files add "Website Q1 2026" --file "menu/*.xml" --target-locale fr-FR
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			params, err := resolveParams(args[0], filePatterns, targetLocales)
			if err != nil {
				return fmt.Errorf("failed to resolve add params: %w", err)
			}
			format, err := cmd.Flags().GetString("output")
			if err != nil {
				return err
			}
			return run(ctx, initializer, params, output.Params{Format: format})
		},
	}

	addCmd.Flags().StringArrayVar(&filePatterns, fileFlag, nil, "File URI glob pattern to add (repeatable, required).")
	addCmd.Flags().StringArrayVar(&targetLocales, targetLocaleFlag, nil, "Target locale to add the files to (repeatable; default all job locales).")
	if err := addCmd.MarkFlagRequired(fileFlag); err != nil {
		rlog.Errorf("failed to mark --%s required: %s", fileFlag, err)
		os.Exit(1)
	}

	return addCmd
}
