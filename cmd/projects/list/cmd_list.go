package list

import (
	"github.com/Smartling/smartling-cli/services/projects"

	"github.com/spf13/cobra"
)

var short bool

func NewListCmd(s projects.Service) *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Lists projects for current account.",
		Long:  `Lists projects for current account.`,
		Run: func(cmd *cobra.Command, args []string) {
			err := s.RunList(short)
			if err != nil {
				// TODO log it
			}
		},
	}
	listCmd.Flags().BoolVarP(&short, "short", "s", false, "Display only project IDs.")
	return listCmd
}
