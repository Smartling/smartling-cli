package add

import (
	rootcmd "github.com/Smartling/smartling-cli/cmd"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	srv "github.com/Smartling/smartling-cli/services/jobs/strings"
)

func resolveParams(jobUIDOrName string, hashcodes, targetLocales []string, moveEnabled bool) (srv.AddParams, error) {
	rlog.Debugf("resolving add params")

	cnf, err := rootcmd.Config()
	if err != nil {
		return srv.AddParams{}, clierror.UIError{
			Operation:   "config",
			Err:         err,
			Description: "failed to read config",
		}
	}

	return srv.AddParams{
		ProjectID:       cnf.ProjectID,
		JobUIDOrName:    jobUIDOrName,
		Hashcodes:       hashcodes,
		TargetLocaleIDs: targetLocales,
		MoveEnabled:     moveEnabled,
	}, nil
}
