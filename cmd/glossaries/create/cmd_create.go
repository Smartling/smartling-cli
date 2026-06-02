package glcreate

import (
	"fmt"

	glossariescmd "github.com/Smartling/smartling-cli/cmd/glossaries"
	"github.com/Smartling/smartling-cli/output"

	"github.com/spf13/cobra"
)

// Flag names accepted by `glossaries create`. Each flag maps to a field on the
// Smartling Glossary Create API request body
// (https://api-reference.smartling.com/#tag/Glossary-API/operation/createGlossary).
const (
	descriptionFlag      = "description"
	verificationModeFlag = "verification-mode"
	localeFlag           = "locale"
	fallbackLocaleFlag   = "fallback-locale"
)

// NewCreateCmd builds the `glossaries create` command.
func NewCreateCmd(initializer glossariescmd.SrvInitializer) *cobra.Command {
	var (
		description      string
		verificationMode bool
		locales          []string
		fallbackLocales  []string
	)

	createCmd := &cobra.Command{
		Use:   "create <glossaryName>",
		Short: "Glossary create",
		Long: `Create a new glossary in the current account.

The glossary name is taken from the positional argument and maps to the API
request's "glossaryName" field. All other API fields are exposed as flags.

Fallback locales are repeatable and use the format
  --fallback-locale <fallbackLocaleId>:<localeId>[,<localeId>...]
e.g. --fallback-locale es:es-MX,es-AR.`,
		Example: `
# Create a glossary named "CLI glossary" with two locales

  smartling-cli glossaries create "CLI glossary" --locale es-ES --locale fr-FR
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			glossaryName := args[0]

			fileConfig, err := glossariescmd.BindFileConfig(cmd)
			if err != nil {
				return err
			}
			params, err := resolveParams(cmd, fileConfig, glossaryName)
			if err != nil {
				return fmt.Errorf("failed to resolve create params: %w", err)
			}

			format, err := cmd.Flags().GetString("output")
			if err != nil {
				return err
			}
			outputParams := output.Params{Format: format}
			return run(ctx, initializer, params, outputParams)
		},
	}

	f := createCmd.Flags()
	f.StringVar(&description, descriptionFlag, "", "Glossary description.")
	f.BoolVar(&verificationMode, verificationModeFlag, false, "Enable verification mode for the glossary.")
	f.StringArrayVar(&locales, localeFlag, nil, "Target locale ID (repeatable → localeIds).")
	f.StringArrayVar(&fallbackLocales, fallbackLocaleFlag, nil, "Fallback locale mapping (repeatable). Format: '<fallbackLocaleId>:<localeId>[,<localeId>...]'.")

	return createCmd
}
