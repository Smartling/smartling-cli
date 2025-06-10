package list

import (
	rootcmd "github.com/Smartling/smartling-cli/cmd"
	projectscmd "github.com/Smartling/smartling-cli/cmd/projects"

	"github.com/spf13/cobra"
)

var short bool

func NewListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Lists projects for current account.",
		Long:  `Lists projects for current account.`,
		Run: func(cmd *cobra.Command, args []string) {
			s, err := projectscmd.GetService()
			if err != nil {
				rootcmd.Logger().Errorf("failed to get project service: %s", err)
				return
			}

			err = s.RunList(short)
			if err != nil {
				rootcmd.Logger().Errorf("failed to run list: %s", err)
				return
			}
		},
	}
	listCmd.Flags().BoolVarP(&short, "short", "s", false, "Display only project IDs.")
	return listCmd
}
