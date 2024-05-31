package main

import (
	"io"
	"os"
	"path/filepath"

	smartling "github.com/Smartling/api-sdk-go"
	"github.com/reconquest/hierr-go"
)

func downloadFile(
	client *smartling.Client,
	project string,
	file smartling.File,
	locale string,
	path string,
	retrievalType smartling.RetrievalType,
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
		request := smartling.FileDownloadRequest{}
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
