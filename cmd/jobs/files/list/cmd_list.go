package list

import (
	rootcmd "github.com/Smartling/smartling-cli/cmd"
	jobscmd "github.com/Smartling/smartling-cli/cmd/jobs"
	"github.com/Smartling/smartling-cli/output"
	srv "github.com/Smartling/smartling-cli/services/jobs"

	"github.com/spf13/cobra"
)

// NewListCmd builds the `jobs files list` command.
func NewListCmd(initializer jobscmd.SrvInitializer) *cobra.Command {
	var (
		limit  uint32
		offset uint32
	)

	filesCmd := &cobra.Command{
		Use:   "list <translationJobUid|translationJobName>",
		Short: "List source files attached to a translation job.",
		Long:  `List the source files attached to a translation job, by UID or name.`,
		Args:  cobra.ExactArgs(1),
		Example: `
# List files for a job by UID

  smartling-cli jobs files list aabbccdd1122

# List files for a job by name in JSON

  smartling-cli jobs files list "Website Q1 2026" --output json
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			cnf, err := rootcmd.Config()
			if err != nil {
				return err
			}

			params := srv.FilesParams{
				ProjectUID:   cnf.ProjectID,
				JobUIDOrName: args[0],
				Limit:        limit,
				Offset:       offset,
			}
			format, err := cmd.Flags().GetString("output")
			if err != nil {
				return err
			}
			return run(ctx, initializer, params, output.Params{Format: format})
		},
	}

	filesCmd.Flags().Uint32Var(&limit, "limit", srv.DefaultListPageLimit, "Maximum number of files to return.")
	filesCmd.Flags().Uint32Var(&offset, "offset", 0, "Offset for pagination.")

	return filesCmd
}
