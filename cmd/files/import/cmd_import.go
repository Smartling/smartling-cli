package importcmd

import (
	"github.com/spf13/cobra"
)

var (
	uri    string
	file   string
	locale string

	published       bool
	postTranslation bool
	typ             string
	overwrite       bool
)

func NewImportCmd() *cobra.Command {
	importCmd := &cobra.Command{
		Use:   "import <uri> <file> <locale>",
		Short: "Imports translations for given original file URI with.",
		Long:  `Imports translations for given original file URI with.`,
		Run: func(cmd *cobra.Command, args []string) {
			uri = args[0]
			file = args[1]
			locale = args[2]
		},
	}

	importCmd.Flags().BoolVar(&published, "published", false, "Translated content will be published.")
	importCmd.Flags().BoolVar(&postTranslation, "post-translation", false, `Translated content will be imported into first step of translation. If there are none, it will be published.`)
	importCmd.Flags().StringVar(&typ, "type", "", "Specify file type. If option is not given, file type will be deduced from extension.")
	importCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite any existing translations.")

	return importCmd
}
