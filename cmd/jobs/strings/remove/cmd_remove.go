package remove

import (
	"fmt"
	"os"

	stringscmd "github.com/Smartling/smartling-cli/cmd/jobs/strings"
	"github.com/Smartling/smartling-cli/output"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/spf13/cobra"
)

const (
	hashcodeFlag = "hashcode"
	localeFlag   = "locale"
)

// NewJobStringsRemoveCmd returns new command to job string remove
func NewJobStringsRemoveCmd(initializer stringscmd.SrvInitializer) *cobra.Command {
	var (
		hashcodes []string
		localeIDs []string
	)
	removeCmd := &cobra.Command{
		Use:   "remove <translationJobUid|translationJobName>",
		Short: "Remove strings from a translation job.",
		Long:  `Detach strings (by hashcode) from an existing translation job, identified by UID or name.`,
		Args:  cobra.ExactArgs(1),
		Example: `
# Remove two strings from a job

  smartling-cli jobs strings remove aabbccdd1122 --hashcode h1 --hashcode h2

# Remove a string from specific locales only

  smartling-cli jobs strings remove "Website Q1 2026" --hashcode h1 --locale fr-FR
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			params, err := resolveParams(args[0], hashcodes, localeIDs)
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

	removeCmd.Flags().StringArrayVar(&hashcodes, hashcodeFlag, nil, "String hashcode to remove (repeatable, required).")
	removeCmd.Flags().StringArrayVar(&localeIDs, localeFlag, nil, "Locale to remove the strings from (repeatable; default all job locales).")
	if err := removeCmd.MarkFlagRequired(hashcodeFlag); err != nil {
		rlog.Errorf("failed to mark --%s required: %s", hashcodeFlag, err)
		os.Exit(1)
	}

	return removeCmd
}
