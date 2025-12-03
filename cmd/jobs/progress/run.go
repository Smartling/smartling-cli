package progress

import (
	"context"

	jobscmd "github.com/Smartling/smartling-cli/cmd/jobs"
	"github.com/Smartling/smartling-cli/output"
	"github.com/Smartling/smartling-cli/output/jobs"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	srv "github.com/Smartling/smartling-cli/services/jobs"
)

func run(ctx context.Context,
	initializer jobscmd.SrvInitializer,
	params srv.ProgressParams,
	outputParams output.Params,
) error {
	rlog.Debugf("running progress with params: %v", params)
	jobSrv, err := initializer.InitJobSrv()
	if err != nil {
		return clierror.UIError{
			Operation:   "init",
			Err:         err,
			Description: "unable to initialize Jobs service",
		}
	}

	progressOutput, err := jobSrv.RunProgress(ctx, params)
	if err != nil {
		return err
	}

	if progressOutput.TranslationJobUID == "" {
		rlog.Infof("no jobs found for given translationJobUid or translationJobName: %s", params.JobIDOrName)
		return nil
	}

	outputFormat := jobs.GetOutputFormat(outputParams.Format)
	outputFormat.FormatAndRender(progressOutput)
	return nil
}
