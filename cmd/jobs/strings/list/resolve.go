package list

import (
	rootcmd "github.com/Smartling/smartling-cli/cmd"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	srv "github.com/Smartling/smartling-cli/services/jobs/strings"
)

func resolveParams(jobUIDOrName, targetLocale string, limit, offset uint32) (srv.ListParams, error) {
	rlog.Debugf("resolving list params")

	cnf, err := rootcmd.Config()
	if err != nil {
		return srv.ListParams{}, clierror.UIError{
			Operation:   "config",
			Err:         err,
			Description: "failed to read config",
		}
	}

	return srv.ListParams{
		ProjectID:      cnf.ProjectID,
		JobUIDOrName:   jobUIDOrName,
		TargetLocaleID: targetLocale,
		Limit:          limit,
		Offset:         offset,
	}, nil
}
