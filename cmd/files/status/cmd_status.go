package status

import (
	filescmd "github.com/Smartling/smartling-cli/cmd/files"
	"github.com/Smartling/smartling-cli/services/files"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/spf13/cobra"
)

var (
	format    string
	directory string
)

func NewStatusCmd() *cobra.Command {
	statusCmd := &cobra.Command{
		Use:   "status <uri>",
		Short: "Shows file translation status.",
		Long:  `Shows file translation status.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				rlog.Error("missing required argument `<uri>`")
				return
			}
			uri := args[0]

			s, err := filescmd.InitFilesSrv()
			if err != nil {
				rlog.Errorf("failed to get files service: %s", err)
				return
			}

			p := files.StatusParams{
				URI:       uri,
				Directory: directory,
				Format:    format,
			}

			err = s.RunStatus(p)
			if err != nil {
				rlog.Errorf("failed to run status: %s", err)
				return
			}
		},
	}

	statusCmd.Flags().StringVar(&format, "format", "", `Specifies format to use for file status output. [default: $FILE_STATUS_FORMAT]`)
	statusCmd.Flags().StringVar(&directory, "directory", "", `Use another directory as reference to check for local files.`)

	return statusCmd
}
