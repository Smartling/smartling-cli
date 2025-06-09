package files

import (
	"github.com/Smartling/smartling-cli/cmd"
	"github.com/Smartling/smartling-cli/services/files"

	"github.com/spf13/cobra"
)

func NewFilesCmd() *cobra.Command {
	filesCmd := &cobra.Command{
		Use:     "files",
		Aliases: []string{"f"},
		Short:   "Used to access various files sub-commands.",
		Long:    `Used to access various files sub-commands.`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return filesCmd
}

func GetService() (*files.Service, error) {
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
	srv := files.NewService(&client, cnf, fileConfig)
	return srv, nil
}
