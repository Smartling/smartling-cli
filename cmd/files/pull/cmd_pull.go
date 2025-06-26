package pull

import (
	"os"

	filescmd "github.com/Smartling/smartling-cli/cmd/files"
	"github.com/Smartling/smartling-cli/services/files"
	"github.com/Smartling/smartling-cli/services/helpers/format"
	"github.com/Smartling/smartling-cli/services/helpers/help"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/spf13/cobra"
)

var (
	uri        string
	source     bool
	progress   string
	retrieve   string
	directory  string
	formatPath string
	locales    []string
)

// NewPullCmd creates a new command to pull files.
func NewPullCmd(initializer filescmd.SrvInitializer) *cobra.Command {
	pullCmd := &cobra.Command{
		Use:   "pull <uri>",
		Short: "Pulls specified files from server.",
		Long: `smartling-cli files pull — downloads translated files from project.

Downloads files from specified project into local directory.

It's possible to download only specific files by file mask, to download source
files with translations, to download file to specific directory or to download
specific locales only.

If special value of "-" is specified as <uri>, then program will expect
to read files list from stdin:

  cat files.txt | smartling-cli files pull -

<uri> ` + help.GlobPattern + `

If --locale flag is not specified, all available locales are downloaded. To
see available locales, use "status" command.

To download files into subdirectory, use --directory option and specify
directory name you want to download into.

To download source file as well as translated files specify --source option.

Files will be downloaded and stored under names used while upload (e.g. File
URI). While downloading translated file suffix "_<locale>" will be appended to
file name before extension. To override file format name, use --format option.
` + help.FormatOption + `
Following variables are available:

  > .FileURI — full file URI in Smartling system;
  > .Locale — locale ID for translated file and empty for source file;


Available options:
  -p --project <project>
    Specify project to use.

  --source
    Download source files along with translated files.

  —d ——directory <dir>
    Download files into specified directory.

  --format <format>
    Specify format for download file nmae.

  --progress <percents>
    Specify minimum of translation progress in percents.
	By default that filter does not apply.

  --retrieve <type>
    Retrieval type according to API specs:
    > pending — returns any translations, including non-published ones);
    > published — returns only published translations;
    > pseudo — returns modified version of original text with certain
               characters transformed;
    > contextMatchingInstrumented — to use with Chrome Context Capture;
` + help.AuthenticationOptions,
		Run: func(_ *cobra.Command, args []string) {
			if len(args) > 0 {
				uri = args[0]
			}

			s, err := initializer.InitFilesSrv()
			if err != nil {
				rlog.Errorf("failed to get files service: %s", err)
				os.Exit(1)
			}

			params := files.PullParams{
				URI:       uri,
				Format:    formatPath,
				Directory: directory,
				Source:    source,
				Locales:   locales,
				Progress:  progress,
				Retrieve:  retrieve,
			}
			err = s.RunPull(params)
			if err != nil {
				rlog.Errorf("failed to run pull: %s", err)
				os.Exit(1)
			}
		},
	}

	pullCmd.Flags().BoolVar(&source, "source", false, `Pulls source file as well.`)
	pullCmd.Flags().StringVar(&progress, "progress", "", `Pulls only translations that are at least specified percent of work complete.`)
	pullCmd.Flags().StringVar(&retrieve, "retrieve", "", `Retrieval type: pending, published, pseudo or contextMatchingInstrumented.`)
	pullCmd.Flags().StringVarP(&directory, "directory", "d", ".", `Download all files to specified directory.`)
	pullCmd.Flags().StringArrayVarP(&locales, "locale", "l", []string{}, `Authorize only specified locales.`)
	pullCmd.Flags().StringVar(&formatPath, "format", "", `Can be used to format path to downloaded files.
                           Note, that single file can be translated in
                           different locales, so format should include locale
                           to create several file paths.
                           Default: `+format.DefaultFilePullFormat)

	return pullCmd
}
