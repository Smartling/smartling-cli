package add

import (
	"fmt"

	stringscmd "github.com/Smartling/smartling-cli/cmd/jobs/strings"
	"github.com/Smartling/smartling-cli/output"

	"github.com/spf13/cobra"
)

const (
	hashcodeFlag     = "hashcode"
	targetLocaleFlag = "target-locale"
	moveEnabledFlag  = "move-enabled"
)

// NewJobStringsAddCmd returns new command to job string add
func NewJobStringsAddCmd(initializer stringscmd.SrvInitializer) *cobra.Command {
	var (
		hashcodes     []string
		targetLocales []string
		moveEnabled   bool
	)
	addCmd := &cobra.Command{
		Use:   "add <translationJobUid|translationJobName>",
		Short: "Add strings to a translation job.",
		Long:  `Assign strings (by hashcode) to an existing translation job, identified by UID or name.`,
		Args:  cobra.ExactArgs(1),
		Example: `
# Add two strings to a job

  smartling-cli jobs strings add aabbccdd1122 --hashcode h1 --hashcode h2

# Add a string for specific locales, moving it if it already belongs to another job

  smartling-cli jobs strings add "Website Q1 2026" --hashcode h1 --target-locale fr-FR --move-enabled
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			params, err := resolveParams(args[0], hashcodes, targetLocales, moveEnabled)
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

	addCmd.Flags().StringArrayVar(&hashcodes, hashcodeFlag, nil, "String hashcode to add (repeatable, required).")
	addCmd.Flags().StringArrayVar(&targetLocales, targetLocaleFlag, nil, "Target locale to add the strings to (repeatable; default all job locales).")
	addCmd.Flags().BoolVar(&moveEnabled, moveEnabledFlag, false, "Move the string into this job if it already belongs to another job for a locale.")

	return addCmd
}
