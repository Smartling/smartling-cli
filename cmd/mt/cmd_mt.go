package mt

import (
	"fmt"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

const (
	outputFormatFlag = "output"
	outputModeFlag   = "output-mode"
)

var (
	outputFormat   string
	allowedOutputs = []string{
		"table",
		"json",
		"simple",
	}
	joinedAllowedOutputs = strings.Join(allowedOutputs, ", ")
	outputMode           string
	allowedOutputModes   = []string{
		"dynamic",
		"static",
	}
	joinedAllowedOutputModes = strings.Join(allowedOutputModes, ", ")
)

// NewMTCmd returns new mt command
func NewMTCmd() *cobra.Command {
	mtCmd := &cobra.Command{
		Use:   "mt",
		Short: "File Machine Translations",
		Long:  `Machine Translations offers a simple way to upload files and execute actions on them without any complex setup required`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
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

	mtCmd.PersistentFlags().StringVar(&outputFormat, outputFormatFlag, "simple", "Output format: "+joinedAllowedOutputs)
	mtCmd.PersistentFlags().StringVar(&outputMode, outputModeFlag, "static", "Output mode: "+joinedAllowedOutputModes)

	return mtCmd
}
