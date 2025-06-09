package projects

import (
	"github.com/Smartling/smartling-cli/cmd"
	"github.com/Smartling/smartling-cli/services/projects"
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

	return projectsCmd
}

func GetService() (*projects.Service, error) {
	client, err := cmd.Client()
	if err != nil {
		return nil, err
	}
	cnf, err := cmd.Config()
	if err != nil {
		return nil, err
	}
	srv := projects.NewService(&client, cnf)
	return srv, nil
}
