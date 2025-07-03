package detect

import (
	"errors"
	"fmt"
	"time"

	mtcmd "github.com/Smartling/smartling-cli/cmd/mt"
	output "github.com/Smartling/smartling-cli/output/mt"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
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

// NewDetectCmd returns new detect command
func NewDetectCmd(initializer mtcmd.SrvInitializer) *cobra.Command {
	detectCmd := &cobra.Command{
		Use:   "detect <file|pattern>",
		Short: "Detect the source language of files using Smartling's File MT API.",
		Long:  `Detect the source language of files using Smartling's File MT API.`,
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

			mtSrv, err := initializer.InitMTSrv()
			if err != nil {
				return clierror.UIError{
					Operation:   "init",
					Err:         err,
					Description: "unable to initialize MT service",
				}
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

			files, err := mtSrv.GetFiles(params.InputDirectory, fileOrPattern)
			if err != nil {
				return clierror.UIError{
					Operation:   "get files",
					Err:         err,
					Description: "unable to get input files",
				}
			}

			outputParams, err := mtcmd.ResolveOutputParams(cmd, fileConfig.MT.FileFormat)
			if err != nil {
				return err
			}
			var dataProvider output.DetectDataProvider
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

			updates := make(chan any)
			var errGroup errgroup.Group
			errGroup.Go(func() error {
				defer func() {
					close(updates)
				}()

				_, err := mtSrv.RunDetect(ctx, files, params, updates)
				if err != nil {
					return clierror.UIError{
						Operation: "run detect",
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

	detectCmd.Flags().StringVar(&fileType, fileTypeFlag, "", "Override automatically detected file type.")
	detectCmd.Flags().StringVar(&inputDirectory, inputDirectoryFlag, ".", "Input directory with files")
	detectCmd.Flags().StringVar(&outputTemplate, outputTemplateFlag, "", `Output format template.
Default: `+output.DefaultDetectTemplate+`
{{.File}} - Original file path
{{.Language}} - Detected language code
{{.Confidence}} - Detection confidence (if available)`)

	return detectCmd
}
