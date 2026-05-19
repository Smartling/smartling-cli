package info

import (
	"os"

	projectscmd "github.com/Smartling/smartling-cli/cmd/projects"
	output "github.com/Smartling/smartling-cli/output/projects"
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
		Run: func(cmd *cobra.Command, _ []string) {
			ctx := cmd.Context()
			s, err := initializer.InitProjectsSrv(ctx)
			if err != nil {
				rlog.Errorf("failed to get project service: %s", err)
				os.Exit(1)
			}
			infoOutput, err := s.RunInfo(ctx)
			if err != nil {
				rlog.Errorf("failed to run info: %s", err)
				os.Exit(1)
			}
			if err := output.RenderTable(infoOutput); err != nil {
				rlog.Errorf("failed to render info output: %s", err)
				os.Exit(1)
			}
		},
	}
	return infoCmd
}
