package importcmd

import (
	"github.com/Smartling/smartling-cli/services/files"

	"github.com/spf13/cobra"
)

func NewImportCmd(s files.Service) *cobra.Command {
	var (
		published       bool
		postTranslation bool
		typ             string
		overwrite       bool
	)

	importCmd := &cobra.Command{
		Use:   "import <uri> <file> <locale>",
		Short: "Imports translations for given original file URI with.",
		Long:  `Imports translations for given original file URI with.`,
		Run: func(cmd *cobra.Command, args []string) {
			uri := args[0]
			file := args[1]
			locale := args[2]

			params := files.ImportParams{
				URI:             uri,
				File:            file,
				Locale:          locale,
				FileType:        typ,
				PostTranslation: postTranslation,
				Overwrite:       overwrite,
			}
			err := s.RunImport(params)
			if err != nil {
				// TODO log it
			}
		},
	}

	importCmd.Flags().BoolVar(&published, "published", false, "Translated content will be published.")
	importCmd.Flags().BoolVar(&postTranslation, "post-translation", false, `Translated content will be imported into first step of translation. If there are none, it will be published.`)
	importCmd.Flags().StringVar(&typ, "type", "", "Specify file type. If option is not given, file type will be deduced from extension.")
	importCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite any existing translations.")

	return importCmd
}
