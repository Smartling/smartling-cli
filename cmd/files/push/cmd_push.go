package push

import (
	filescmd "github.com/Smartling/smartling-cli/cmd/files"
	"github.com/Smartling/smartling-cli/services/files"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/spf13/cobra"
)

// NewPushCmd creates a new command to upload files to the Smartling platform.
func NewPushCmd(initializer filescmd.SrvInitializer) *cobra.Command {
	var (
		authorize  bool
		locales    []string
		branch     string
		fileType   string
		directory  string
		directives []string
	)

	pushCmd := &cobra.Command{
		Use:   "push <file> <uri>",
		Short: "Uploads specified file into Smartling platform.",
		Long:  `Uploads specified file into Smartling platform.`,
		Run: func(_ *cobra.Command, args []string) {
			var (
				file string
				uri  string
			)
			if len(args) > 0 {
				file = args[0]
			}
			if len(args) > 1 {
				uri = args[1]
			}

			s, err := initializer.InitFilesSrv()
			if err != nil {
				rlog.Errorf("failed to get files service: %s", err)
				return
			}

			p := files.PushParams{
				URI:        uri,
				File:       file,
				Branch:     branch,
				Locales:    locales,
				Authorize:  authorize,
				Directory:  directory,
				FileType:   fileType,
				Directives: directives,
			}

			if err := s.RunPush(p); err != nil {
				rlog.Errorf("failed to run push: %s", err)
				return
			}
		},
	}

	pushCmd.Flags().BoolVarP(&authorize, "authorize", "z", false, `Automatically authorize all locales in specified file. Incompatible with -l option.`)
	pushCmd.Flags().StringSliceVarP(&locales, "locales", "l", []string{}, `Authorize only specified locales.`)
	pushCmd.Flags().StringVarP(&branch, "branch", "b", "", `Prepend specified text to the file uri.`)
	pushCmd.Flags().StringVarP(&fileType, "type", "t", "", `Specifies file type which will be used instead of automatically deduced from extension.`)
	pushCmd.Flags().StringSliceVarP(&directives, "directive", "r", []string{}, `Specifies one or more directives to use in push request.`)

	return pushCmd
}
