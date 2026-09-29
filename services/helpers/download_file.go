package helpers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	sdk "github.com/Smartling/api-sdk-go"
	sdkfile "github.com/Smartling/api-sdk-go/helpers/sm_file"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	"github.com/reconquest/hierr-go"
)

// DownloadFile downloads a file.
func DownloadFile(
	ctx context.Context,
	client sdk.APIClient,
	project string,
	file sdkfile.File,
	locale string,
	path string,
	retrievalType sdk.RetrievalType,
) error {
	var (
		reader io.ReadCloser
		err    error
	)

	if locale == "" {
		reader, err = client.DownloadFile(ctx, project, file.FileURI)
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

		reader, err = client.DownloadTranslation(ctx, project, locale, request)
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
	defer func() {
		if err := reader.Close(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, 0o755)
	if err != nil {
		return hierr.Errorf(
			err,
			`unable to create dirs hierarchy "%s" for downloaded file`,
			path,
		)
	}

	return writeFileAtomically(dir, path, reader)
}

// writeFileAtomically writes to a temp file and renames it into place, so an
// interrupted download never leaves a truncated file at path.
func writeFileAtomically(dir, path string, reader io.Reader) error {
	writer, err := os.CreateTemp(dir, "."+filepath.Base(path)+".*.tmp")
	if err != nil {
		return hierr.Errorf(
			err,
			`unable to create temporary file for "%s"`,
			path,
		)
	}
	defer func() {
		if err := os.Remove(writer.Name()); err != nil && !errors.Is(err, fs.ErrNotExist) {
			rlog.Error(err.Error())
		}
	}()

	_, err = io.Copy(writer, reader)
	if err != nil {
		if err := writer.Close(); err != nil {
			rlog.Error(err.Error())
		}
		return hierr.Errorf(
			err,
			`unable to write file contents into "%s"`,
			path,
		)
	}

	err = writer.Close()
	if err != nil {
		return hierr.Errorf(
			err,
			`unable to close output file "%s"`,
			path,
		)
	}

	mode := os.FileMode(0o644)
	if info, err := os.Stat(path); err == nil {
		mode = info.Mode().Perm()
	}
	err = os.Chmod(writer.Name(), mode)
	if err != nil {
		return hierr.Errorf(
			err,
			`unable to set permissions on output file "%s"`,
			path,
		)
	}

	err = os.Rename(writer.Name(), path)
	if err != nil {
		return hierr.Errorf(
			err,
			`unable to move downloaded file into "%s"`,
			path,
		)
	}

	return nil
}
