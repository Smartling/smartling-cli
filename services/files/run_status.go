package files

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/Smartling/smartling-cli/services/helpers/format"
	globfiles "github.com/Smartling/smartling-cli/services/helpers/glob_files"
	"github.com/Smartling/smartling-cli/services/helpers/progress"
	"github.com/Smartling/smartling-cli/services/helpers/table"

	sdk "github.com/Smartling/api-sdk-go"
)

type StatusParams struct {
	URI       string
	Directory string
	Format    string
}

func (s Service) RunStatus(params StatusParams) error {
	defaultFormat := params.Format
	if defaultFormat == "" {
		defaultFormat = format.DefaultFileStatusFormat
	}

	projectID := s.Config.ProjectID
	info, err := s.Client.GetProjectDetails(projectID)
	if err != nil {
		return err
	}

	files, err := globfiles.Remote(s.Client, projectID, params.URI)
	if err != nil {
		return err
	}

	var tableWriter = table.NewTableWriter(os.Stdout)

	var progress = progress.Progress{
		Total: len(files),
	}

	for _, file := range files {
		status, err := s.Client.GetFileStatus(projectID, file.FileURI)
		if err != nil {
			return err
		}

		progress.Increment()
		progress.Flush()

		translations := status.Items

		translations = append(
			[]sdk.FileStatusTranslation{
				{
					CompletedStringCount: status.TotalStringCount,
					CompletedWordCount:   status.TotalWordCount,
				},
			},
			translations...,
		)

		for _, translation := range translations {
			path, err := format.ExecuteFileFormat(
				s.Config,
				file,
				defaultFormat,
				format.UsePullFormat,
				map[string]interface{}{
					"FileURI": file.FileURI,
					"Locale":  translation.LocaleID,
				},
			)
			if err != nil {
				return err
			}

			path = filepath.Join(params.Directory, path)

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

			writeFileStatus(tableWriter, map[string]string{
				"Path":     path,
				"Locale":   locale,
				"State":    state,
				"Progress": progress,
				"Strings":  fmt.Sprint(translation.CompletedStringCount),
				"Words":    fmt.Sprint(translation.CompletedWordCount),
			})
		}
	}

	err = table.Render(tableWriter)
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

func isFileExists(path string) bool {
	// we don't care about any other errors there, just return false if stat
	// failed for whatever reason
	_, err := os.Stat(path)
	return err == nil
}
