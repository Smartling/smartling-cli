package remove

import (
	"fmt"

	localescmd "github.com/Smartling/smartling-cli/cmd/jobs/locales"
	"github.com/Smartling/smartling-cli/output"

	"github.com/spf13/cobra"
)

// NewJobLocalesRemoveCmd returns new command to job locale remove
func NewJobLocalesRemoveCmd(initializer localescmd.SrvInitializer) *cobra.Command {
	removeCmd := &cobra.Command{
		Use:   "remove <translationJobUid|translationJobName> <targetLocaleId>",
		Short: "Remove a target locale from a translation job.",
		Long:  `Detach a target locale from an existing translation job, identified by UID or name.`,
		Args:  cobra.ExactArgs(2),
		Example: `
# Remove a locale from a job by UID

  smartling-cli jobs locales remove aabbccdd1122 fr-FR

# Remove a locale from a job by name

  smartling-cli jobs locales remove "Website Q1 2026" fr-FR
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			params, err := resolveParams(args[0], args[1])
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

	return removeCmd
}
