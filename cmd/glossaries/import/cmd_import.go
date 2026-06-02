package glimport

import (
	"fmt"

	glossariescmd "github.com/Smartling/smartling-cli/cmd/glossaries"
	"github.com/Smartling/smartling-cli/output"

	"github.com/spf13/cobra"
)

// Flag names accepted by `glossaries import`. Each flag maps to a field on the
// Smartling Glossary Import API request body
// (https://api-reference.smartling.com/#tag/Glossary-API/operation/importGlossary).
const (
	archiveModeFlag = "archive-mode"
	mediaTypeFlag   = "media-type"
)

// NewImportCmd builds the `glossaries import` command.
func NewImportCmd(initializer glossariescmd.SrvInitializer) *cobra.Command {
	var (
		archiveMode bool
		mediaType   string
	)

	importCmd := &cobra.Command{
		Use:   "import <glossaryUID|glossaryName> <inFile>",
		Short: "Glossary import process",
		Long: `Upload a CSV, XLSX, or TBX file into an existing glossary.

Validates and uploads the file, waits for the server to confirm the import,
then polls until the import reaches SUCCESSFUL or FAILED status. New entries
are created.`,
		Example: `
# Import a CSV file into a glossary (media type derived from .csv extension)

  smartling-cli glossaries import "CLI glossary" ./terms.csv

# Import a TBX file and archive entries that are missing from the file

  smartling-cli glossaries import "CLI glossary" ./terms.tbx --archive-mode

# Import by glossary UID instead of name

  smartling-cli glossaries import 03e37fc4-842b-4cdb-b19d-79b13d6edbd2 ./terms.xlsx

# Override the auto-derived media type

  smartling-cli glossaries import "CLI glossary" ./terms.dat --media-type text/csv
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			glossaryUIDOrName, inFile := args[0], args[1]

			fileConfig, err := glossariescmd.BindFileConfig(cmd)
			if err != nil {
				return err
			}
			params, err := resolveParams(cmd, fileConfig, glossaryUIDOrName, inFile)
			if err != nil {
				return fmt.Errorf("failed to resolve import params: %w", err)
			}

			format, err := cmd.Flags().GetString("output")
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
