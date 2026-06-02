package gllist

import (
	"fmt"

	glossariescmd "github.com/Smartling/smartling-cli/cmd/glossaries"
	"github.com/Smartling/smartling-cli/output"

	"github.com/spf13/cobra"
)

const nameFlag = "name"

// NewListCmd builds the `glossaries list` command.
func NewListCmd(initializer glossariescmd.SrvInitializer) *cobra.Command {
	var name string

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List glossaries in the current account",
		Long: `List glossaries under the current account. When --name is supplied it is
passed to the API as the "glossaryName" filter; otherwise all glossaries are
returned.`,
		Example: `
# List glossaries

  smartling-cli glossaries list

# List glossaries with name "CLI"

  smartling-cli glossaries list --name "CLI"

# List glossaries with name "CLI" and output in table

  smartling-cli glossaries list --name "CLI" --output table
`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			params, err := resolveParams(cmd, name)
			if err != nil {
				return fmt.Errorf("failed to resolve list params: %w", err)
			}

			format, err := cmd.Parent().PersistentFlags().GetString("output")
			if err != nil {
				return err
			}
			outputParams := output.Params{Format: format}
			return run(ctx, initializer, params, outputParams)
		},
	}

	listCmd.Flags().StringVar(&name, nameFlag, "", "Filter glossaries by name (maps to the API `glossaryName` query parameter).")

	return listCmd
}
