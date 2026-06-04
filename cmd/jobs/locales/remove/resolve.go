package remove

import (
	rootcmd "github.com/Smartling/smartling-cli/cmd"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	srv "github.com/Smartling/smartling-cli/services/jobs/locales"
)

func resolveParams(jobUIDOrName, targetLocaleID string) (srv.RemoveParams, error) {
	rlog.Debugf("resolving remove params")

	cnf, err := rootcmd.Config()
	if err != nil {
		return srv.RemoveParams{}, clierror.UIError{
			Operation:   "config",
			Err:         err,
			Description: "failed to read config",
		}
	}

	return srv.RemoveParams{
		ProjectID:      cnf.ProjectID,
		JobUIDOrName:   jobUIDOrName,
		TargetLocaleID: targetLocaleID,
	}, nil
}
