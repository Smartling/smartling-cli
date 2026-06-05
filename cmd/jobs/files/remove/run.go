package remove

import (
	"context"
	"errors"
	"fmt"
	"strings"

	filescmd "github.com/Smartling/smartling-cli/cmd/jobs/files"
	"github.com/Smartling/smartling-cli/output"
	"github.com/Smartling/smartling-cli/output/static"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	srv "github.com/Smartling/smartling-cli/services/jobs/files"

	jobapi "github.com/Smartling/api-sdk-go/api/job"
)

func run(ctx context.Context,
	initializer filescmd.SrvInitializer,
	params srv.RemoveParams,
	outputParams output.Params,
) error {
	rlog.Debugf("running jobs files remove with params: %v", params)
	filesSrv, err := initializer.InitJobFilesSrv(ctx)
	if err != nil {
		return clierror.UIError{
			Operation:   "init",
			Err:         err,
			Description: "unable to initialize Job Files service",
		}
	}

	removeOutput, err := filesSrv.RunRemove(ctx, params)
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

	static.GetOutputFormat[srv.MutateOutput](outputParams.Format).FormatAndRender(removeOutput)

	if failed := removeOutput.FailedFileURIs(); len(failed) > 0 {
		return clierror.UIError{
			Operation:   "remove files",
			Err:         fmt.Errorf("%d of %d file(s) failed to remove", len(failed), len(removeOutput.Files)),
			Description: fmt.Sprintf("failed files: %s", strings.Join(failed, ", ")),
		}
	}
	return nil
}
