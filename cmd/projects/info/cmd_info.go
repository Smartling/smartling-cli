package info

import (
	projectscmd "github.com/Smartling/smartling-cli/cmd/projects"

	"github.com/spf13/cobra"
)

func NewInfoCmd() *cobra.Command {
	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Get project details about specific project.",
		Long:  `Get project details about specific project.`,
		Run: func(cmd *cobra.Command, args []string) {
			s, err := projectscmd.GetService()
			if err != nil {
				// TODO log it
			}
			err = s.RunInfo()
			if err != nil {
				// TODO log it
			}
		},
	}
	return infoCmd
}
