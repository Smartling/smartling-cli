package info

import (
	"os"

	projectscmd "github.com/Smartling/smartling-cli/cmd/projects"
	"github.com/Smartling/smartling-cli/services/helpers/help"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/spf13/cobra"
)

// NewInfoCmd creates a new command to get project details.
func NewInfoCmd(initializer projectscmd.SrvInitializer) *cobra.Command {
	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Get project details about specific project.",
		Long: `smartling-cli projects info — show detailed project info.

Displays detailed information for specific project.

Project should be specified either in config or via --project option.

Available options:` + help.AuthenticationOptions,
		Example: `
# View project information

  smartling-cli projects info

`,
		Run: func(_ *cobra.Command, _ []string) {
			s, err := initializer.InitProjectsSrv()
			if err != nil {
				rlog.Errorf("failed to get project service: %s", err)
				os.Exit(1)
			}
			err = s.RunInfo()
			if err != nil {
				rlog.Errorf("failed to run info: %s", err)
				os.Exit(1)
			}
		},
	}
	return infoCmd
}
