package progress

import (
	"context"
	"errors"
	"fmt"

	jobapi "github.com/Smartling/api-sdk-go/api/job"
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
	jobSrv, err := initializer.InitJobSrv(ctx)
	if err != nil {
		return clierror.UIError{
			Operation:   "init",
			Err:         err,
			Description: "unable to initialize Jobs service",
		}
	}

	progressOutput, err := jobSrv.RunProgress(ctx, params)
	if err != nil {
		if errors.Is(err, jobapi.ErrNotFound) {
			return clierror.UIError{
				Operation:   "find job",
				Err:         err,
				Description: fmt.Sprintf("no job found for %q", params.JobUIDOrName),
			}
		}
		return err
	}

	outputFormat := jobs.GetOutputFormat(outputParams.Format)
	outputFormat.FormatAndRender(progressOutput)
	return nil
}
