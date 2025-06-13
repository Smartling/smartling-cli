package delete

import (
	"github.com/Smartling/smartling-cli/cmd/files"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/spf13/cobra"
)

var uri string

// NewDeleteCmd creates a new command to delete files.
func NewDeleteCmd(initializer files.SrvInitializer) *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:   "delete <uri>",
		Short: "Deletes given file from Smartling.",
		Long:  `Deletes given file from Smartling. This operation can not be undone, so use with care.`,
		Run: func(_ *cobra.Command, args []string) {
			if len(args) > 0 {
				uri = args[0]
			}

			s, err := initializer.InitFilesSrv()
			if err != nil {
				rlog.Errorf("failed to get files service: %s", err)
				return
			}

			err = s.RunDelete(uri)
			if err != nil {
				rlog.Errorf("failed to run delete: %s", err)
				return
			}
		},
	}

	return deleteCmd
}
