package list

import (
	filescmd "github.com/Smartling/smartling-cli/cmd/files"
	"github.com/Smartling/smartling-cli/services/helpers/format"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/spf13/cobra"
)

var (
	formatType string
	short      bool
)

// NewListCmd creates a new command to list files.
func NewListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list <uri>",
		Short: "Lists files from specified project.",
		Long:  `Lists files from specified project.`,
		Run: func(_ *cobra.Command, args []string) {
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

			err = s.RunList(formatType, short, uri)
			if err != nil {
				rlog.Errorf("failed to run list: %s", err)
				return
			}
		},
	}
	listCmd.Flags().BoolVarP(&short, "short", "s", false, "Display only project IDs.")
	listCmd.Flags().StringVar(&formatType, "format", "", `Can be used to format path to downloaded files.
                           Note, that single file can be translated in
                           different locales, so format should include locale
                           to create several file paths.
                           Default: `+format.DefaultFilesListFormat)
	return listCmd
}
