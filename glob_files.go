package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Smartling/smartling-cli/services/helpers/cli_error"

	"github.com/gobwas/glob"
	"github.com/reconquest/hierr-go"
)

func globFilesLocallyFunc(
	directory string,
	base string,
	mask string,
) ([]string, error) {
	if strings.HasPrefix(base, "/") {
		directory = base
	} else {
		directory = filepath.Join(directory, base)
	}

	pattern, err := glob.Compile(mask, '/')
	if err != nil {
		return nil, clierror.NewError(
			err,
			"Search file pattern is malformed. Check out help for more "+
				"information about search patterns.",
		)
	}

	if _, err := os.Stat(filepath.Join(directory, mask)); err == nil {
		return []string{filepath.Join(directory, mask)}, nil
	}

	var result []string

	err = filepath.Walk(
		directory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			path = strings.TrimPrefix(path, directory)
			path = strings.TrimPrefix(path, "/")

			if pattern.Match(path) {
				result = append(
					result,
					filepath.Join(directory, path),
				)
			}

			return nil
		},
	)

	if err != nil {
		return nil, hierr.Errorf(
			err,
			`unable to walk down files in dir "%s"`,
			directory,
		)
	}

	return result, nil
}

var globFilesLocally = globFilesLocallyFunc
