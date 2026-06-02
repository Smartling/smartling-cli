package glcreate

import (
	"strings"

	rootcmd "github.com/Smartling/smartling-cli/cmd"
	glossariescmd "github.com/Smartling/smartling-cli/cmd/glossaries"
	"github.com/Smartling/smartling-cli/cmd/helpers/resolve"
	srv "github.com/Smartling/smartling-cli/services/glossary"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/spf13/cobra"
)

func resolveParams(cmd *cobra.Command, fileConfig glossariescmd.FileConfig, glossaryName string) (srv.CreateParams, error) {
	rlog.Debugf("resolving create params")

	cnf, err := rootcmd.Config()
	if err != nil {
		return srv.CreateParams{}, clierror.UIError{
			Operation:   "config",
			Err:         err,
			Description: "failed to read config",
		}
	}

	accountUID, err := resolve.FallbackAccount(cmd.Root().PersistentFlags().Lookup("account"), cnf.AccountID)
	if err != nil {
		return srv.CreateParams{}, err
	}

	cfg := fileConfig.Glossaries.Create

	description := resolve.FallbackString(cmd.Flags().Lookup(descriptionFlag), resolve.StringParam{FlagName: descriptionFlag})
	verificationMode := resolve.FallbackBool(cmd.Flags().Lookup(verificationModeFlag), resolve.BoolParam{FlagName: verificationModeFlag, Config: &cfg.VerificationMode})
	localeIDs := resolve.FallbackStringArray(cmd, localeFlag, cfg.LocaleIDs)
	rawFallbacks := resolve.FallbackStringArray(cmd, fallbackLocaleFlag, cfg.FallbackLocales)

	fallbacks, err := parseFallbackLocales(rawFallbacks)
	if err != nil {
		return srv.CreateParams{}, err
	}

	return srv.CreateParams{
		AccountUID:       accountUID,
		GlossaryName:     glossaryName,
		Description:      description,
		VerificationMode: verificationMode,
		LocaleIDs:        localeIDs,
		FallbackLocales:  fallbacks,
	}, nil
}

func parseFallbackLocales(raws []string) ([]srv.FallbackLocale, error) {
	fallbacks := make([]srv.FallbackLocale, 0, len(raws))
	for _, raw := range raws {
		parts := strings.SplitN(raw, ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, clierror.UIError{
				Operation:   "parse",
				Description: "fallback locale must use format '<fallbackLocaleId>:<localeId>[,<localeId>...]', got: " + raw,
			}
		}
		fallbacks = append(fallbacks, srv.FallbackLocale{
			FallbackLocaleID: parts[0],
			LocaleIDs:        strings.Split(parts[1], ","),
		})
	}
	return fallbacks, nil
}
