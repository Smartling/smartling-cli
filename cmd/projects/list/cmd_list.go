package list

import (
	projectscmd "github.com/Smartling/smartling-cli/cmd/projects"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/spf13/cobra"
)

var short bool

// NewListCmd creates a new command to list projects.
func NewListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Lists projects for current account.",
		Long:  `Lists projects for current account.`,
		Run: func(_ *cobra.Command, _ []string) {
			s, err := projectscmd.GetService()
			if err != nil {
				rlog.Errorf("failed to get project service: %s", err)
				return
			}

			err = s.RunList(short)
			if err != nil {
				rlog.Errorf("failed to run list: %s", err)
				return
			}
		},
	}
	listCmd.Flags().BoolVarP(&short, "short", "s", false, "Display only project IDs.")
	return listCmd
}
