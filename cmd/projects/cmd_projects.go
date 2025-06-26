package projects

import (
	"os"

	"github.com/Smartling/smartling-cli/cmd"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	"github.com/Smartling/smartling-cli/services/projects"

	"github.com/spf13/cobra"
)

// NewProjectsCmd creates a new projects command.
func NewProjectsCmd() *cobra.Command {
	projectsCmd := &cobra.Command{
		Use:     "projects",
		Aliases: []string{"p"},
		Short:   "Used to access various projects sub-commands.",
		Long:    `Used to access various projects sub-commands.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 && cmd.Flags().NFlag() == 0 {
				if err := cmd.Help(); err != nil {
					rlog.Error(err.Error())
					os.Exit(1)
				}
				return
			}
		},
	}

	return projectsCmd
}

// SrvInitializer defines projects service initializer
type SrvInitializer interface {
	InitProjectsSrv() (projects.Service, error)
}

// NewSrvInitializer returns new SrvInitializer implementation
func NewSrvInitializer() SrvInitializer {
	return srvInitializer{}
}

type srvInitializer struct{}

// InitProjectsSrv returns a new instance of projects service.
func (s srvInitializer) InitProjectsSrv() (projects.Service, error) {
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
