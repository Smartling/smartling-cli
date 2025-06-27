package detect

import (
	"fmt"
	rootcmd "github.com/Smartling/smartling-cli/cmd"
	mtcmd "github.com/Smartling/smartling-cli/cmd/mt"
	output "github.com/Smartling/smartling-cli/output/mt"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	srv "github.com/Smartling/smartling-cli/services/mt"
	"os"

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
				os.Exit(1)
			}
			var fileOrPattern string
			if len(args) == 1 {
				fileOrPattern = args[0]
			}

			mtSrv, listAllFilesFn, err := initializer.InitMTSrv()
			if err != nil {
				rlog.Errorf("unable to initialize MT service: %w", err)
				os.Exit(1)
			}

			ctx := cmd.Context()

			params, err := resolveParams(cmd, fileOrPattern)
			if err != nil {
				rlog.Error(err)
				os.Exit(1)
			}
			out, err := mtSrv.RunDetect(ctx, params, listAllFilesFn)
			if err != nil {
				rlog.Errorf("unable to run detect: %w", err)
				os.Exit(1)
			}

			outputFormat, err := cmd.Parent().PersistentFlags().GetString("output")
			if err != nil {
				rlog.Errorf("unable to get output: %w", err)
				os.Exit(1)
			}

			fileConfig, err := mtcmd.BindFileConfig(cmd)
			if err != nil {
				rlog.Errorf("unable to bind config: %w", err)
				os.Exit(1)
			}
			outTemplate := resolveOutputTemplate(cmd, fileConfig)
			err = output.RenderDetect(out, outputFormat, outTemplate)
			if err != nil {
				rlog.Errorf("unable to render detect: %w", err)
				os.Exit(1)
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

func resolveParams(cmd *cobra.Command, fileOrPattern string) (srv.DetectParams, error) {
	cnf, err := rootcmd.Config()
	if err != nil {
		return srv.DetectParams{}, fmt.Errorf("unable to read config: %w", err)
	}
	params := srv.DetectParams{
		FileType:      resolveFileType(cmd),
		FileOrPattern: fileOrPattern,
		URI:           "",
	}
	params.AccountUID, err = resolveAccountUID(cmd, cnf.AccountID)
	if err != nil {
		return srv.DetectParams{}, fmt.Errorf("unable to resolve AccountUID: %w", err)
	}
	params.ProjectID, err = resolveProjectID(cmd, cnf.ProjectID)
	if err != nil {
		return srv.DetectParams{}, fmt.Errorf("unable to resolve ProjectID: %w", err)
	}
	return params, nil
}
