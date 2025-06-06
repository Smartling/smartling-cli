package files

import (
	"fmt"

	"github.com/Smartling/smartling-cli/services/helpers/cli_error"
	globfiles "github.com/Smartling/smartling-cli/services/helpers/glob_files"
	"github.com/Smartling/smartling-cli/services/helpers/reader"

	"github.com/Smartling/api-sdk-go"
	"github.com/reconquest/hierr-go"
)

func (s Service) RunDelete(uri string) error {
	projectID := s.Config.ProjectID
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
		files, err = globfiles.Remote(s.Client, projectID, uri)
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
		err := s.Client.DeleteFile(projectID, file.FileURI)
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
