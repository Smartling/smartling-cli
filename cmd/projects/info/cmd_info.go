package info

import (
	projectscmd "github.com/Smartling/smartling-cli/cmd/projects"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/spf13/cobra"
)

// NewInfoCmd creates a new command to get project details.
func NewInfoCmd() *cobra.Command {
	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Get project details about specific project.",
		Long:  `Get project details about specific project.`,
		Run: func(_ *cobra.Command, _ []string) {
			s, err := projectscmd.GetService()
			if err != nil {
				rlog.Errorf("failed to get project service: %s", err)
				return
			}
			err = s.RunInfo()
			if err != nil {
				rlog.Errorf("failed to run info: %s", err)
				return
			}
		},
	}
	return infoCmd
}
