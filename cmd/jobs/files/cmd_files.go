package jobfiles

import (
	"errors"
	"fmt"

	rootcmd "github.com/Smartling/smartling-cli/cmd"
	jobscmd "github.com/Smartling/smartling-cli/cmd/jobs"
	"github.com/Smartling/smartling-cli/output"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	srv "github.com/Smartling/smartling-cli/services/jobs"

	"github.com/spf13/cobra"
)

// NewFilesCmd builds the `jobs files` command.
func NewFilesCmd(initializer jobscmd.SrvInitializer) *cobra.Command {
	var (
		limit  uint32
		offset uint32
	)

	filesCmd := &cobra.Command{
		Use:   "files <translationJobUid|translationJobName>",
		Short: "List source files attached to a translation job.",
		Long:  `List the source files attached to a translation job, by UID or name.`,
		Example: `
# List files for a job by UID

  smartling-cli jobs files aabbccdd1122

# List files for a job by name in JSON

  smartling-cli jobs files "Website Q1 2026" --output json
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if len(args) != 1 {
				return clierror.UIError{
					Operation:   "check args",
					Err:         errors.New("wrong argument quantity"),
					Description: fmt.Sprintf("expected one argument, got: %d", len(args)),
				}
			}

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

	filesCmd.Flags().Uint32Var(&limit, "limit", 500, "Maximum number of files to return.")
	filesCmd.Flags().Uint32Var(&offset, "offset", 0, "Offset for pagination.")

	return filesCmd
}
