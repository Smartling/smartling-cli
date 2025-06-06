package files

import (
	"fmt"

	"github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/config"
	globfiles "github.com/Smartling/smartling-cli/services/helpers/glob_files"
	"github.com/Smartling/smartling-cli/services/helpers/reader"

	"github.com/Smartling/api-sdk-go"
	"github.com/reconquest/hierr-go"
)

type DeleteParams struct {
	URI    string
	Config config.Config
}

func RunDelete(client *smartling.Client, params DeleteParams) error {
	projectID := params.Config.ProjectID
	var (
		err   error
		files []smartling.File
	)
	if params.URI == "-" {
		files, err = reader.ReadFilesFromStdin()
		if err != nil {
			return err
		}
	} else {
		files, err = globfiles.Remote(client, projectID, params.URI)
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
		err := client.DeleteFile(projectID, file.FileURI)
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
