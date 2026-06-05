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

	failed := removeOutput.FailedFileURIs()
	if len(failed) > 0 || len(removeOutput.Unmatched) > 0 {
		var parts []string
		if len(failed) > 0 {
			parts = append(parts, fmt.Sprintf("failed to remove: %s", strings.Join(failed, ", ")))
		}
		if len(removeOutput.Unmatched) > 0 {
			parts = append(parts, fmt.Sprintf("no files matched: %s", strings.Join(removeOutput.Unmatched, ", ")))
		}
		return clierror.UIError{
			Operation:   "remove files",
			Err:         errors.New("some --file patterns did not fully apply"),
			Description: strings.Join(parts, "; "),
		}
	}
	return nil
}
