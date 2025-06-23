package initialize

import (
	sdk "github.com/Smartling/api-sdk-go"
	rootcmd "github.com/Smartling/smartling-cli/cmd"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	"github.com/Smartling/smartling-cli/services/init"

	"github.com/spf13/cobra"
)

var (
	dryRun bool
)

// NewInitCmd creates a new command to initialize the Smartling CLI.
func NewInitCmd(srvInitializer SrvInitializer) *cobra.Command {
	initCmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Prepares project to work with Smartling",
		Long: `Prepares project to work with Smartling,
essentially, assisting user in creating
configuration file.`,
		Run: func(_ *cobra.Command, _ []string) {
			s, err := srvInitializer.InitSrv()
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

// SrvInitializer defines service initializer
type SrvInitializer interface {
	InitSrv() (initialize.Service, error)
}

// NewSrvInitializer returns new SrvInitializer implementation
func NewSrvInitializer() SrvInitializer {
	return srvInitializer{}
}

type srvInitializer struct{}

// Init initializes and returns a new instance of the init service.
func (s srvInitializer) InitSrv() (initialize.Service, error) {
	/*client, err := rootcmd.Client()
	if err != nil {
		return nil, err
	}
	cnf, err := rootcmd.Config()
	if err != nil {
		return nil, err
	}
	smClient := sdk.NewClient(cnf.UserID, cnf.Secret)
	srv := initialize.NewService(smClient, cnf) */
	return initialize.Service{}, nil
}
