package remove

import (
	"context"
	"errors"
	"fmt"

	localescmd "github.com/Smartling/smartling-cli/cmd/jobs/locales"
	"github.com/Smartling/smartling-cli/output"
	"github.com/Smartling/smartling-cli/output/static"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	srv "github.com/Smartling/smartling-cli/services/jobs/locales"

	jobapi "github.com/Smartling/api-sdk-go/api/job"
)

func run(ctx context.Context,
	initializer localescmd.SrvInitializer,
	params srv.RemoveParams,
	outputParams output.Params,
) error {
	rlog.Debugf("running jobs locales remove with params: %v", params)
	localeSrv, err := initializer.InitJobLocalesSrv(ctx)
	if err != nil {
		return clierror.UIError{
			Operation:   "init",
			Err:         err,
			Description: "unable to initialize Job Locales service",
		}
	}

	removeOutput, err := localeSrv.RunRemove(ctx, params)
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

	static.GetOutputFormat[srv.Output](outputParams.Format).FormatAndRender(removeOutput)
	return nil
}
