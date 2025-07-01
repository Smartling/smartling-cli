package detect

import (
	"fmt"
	"os"

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
	inputDirectoryFlag = "input-directory"
)

var (
	fileType       string
	outputTemplate string
	inputDirectory string
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

			mtSrv, _, err := initializer.InitMTSrv()
			if err != nil {
				rlog.Errorf("unable to initialize MT service: %s", err)
				os.Exit(1)
			}

			ctx := cmd.Context()

			fileConfig, err := mtcmd.BindFileConfig(cmd)
			if err != nil {
				rlog.Errorf("unable to bind config: %s", err)
				os.Exit(1)
			}

			params, err := resolveParams(cmd, fileConfig, fileOrPattern)
			if err != nil {
				rlog.Error(err)
				os.Exit(1)
			}

			files, err := mtSrv.GetFiles(params.InputDirectory, fileOrPattern)
			if err != nil {
				rlog.Error(err)
				os.Exit(1)
			}

			outputFormat, err := cmd.Parent().PersistentFlags().GetString("output")
			if err != nil {
				rlog.Errorf("unable to get output: %s", err)
				os.Exit(1)
			}

			outTemplate := resolveOutputTemplate(cmd, fileConfig)
			program, cellCoords, err := output.RenderDetectFiles(files, outputFormat, outTemplate)
			if err != nil {
				rlog.Errorf("unable to render detect: %s", err)
				os.Exit(1)
			}

			updates := make(chan srv.DetectUpdates)

			go func() {
				_, err := mtSrv.RunDetect(ctx, files, params, updates)
				if err != nil {
					rlog.Errorf("unable to run detect: %s", err)
					os.Exit(1)
				}
			}()

			go func() {
				for update := range updates {
					updateRow := output.DetectUpdateRow{
						Coords:  cellCoords,
						Updates: update,
					}
					program.Send(updateRow)
				}
				program.Quit()
			}()

			if _, err := program.Run(); err != nil {
				rlog.Errorf("unable to program run: %w", err)
				os.Exit(1)
			}
		},
	}

	detectCmd.Flags().StringVar(&fileType, fileTypeFlag, "", "Override automatically detected file type.")
	detectCmd.Flags().StringVar(&inputDirectory, inputDirectoryFlag, ".", "Input directory with files")
	detectCmd.Flags().StringVar(&outputTemplate, outputTemplateFlag, "", `Output format template.
Default: `+output.DefaultDetectTemplate+`
{{.File}} - Original file path
{{.Language}} - Detected language code
{{.Confidence}} - Detection confidence (if available)`)

	return detectCmd
}

func resolveParams(cmd *cobra.Command, fileConfig mtcmd.FileConfig, fileOrPattern string) (srv.DetectParams, error) {
	cnf, err := rootcmd.Config()
	if err != nil {
		return srv.DetectParams{}, fmt.Errorf("unable to read config: %w", err)
	}
	params := srv.DetectParams{
		FileType:       resolveFileType(cmd),
		InputDirectory: resolveInputDirectory(cmd, fileConfig),
		FileOrPattern:  fileOrPattern,
		URI:            "",
	}
	params.AccountUID, err = resolveAccountUID(cmd, cnf.AccountID)
	if err != nil {
		return srv.DetectParams{}, fmt.Errorf("unable to resolve AccountUID: %w", err)
	}
	return params, nil
}
