package delete

import (
	"github.com/Smartling/smartling-cli/services/files"

	"github.com/spf13/cobra"
)

var uri string

func NewDeleteCmd(s *files.Service) *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:   "delete <uri>",
		Short: "Deletes given file from Smartling.",
		Long:  `Deletes given file from Smartling. This operation can not be undone, so use with care.`,
		Run: func(cmd *cobra.Command, args []string) {
			uri = args[0]

			err := s.RunDelete(uri)
			if err != nil {
				// TODO log it
			}
		},
	}

	return deleteCmd
}
