package translate

import (
	rootcmd "github.com/Smartling/smartling-cli/cmd"
	mtcmd "github.com/Smartling/smartling-cli/cmd/mt"
	output "github.com/Smartling/smartling-cli/output/mt"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	srv "github.com/Smartling/smartling-cli/services/mt"

	"github.com/spf13/cobra"
)

const (
	sourceLocaleFlag          = "source-locale"
	detectLanguageFlag        = "detect-language"
	targetLocaleFlag          = "target-locale"
	outputDirectoryFlag       = "output-directory"
	directiveFlag             = "directive"
	progressFlag              = "progress"
	overrideTypeDetectionFlag = "type"
	outputTemplateFlag        = "format"
)

var (
	sourceLocale          string
	detectLanguage        bool
	targetLocales         []string
	outputDirectory       string
	directive             []string
	progress              bool
	overrideTypeDetection bool
	outputTemplate        string
)

// NewTranslateCmd ...
func NewTranslateCmd(initializer mtcmd.SrvInitializer) *cobra.Command {
	translateCmd := &cobra.Command{
		Use:   "translate <file|pattern>",
		Short: "Translate files using Smartling's File Machine Translation API.",
		Long:  `Translate files using Smartling's File Machine Translation API.`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			if len(args) > 1 {
				rlog.Errorf("expected one argument, got: %d", len(args))
				return
			}
			fileOrPattern := args[0]

			mtSrv, _, err := initializer.InitMTSrv()
			if err != nil {
				rlog.Errorf("unable to initialize MT service: %w", err)
				return
			}

			cnf, err := rootcmd.Config()
			if err != nil {
				rlog.Errorf("unable to read config: %w", err)
				return
			}

			fileConfig, err := mtcmd.BindFileConfig(cmd)
			if err != nil {
				rlog.Errorf("unable to bind config: %w", err)
				return
			}

			params := srv.TranslateParams{
				SourceLocale:          resolveSourceLocale(cmd, fileConfig),
				DetectLanguage:        resolveDetectLanguage(cmd),
				TargetLocales:         resolveTargetLocale(cmd, fileConfig),
				OutputDirectory:       resolveOutputDirectory(cmd, fileConfig),
				Progress:              resolveProgress(cmd),
				OverrideTypeDetection: resolveOverrideTypeDetection(cmd),
				FileOrPattern:         []string{fileOrPattern},
				URI:                   "",
			}
			params.Directives, err = resolveDirectives(cmd, fileConfig)
			if err != nil {
				rlog.Errorf("unable to resolve directives: %w", err)
			}
			params.AccountUID, err = resolveAccountUID(cmd, cnf.AccountID)
			if err != nil {
				rlog.Errorf("unable to resolve AccountUID: %w", err)
			}
			params.ProjectID, err = resolveProjectID(cmd, cnf.ProjectID)
			if err != nil {
				rlog.Errorf("unable to resolve ProjectID: %w", err)
			}
			out, err := mtSrv.RunTranslate(ctx, params)
			if err != nil {
				rlog.Errorf("unable to run translate: %w", err)
				return
			}

			outFormat, err := cmd.Parent().PersistentFlags().GetString("output")
			if err != nil {
				rlog.Errorf("unable to get output: %w", err)
				return
			}
			outTemplate := resolveOutputTemplate(cmd, fileConfig)
			err = output.RenderTranslate(out, outFormat, outTemplate)
			if err != nil {
				rlog.Errorf("unable to render translate: %w", err)
				return
			}
		},
	}

	translateCmd.Flags().StringVar(&sourceLocale, sourceLocaleFlag, "", "Explicitly specify source language")
	translateCmd.Flags().BoolVar(&detectLanguage, detectLanguageFlag, false, "Auto-detect source language")
	translateCmd.Flags().StringArrayVar(&targetLocales, targetLocaleFlag, nil, "Target language(s). Can be specified multiple times")
	translateCmd.Flags().BoolVar(&overrideTypeDetection, overrideTypeDetectionFlag, false, "Override automatically detected file type")
	translateCmd.Flags().StringVar(&outputDirectory, outputDirectoryFlag, "", "Output directory for translated files")
	translateCmd.Flags().StringVar(&outputTemplate, outputTemplateFlag, "", `Translated file naming template.
Default: `+output.DefaultTranslateTemplate+`
{{.File}} - Original file path
{{.Locale}} - Target locale
{{name .File}} - File name without extension
{{ext .File}} - File extension
{{dir .File}} - Directory path`)
	translateCmd.Flags().StringArrayVar(&directive, directiveFlag, nil, "Smartling directive. Can be specified multiple times")
	translateCmd.Flags().BoolVar(&progress, progressFlag, true, "Display progress")

	if err := translateCmd.MarkFlagRequired(targetLocaleFlag); err != nil {
		rlog.Errorf("failed to mark flag required: %v", err)
	}

	return translateCmd
}
