package list

import (
	"github.com/spf13/cobra"
)

var short bool

func NewListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Lists projects for current account.",
		Long:  `Lists projects for current account.`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	listCmd.Flags().BoolVarP(&short, "short", "s", false, "Display only project IDs.")
	return listCmd
}
