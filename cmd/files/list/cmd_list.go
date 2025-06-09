package list

import (
	"github.com/Smartling/smartling-cli/services/files"

	"github.com/spf13/cobra"
)

var (
	format string
	short  bool
)

func NewListCmd(s *files.Service) *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list <uri>",
		Short: "Lists files from specified project.",
		Long:  `Lists files from specified project.`,
		Run: func(cmd *cobra.Command, args []string) {
			uri := args[0]

			err := s.RunList(format, short, uri)
			if err != nil {
				// TODO log it
			}
		},
	}
	listCmd.Flags().BoolVarP(&short, "short", "s", false, "Display only project IDs.")
	listCmd.Flags().StringVar(&format, "format", "", `Can be used to format path to downloaded files.
                           Note, that single file can be translated in
                           different locales, so format should include locale
                           to create several file paths.
                           [default: $FILE_PULL_FORMAT]`)
	return listCmd
}
