package list

import (
	"os"

	filescmd "github.com/Smartling/smartling-cli/cmd/files"
	"github.com/Smartling/smartling-cli/services/helpers/format"
	"github.com/Smartling/smartling-cli/services/helpers/help"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/spf13/cobra"
)

var (
	formatType string
	short      bool
)

// NewListCmd creates a new command to list files.
func NewListCmd(initializer filescmd.SrvInitializer) *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list <uri>",
		Short: "Lists files from specified project.",
		Long: `smartling-cli files list — list files from project.

Lists all files from project or only files which matches specified uri.

Note, that by default listing is limited to 500 items in Smartling API,
so several requests may be needed to obtain full file list, which will
take some time.

List command will output following fields in tabular format by default:

  > File URI;
  > Last uploaded date;
  > File Type;
` + help.FormatOption + `
Following variables are available:

  > .FileURI — full file URI in Smartling system;
  > .FileType — internal Smartling file type;
  > .LastUploaded — timestamp when file was last uploaded;
  > .HasInstructions — true/false if file has translation instructions;

` + "`<uri>` " + help.GlobPattern + `


Available options:
  -p --project <project>
    Specify project to use.

  -s --short
    List only file URIs.

  --format <format>
    Override default listing format.
` + help.AuthenticationOptions,
		Example: `
# List project files

smartling-cli files list
`,
		Run: func(_ *cobra.Command, args []string) {
			var uri string
			if len(args) > 0 {
				uri = args[0]
			}

			s, err := initializer.InitFilesSrv()
			if err != nil {
				rlog.Errorf("failed to get files service: %s", err)
				os.Exit(1)
			}

			err = s.RunList(formatType, short, uri)
			if err != nil {
				rlog.Errorf("failed to run list: %s", err)
				os.Exit(1)
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
