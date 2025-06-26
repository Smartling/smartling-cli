package initialize

import (
	"os"

	rootcmd "github.com/Smartling/smartling-cli/cmd"
	"github.com/Smartling/smartling-cli/services/helpers/help"
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
		Long: `smartling-cli init — create config file interactively.

Walk down common config file parameters and fill them through dialog.

Init process will inspect if config file already exists and if it is, it will
be loaded as default values, so init can be used sequentially without config
is lost.

Options like --user, --secret, --account and --project can be used to specify
config values prior dialog:

  smartling-cli init --user=your_user_id

Also, --dry-run option can be used to just look at resulting config without
overwritting anything:

  smartling-cli init --dry-run

By default, smartling.yml file in the local directory will be used as target
config file, but it can be overriden by using --config option:

  smartling-cli init --config=/path/to/project/smartling.yml


Available options:
  -c --config <file>
    Specify config file to operate on. Default: smartling.yml

  --dry-run
    Do not overwrite config file, only output to stdout.

Default config values can be passed via following options:` +
			help.AuthenticationOptions + `
  -p --project <project>
    Specify default project.`,
		Run: func(_ *cobra.Command, _ []string) {
			s, err := srvInitializer.InitSrv()
			if err != nil {
				rlog.Errorf("failed to get init service: %s", err)
				os.Exit(1)
			}
			err = s.RunInit(dryRun)
			if err != nil {
				rlog.Errorf("failed to run init: %s", err)
				os.Exit(1)
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
	client, err := rootcmd.Client()
	if err != nil {
		return nil, err
	}
	cnf, err := rootcmd.Config()
	if err != nil {
		return nil, err
	}
	srv := initialize.NewService(&client, cnf)
	smClient := sdk.NewHttpAPIClient(cnf.UserID, cnf.Secret)
	srv := initialize.NewService(smClient, cnf)
	return srv, nil
}
