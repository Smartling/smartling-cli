package translate

import (
	"github.com/Smartling/smartling-cli/services/helpers/format"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/spf13/cobra"
)

var (
	sourceLocale   string
	detectLanguage string
	targetLocale   string
	directory      string
	directive      string
	progress       string
	fileType       string
	formatPath     string
)

// NewTranslateCmd ...
func NewTranslateCmd() *cobra.Command {
	translateCmd := &cobra.Command{
		Use:   "translate <file|pattern>",
		Short: "Translate files using Smartling's File Machine Translation API.",
		Long:  `Translate files using Smartling's File Machine Translation API.`,

		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				rlog.Error("<file|pattern> argument expected")
				return
			}
			if len(args) > 1 {
				rlog.Errorf("expected one argument, got: %d", len(args))
				return
			}

			//output, _ := cmd.Parent().PersistentFlags().GetString("output")

		},
	}

	translateCmd.Flags().StringVar(&sourceLocale, "source-locale", "", "Explicitly specify source language")
	translateCmd.Flags().StringVar(&detectLanguage, "detect-language", "", "Auto-detect source language")
	translateCmd.Flags().StringVar(&targetLocale, "target-locale", "", "Override file type detection")
	translateCmd.Flags().StringVar(&fileType, "type", "", "Override automatically detected file type.")
	translateCmd.Flags().StringVar(&directory, "directory", "", "Output directory for translated files")
	translateCmd.Flags().StringVar(&formatPath, "format", "", `Translated file naming template.
Default: `+format.DefaultFilePullFormat+`
{{.File}} - Original file path
{{.Locale}} - Target locale
{{name .File}} - File name without extension
{{ext .File}} - File extension
{{dir .File}} - Directory path`)
	translateCmd.Flags().StringVar(&directive, "directive", "", "Smartling directive. Can be specified multiple times")
	translateCmd.Flags().StringVar(&progress, "progress", "", "Override automatically detected file type.")

	if err := translateCmd.MarkFlagRequired("target-locale"); err != nil {
		rlog.Errorf("failed to mark flag required: %v", err)
	}

	return translateCmd
}
