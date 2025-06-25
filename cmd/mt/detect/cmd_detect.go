package detect

import (
	rootcmd "github.com/Smartling/smartling-cli/cmd"
	mtcmd "github.com/Smartling/smartling-cli/cmd/mt"
	output "github.com/Smartling/smartling-cli/output/mt"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	srv "github.com/Smartling/smartling-cli/services/mt"

	"github.com/spf13/cobra"
)

const (
	fileTypeFlag       = "type"
	outputTemplateFlag = "format"
)

var (
	fileType       string
	outputTemplate string
)

// NewDetectCmd ...
func NewDetectCmd(initializer mtcmd.SrvInitializer) *cobra.Command {
	detectCmd := &cobra.Command{
		Use:   "detect <file|pattern>",
		Short: "Detect the source language of files using Smartling's File MT API.",
		Long:  `Detect the source language of files using Smartling's File MT API.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				rlog.Errorf("expected one argument, got: %d", len(args))
				return
			}
			var fileOrPattern string
			if len(args) == 1 {
				fileOrPattern = args[0]
			}

			mtSrv, listAllFilesFn, err := initializer.InitMTSrv()
			if err != nil {
				rlog.Errorf("unable to initialize MT service: %w", err)
				return
			}

			ctx := cmd.Context()

			cnf, err := rootcmd.Config()
			if err != nil {
				rlog.Errorf("unable to read config: %w", err)
				return
			}
			params := srv.DetectParams{
				FileType:      resolveFileType(cmd),
				FileOrPattern: fileOrPattern,
				URI:           "",
			}
			params.AccountUID, err = resolveAccountUID(cmd, cnf.AccountID)
			if err != nil {
				rlog.Errorf("unable to resolve AccountUID: %w", err)
			}
			params.ProjectID, err = resolveProjectID(cmd, cnf.ProjectID)
			if err != nil {
				rlog.Errorf("unable to resolve ProjectID: %w", err)
			}
			out, err := mtSrv.RunDetect(ctx, params, listAllFilesFn)
			if err != nil {
				rlog.Errorf("unable to run detect: %w", err)
				return
			}

			outputFormat, err := cmd.Parent().PersistentFlags().GetString("output")
			if err != nil {
				rlog.Errorf("unable to get output: %w", err)
				return
			}

			fileConfig, err := mtcmd.BindFileConfig(cmd)
			if err != nil {
				rlog.Errorf("unable to bind config: %w", err)
				return
			}
			outTemplate := resolveOutputTemplate(cmd, fileConfig)
			err = output.RenderDetect(out, outputFormat, outTemplate)
			if err != nil {
				rlog.Errorf("unable to render detect: %w", err)
				return
			}

		},
	}

	detectCmd.Flags().StringVar(&fileType, fileTypeFlag, "", "Override automatically detected file type.")
	detectCmd.Flags().StringVar(&outputTemplate, outputTemplateFlag, "", `Output format template.
Default: `+output.DefaultDetectTemplate+`
{{.File}} - Original file path
{{.Language}} - Detected language code
{{.Confidence}} - Detection confidence (if available)`)

	return detectCmd
}
