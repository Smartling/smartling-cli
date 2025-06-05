package main

import (
	"fmt"
	"github.com/Smartling/smartling-cli/services/helpers/config"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/Smartling/api-sdk-go"
)

func doFilesStatus(
	client *smartling.Client,
	config config.Config,
	args map[string]interface{},
) error {
	var (
		project   = config.ProjectID
		uri, _    = args["<uri>"].(string)
		directory = args["--directory"].(string)

		defaultFormat, _ = args["--format"].(string)
	)

	if defaultFormat == "" {
		defaultFormat = defaultFileStatusFormat
	}

	info, err := client.GetProjectDetails(project)
	if err != nil {
		return err
	}

	files, err := globFilesRemote(client, project, uri)
	if err != nil {
		return err
	}

	var table = NewTableWriter(os.Stdout)

	var progress = Progress{
		Total: len(files),
	}

	for _, file := range files {
		status, err := client.GetFileStatus(project, file.FileURI)
		if err != nil {
			return err
		}

		progress.Increment()
		progress.Flush()

		translations := status.Items

		translations = append(
			[]smartling.FileStatusTranslation{
				{
					CompletedStringCount: status.TotalStringCount,
					CompletedWordCount:   status.TotalWordCount,
				},
			},
			translations...,
		)

		for _, translation := range translations {
			path, err := executeFileFormat(
				config,
				file,
				defaultFormat,
				usePullFormat,
				map[string]interface{}{
					"FileURI": file.FileURI,
					"Locale":  translation.LocaleID,
				},
			)
			if err != nil {
				return err
			}

			path = filepath.Join(directory, path)

			var (
				locale   = info.SourceLocaleID
				state    = "source"
				progress = "source"
			)

			if translation.LocaleID != "" {
				locale = translation.LocaleID
				state = "remote"
				if status.TotalStringCount > 0 {
					progress = fmt.Sprintf(
						"%d%%",
						int(
							100*
								float64(translation.CompletedStringCount)/
								float64(status.TotalStringCount),
						),
					)
				} else {
					progress = "-"
				}
			}

			if !isFileExists(path) {
				state = "missing"
			}

			writeFileStatus(table, map[string]string{
				"Path":     path,
				"Locale":   locale,
				"State":    state,
				"Progress": progress,
				"Strings":  fmt.Sprint(translation.CompletedStringCount),
				"Words":    fmt.Sprint(translation.CompletedWordCount),
			})
		}
	}

	err = RenderTable(table)
	if err != nil {
		return err
	}

	return nil
}

func writeFileStatus(table *tabwriter.Writer, row map[string]string) {
	fmt.Fprintf(
		table,
		"%s\t%s\t%s\t%s\t%s\t%s\n",
		row["Path"],
		row["Locale"],
		row["State"],
		row["Progress"],
		row["Strings"],
		row["Words"],
	)
}
