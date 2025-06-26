package files

import (
	"fmt"
	"io"
	"os"

	"github.com/Smartling/smartling-cli/services/helpers/format"
	globfiles "github.com/Smartling/smartling-cli/services/helpers/glob_files"
	"github.com/Smartling/smartling-cli/services/helpers/table"

	"github.com/reconquest/hierr-go"
)

// RunList retrieves and outputs a list of files.
func (s service) RunList(formatType string, short bool, uri string) error {
	if formatType == "" {
		formatType = format.DefaultFilesListFormat
	}

	format, err := format.Compile(formatType)
	if err != nil {
		return err
	}

	files, err := globfiles.Remote(s.Client, s.Config.ProjectID, uri)
	if err != nil {
		return err
	}

	tableWriter := table.NewTableWriter(os.Stdout)

	for _, file := range files {
		if short {
			fmt.Fprintf(tableWriter, "%s\n", file.FileURI)
		} else {
			row, err := format.Execute(file)
			if err != nil {
				return err
			}

			_, err = io.WriteString(tableWriter, row)
			if err != nil {
				return hierr.Errorf(
					err,
					"unable to write row to output table",
				)
			}
		}
	}

	err = table.Render(tableWriter)
	if err != nil {
		return err
	}

	return nil
}
