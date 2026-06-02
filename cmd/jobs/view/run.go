package jobview

import (
	"context"
	"errors"
	"fmt"

	jobscmd "github.com/Smartling/smartling-cli/cmd/jobs"
	"github.com/Smartling/smartling-cli/output"
	"github.com/Smartling/smartling-cli/output/static"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	srv "github.com/Smartling/smartling-cli/services/jobs"

	jobapi "github.com/Smartling/api-sdk-go/api/job"
)

func run(ctx context.Context,
	initializer jobscmd.SrvInitializer,
	params srv.ViewParams,
	outputParams output.Params,
) error {
	rlog.Debugf("running jobs view with params: %v", params)
	jobSrv, err := initializer.InitJobSrv(ctx)
	if err != nil {
		return clierror.UIError{
			Operation:   "init",
			Err:         err,
			Description: "unable to initialize Jobs service",
		}
	}

	viewOutput, err := jobSrv.RunView(ctx, params)
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

	outputFormat := static.GetOutputFormat[srv.ViewOutput](outputParams.Format)
	outputFormat.FormatAndRender(viewOutput)
	return nil
}
