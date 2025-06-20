package detect

import (
	"github.com/Smartling/smartling-cli/services/helpers/format"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/spf13/cobra"
)

var (
	fileType   string
	formatPath string
)

// NewDetectCmd ...
func NewDetectCmd() *cobra.Command {
	detectCmd := &cobra.Command{
		Use:   "detect <file|pattern>",
		Short: "Detect the source language of files using Smartling's File MT API.",
		Long:  `Detect the source language of files using Smartling's File MT API.`,
		Run: func(_ *cobra.Command, args []string) {
			if len(args) == 0 {
				rlog.Error("<file|pattern> argument expected")
				return
			}
			if len(args) > 1 {
				rlog.Errorf("expected one argument, got: %d", len(args))
				return
			}

		},
	}

	detectCmd.Flags().StringVar(&fileType, "type", "", "Override automatically detected file type.")
	detectCmd.Flags().StringVar(&formatPath, "format", "", `Output format template.
Default: `+format.DefaultFilePullFormat+`
{{.File}} - Original file path
{{.Language}} - Detected language code
{{.Confidence}} - Detection confidence (if available)`)

	return detectCmd
}
