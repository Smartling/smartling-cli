package importcmd

import (
	"os"

	filescmd "github.com/Smartling/smartling-cli/cmd/files"
	"github.com/Smartling/smartling-cli/services/files"
	"github.com/Smartling/smartling-cli/services/helpers/help"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/spf13/cobra"
)

// NewImportCmd creates a new command to import translations.
func NewImportCmd(initializer filescmd.SrvInitializer) *cobra.Command {
	var (
		published       bool
		postTranslation bool
		fileType        string
		overwrite       bool
	)

	importCmd := &cobra.Command{
		Use:   "import <uri> <file> <locale>",
		Short: "Imports translations for given original file URI with.",
		Long: `smartling-cli files import — import file translations.

Import pre-existent file translations into Smartling. Note, that
original file should be pushed prior file translations are imported.

Either --published or --post-translation should present to specify state
of imported translation.  Value indicates the workflow state to import the
translations into. Content will be imported into the language's default
workflow.

--overwrite option can be used to replace existent translations.

Available options:
  --published
    The translated content is published.

  --post-translation
   The translated content is imported into the first step after translation
   If there are none, it will be published.

  --overwrite
    Overwrite existing translations.
` + help.AuthenticationOptions,
		Run: func(_ *cobra.Command, args []string) {
			var (
				uri    string
				file   string
				locale string
			)
			if len(args) > 0 {
				uri = args[0]
			}
			if len(args) > 1 {
				file = args[1]
			}
			if len(args) > 2 {
				locale = args[2]
			}

			s, err := initializer.InitFilesSrv()
			if err != nil {
				rlog.Errorf("failed to get files service: %s", err)
				os.Exit(1)
			}

			params := files.ImportParams{
				URI:             uri,
				File:            file,
				Locale:          locale,
				FileType:        fileType,
				PostTranslation: postTranslation,
				Overwrite:       overwrite,
			}
			err = s.RunImport(params)
			if err != nil {
				rlog.Errorf("failed to run import: %s", err)
				os.Exit(1)
			}
		},
	}

	importCmd.Flags().BoolVar(&published, "published", false, "Translated content will be published.")
	importCmd.Flags().BoolVar(&postTranslation, "post-translation", false, `Translated content will be imported into first step of translation. If there are none, it will be published.`)
	importCmd.Flags().StringVar(&fileType, "type", "", "Specify file type. If option is not given, file type will be deduced from extension.")
	importCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite any existing translations.")

	return importCmd
}
