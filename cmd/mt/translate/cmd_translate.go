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
		Example: `
# Translate with automatic language detection

smartling-cli mt translate document.txt --target-locale es-ES

# Translate with explicit source language

smartling-cli mt translate document.txt --source-locale en --target-locale fr-FR

`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if len(args) != 1 {
				return clierror.UIError{
					Operation:   "check args",
					Err:         errors.New("wrong argument quantity"),
					Description: fmt.Sprintf("expected one argument, got: %d", len(args)),
				}
			}
			fileOrPattern := args[0]

			fileConfig, err := mtcmd.BindFileConfig(cmd)
			if err != nil {
				return clierror.UIError{
					Operation:   "bind",
					Err:         err,
					Description: "unable to bind config",
				}
			}

			outputParams, err := mtcmd.ResolveOutputParams(cmd, fileConfig.MT.FileFormat)
			if err != nil {
				return err
			}

			params, err := resolveParams(cmd, fileConfig)
			if err != nil {
				return clierror.UIError{
					Operation: "resolve params",
					Err:       err,
				}
			}

			return run(ctx, initializer, params, fileOrPattern, outputParams)
		},
	}

	translateCmd.Flags().StringVar(&sourceLocale, sourceLocaleFlag, "", "Explicitly specify source language")
	translateCmd.Flags().StringArrayVarP(&targetLocales, targetLocaleFlag, "l", nil, `Target language(s). Can be specified multiple times.
Example: Specifying two target locales
smartling-cli mt translate --target-locale fr --target-locale es-ES`)
	translateCmd.Flags().StringVar(&overrideFileType, overrideFileTypeFlag, "", `Override the automatically detected file type. 
A complete list of supported types can be found in the API documentation:
https://api-reference.smartling.com/#tag/File-Machine-Translations-(MT)/operation/fileUpload`)
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
