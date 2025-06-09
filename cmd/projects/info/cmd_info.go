package info

import (
	"github.com/Smartling/smartling-cli/services/projects"

	"github.com/spf13/cobra"
)

func NewInfoCmd(s *projects.Service) *cobra.Command {
	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Get project details about specific project.",
		Long:  `Get project details about specific project.`,
		Run: func(cmd *cobra.Command, args []string) {
			err := s.RunInfo()
			if err != nil {
				// TODO log it
			}
		},
	}
	return infoCmd
}
