package list

import (
	"os"

	projectscmd "github.com/Smartling/smartling-cli/cmd/projects"
	"github.com/Smartling/smartling-cli/services/helpers/help"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/spf13/cobra"
)

var short bool

// NewListCmd creates a new command to list projects.
func NewListCmd(initializer projectscmd.SrvInitializer) *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Lists projects for current account.",
		Long: `smartling-cli projects list — list projects from account.

Command will list projects from specified account in tabular format with
following information:

  > Project ID
  > Project Description
  > Project Source Locale ID

Only project IDs will be listed if --short option is specified.

Note, that you should specify account ID either in config file or via --account
option to be able to see projects list.


Available options:
  -s --short
    List only project IDs.
` + help.AuthenticationOptions,
		Run: func(_ *cobra.Command, _ []string) {
			s, err := initializer.InitProjectsSrv()
			if err != nil {
				rlog.Errorf("failed to get project service: %s", err)
				os.Exit(1)
			}

			err = s.RunList(short)
			if err != nil {
				rlog.Errorf("failed to run list: %s", err)
				os.Exit(1)
			}
		},
	}
	listCmd.Flags().BoolVarP(&short, "short", "s", false, "Display only project IDs.")
	return listCmd
}
