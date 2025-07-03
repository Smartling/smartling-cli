package detect

import (
	"errors"
	"fmt"
	"time"

	rootcmd "github.com/Smartling/smartling-cli/cmd"
	"github.com/Smartling/smartling-cli/cmd/helpers/resolve"
	mtcmd "github.com/Smartling/smartling-cli/cmd/mt"
	output "github.com/Smartling/smartling-cli/output/mt"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	srv "github.com/Smartling/smartling-cli/services/mt"

	api "github.com/Smartling/api-sdk-go/api/mt"
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

			render, err := mtcmd.InitRender(cmd, fileConfig.MT.FileFormat, files)
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
