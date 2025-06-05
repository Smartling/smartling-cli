package info

import (
	"github.com/spf13/cobra"
)

func NewInfoCmd() *cobra.Command {
	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Get project details about specific project.",
		Long:  `Get project details about specific project.`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	return infoCmd
}
