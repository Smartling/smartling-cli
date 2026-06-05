package findbystrings

import (
	rootcmd "github.com/Smartling/smartling-cli/cmd"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	srv "github.com/Smartling/smartling-cli/services/jobs"
)

func resolveParams(hashcodes, locales []string) (srv.FindByStringsParams, error) {
	rlog.Debugf("resolving find-by-strings params")

	cnf, err := rootcmd.Config()
	if err != nil {
		return srv.FindByStringsParams{}, clierror.UIError{
			Operation:   "config",
			Err:         err,
			Description: "failed to read config",
		}
	}

	return srv.FindByStringsParams{
		ProjectUID: cnf.ProjectID,
		Hashcodes:  hashcodes,
		LocaleIDs:  locales,
	}, nil
}
