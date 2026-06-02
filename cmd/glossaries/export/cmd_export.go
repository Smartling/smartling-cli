package glexport

import (
	"strings"

	glossariescmd "github.com/Smartling/smartling-cli/cmd/glossaries"
	"github.com/Smartling/smartling-cli/output"
	glossarysvc "github.com/Smartling/smartling-cli/services/glossary"

	"github.com/spf13/cobra"
)

const (
	fileTypeFlag   = "file-type"
	tbxVersionFlag = "tbx-version"

	focusLocaleFlag = "focus-locale"
	localeFlag      = "locale"
	skipEntriesFlag = "skip-entries"

	filterQueryFlag                      = "filter-query"
	filterLocaleFlag                     = "filter-locale"
	filterEntryUIDFlag                   = "filter-entry-uid"
	filterEntryStateFlag                 = "filter-entry-state"
	filterMissingTranslationLocaleFlag   = "filter-missing-translation-locale"
	filterPresentTranslationLocaleFlag   = "filter-present-translation-locale"
	filterDntLocaleFlag                  = "filter-dnt-locale"
	filterReturnFallbackTranslationsFlag = "filter-return-fallback-translations"
	filterLabelsTypeFlag                 = "filter-labels-type"
	filterDntTermSetFlag                 = "filter-dnt-term-set"

	filterCreatedLevelFlag = "filter-created-level"
	filterCreatedTypeFlag  = "filter-created-type"
	createdDateFlag        = "filter-created-date"

	filterLastModifiedLevelFlag = "filter-last-modified-level"
	filterLastModifiedTypeFlag  = "filter-last-modified-type"
	lastModifiedDateFlag        = "filter-last-modified-date"

	filterCreatedByLevelFlag  = "filter-created-by-level"
	filterCreatedByUserIDFlag = "filter-created-by-user-id"

	filterLastModifiedByLevelFlag  = "filter-last-modified-by-level"
	filterLastModifiedByUserIDFlag = "filter-last-modified-by-user-id"
)

// NewExportCmd builds the `glossaries export` command.
func NewExportCmd(initializer glossariescmd.SrvInitializer) *cobra.Command {
	var (
		fileType   string
		tbxVersion string

		focusLocale string
		locales     []string
		skipEntries bool

		fQuery                      string
		fLocaleIDs                  []string
		fEntryUIDs                  []string
		fEntryState                 string
		fMissingTranslationLocaleID string
		fPresentTranslationLocaleID string
		fDntLocaleID                string
		fReturnFallback             bool
		fLabelsType                 string
		fDntTermSet                 bool

		fCreatedLevel   string
		fCreatedType    string
		fCreatedDateStr string
		fLastModLevel   string
		fLastModType    string
		fLastModDateStr string
		fCreatedByLevel string
		fCreatedByUsers []string
		fLastModByLevel string
		fLastModByUsers []string
	)

	exportCmd := &cobra.Command{
		Use:   "export <glossaryUID|glossaryName> [outFile]",
		Short: "Export a glossary to a file",
		Long: `Export a glossary's entries to a file.

The glossary is identified by its UID or name (the first positional argument).
Choose the format with --file-type: csv, xlsx, or tbx. For tbx you must also
pass --tbx-version (v2 or v3).

By default every locale in the glossary is exported; pass one or more --locale
flags to limit the export to specific locales. The --filter-* flags narrow the
exported entries (by query, entry state, labels, creation/modification date,
and more).

If outFile is omitted the file is written to "<glossaryUID>.<file-type>" in the
current directory.`,
		Example: `
# Export a glossary as CSV to the server-suggested filename

  smartling-cli glossaries export "CLI glossary" --file-type csv

# Export to a specific output file

  smartling-cli glossaries export "CLI glossary" terms.csv --file-type csv

# Export as TBX v3 (requires --tbx-version)

  smartling-cli glossaries export "CLI glossary" terms.tbx --file-type tbx --tbx-version v3

# Export only a subset of locales

  smartling-cli glossaries export "CLI glossary" terms.xlsx --file-type xlsx --locale es-ES --locale fr-FR

# Export entries matching a free-text filter

  smartling-cli glossaries export "CLI glossary" terms.csv --file-type csv --filter-query "checkout"
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			var glossaryUIDOrName, outFile string
			if len(args) > 0 {
				glossaryUIDOrName = args[0]
			}
			if len(args) > 1 {
				outFile = args[1]
			}

			fileConfig, err := glossariescmd.BindFileConfig(cmd)
			if err != nil {
				return err
			}

			params, err := resolveParams(cmd, fileConfig, glossaryUIDOrName, outFile)
			if err != nil {
				return err
			}

			format, err := cmd.Flags().GetString("output")
			if err != nil {
				return err
			}
			outputParams := output.Params{Format: format}
			return run(ctx, initializer, params, outputParams)
		},
	}

	f := exportCmd.Flags()

	f.StringVar(&fileType, fileTypeFlag, "", "File type. Allowed: "+strings.Join(glossarysvc.AllowedExportFileTypes, ", "))
	f.StringVar(&tbxVersion, tbxVersionFlag, "", "TBX version, required when --file-type=TBX. Allowed: "+strings.Join(glossarysvc.AllowedExportTbxVersions, ", "))

	f.StringVar(&focusLocale, focusLocaleFlag, "", "Locale ID to use as the export focus locale.")
	f.StringArrayVar(&locales, localeFlag, nil, "Target locale ID to include in the export (repeatable).")
	f.BoolVar(&skipEntries, skipEntriesFlag, false, "Skip glossary entries in the export.")

	f.StringVar(&fQuery, filterQueryFlag, "", "Filter: free-text query to match entries.")
	f.StringArrayVar(&fLocaleIDs, filterLocaleFlag, nil, "Filter: locale ID to match (repeatable → filter.localeIds).")
	f.StringArrayVar(&fEntryUIDs, filterEntryUIDFlag, nil, "Filter: entry UID to match (repeatable → filter.entryUids).")
	f.StringVar(&fEntryState, filterEntryStateFlag, "", "Filter: entry state to match.")
	f.StringVar(&fMissingTranslationLocaleID, filterMissingTranslationLocaleFlag, "", "Filter: locale ID that must be missing a translation.")
	f.StringVar(&fPresentTranslationLocaleID, filterPresentTranslationLocaleFlag, "", "Filter: locale ID that must have a translation.")
	f.StringVar(&fDntLocaleID, filterDntLocaleFlag, "", "Filter: DNT (do-not-translate) locale ID.")
	f.BoolVar(&fReturnFallback, filterReturnFallbackTranslationsFlag, false, "Filter: include fallback translations in the result.")
	f.StringVar(&fLabelsType, filterLabelsTypeFlag, "", "Filter: labels.type to match.")
	f.BoolVar(&fDntTermSet, filterDntTermSetFlag, false, "Filter: restrict to entries whose DNT term-set flag is set.")

	f.StringVar(&fCreatedLevel, filterCreatedLevelFlag, "", "Filter: created.level.")
	f.StringVar(&fCreatedType, filterCreatedTypeFlag, "", "Filter: created.type (e.g. AFTER, BEFORE).")
	f.StringVar(&fCreatedDateStr, createdDateFlag, "", "Filter: created.date in RFC3339 (e.g. 2026-01-02T15:04:05Z).")

	f.StringVar(&fLastModLevel, filterLastModifiedLevelFlag, "", "Filter: lastModified.level.")
	f.StringVar(&fLastModType, filterLastModifiedTypeFlag, "", "Filter: lastModified.type.")
	f.StringVar(&fLastModDateStr, lastModifiedDateFlag, "", "Filter: lastModified.date in RFC3339.")

	f.StringVar(&fCreatedByLevel, filterCreatedByLevelFlag, "", "Filter: createdBy.level.")
	f.StringArrayVar(&fCreatedByUsers, filterCreatedByUserIDFlag, nil, "Filter: createdBy.userIds entry (repeatable).")
	f.StringVar(&fLastModByLevel, filterLastModifiedByLevelFlag, "", "Filter: lastModifiedBy.level.")
	f.StringArrayVar(&fLastModByUsers, filterLastModifiedByUserIDFlag, nil, "Filter: lastModifiedBy.userIds entry (repeatable).")

	return exportCmd
}
