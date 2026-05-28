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
		Short: "glossary",
		Long:  `glossary`,
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
