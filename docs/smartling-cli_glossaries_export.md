## smartling-cli glossaries export

Export a glossary to a file

### Synopsis

Export a glossary's entries to a file.

The glossary is identified by its UID or name (the first positional argument).
Choose the format with --file-type: csv, xlsx, or tbx. For tbx you must also
pass --tbx-version (v2 or v3).

By default every locale in the glossary is exported; pass one or more --locale
flags to limit the export to specific locales. The --filter-* flags narrow the
exported entries (by query, entry state, labels, creation/modification date,
and more).

If outFile is omitted the file is written to "<glossaryUID>.<file-type>" in the
current directory.

```
smartling-cli glossaries export <glossaryUID|glossaryName> [outFile] [flags]
```

### Examples

```

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

```

### Options

```
      --file-type string                              File type. Allowed: csv, xlsx, tbx
      --filter-created-by-level string                Filter: createdBy.level.
      --filter-created-by-user-id stringArray         Filter: createdBy.userIds entry (repeatable).
      --filter-created-date string                    Filter: created.date in RFC3339 (e.g. 2026-01-02T15:04:05Z).
      --filter-created-level string                   Filter: created.level.
      --filter-created-type string                    Filter: created.type (e.g. AFTER, BEFORE).
      --filter-dnt-locale string                      Filter: DNT (do-not-translate) locale ID.
      --filter-dnt-term-set                           Filter: restrict to entries whose DNT term-set flag is set.
      --filter-entry-state string                     Filter: entry state to match.
      --filter-entry-uid stringArray                  Filter: entry UID to match (repeatable → filter.entryUids).
      --filter-labels-type string                     Filter: labels.type to match.
      --filter-last-modified-by-level string          Filter: lastModifiedBy.level.
      --filter-last-modified-by-user-id stringArray   Filter: lastModifiedBy.userIds entry (repeatable).
      --filter-last-modified-date string              Filter: lastModified.date in RFC3339.
      --filter-last-modified-level string             Filter: lastModified.level.
      --filter-last-modified-type string              Filter: lastModified.type.
      --filter-locale stringArray                     Filter: locale ID to match (repeatable → filter.localeIds).
      --filter-missing-translation-locale string      Filter: locale ID that must be missing a translation.
      --filter-present-translation-locale string      Filter: locale ID that must have a translation.
      --filter-query string                           Filter: free-text query to match entries.
      --filter-return-fallback-translations           Filter: include fallback translations in the result.
      --focus-locale string                           Locale ID to use as the export focus locale.
  -h, --help                                          help for export
      --locale stringArray                            Target locale ID to include in the export (repeatable).
      --skip-entries                                  Skip glossary entries in the export.
      --tbx-version string                            TBX version, required when --file-type=TBX. Allowed: v2, v3
```

### Options inherited from parent commands

```
  -a, --account string               Account ID to operate on.
                                     This option overrides config value "account_id".
  -c, --config string                Config file in YAML format.
                                     By default CLI will look for file named
                                     "smartling.yml" in current directory and in all
                                     intermediate parents, emulating git behavior.
  -k, --insecure                     Skip HTTPS certificate validation.
      --operation-directory string   Sets directory to operate on, usually, to store or to
                                     read files.  Depends on command. (default ".")
      --output string                Output format: table, json, simple (default "simple")
  -p, --project string               Project ID to operate on.
                                     This option overrides config value "project_id".
      --proxy string                 Use specified URL as proxy server.
      --secret string                Token Secret which will be used for authentication.
                                     This option overrides config value "secret".
      --show-config                  Print the resolved account, project, user, and config file path
                                     to stderr before the command runs.
      --smartling-url string         Specify base Smartling URL, merely for testing
                                     purposes.
      --user string                  User ID which will be used for authentication.
                                     This option overrides config value "user_id".
  -v, --verbose count                Verbose logging
```

### SEE ALSO

* [smartling-cli glossaries](smartling-cli_glossaries.md)	 - Manage Smartling glossaries

###### Auto generated by spf13/cobra on 5-Jun-2026
