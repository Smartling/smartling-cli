package main

import (
	"fmt"
	"github.com/Smartling/smartling-cli/services/helpers/config"
	"github.com/Smartling/smartling-cli/services/helpers/format"
	globfiles "github.com/Smartling/smartling-cli/services/helpers/glob_files"
	table2 "github.com/Smartling/smartling-cli/services/helpers/table"
	"io"
	"os"

	"github.com/Smartling/api-sdk-go"
	"github.com/reconquest/hierr-go"
)

func doFilesList(
	client *smartling.Client,
	config config.Config,
	args map[string]interface{},
) error {
	var (
		project = config.ProjectID
		short   = args["--short"].(bool)
		uri, _  = args["<uri>"].(string)
	)

	if args["--format"] == nil {
		args["--format"] = format.DefaultFilesListFormat
	}

	format, err := format.Compile(args["--format"].(string))
	if err != nil {
		return err
	}

	files, err := globfiles.Remote(client, project, uri)
	if err != nil {
		return err
	}

	table := table2.NewTableWriter(os.Stdout)

	for _, file := range files {
		if short {
			fmt.Fprintf(table, "%s\n", file.FileURI)
		} else {
			row, err := format.Execute(file)
			if err != nil {
				return err
			}

			_, err = io.WriteString(table, row)
			if err != nil {
				return hierr.Errorf(
					err,
					"unable to write row to output table",
				)
			}
		}
	}

	err = table2.Render(table)
	if err != nil {
		return err
	}

	return nil
}
