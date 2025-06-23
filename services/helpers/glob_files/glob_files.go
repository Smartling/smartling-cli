package globfiles

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Smartling/smartling-cli/services/helpers/cli_error"

	sdkerror "github.com/Smartling/api-sdk-go/helpers/sm_error"
	sdkfile "github.com/Smartling/api-sdk-go/helpers/sm_file"
	"github.com/gobwas/glob"
	"github.com/reconquest/hierr-go"
)

// Remote searches for files matching a specified glob pattern on the remote server.
func Remote(
	listAllFilesFn ListFilesFn,
	project string,
	uri string,
) ([]sdkfile.File, error) {
	if uri == "" {
		uri = "**"
	}

	pattern, err := glob.Compile(uri, '/')
	if err != nil {
		return nil, clierror.NewError(
			err,
			"Search file URI is malformed. Check out help for more "+
				"information about search patterns.",
		)
	}

	request := sdkfile.FilesListRequest{}

	files, err := listAllFilesFn(project, request)
	if err != nil {
		if _, ok := err.(sdkerror.NotFoundError); ok {
			return nil, clierror.ProjectNotFoundError{}
		}

		return nil, hierr.Errorf(
			err,
			`unable to list files in project "%s"`,
			project,
		)
	}

	result := []sdkfile.File{}

	for _, file := range files {
		if pattern.Match(file.FileURI) {
			result = append(result, file)
		}
	}

	if len(result) == 0 {
		return nil, clierror.NewError(
			fmt.Errorf(
				"no files found on the remote server matching provided pattern",
			),

			"Check that file URI pattern is correct.",
		)
	}

	return result, nil
}

// GetDirectoryFromPattern extracts the directory and file mask from a glob pattern.
func GetDirectoryFromPattern(mask string) (string, string) {
	matches := regexp.MustCompile(`^([^*?{}\[\]]+)/(.+)$`).FindStringSubmatch(
		mask,
	)

	if len(matches) < 2 {
		return "", mask
	}

	return matches[1], matches[2]
}

// LocallyFunc searches for files matching a specified glob pattern.
func LocallyFunc(
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

// LocallyFn searches for files matching a specified glob pattern.
func LocallyFn(
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

// GlobFilesLocally searches for files matching a specified glob pattern.
var GlobFilesLocally = LocallyFn

// ListFilesFn is function to list files
type ListFilesFn func(projectID string, request sdkfile.FilesListRequest) ([]sdkfile.File, error)
