package translate

import (
	"context"
	"fmt"
	"time"

	mtcmd "github.com/Smartling/smartling-cli/cmd/mt"
	output "github.com/Smartling/smartling-cli/output/mt"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	srv "github.com/Smartling/smartling-cli/services/mt"

	"golang.org/x/sync/errgroup"
)

func run(ctx context.Context,
	initializer mtcmd.SrvInitializer,
	params srv.TranslateParams,
	fileOrPattern string,
	outputParams output.OutputParams) error {
	rlog.Debugf("running translate with params: %v", params)
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
	render := output.InitRender(outputParams, dataProvider, files, uint8(len(params.TargetLocales)))
	renderRun := make(chan struct{})
	var runGroup errgroup.Group
	runGroup.Go(func() error {
		close(renderRun)
		if err = render.Run(); err != nil {
			return clierror.UIError{
				Operation: "render run",
				Err:       err,
				Fields: map[string]string{
					"render": fmt.Sprintf("%T", render),
				},
				Description: "unable to run render",
			}
		}
		return nil
	})
	<-renderRun
	time.Sleep(time.Second)

	updates := make(chan any)
	var updateGroup errgroup.Group
	updateGroup.Go(func() error {
		defer func() {
			close(updates)
		}()
		_, err := mtSrv.RunTranslate(ctx, params, files, updates)
		if err != nil {
			return err
		}
		return nil
	})

	updateGroup.Go(func() error {
		return render.Update(updates)
	})

	if err := updateGroup.Wait(); err != nil {
		return err
	}
	if err := runGroup.Wait(); err != nil {
		return err
	}
	render.End()
	return nil
}
