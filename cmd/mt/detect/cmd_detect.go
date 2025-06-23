package detect

import (
	sdk "github.com/Smartling/api-sdk-go/api/mt"
	rootcmd "github.com/Smartling/smartling-cli/cmd"
	mtsrv "github.com/Smartling/smartling-cli/cmd/mt"
	output "github.com/Smartling/smartling-cli/output/mt"
	"github.com/Smartling/smartling-cli/services/helpers/format"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	"github.com/Smartling/smartling-cli/services/mt"

	"github.com/spf13/cobra"
)

var (
	fileType   string
	formatPath string
)

// NewDetectCmd ...
func NewDetectCmd(initializer mtsrv.SrvInitializer) *cobra.Command {
	detectCmd := &cobra.Command{
		Use:   "detect <file|pattern>",
		Short: "Detect the source language of files using Smartling's File MT API.",
		Long:  `Detect the source language of files using Smartling's File MT API.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				rlog.Error("<file|pattern> argument expected")
				return
			}
			if len(args) > 1 {
				rlog.Errorf("expected one argument, got: %d", len(args))
				return
			}
			fileOrPattern := args[0]

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
			params := mt.DetectParams{
				FileType:      fileType,
				FormatPath:    formatPath,
				FileOrPattern: fileOrPattern,
				ProjectID:     cnf.ProjectID,
				AccountUID:    sdk.AccountUID(cnf.AccountID),
				URI:           "",
			}
			out, err := mtSrv.RunDetect(ctx, params, listAllFilesFn)
			if err != nil {
				rlog.Errorf("unable to run detect: %w", err)
				return
			}

			err = output.RenderDetect(out, formatPath)
			if err != nil {
				rlog.Errorf("unable to render detect: %w", err)
				return
			}

		},
	}

	detectCmd.Flags().StringVar(&fileType, "type", "", "Override automatically detected file type.")
	detectCmd.Flags().StringVar(&formatPath, "format", "", `Output format template.
Default: `+format.DefaultFilePullFormat+`
{{.File}} - Original file path
{{.Language}} - Detected language code
{{.Confidence}} - Detection confidence (if available)`)

	return detectCmd
}
