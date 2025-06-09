package initialize

import (
	"github.com/Smartling/smartling-cli/cmd"
	"github.com/Smartling/smartling-cli/services/init"

	"github.com/spf13/cobra"
)

var (
	dryRun bool
)

func NewInitCmd() *cobra.Command {
	initCmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Prepares project to work with Smartling",
		Long: `Prepares project to work with Smartling,
essentially, assisting user in creating
configuration file.`,
		Run: func(cmd *cobra.Command, args []string) {
			s, err := GetService()
			if err != nil {
				// TODO log it
			}
			err = s.RunInit(dryRun)
			if err != nil {
				// TODO log it
			}
		},
	}
	initCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Do not actually write file, just output it on stdout.")

	return initCmd
}

func GetService() (*initialize.Service, error) {
	client, err := cmd.Client()
	if err != nil {
		return nil, err
	}
	cnf, err := cmd.Config()
	if err != nil {
		return nil, err
	}
	srv := initialize.NewService(&client, cnf, cmd.CLIClientConfig())
	return srv, nil
}
