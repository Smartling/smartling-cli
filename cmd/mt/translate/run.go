package translate

import (
	"context"
	"fmt"
	"time"

	mtcmd "github.com/Smartling/smartling-cli/cmd/mt"
	output "github.com/Smartling/smartling-cli/output/mt"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	srv "github.com/Smartling/smartling-cli/services/mt"

	"golang.org/x/sync/errgroup"
)

func run(ctx context.Context,
	initializer mtcmd.SrvInitializer,
	params srv.TranslateParams,
	fileOrPattern string,
	outputParams output.OutputParams) error {
	mtSrv, err := initializer.InitMTSrv()
	if err != nil {
		return clierror.UIError{
			Operation:   "init",
			Err:         err,
			Description: "unable to initialize MT service",
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
	var dataProvider output.TranslateDataProvider
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
		_, err := mtSrv.RunTranslate(ctx, params, files, updates)
		if err != nil {
			return clierror.UIError{
				Operation: "run translate",
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
}
