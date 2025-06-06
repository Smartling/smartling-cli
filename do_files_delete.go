package main

import (
	"fmt"

	"github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/config"
	globfiles "github.com/Smartling/smartling-cli/services/helpers/glob_files"
	"github.com/Smartling/smartling-cli/services/helpers/reader"

	"github.com/Smartling/api-sdk-go"
	"github.com/reconquest/hierr-go"
)

func doFilesDelete(
	client *smartling.Client,
	config config.Config,
	args map[string]interface{},
) error {
	var (
		project = config.ProjectID
		uri     = args["<uri>"].(string)
	)

	var (
		err   error
		files []smartling.File
	)

	if uri == "-" {
		files, err = reader.ReadFilesFromStdin()
		if err != nil {
			return err
		}
	} else {
		files, err = globfiles.Remote(client, project, uri)
		if err != nil {
			return err
		}
	}

	if len(files) == 0 {
		return clierror.NewError(
			fmt.Errorf("no files match specified pattern"),

			`Check files list on remote server and your pattern according `+
				`to help.`,
		)
	}

	for _, file := range files {
		err := client.DeleteFile(project, file.FileURI)
		if err != nil {
			return hierr.Errorf(
				err,
				`unable to delete file "%s"`,
				file.FileURI,
			)
		}

		fmt.Printf("%s deleted\n", file.FileURI)
	}

	return nil
}
