package detect

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

			outFormat, err := cmd.Parent().PersistentFlags().GetString("output")
			if err != nil {
				rlog.Errorf("unable to get output: %s", err)
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

			var dataProvider output.DetectDataProvider
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

			updates := make(chan any)
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer func() {
					close(updates)
					wg.Done()
				}()

				_, err := mtSrv.RunDetect(ctx, files, params, updates)
				if err != nil {
					rlog.Errorf("unable to run detect: %s", err)
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
	fileTypeParam := resolve.FallbackString(cmd.Flags().Lookup(fileTypeFlag), resolve.StringParam{
		FlagName: fileTypeFlag,
	})
	inputDirectoryParam := resolve.FallbackString(cmd.Flags().Lookup(inputDirectoryFlag), resolve.StringParam{
		FlagName: inputDirectoryFlag,
		Config:   fileConfig.MT.InputDirectory,
	})
	cnf, err := rootcmd.Config()
	if err != nil {
		return srv.DetectParams{}, fmt.Errorf("unable to read config: %w", err)
	}
	var accountIDConfig *string
	if cnf.AccountID != "" {
		accountIDConfig = &cnf.AccountID
	}
	accountUIDParam := resolve.FallbackString(cmd.Root().PersistentFlags().Lookup("account"), resolve.StringParam{
		FlagName: "account",
		Config:   accountIDConfig,
	})
	return srv.DetectParams{
		FileType:       fileTypeParam,
		InputDirectory: inputDirectoryParam,
		FileOrPattern:  fileOrPattern,
		AccountUID:     api.AccountUID(accountUIDParam),
	}, nil
}
