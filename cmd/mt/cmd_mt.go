package mt

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/Smartling/smartling-cli/cmd"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

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
		PersistentPreRunE: func(c *cobra.Command, args []string) error {
			ctx := c.Context()
			if !slices.Contains(allowedOutputs, outputFormat) {
				return fmt.Errorf("invalid output: %s (allowed: %s)", outputFormat, joinedAllowedOutputs)
			}
			if !slices.Contains(allowedOutputModes, outputMode) {
				return fmt.Errorf("invalid output-mode: %s (allowed: %s)", outputMode, joinedAllowedOutputModes)
			}
			if err := cmd.ShowConfigBanner(ctx); err != nil {
				rlog.Error(err)
				os.Exit(1)
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
