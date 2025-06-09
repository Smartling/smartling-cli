package projects

import (
	"github.com/Smartling/smartling-cli/cmd/projects/info"
	"github.com/Smartling/smartling-cli/cmd/projects/list"
	"github.com/Smartling/smartling-cli/cmd/projects/locales"
	"github.com/Smartling/smartling-cli/services/projects"

	"github.com/spf13/cobra"
)

func NewProjectsCmd(s *projects.Service) *cobra.Command {
	projectsCmd := &cobra.Command{
		Use:     "projects",
		Aliases: []string{"f"},
		Short:   "Used to access various projects sub-commands.",
		Long:    `Used to access various projects sub-commands.`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	projectsCmd.AddCommand(list.NewListCmd(s))
	projectsCmd.AddCommand(info.NewInfoCmd(s))
	projectsCmd.AddCommand(locales.NewLocatesCmd(s))

	return projectsCmd
}
