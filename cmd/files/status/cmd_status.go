package status

import (
	filescmd "github.com/Smartling/smartling-cli/cmd/files"
	"github.com/Smartling/smartling-cli/services/files"
	"github.com/Smartling/smartling-cli/services/helpers/format"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/spf13/cobra"
)

var (
	formatType string
	directory  string
)

// NewStatusCmd creates a new command to show file translation status.
func NewStatusCmd(initializer filescmd.SrvInitializer) *cobra.Command {
	statusCmd := &cobra.Command{
		Use:   "status <uri>",
		Short: "Shows file translation status.",
		Long:  `Shows file translation status.`,
		Run: func(_ *cobra.Command, args []string) {
			var uri string
			if len(args) > 0 {
				uri = args[0]
			}

			s, err := initializer.InitFilesSrv()
			if err != nil {
				rlog.Errorf("failed to get files service: %s", err)
				return
			}

			p := files.StatusParams{
				URI:       uri,
				Directory: directory,
				Format:    formatType,
			}
			err = s.RunStatus(p)
			if err != nil {
				rlog.Errorf("failed to run status: %s", err)
				return
			}
		},
	}

	statusCmd.Flags().StringVar(&formatType, "format", "", `Specifies format to use for file status output. 
								Default: `+format.DefaultFileStatusFormat)
	statusCmd.Flags().StringVar(&directory, "directory", ".", `Use another directory as reference to check for local files.`)

	return statusCmd
}
