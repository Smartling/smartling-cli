package jobview

import (
	rootcmd "github.com/Smartling/smartling-cli/cmd"
	jobscmd "github.com/Smartling/smartling-cli/cmd/jobs"
	"github.com/Smartling/smartling-cli/output"
	srv "github.com/Smartling/smartling-cli/services/jobs"

	"github.com/spf13/cobra"
)

// NewViewCmd builds the `jobs view` command.
func NewViewCmd(initializer jobscmd.SrvInitializer) *cobra.Command {
	viewCmd := &cobra.Command{
		Use:   "view <translationJobUid|translationJobName>",
		Short: "Show full details of a translation job.",
		Long:  `Retrieve full details of a single translation job by UID or name.`,
		Args:  cobra.ExactArgs(1),
		Example: `
# View a job by UID

  smartling-cli jobs view aabbccdd1122

# View a job by name

  smartling-cli jobs view "Website Q1 2026" --output json
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			cnf, err := rootcmd.Config()
			if err != nil {
				return err
			}

			params := srv.ViewParams{
				ProjectUID:   cnf.ProjectID,
				JobUIDOrName: args[0],
			}
			format, err := cmd.Flags().GetString("output")
			if err != nil {
				return err
			}
			return run(ctx, initializer, params, output.Params{Format: format})
		},
	}

	return viewCmd
}
