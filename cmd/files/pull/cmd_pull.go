package pull

import (
	filescmd "github.com/Smartling/smartling-cli/cmd/files"
	"github.com/Smartling/smartling-cli/services/files"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/spf13/cobra"
)

var (
	uri       string
	source    bool
	progress  string
	retrieve  string
	directory string
	format    string
)

func NewPullCmd() *cobra.Command {
	pullCmd := &cobra.Command{
		Use:   "pull <uri>",
		Short: "Pulls specified files from server.",
		Long:  `Pulls specified files from server.`,
		Run: func(cmd *cobra.Command, args []string) {
			uri = args[0]

			s, err := filescmd.InitFilesSrv()
			if err != nil {
				rlog.Errorf("failed to get files service: %s", err)
				return
			}

			params := files.PullParams{
				URI:       uri,
				Format:    format,
				Directory: directory,
				Source:    source,
				Locales:   nil,
				Progress:  "",
				Retrieve:  "",
			}
			err = s.RunPull(params)
			if err != nil {
				rlog.Errorf("failed to run pull: %s", err)
				return
			}
		},
	}

	pullCmd.Flags().BoolVar(&source, "source", false, `Pulls source file as well.`)
	pullCmd.Flags().StringVar(&progress, "progress", "", `Pulls only translations that are at least specified percent of work complete.`)
	pullCmd.Flags().StringVar(&retrieve, "retrieve", "", `Retrieval type: pending, published, pseudo or contextMatchingInstrumented.`)
	pullCmd.Flags().StringVarP(&directory, "directory", "d", "", `Download all files to specified directory.`)
	pullCmd.Flags().StringVar(&format, "format", "", `Can be used to format path to downloaded files.
                           Note, that single file can be translated in
                           different locales, so format should include locale
                           to create several file paths.
                           [default: $FILE_PULL_FORMAT]`)

	return pullCmd
}
