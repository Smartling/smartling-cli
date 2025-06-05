package projects

import (
	"github.com/Smartling/smartling-cli/cmd/projects/info"
	"github.com/Smartling/smartling-cli/cmd/projects/list"
	"github.com/Smartling/smartling-cli/cmd/projects/locales"

	"github.com/spf13/cobra"
)

func NewProjectsCmd() *cobra.Command {
	projectsCmd := &cobra.Command{
		Use:     "projects",
		Aliases: []string{"f"},
		Short:   "Used to access various projects sub-commands.",
		Long:    `Used to access various projects sub-commands.`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	projectsCmd.AddCommand(list.NewListCmd())
	projectsCmd.AddCommand(info.NewInfoCmd())
	projectsCmd.AddCommand(locales.NewLocatesCmd())

	return projectsCmd
}
