package translate

import (
	"errors"
	"fmt"
	"time"

	"github.com/Smartling/smartling-cli/cmd/helpers/resolve"
	mtcmd "github.com/Smartling/smartling-cli/cmd/mt"
	output "github.com/Smartling/smartling-cli/output/mt"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
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

			mtSrv, err := initializer.InitMTSrv()
			if err != nil {
				output.RenderAndExitIfErr(clierror.UIError{
					Operation:   "init",
					Err:         err,
					Description: "unable to initialize MT service",
				})
			}

			fileConfig, err := mtcmd.BindFileConfig(cmd)
			if err != nil {
				output.RenderAndExitIfErr(clierror.UIError{
					Operation:   "bind",
					Err:         err,
					Description: "unable to bind config",
				})
			}

			inputDirectoryParam := resolve.FallbackString(cmd.Flags().Lookup(inputDirectoryFlag), resolve.StringParam{
				FlagName: inputDirectoryFlag,
				Config:   fileConfig.MT.InputDirectory,
			})
			files, err := mtSrv.GetFiles(inputDirectoryParam, fileOrPattern)
			if err != nil {
				output.RenderAndExitIfErr(clierror.UIError{
					Operation:   "get files",
					Err:         err,
					Description: "unable to get input files",
				})
			}

			outputParams, err := mtcmd.ResolveOutputParams(cmd, fileConfig.MT.FileFormat)
			if err != nil {
				return err
			}
			var dataProvider output.TranslateDataProvider
			render, err := mtcmd.InitRender(outputParams, dataProvider, files)
			if err != nil {
				return err
			}
			renderRun := make(chan struct{})
			go func() {
				close(renderRun)
				if err = render.Run(); err != nil {
					output.RenderAndExitIfErr(clierror.UIError{
						Operation: "render run",
						Err:       err,
						Fields: map[string]string{
							"render": fmt.Sprintf("%T", render),
						},
						Description: "unable to run render",
					})
				}
			}()
			<-renderRun
			time.Sleep(time.Second)

			params, err := resolveParams(cmd, fileConfig)
			if err != nil {
				output.RenderAndExitIfErr(clierror.UIError{
					Operation: "resolve params",
					Err:       err,
				})
			}
			params.InputDirectory = inputDirectoryParam

			updates := make(chan any)
			var errGroup errgroup.Group

			errGroup.Go(func() error {
				defer func() {
					close(updates)
				}()
				_, err := mtSrv.RunTranslate(ctx, params, files, updates)
				if err != nil {
					return clierror.UIError{
						Operation: "run translate",
						Err:       err,
					}
				}
				return nil
			})

			errGroup.Go(func() error {
				render.Update(updates)
				return nil
			})

			if err := errGroup.Wait(); err != nil {
				return err
			}
			render.End()
			return nil
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
