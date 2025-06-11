package initialize

import (
	rootcmd "github.com/Smartling/smartling-cli/cmd"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	"github.com/Smartling/smartling-cli/services/init"

	"github.com/spf13/cobra"
)

var (
	dryRun bool
)

// NewInitCmd creates a new command to initialize the Smartling CLI.
func NewInitCmd() *cobra.Command {
	initCmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Prepares project to work with Smartling",
		Long: `Prepares project to work with Smartling,
essentially, assisting user in creating
configuration file.`,
		Run: func(_ *cobra.Command, _ []string) {
			s, err := GetService()
			if err != nil {
				rlog.Errorf("failed to get init service: %s", err)
				return
			}
			err = s.RunInit(dryRun)
			if err != nil {
				rlog.Errorf("failed to run init: %s", err)
				return
			}
		},
	}
	initCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Do not actually write file, just output it on stdout.")

	return initCmd
}

// GetService initializes and returns a new instance of the init service.
func GetService() (*initialize.Service, error) {
	client, err := rootcmd.Client()
	if err != nil {
		return nil, err
	}
	cnf, err := rootcmd.Config()
	if err != nil {
		return nil, err
	}
	srv := initialize.NewService(&client, cnf)
	return srv, nil
}
