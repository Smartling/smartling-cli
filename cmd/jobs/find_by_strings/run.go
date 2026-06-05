package findbystrings

import (
	"context"

	jobscmd "github.com/Smartling/smartling-cli/cmd/jobs"
	"github.com/Smartling/smartling-cli/output"
	"github.com/Smartling/smartling-cli/output/static"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	srv "github.com/Smartling/smartling-cli/services/jobs"
)

func run(ctx context.Context,
	initializer jobscmd.SrvInitializer,
	params srv.FindByStringsParams,
	outputParams output.Params,
) error {
	rlog.Debugf("running jobs find-by-strings with params: %v", params)
	jobSrv, err := initializer.InitJobSrv(ctx)
	if err != nil {
		return clierror.UIError{
			Operation:   "init",
			Err:         err,
			Description: "unable to initialize Jobs service",
		}
	}

	findOutput, err := jobSrv.RunFindByStrings(ctx, params)
	if err != nil {
		return err
	}

	outputFormat := static.GetOutputFormat[srv.FindByStringsOutput](outputParams.Format)
	outputFormat.FormatAndRender(findOutput)
	return nil
}
