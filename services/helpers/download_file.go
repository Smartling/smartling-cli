package helpers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	sdk "github.com/Smartling/api-sdk-go"
	sdkfile "github.com/Smartling/api-sdk-go/helpers/sm_file"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	"github.com/reconquest/hierr-go"
)

var (
	downloadIdleTimeout = 2 * time.Minute
	errDownloadStalled  = errors.New("download stalled")
)

// DownloadFile downloads a file. The download is aborted if the response body
// delivers no data for downloadIdleTimeout.
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

	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

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
			rlog.Error(err.Error())
		}
	}()

	err = os.MkdirAll(filepath.Dir(path), 0o755)
	if err != nil {
		return hierr.Errorf(
			err,
			`unable to create dirs hierarchy "%s" for downloaded file`,
			path,
		)
	}

	timer := time.AfterFunc(downloadIdleTimeout, func() {
		cancel(errDownloadStalled)
	})
	defer timer.Stop()

	err = writeFileAtomically(path, &idleTimeoutReader{reader: reader, timer: timer})
	if err != nil && errors.Is(context.Cause(ctx), errDownloadStalled) {
		return fmt.Errorf(
			`no data received for %s while downloading into "%s": %w`,
			downloadIdleTimeout,
			path,
			errDownloadStalled,
		)
	}
	return err
}

type idleTimeoutReader struct {
	reader io.Reader
	timer  *time.Timer
}

func (r *idleTimeoutReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	if n > 0 {
		r.timer.Reset(downloadIdleTimeout)
	}
	return n, err
}

// writeFileAtomically writes to a temp file and renames it into place, so an
// interrupted download never leaves a truncated file at path. If path is a
// symlink, its target is replaced and the link is kept.
func writeFileAtomically(path string, reader io.Reader) error {
	target, err := resolveSymlinks(path)
	if err != nil {
		return err
	}

	writer, err := os.CreateTemp(filepath.Dir(target), "."+filepath.Base(target)+".*.tmp")
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

	mode, err := outputFileMode(target)
	if err != nil {
		return err
	}
	err = os.Chmod(writer.Name(), mode)
	if err != nil {
		return hierr.Errorf(
			err,
			`unable to set permissions on output file "%s"`,
			path,
		)
	}

	err = os.Rename(writer.Name(), target)
	if err != nil {
		return hierr.Errorf(
			err,
			`unable to move downloaded file into "%s"`,
			path,
		)
	}

	return nil
}

func resolveSymlinks(path string) (string, error) {
	resolved, err := filepath.EvalSymlinks(path)
	if errors.Is(err, fs.ErrNotExist) {
		return path, nil
	}
	if err != nil {
		return "", hierr.Errorf(err, `unable to resolve output path "%s"`, path)
	}
	return resolved, nil
}

func outputFileMode(path string) (os.FileMode, error) {
	info, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return 0o644, nil
	}
	if err != nil {
		return 0, hierr.Errorf(err, `unable to stat output file "%s"`, path)
	}
	return info.Mode().Perm(), nil
}
