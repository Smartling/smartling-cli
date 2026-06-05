package add

import (
	"context"
	"errors"
	"fmt"

	stringscmd "github.com/Smartling/smartling-cli/cmd/jobs/strings"
	"github.com/Smartling/smartling-cli/output"
	"github.com/Smartling/smartling-cli/output/static"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	srv "github.com/Smartling/smartling-cli/services/jobs/strings"

	jobapi "github.com/Smartling/api-sdk-go/api/job"
)

func run(ctx context.Context,
	initializer stringscmd.SrvInitializer,
	params srv.AddParams,
	outputParams output.Params,
) error {
	rlog.Debugf("running jobs strings add with params: %v", params)
	stringsSrv, err := initializer.InitJobStringsSrv(ctx)
	if err != nil {
		return clierror.UIError{
			Operation:   "init",
			Err:         err,
			Description: "unable to initialize Job Strings service",
		}
	}

	addOutput, err := stringsSrv.RunAdd(ctx, params)
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

	static.GetOutputFormat[srv.MutateOutput](outputParams.Format).FormatAndRender(addOutput)
	return nil
}
