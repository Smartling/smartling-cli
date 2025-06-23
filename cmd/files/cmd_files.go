package files

import (
	"github.com/Smartling/smartling-cli/cmd"
	"github.com/Smartling/smartling-cli/services/files"

	"github.com/spf13/cobra"
)

// NewFilesCmd creates a new command to access various files sub-commands.
func NewFilesCmd() *cobra.Command {
	filesCmd := &cobra.Command{
		Use:     "files",
		Aliases: []string{"f"},
		Short:   "Used to access various files sub-commands.",
		Long:    `Used to access various files sub-commands.`,
		Run: func(_ *cobra.Command, _ []string) {

		},
	}

	return filesCmd
}

// SrvInitializer defines files service initializer
type SrvInitializer interface {
	InitFilesSrv() (files.Service, error)
}

// NewSrvInitializer returns new SrvInitializer implementation
func NewSrvInitializer() SrvInitializer {
	return srvInitializer{}
}

type srvInitializer struct{}

// InitFilesSrv initializes `files` service with the client and configuration.
func (i srvInitializer) InitFilesSrv() (files.Service, error) {
	client, err := cmd.Client()
	if err != nil {
		return nil, err
	}
	cnf, err := cmd.Config()
	if err != nil {
		return nil, err
	}
	fileConfig, err := cnf.GetFileConfig(cmd.ConfigFile())
	if err != nil {
		return nil, err
	}
	srv := files.NewService(client, cnf, fileConfig)

	return srv, nil
}
