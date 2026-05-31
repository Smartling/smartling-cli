package glimport

import (
	"fmt"

	glossarycmd "github.com/Smartling/smartling-cli/cmd/glossary"
	"github.com/Smartling/smartling-cli/output"

	"github.com/spf13/cobra"
)

// Flag names accepted by `glossary import`. Each flag maps to a field on the
// Smartling Glossary Import API request body
// (https://api-reference.smartling.com/#tag/Glossary-API/operation/importGlossary).
const (
	archiveModeFlag = "archive-mode"
	mediaTypeFlag   = "media-type"
)

// NewImportCmd builds the `glossary import` command.
func NewImportCmd(initializer glossarycmd.SrvInitializer) *cobra.Command {
	var (
		archiveMode bool
		mediaType   string
	)

	importCmd := &cobra.Command{
		Use:   "import <glossaryUID|glossaryName> <inFile>",
		Short: "Glossary import process",
		Long: `Import a glossary file (CSV/XLSX/TBX) into an existing glossary.

The first argument is the glossary UID or name; the second is the local file
to upload. The file extension is used to derive a media type unless
--media-type is supplied explicitly.`,
		Example: `
# Import a CSV file into a glossary (media type derived from .csv extension)

  smartling-cli glossary import "CLI glossary" ./terms.csv

# Import a TBX file and archive entries that are missing from the file

  smartling-cli glossary import "CLI glossary" ./terms.tbx --archive-mode

# Import by glossary UID instead of name

  smartling-cli glossary import 03e37fc4-842b-4cdb-b19d-79b13d6edbd2 ./terms.xlsx

# Override the auto-derived media type

  smartling-cli glossary import "CLI glossary" ./terms.dat --media-type text/csv
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			glossaryUIDOrName, inFile := args[0], args[1]

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

	importCmd.Flags().BoolVar(&archiveMode, archiveModeFlag, false, "Archive entries that are missing from the imported file.")
	importCmd.Flags().StringVar(&mediaType, mediaTypeFlag, "", `Override the media type. Must be one of "text/csv", "text/xml", or "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet". By default derived from the file extension.`)

	return importCmd
}
