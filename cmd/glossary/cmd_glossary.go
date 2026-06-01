package glossary

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Smartling/smartling-cli/cmd"

	"github.com/spf13/cobra"
)

const (
	outputFormatFlag = "output"
)

var (
	outputFormat   string
	allowedOutputs = []string{
		"table",
		"json",
		"simple",
	}
	joinedAllowedOutputs = strings.Join(allowedOutputs, ", ")
)

// NewGlossaryCmd  builds the `glossary` command.
func NewGlossaryCmd() *cobra.Command {
	glossaryCmd := &cobra.Command{
		Use:   "glossary",
		Short: "Manage Smartling glossaries",
		Long: `Manage Smartling glossaries from the command line.

A glossary is a list of the terms and phrases your brand uses in a specific
way, along with instructions on how translation resources should treat them -
definitions, part of speech, notes, term variations, and a "do not translate"
(DNT) flag. Glossaries give translators and machine translation a shared,
consistent understanding of your terminology across every locale.`,
		PersistentPreRunE: func(c *cobra.Command, args []string) error {
			if err := cmd.RunRootPersistentPreRun(c); err != nil {
				return err
			}
			if !slices.Contains(allowedOutputs, outputFormat) {
				return fmt.Errorf("invalid output: %s (allowed: %s)", outputFormat, joinedAllowedOutputs)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 && cmd.Flags().NFlag() == 0 {
				return cmd.Help()
			}
			return nil
		},
	}

	glossaryCmd.PersistentFlags().StringVar(&outputFormat, outputFormatFlag, "simple", "Output format: "+joinedAllowedOutputs)

	return glossaryCmd
}
