package push

import (
	"github.com/spf13/cobra"
)

var (
	file      string
	uri       string
	authorize bool
	locale    string
	branch    string
	typ       string
	directive string
)

func NewPushCmd() *cobra.Command {
	pushCmd := &cobra.Command{
		Use:   "push <file> <uri>",
		Short: "Uploads specified file into Smartling platform.",
		Long:  `Uploads specified file into Smartling platform.`,
		Run: func(cmd *cobra.Command, args []string) {
			file = args[0]
			uri = args[1]
		},
	}

	pushCmd.Flags().BoolVarP(&authorize, "authorize", "-z", false, `Automatically authorize all locales in specified file. Incompatible with -l option.`)
	pushCmd.Flags().StringVarP(&locale, "locale", "-l", "", `Authorize only specified locales.`)
	pushCmd.Flags().StringVarP(&branch, "branch", "-b", "", `Prepend specified text to the file uri.`)
	pushCmd.Flags().StringVarP(&typ, "type", "-t", "", `Specifies file type which will be used instead of automatically deduced from extension.`)
	pushCmd.Flags().StringVarP(&directive, "directive", "-r", "", `Specifies one or more directives to use in push request.`)

	return pushCmd
}
