package mt

import (
	"context"
	"fmt"
	"io"
	"os"

	globfiles "github.com/Smartling/smartling-cli/services/helpers/glob_files"
	"github.com/Smartling/smartling-cli/services/helpers/table"

	sdkfile "github.com/Smartling/api-sdk-go/helpers/sm_file"
	"github.com/reconquest/hierr-go"
)

// DetectParams is the parameters for the RunDetect method.
type DetectParams struct {
	FileType      string
	FormatPath    string
	FileOrPattern string
}

func (s service) RunDetect(ctx context.Context, p DetectParams) (DetectOutput, error) {

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
	s.translationControl

	s.translationControl.DetectionProgress()

	sdkfile.File{}
	return DetectOutput{}, nil
}

type DetectOutput struct {
	File       string
	Language   string
	Confidence string
}
