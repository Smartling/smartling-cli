package detect

import (
	"errors"
	"fmt"

	mtcmd "github.com/Smartling/smartling-cli/cmd/mt"
	output "github.com/Smartling/smartling-cli/output/mt"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"

	"github.com/spf13/cobra"
)

const (
	fileTypeFlag       = "type"
	outputTemplateFlag = "format"
	inputDirectoryFlag = "input-directory"
	shortFlag          = "short"
)

var (
	fileType       string
	outputTemplate string
	inputDirectory string
	short          bool
)

// NewDetectCmd returns new detect command
func NewDetectCmd(initializer mtcmd.SrvInitializer) *cobra.Command {
	detectCmd := &cobra.Command{
		Use:   "detect <file|pattern>",
		Short: "Detect the source language of files using Smartling's File MT API.",
		Long:  `Detect the source language of files using Smartling's File MT API.`,
		Example: `
# Detect file language

  smartling-cli mt detect document.txt

`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return clierror.UIError{
					Operation:   "check args",
					Err:         errors.New("wrong argument quantity"),
					Description: fmt.Sprintf("expected one argument, got: %d", len(args)),
				}
			}
			var fileOrPattern string
			if len(args) == 1 {
				fileOrPattern = args[0]
			}

			ctx := cmd.Context()

			fileConfig, err := mtcmd.BindFileConfig(cmd)
			if err != nil {
				return clierror.UIError{
					Operation:   "bind",
					Err:         err,
					Description: "unable to bind config",
				}
			}

			params, err := resolveParams(cmd, fileConfig, fileOrPattern)
			if err != nil {
				return clierror.UIError{
					Operation: "resolve params",
					Err:       err,
				}
			}

			outputParams, err := mtcmd.ResolveOutputParams(cmd, fileConfig.MT.FileFormat)
			if err != nil {
				return err
			}
			if short {
				outputParams.Template = output.DefaultShortDetectTemplate
			}

			return run(ctx, initializer, params, outputParams)
		},
	}

	detectCmd.Flags().StringVar(&fileType, fileTypeFlag, "", `Override the automatically detected file type. 
A complete list of supported types can be found in the API documentation:
https://api-reference.smartling.com/#tag/File-Machine-Translations-(MT)/operation/fileUpload`)
	detectCmd.Flags().StringVar(&inputDirectory, inputDirectoryFlag, ".", "Input directory with files")
	detectCmd.Flags().BoolVarP(&short, shortFlag, "s", false, "Output only detected languages.")
	detectCmd.Flags().StringVar(&outputTemplate, outputTemplateFlag, output.DefaultDetectTemplate, `Output format template.
Default: `+output.DefaultDetectTemplate+`
{{.File}} - Original file path
{{.Language}} - Detected language code`)

	return detectCmd
}
