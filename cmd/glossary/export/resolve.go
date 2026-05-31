package glexport

import (
	"strings"

	rootcmd "github.com/Smartling/smartling-cli/cmd"
	glossarycmd "github.com/Smartling/smartling-cli/cmd/glossary"
	"github.com/Smartling/smartling-cli/cmd/helpers/resolve"
	srv "github.com/Smartling/smartling-cli/services/glossary"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/spf13/cobra"
)

func resolveParams(cmd *cobra.Command, fileConfig glossarycmd.FileConfig, glossaryUIDOrName, outFile string) (srv.ExportParams, error) {
	rlog.Debugf("resolving params")

	cnf, err := rootcmd.Config()
	if err != nil {
		return srv.ExportParams{}, clierror.UIError{
			Operation:   "config",
			Err:         err,
			Description: "failed to read config",
		}
	}

	accountUID, err := resolve.FallbackAccount(cmd.Root().PersistentFlags().Lookup("account"), cnf.AccountID)
	if err != nil {
		return srv.ExportParams{}, err
	}

	cfg := fileConfig.Glossary.Export

	fileType := resolve.FallbackString(cmd.Flags().Lookup(fileTypeFlag), resolve.StringParam{FlagName: fileTypeFlag, Config: &cfg.FileType})
	fileType = strings.ToLower(fileType)
	tbxVersion := resolve.FallbackString(cmd.Flags().Lookup(tbxVersionFlag), resolve.StringParam{FlagName: tbxVersionFlag, Config: &cfg.TbxVersion})
	tbxVersion = strings.ToLower(tbxVersion)
	params := srv.ExportParams{
		AccountUID:        accountUID,
		GlossaryUIDOrName: glossaryUIDOrName,
		OutFile:           outFile,

		FileType:   fileType,
		TbxVersion: tbxVersion,

		FocusLocaleID: resolve.FallbackString(cmd.Flags().Lookup(focusLocaleFlag), resolve.StringParam{FlagName: focusLocaleFlag, Config: &cfg.FocusLocaleID}),
		LocaleIDs:     resolve.FallbackStringArray(cmd, localeFlag, cfg.LocaleIDs),
		SkipEntries:   resolve.FallbackBool(cmd.Flags().Lookup(skipEntriesFlag), resolve.BoolParam{FlagName: skipEntriesFlag, Config: &cfg.SkipEntries}),

		Filter: srv.ExportFilter{
			Query:                      resolve.FallbackString(cmd.Flags().Lookup(filterQueryFlag), resolve.StringParam{FlagName: filterQueryFlag, Config: &cfg.Filter.Query}),
			LocaleID:                   resolve.FallbackStringArray(cmd, filterLocaleFlag, cfg.Filter.LocaleIDs),
			EntryUIDs:                  resolve.FallbackStringArray(cmd, filterEntryUIDFlag, cfg.Filter.EntryUIDs),
			EntryState:                 resolve.FallbackString(cmd.Flags().Lookup(filterEntryStateFlag), resolve.StringParam{FlagName: filterEntryStateFlag, Config: &cfg.Filter.EntryState}),
			MissingTranslationLocaleID: resolve.FallbackString(cmd.Flags().Lookup(filterMissingTranslationLocaleFlag), resolve.StringParam{FlagName: filterMissingTranslationLocaleFlag, Config: &cfg.Filter.MissingTranslationLocaleID}),
			PresentTranslationLocaleID: resolve.FallbackString(cmd.Flags().Lookup(filterPresentTranslationLocaleFlag), resolve.StringParam{FlagName: filterPresentTranslationLocaleFlag, Config: &cfg.Filter.PresentTranslationLocaleID}),
			DntLocaleID:                resolve.FallbackString(cmd.Flags().Lookup(filterDntLocaleFlag), resolve.StringParam{FlagName: filterDntLocaleFlag, Config: &cfg.Filter.DntLocaleID}),
			ReturnFallbackTranslations: resolve.FallbackBool(cmd.Flags().Lookup(filterReturnFallbackTranslationsFlag), resolve.BoolParam{FlagName: filterReturnFallbackTranslationsFlag, Config: &cfg.Filter.ReturnFallbackTranslations}),
			LabelsType:                 resolve.FallbackString(cmd.Flags().Lookup(filterLabelsTypeFlag), resolve.StringParam{FlagName: filterLabelsTypeFlag, Config: &cfg.Filter.LabelsType}),
			DntTermSet:                 resolve.FallbackBool(cmd.Flags().Lookup(filterDntTermSetFlag), resolve.BoolParam{FlagName: filterDntTermSetFlag, Config: &cfg.Filter.DntTermSet}),
			Created: srv.Created{
				Level: resolve.FallbackString(cmd.Flags().Lookup(filterCreatedLevelFlag), resolve.StringParam{FlagName: filterCreatedLevelFlag, Config: &cfg.Filter.Created.Level}),
				Type:  resolve.FallbackString(cmd.Flags().Lookup(filterCreatedTypeFlag), resolve.StringParam{FlagName: filterCreatedTypeFlag, Config: &cfg.Filter.Created.Type}),
			},
			CreatedBy: srv.CreatedBy{
				Level:   resolve.FallbackString(cmd.Flags().Lookup(filterCreatedByLevelFlag), resolve.StringParam{FlagName: filterCreatedByLevelFlag, Config: &cfg.Filter.CreatedBy.Level}),
				UserIDs: resolve.FallbackStringArray(cmd, filterCreatedByUserIDFlag, cfg.Filter.CreatedBy.UserIDs),
			},
			LastModified: srv.LastModified{
				Level: resolve.FallbackString(cmd.Flags().Lookup(filterLastModifiedLevelFlag), resolve.StringParam{FlagName: filterLastModifiedLevelFlag, Config: &cfg.Filter.LastModified.Level}),
				Type:  resolve.FallbackString(cmd.Flags().Lookup(filterLastModifiedTypeFlag), resolve.StringParam{FlagName: filterLastModifiedTypeFlag, Config: &cfg.Filter.LastModified.Type}),
			},
			LastModifiedBy: srv.LastModifiedBy{
				Level:   resolve.FallbackString(cmd.Flags().Lookup(filterLastModifiedByLevelFlag), resolve.StringParam{FlagName: filterLastModifiedByLevelFlag, Config: &cfg.Filter.LastModifiedBy.Level}),
				UserIDs: resolve.FallbackStringArray(cmd, filterLastModifiedByUserIDFlag, cfg.Filter.LastModifiedBy.UserIDs),
			},
		},
	}

	if params.Filter.Created.Date, err = resolve.FallbackDate(cmd, createdDateFlag, cfg.Filter.Created.Date); err != nil {
		return srv.ExportParams{}, err
	}
	if params.Filter.LastModified.Date, err = resolve.FallbackDate(cmd, lastModifiedDateFlag, cfg.Filter.LastModified.Date); err != nil {
		return srv.ExportParams{}, err
	}

	return params, nil
}
