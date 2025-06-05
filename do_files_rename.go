package main

import (
	"github.com/Smartling/api-sdk-go"
	"github.com/Smartling/smartling-cli/services/helpers/config"
	"github.com/reconquest/hierr-go"
)

func doFilesRename(
	client *smartling.Client,
	config config.Config,
	args map[string]interface{},
) error {
	var (
		project = config.ProjectID
		oldURI  = args["<old-uri>"].(string)
		newURI  = args["<new-uri>"].(string)
	)

	err := client.RenameFile(project, oldURI, newURI)
	if err != nil {
		return hierr.Errorf(
			err,
			`unable to rename file "%s" -> "%s"`,
			oldURI,
			newURI,
		)
	}

	return nil
}
