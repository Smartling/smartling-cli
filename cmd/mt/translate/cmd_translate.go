package translate

import (
	"errors"
	"fmt"
	mtcmd "github.com/Smartling/smartling-cli/cmd/mt"
	output "github.com/Smartling/smartling-cli/output/mt"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/spf13/cobra"
)

const (
	sourceLocaleFlag     = "source-locale"
	detectLanguageFlag   = "detect-language"
	targetLocaleFlag     = "target-locale"
	inputDirectoryFlag   = "input-directory"
	outputDirectoryFlag  = "output-directory"
	directiveFlag        = "directive"
	progressFlag         = "progress"
	overrideFileTypeFlag = "type"
	outputTemplateFlag   = "format"
)

var (
	sourceLocale     string
	detectLanguage   bool
	targetLocales    []string
	inputDirectory   string
	outputDirectory  string
	directive        []string
	progress         bool
	overrideFileType string
	outputTemplate   string
)

// NewTranslateCmd retutns new translate command
func NewTranslateCmd(initializer mtcmd.SrvInitializer) *cobra.Command {
	translateCmd := &cobra.Command{
		Use:   "translate <file|pattern>",
		Short: "Translate files using Smartling's File Machine Translation API.",
		Long:  `Translate files using Smartling's File Machine Translation API.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if len(args) != 1 {
				output.RenderAndExitIfErr(clierror.UIError{
					Operation:   "check args",
					Err:         errors.New("wrong argument quantity"),
					Description: fmt.Sprintf("expected one argument, got: %d", len(args)),
				})
			}
			fileOrPattern := args[0]

			fileConfig, err := mtcmd.BindFileConfig(cmd)
			if err != nil {
				output.RenderAndExitIfErr(clierror.UIError{
					Operation:   "bind",
					Err:         err,
					Description: "unable to bind config",
				})
			}

			outputParams, err := mtcmd.ResolveOutputParams(cmd, fileConfig.MT.FileFormat)
			if err != nil {
				return err
			}

			params, err := resolveParams(cmd, fileConfig)
			if err != nil {
				output.RenderAndExitIfErr(clierror.UIError{
					Operation: "resolve params",
					Err:       err,
				})
			}

			return run(ctx, initializer, params, fileOrPattern, outputParams)
		},
	}

	translateCmd.Flags().StringVar(&sourceLocale, sourceLocaleFlag, "", "Explicitly specify source language")
	translateCmd.Flags().BoolVar(&detectLanguage, detectLanguageFlag, false, "Auto-detect source language")
	translateCmd.Flags().StringArrayVar(&targetLocales, targetLocaleFlag, nil, "Target language(s). Can be specified multiple times")
	translateCmd.Flags().StringVar(&overrideFileType, overrideFileTypeFlag, "", "Set file type to override automatically detected file type. More info: https://help.smartling.com/hc/en-us/articles/360007998893--Supported-File-Types")
	translateCmd.Flags().StringVar(&inputDirectory, inputDirectoryFlag, ".", "Input directory with files")
	translateCmd.Flags().StringVar(&outputDirectory, outputDirectoryFlag, ".", "Output directory for translated files")
	translateCmd.Flags().StringVar(&outputTemplate, outputTemplateFlag, output.DefaultTranslateTemplate, `Translated file naming template.
Default: `+output.DefaultTranslateTemplate+`
{{.File}} - Original file path
{{.Locale}} - Target locale
{{name .File}} - File name without extension
{{ext .File}} - File extension
{{dir .Directory}} - Directory path`)
	translateCmd.Flags().StringArrayVar(&directive, directiveFlag, nil, "Smartling directive. Can be specified multiple times")
	translateCmd.Flags().BoolVar(&progress, progressFlag, true, "Display progress")

	if err := translateCmd.MarkFlagRequired(targetLocaleFlag); err != nil {
		output.RenderAndExitIfErr(clierror.UIError{
			Operation:   "MarkFlagRequired",
			Err:         err,
			Description: "failed to mark " + targetLocaleFlag + " flag required",
		})
	}

	return translateCmd
}
