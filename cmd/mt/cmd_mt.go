package mt

import (
	"fmt"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

var (
	output         string
	noProgress     bool
	allowedOutputs = []string{
		"table",
		"json",
		"simple",
	}
	joinedAllowedOutputs = strings.Join(allowedOutputs, ", ")
)

// NewMTCmd ...
func NewMTCmd() *cobra.Command {
	mtCmd := &cobra.Command{
		Use:   "mt",
		Short: "mt...",
		Long:  `mt...`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !slices.Contains(allowedOutputs, output) {
				return fmt.Errorf("invalid output: %s (allowed: %s)", output, joinedAllowedOutputs)
			}
			return nil
		},
		Run: func(_ *cobra.Command, _ []string) {

		},
	}

	mtCmd.PersistentFlags().StringVar(&output, "output", "simple", "Output format: "+joinedAllowedOutputs)
	mtCmd.PersistentFlags().BoolVar(&noProgress, "no-progress", false, "Disable progress indicators")

	return mtCmd
}
