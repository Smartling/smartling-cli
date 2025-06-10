package info

import (
	rootcmd "github.com/Smartling/smartling-cli/cmd"
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
				rootcmd.Logger().Errorf("failed to get project service: %s", err)
				return
			}
			err = s.RunInfo()
			if err != nil {
				rootcmd.Logger().Errorf("failed to run info: %s", err)
				return
			}
		},
	}
	return infoCmd
}
