package helpers

import (
	"io"
	"os"
	"path/filepath"

	sdk "github.com/Smartling/api-sdk-go"
	sdkfile "github.com/Smartling/api-sdk-go/helpers/sm_file"
	"github.com/reconquest/hierr-go"
)

// DownloadFile downloads a file.
func DownloadFile(
	client sdk.ClientInterface,
	project string,
	file sdkfile.File,
	locale string,
	path string,
	retrievalType sdk.RetrievalType,
) error {
	var (
		reader io.Reader
		err    error
	)

	if locale == "" {
		reader, err = client.DownloadFile(project, file.FileURI)
		if err != nil {
			return hierr.Errorf(
				err,
				`unable to download original file "%s" from project "%s"`,
				file.FileURI,
				project,
			)
		}
	} else {
		request := sdk.FileDownloadRequest{}
		request.FileURI = file.FileURI
		request.Type = retrievalType

		reader, err = client.DownloadTranslation(project, locale, request)
		if err != nil {
			return hierr.Errorf(
				err,
				`unable to download file "%s" from project "%s" (locale "%s")`,
				file.FileURI,
				project,
				locale,
			)
		}
	}

	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return hierr.Errorf(
			err,
			`unable to create dirs hierarchy "%s" for downloaded file`,
			path,
		)
	}

	writer, err := os.Create(path)
	if err != nil {
		return hierr.Errorf(
			err,
			`unable to create output file "%s"`,
			path,
		)
	}

	defer writer.Close()

	_, err = io.Copy(writer, reader)
	if err != nil {
		return hierr.Errorf(
			err,
			`unable to write file contents into "%s"`,
			path,
		)
	}

	return nil
}
