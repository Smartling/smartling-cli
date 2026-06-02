package gllist

import (
	rootcmd "github.com/Smartling/smartling-cli/cmd"
	"github.com/Smartling/smartling-cli/cmd/helpers/resolve"
	srv "github.com/Smartling/smartling-cli/services/glossary"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/spf13/cobra"
)

func resolveParams(cmd *cobra.Command, name string) (srv.ListParams, error) {
	rlog.Debugf("resolving list params")

	cnf, err := rootcmd.Config()
	if err != nil {
		return srv.ListParams{}, clierror.UIError{
			Operation:   "config",
			Err:         err,
			Description: "failed to read config",
		}
	}

	accountUID, err := resolve.FallbackAccount(cmd.Root().PersistentFlags().Lookup("account"), cnf.AccountID)
	if err != nil {
		return srv.ListParams{}, err
	}

	return srv.ListParams{
		AccountUID: accountUID,
		Name:       name,
	}, nil
}
