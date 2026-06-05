package findbystrings

import (
	"fmt"
	"os"

	jobscmd "github.com/Smartling/smartling-cli/cmd/jobs"
	"github.com/Smartling/smartling-cli/output"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/spf13/cobra"
)

const (
	hashcodeFlag = "hashcode"
	localeFlag   = "locale"
)

// NewFindByStringsCmd builds the `jobs find-by-strings` command.
func NewFindByStringsCmd(initializer jobscmd.SrvInitializer) *cobra.Command {
	var (
		hashcodes []string
		locales   []string
	)
	findCmd := &cobra.Command{
		Use:   "find-by-strings",
		Short: "Find jobs that contain specific strings in specific locales.",
		Long: `Find the translation jobs that contain the given strings (by hashcode) in
the given locales. Results are reported as one row per hashcode+locale match.`,
		Example: `
# Find jobs containing two strings

  smartling-cli jobs find-by-strings --hashcode h1 --hashcode h2

# Restrict the search to specific locales

  smartling-cli jobs find-by-strings --hashcode h1 --locale fr-FR --locale de-DE
`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			params, err := resolveParams(hashcodes, locales)
			if err != nil {
				return fmt.Errorf("failed to resolve find-by-strings params: %w", err)
			}

			format, err := cmd.Flags().GetString("output")
			if err != nil {
				return err
			}
			return run(ctx, initializer, params, output.Params{Format: format})
		},
	}

	findCmd.Flags().StringArrayVar(&hashcodes, hashcodeFlag, nil, "String hashcode to search for (repeatable, required).")
	findCmd.Flags().StringArrayVar(&locales, localeFlag, nil, "Locale to restrict the search to (repeatable; default all locales).")
	if err := findCmd.MarkFlagRequired(hashcodeFlag); err != nil {
		rlog.Errorf("failed to mark --%s required: %s", hashcodeFlag, err)
		os.Exit(1)
	}

	return findCmd
}
