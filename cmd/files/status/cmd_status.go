package status

import (
	filescmd "github.com/Smartling/smartling-cli/cmd/files"
	"github.com/Smartling/smartling-cli/services/files"

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
			uri := args[0]

			s, err := filescmd.GetService()
			if err != nil {
				// TODO log it
			}

			p := files.StatusParams{
				URI:       uri,
				Directory: directory,
				Format:    format,
			}

			err = s.RunStatus(p)
			if err != nil {
				// TODO log it
			}
		},
	}

	statusCmd.Flags().StringVar(&format, "format", "", `Specifies format to use for file status output. [default: $FILE_STATUS_FORMAT]`)
	statusCmd.Flags().StringVar(&directory, "directory", "", `Use another directory as reference to check for local files.`)

	return statusCmd
}
