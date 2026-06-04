package add

import (
	"fmt"

	localescmd "github.com/Smartling/smartling-cli/cmd/jobs/locales"
	"github.com/Smartling/smartling-cli/output"

	"github.com/spf13/cobra"
)

// NewJobLocalesAddCmd returns new command to job locale add
func NewJobLocalesAddCmd(initializer localescmd.SrvInitializer) *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "add <translationJobUid|translationJobName> <targetLocaleId>",
		Short: "Add a target locale to a translation job.",
		Long:  `Attach a target locale to an existing translation job, identified by UID or name.`,
		Args:  cobra.ExactArgs(2),
		Example: `
# Add a locale to a job by UID

  smartling-cli jobs locales add aabbccdd1122 fr-FR

# Add a locale to a job by name

  smartling-cli jobs locales add "Website Q1 2026" fr-FR
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			params, err := resolveParams(args[0], args[1])
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

	return addCmd
}
