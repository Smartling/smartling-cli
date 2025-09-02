package delete

import (
	"os"

	"github.com/Smartling/smartling-cli/cmd/files"
	"github.com/Smartling/smartling-cli/services/helpers/help"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/spf13/cobra"
)

var uri string

// NewDeleteCmd creates a new command to delete files.
func NewDeleteCmd(initializer files.SrvInitializer) *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:   "delete <uri>",
		Short: "Deletes given file from Smartling.",
		Long: `smartling-cli files delete — removes files from project.

Removes files from project according to specified pattern.

<uri> ` + help.GlobPattern + ` </uri>

If special value of "-" is specified as <uri>, then program will expect
to read files list from stdin:

  cat files.txt | smartling-cli files delete -

Available options:
  -p --project <project>
    Specify project to use.
` + help.AuthenticationOptions,
		Run: func(_ *cobra.Command, args []string) {
			if len(args) > 0 {
				uri = args[0]
			}

			s, err := initializer.InitFilesSrv()
			if err != nil {
				rlog.Errorf("failed to get files service: %s", err)
				os.Exit(1)
			}

			err = s.RunDelete(uri)
			if err != nil {
				rlog.Errorf("failed to run delete: %s", err)
				os.Exit(1)
			}
		},
	}

	return deleteCmd
}
