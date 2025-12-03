package progress

import (
	"errors"
	"fmt"

	rootcmd "github.com/Smartling/smartling-cli/cmd"
	"github.com/Smartling/smartling-cli/cmd/helpers/resolve"
	jobscmd "github.com/Smartling/smartling-cli/cmd/jobs"
	"github.com/Smartling/smartling-cli/output"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	srv "github.com/Smartling/smartling-cli/services/jobs"

	"github.com/spf13/cobra"
)

// NewProgressCmd returns new progress command
func NewProgressCmd(initializer jobscmd.SrvInitializer) *cobra.Command {
	progressCmd := &cobra.Command{
		Use:   "progress <translationJobUid|translationJobName>",
		Short: "Get job progress by the translationJobUid or translationJobName.",
		Long:  `Get job progress by the translationJobUid or translationJobName.`,
		Example: `
# Get job progress

  smartling-cli jobs progress aabbccdd

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
			idOrName := args[0]

			cnf, err := rootcmd.Config()
			if err != nil {
				return err
			}

			accountUID, err := resolve.FallbackAccount(cmd.Root().PersistentFlags().Lookup("account"), cnf.AccountID)
			if err != nil {
				return err
			}

			params := srv.ProgressParams{
				AccountUID:  accountUID,
				ProjectUID:  cnf.ProjectID,
				JobIDOrName: idOrName,
			}
			format, err := cmd.Parent().PersistentFlags().GetString("output")
			if err != nil {
				return err
			}
			outputParams := output.Params{Format: format}
			return run(ctx, initializer, params, outputParams)
		},
	}

	return progressCmd
}
