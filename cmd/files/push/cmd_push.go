package push

import (
	rootcmd "github.com/Smartling/smartling-cli/cmd"
	filescmd "github.com/Smartling/smartling-cli/cmd/files"
	"github.com/Smartling/smartling-cli/services/files"

	"github.com/spf13/cobra"
)

func NewPushCmd() *cobra.Command {
	var (
		file       string
		uri        string
		authorize  bool
		locales    []string
		branch     string
		typ        string
		directory  string
		directives []string
	)

	pushCmd := &cobra.Command{
		Use:   "push <file> <uri>",
		Short: "Uploads specified file into Smartling platform.",
		Long:  `Uploads specified file into Smartling platform.`,
		Run: func(cmd *cobra.Command, args []string) {
			file = args[0]
			uri = args[1]

			s, err := filescmd.InitFilesSrv()
			if err != nil {
				rootcmd.Logger().Errorf("failed to get files service: %s", err)
				return
			}

			p := files.PushParams{
				URI:        uri,
				File:       file,
				Branch:     branch,
				Locales:    locales,
				Authorize:  authorize,
				Directory:  directory,
				FileType:   typ,
				Directives: directives,
			}

			if err := s.RunPush(p); err != nil {
				rootcmd.Logger().Errorf("failed to run push: %s", err)
				return
			}
		},
	}

	pushCmd.Flags().BoolVarP(&authorize, "authorize", "z", false, `Automatically authorize all locales in specified file. Incompatible with -l option.`)
	pushCmd.Flags().StringSliceVarP(&locales, "locales", "l", []string{}, `Authorize only specified locales.`)
	pushCmd.Flags().StringVarP(&branch, "branch", "b", "", `Prepend specified text to the file uri.`)
	pushCmd.Flags().StringVarP(&typ, "type", "t", "", `Specifies file type which will be used instead of automatically deduced from extension.`)
	pushCmd.Flags().StringSliceVarP(&directives, "directives", "r", []string{}, `Specifies one or more directives to use in push request.`)

	return pushCmd
}
