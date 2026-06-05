package list

import (
	"fmt"

	stringscmd "github.com/Smartling/smartling-cli/cmd/jobs/strings"
	"github.com/Smartling/smartling-cli/output"

	"github.com/spf13/cobra"
)

const (
	targetLocaleFlag = "target-locale"
	limitFlag        = "limit"
	offsetFlag       = "offset"
)

// NewJobStringsListCmd returns new command to job string list
func NewJobStringsListCmd(initializer stringscmd.SrvInitializer) *cobra.Command {
	var (
		targetLocale string
		limit        uint32
		offset       uint32
	)
	listCmd := &cobra.Command{
		Use:   "list <translationJobUid|translationJobName>",
		Short: "List the strings on a translation job.",
		Long:  `List the strings (by hashcode and target locale) assigned to a translation job, identified by UID or name.`,
		Args:  cobra.ExactArgs(1),
		Example: `
# List a job's strings

  smartling-cli jobs strings list aabbccdd1122

# List strings for one locale as a table

  smartling-cli jobs strings list "Website Q1 2026" --target-locale fr-FR --output table
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			params, err := resolveParams(args[0], targetLocale, limit, offset)
			if err != nil {
				return fmt.Errorf("failed to resolve list params: %w", err)
			}
			format, err := cmd.Flags().GetString("output")
			if err != nil {
				return err
			}
			return run(ctx, initializer, params, output.Params{Format: format})
		},
	}

	listCmd.Flags().StringVar(&targetLocale, targetLocaleFlag, "", "Filter strings by target locale.")
	listCmd.Flags().Uint32Var(&limit, limitFlag, 0, "Maximum number of strings to return.")
	listCmd.Flags().Uint32Var(&offset, offsetFlag, 0, "Number of strings to skip.")

	return listCmd
}
