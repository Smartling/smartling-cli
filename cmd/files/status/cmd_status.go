package status

import (
	filescmd "github.com/Smartling/smartling-cli/cmd/files"
	"github.com/Smartling/smartling-cli/services/files"
	"github.com/Smartling/smartling-cli/services/helpers/format"
	"github.com/Smartling/smartling-cli/services/helpers/help"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/spf13/cobra"
)

var (
	formatType string
	directory  string
)

// NewStatusCmd creates a new command to show file translation status.
func NewStatusCmd(initializer filescmd.SrvInitializer) *cobra.Command {
	statusCmd := &cobra.Command{
		Use:   "status <uri>",
		Short: "Shows file translation status.",
		Long: `smartling-cli files status — show files status from project.

Lists all files from project along with their translation progress into
different locales.

Status command will check, if files are missing locally or not.

Command will list projects from specified account in tabular format with
following information:

  > File URI
  > File Locale
  > File Status on Local System
  > Translation Progress
  > Strings Count
  > Words Count

If no <uri> is specified, all files will be listed.

To list files status from specific directory, --directory option can be used.

To override default file name format --format can be used.
` + help.FormatOption + `
Following variables are available:

  > .FileURI — full file URI in Smartling system;
  > .Locale — locale ID for translated file and empty for source file;

<uri> ` + help.GlobPattern + `


Available options:
  -p --project <project>
    Specify project to use.

  --directory <directory>
    Check files in specific directory instead of local directory.

  --format <format>
    Specify format for listing file names.
` + help.AuthenticationOptions,
		Run: func(_ *cobra.Command, args []string) {
			var uri string
			if len(args) > 0 {
				uri = args[0]
			}

			s, err := initializer.InitFilesSrv()
			if err != nil {
				rlog.Errorf("failed to get files service: %s", err)
				return
			}

			p := files.StatusParams{
				URI:       uri,
				Directory: directory,
				Format:    formatType,
			}
			err = s.RunStatus(p)
			if err != nil {
				rlog.Errorf("failed to run status: %s", err)
				return
			}
		},
	}

	statusCmd.Flags().StringVar(&formatType, "format", "", `Specifies format to use for file status output. 
								Default: `+format.DefaultFileStatusFormat)
	statusCmd.Flags().StringVar(&directory, "directory", ".", `Use another directory as reference to check for local files.`)

	return statusCmd
}
