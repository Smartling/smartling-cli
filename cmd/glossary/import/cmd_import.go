package glimport

import (
	"fmt"

	glossarycmd "github.com/Smartling/smartling-cli/cmd/glossary"

	"github.com/Smartling/smartling-cli/output"
	"github.com/spf13/cobra"
)

// NewImportCmd builds the `glossary import` command.
func NewImportCmd(initializer glossarycmd.SrvInitializer) *cobra.Command {
	var (
		published       bool
		postTranslation bool
		fileType        string
		overwrite       bool
	)

	importCmd := &cobra.Command{
		Use:   "import <glossaryUID|glossaryName> [inFile]",
		Short: "Glossary import process",
		Long:  `Glossary import process`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			var glossaryUIDOrName, inFile string
			if len(args) > 0 {
				glossaryUIDOrName = args[0]
			}
			if len(args) > 1 {
				inFile = args[1]
			}

			fileConfig, err := glossarycmd.BindFileConfig(cmd)
			if err != nil {
				return err
			}
			params, err := resolveParams(cmd, fileConfig, glossaryUIDOrName, inFile)
			if err != nil {
				return fmt.Errorf("failed to resolve import params: %w", err)
			}

			format, err := cmd.Parent().PersistentFlags().GetString("output")
			if err != nil {
				return err
			}
			outputParams := output.Params{Format: format}
			return run(ctx, initializer, params, outputParams)
		},
	}

	importCmd.Flags().BoolVar(&published, "published", false, "Translated content will be published.")
	importCmd.Flags().BoolVar(&postTranslation, "post-translation", false, `Translated content will be imported into first step of translation. If there are none, it will be published.`)
	importCmd.Flags().StringVar(&fileType, "type", "", "Specify file type. If option is not given, file type will be deduced from extension.")
	importCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite any existing translations.")

	return importCmd
}
