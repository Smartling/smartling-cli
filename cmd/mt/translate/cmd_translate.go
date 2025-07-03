package translate

import (
	"fmt"
	"os"
	"sync"
	"time"

	api "github.com/Smartling/api-sdk-go/api/mt"
	rootcmd "github.com/Smartling/smartling-cli/cmd"
	"github.com/Smartling/smartling-cli/cmd/helpers/resolve"
	mtcmd "github.com/Smartling/smartling-cli/cmd/mt"
	output "github.com/Smartling/smartling-cli/output/mt"
	"github.com/Smartling/smartling-cli/services/helpers/config"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	srv "github.com/Smartling/smartling-cli/services/mt"

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

// NewTranslateCmd ...
func NewTranslateCmd(initializer mtcmd.SrvInitializer) *cobra.Command {
	translateCmd := &cobra.Command{
		Use:   "translate <file|pattern>",
		Short: "Translate files using Smartling's File Machine Translation API.",
		Long:  `Translate files using Smartling's File Machine Translation API.`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			if len(args) != 1 {
				rlog.Errorf("expected one argument, got: %d", len(args))
				os.Exit(1)
			}
			fileOrPattern := args[0]

			mtSrv, _, err := initializer.InitMTSrv()
			if err != nil {
				rlog.Errorf("unable to initialize MT service: %w", err)
				os.Exit(1)
			}

			cnf, err := rootcmd.Config()
			if err != nil {
				rlog.Errorf("unable to read config: %w", err)
				os.Exit(1)
			}

			fileConfig, err := mtcmd.BindFileConfig(cmd)
			if err != nil {
				rlog.Errorf("unable to bind config: %w", err)
				os.Exit(1)
			}

			inputDirectoryParam := resolve.FallbackString(cmd.Flags().Lookup(inputDirectoryFlag), resolve.StringParam{
				FlagName: inputDirectoryFlag,
				Config:   fileConfig.MT.InputDirectory,
			})
			files, err := mtSrv.GetFiles(inputDirectoryParam, fileOrPattern)
			if err != nil {
				rlog.Error(err)
				os.Exit(1)
			}

			outFormat, err := cmd.Parent().PersistentFlags().GetString("output")
			if err != nil {
				rlog.Errorf("unable to get output: %w", err)
				os.Exit(1)
			}
			outTemplate := resolve.FallbackString(cmd.Flags().Lookup(outputTemplateFlag), resolve.StringParam{
				FlagName: outputTemplateFlag,
				Config:   fileConfig.MT.FileFormat,
			})

			var render output.Renderer = &output.Static{}
			outMode, err := cmd.Parent().PersistentFlags().GetString("output-mode")
			if err != nil {
				rlog.Errorf("unable to get output mode: %s", err)
				os.Exit(1)
			}
			if outMode == "dynamic" {
				render = &output.Dynamic{}
			}

			var dataProvider output.TranslateDataProvider
			render.Init(dataProvider, files, outFormat, outTemplate)
			renderRun := make(chan struct{})
			go func() {
				close(renderRun)
				if err = render.Run(); err != nil {
					rlog.Error(err)
					os.Exit(1)
				}
			}()
			<-renderRun
			time.Sleep(time.Second)

			params, err := resolveParams(cmd, fileConfig, cnf)
			if err != nil {
				rlog.Error(err)
				os.Exit(1)
			}
			params.InputDirectory = inputDirectoryParam

			updates := make(chan any)
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer func() {
					close(updates)
					wg.Done()
				}()
				_, err := mtSrv.RunTranslate(ctx, params, files, updates)
				if err != nil {
					rlog.Errorf("unable to run translate: %w", err)
					os.Exit(1)
				}
			}()

			go func() {
				defer wg.Done()
				render.Update(updates)
			}()

			wg.Wait()
			render.End()
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
		rlog.Errorf("failed to mark flag required: %v", err)
	}

	return translateCmd
}

func resolveParams(cmd *cobra.Command, fileConfig mtcmd.FileConfig, cnf config.Config) (srv.TranslateParams, error) {
	var err error

	sourceLocaleParam := resolve.FallbackString(cmd.Flags().Lookup(sourceLocaleFlag), resolve.StringParam{
		FlagName: sourceLocaleFlag,
		Config:   fileConfig.MT.DefaultSourceLocale,
	})
	detectLanguageParam := resolve.FallbackBool(cmd.Flags().Lookup(detectLanguageFlag), resolve.BoolParam{
		FlagName: detectLanguageFlag,
	})
	outputDirectoryParam := resolve.FallbackString(cmd.Flags().Lookup(outputDirectoryFlag), resolve.StringParam{
		FlagName: outputDirectoryFlag,
		Config:   fileConfig.MT.OutputDirectory,
	})
	progressParam := resolve.FallbackBool(cmd.Flags().Lookup(progressFlag), resolve.BoolParam{
		FlagName: progressFlag,
	})
	overrideFileTypeParam := resolve.FallbackString(cmd.Flags().Lookup(overrideFileTypeFlag), resolve.StringParam{
		FlagName: overrideFileTypeFlag,
	})

	var accountIDConfig *string
	if cnf.AccountID != "" {
		accountIDConfig = &cnf.AccountID
	}
	accountUIDParam := resolve.FallbackString(cmd.Root().PersistentFlags().Lookup("account"), resolve.StringParam{
		FlagName: "account",
		Config:   accountIDConfig,
	})
	params := srv.TranslateParams{
		SourceLocale:     sourceLocaleParam,
		DetectLanguage:   detectLanguageParam,
		TargetLocales:    resolveTargetLocale(cmd, fileConfig),
		OutputDirectory:  outputDirectoryParam,
		Progress:         progressParam,
		OverrideFileType: overrideFileTypeParam,
		AccountUID:       api.AccountUID(accountUIDParam),
	}
	params.Directives, err = resolveDirectives(cmd, fileConfig)
	if err != nil {
		return srv.TranslateParams{}, fmt.Errorf("unable to resolve directives: %w", err)
	}

	return params, nil
}
