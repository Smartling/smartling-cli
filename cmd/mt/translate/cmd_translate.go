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
	sourceLocaleFlag   = "source-locale"
	detectLanguageFlag = "detect-language"
	targetLocaleFlag   = "target-locale"
	directoryFlag      = "directory"
	directiveFlag      = "directive"
	progressFlag       = "progress"
	fileTypeFlag       = "type"
	outputFormatFlag   = "format"

	defaultOutputFormat = "{{name .File}}_{{.Locale}}{{ext .File}}"
)

var (
	sourceLocale   string
	detectLanguage string
	targetLocale   []string
	directory      string
	directive      []string
	progress       bool
	fileType       string
	outputFormat   string
)

// NewTranslateCmd ...
func NewTranslateCmd(initializer mtcmd.SrvInitializer, fileConfig mtcmd.FileConfig) *cobra.Command {
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

			//output, _ := cmd.Parent().PersistentFlags().GetString("output")

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
			params := srv.TranslateParams{
				SourceLocale:   resolveSourceLocale(cmd, fileConfig),
				DetectLanguage: resolveDetectLanguage(cmd),
				TargetLocale:   resolveTargetLocale(cmd, fileConfig),
				Directory:      resolveDirectory(cmd, fileConfig),
				Progress:       resolveProgress(cmd),
				FileType:       resolveFileType(cmd),
				OutputFormat:   resolveOutputFormat(cmd, fileConfig),
				FileOrPattern:  fileOrPattern,
				URI:            "",
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

			err = output.RenderTranslate(out, outputFormat)
			if err != nil {
				rlog.Errorf("unable to render translate: %w", err)
				return
			}
		},
	}

	translateCmd.Flags().StringVar(&sourceLocale, sourceLocaleFlag, "", "Explicitly specify source language")
	translateCmd.Flags().StringVar(&detectLanguage, detectLanguageFlag, "", "Auto-detect source language")
	translateCmd.Flags().StringArrayVar(&targetLocale, targetLocaleFlag, nil, "Target language(s). Can be specified multiple times")
	translateCmd.Flags().StringVar(&fileType, fileTypeFlag, "", "Override automatically detected file type")
	translateCmd.Flags().StringVar(&directory, directoryFlag, "", "Output directory for translated files")
	translateCmd.Flags().StringVar(&outputFormat, outputFormatFlag, "", `Translated file naming template.
Default: `+defaultOutputFormat+`
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
