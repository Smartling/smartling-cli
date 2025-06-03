package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Smartling/api-sdk-go"
	"github.com/reconquest/hierr-go"
)

func doFilesPush(
	client smartling.ClientInterface,
	config Config,
	args map[string]interface{},
) error {
	var (
		failedFiles   []string
		project       = config.ProjectID
		result        error
		file, _       = args["<file>"].(string)
		uri, useURI   = args["<uri>"].(string)
		branch, _     = args["--branch"].(string)
		locales, _    = args["--locale"].([]string)
		authorize     = args["--authorize"].(bool)
		directory     = args["--directory"].(string)
		fileType, _   = args["--type"].(string)
		directives, _ = args["--directive"].([]string)
	)

	if branch == "@auto" {
		var err error

		branch, err = getGitBranch()
		if err != nil {
			return hierr.Errorf(
				err,
				"unable to autodetect branch name",
			)
		}

		logger.Infof("autodetected branch name: %s", branch)
	}

	if branch != "" {
		branch = strings.TrimSuffix(branch, "/") + "/"
	}

	patterns := []string{}

	if file != "" {
		patterns = append(patterns, file)
	} else {
		for pattern, section := range config.Files {
			if section.Push.Type != "" {
				patterns = append(patterns, pattern)
			}
		}
	}

	files := []string{}

	for _, pattern := range patterns {
		base, pattern := getDirectoryFromPattern(pattern)
		chunk, err := globFilesLocally(
			directory,
			base,
			pattern,
		)
		if err != nil {
			return NewError(
				hierr.Errorf(
					err,
					`unable to find matching files to upload`,
				),

				`Check, that specified pattern is valid and refer to help for`+
					` more information about glob patterns.`,
			)
		}

		files = append(files, chunk...)
	}

	if len(files) == 0 {
		return NewError(
			fmt.Errorf(`no files found by specified patterns`),

			`Check command line pattern if any and configuration file for`+
				` more patterns to search for.`,
		)
	}

	if uri != "" && len(files) > 1 {
		return NewError(
			fmt.Errorf(
				`more than one file is matching speciifed pattern and <uri>`+
					` is specified too`,
			),

			`Either remove <uri> argument or make sure that only one file`+
				` is matching mask.`,
		)
	}

	base, err := filepath.Abs(config.path)
	if err != nil {
		return NewError(
			hierr.Errorf(
				err,
				`unable to resolve absolute path to config`,
			),

			`It's internal error, please, contact developer for more info`,
		)
	}

	base = filepath.Dir(base)

	for _, file := range files {
		name, err := filepath.Abs(file)
		if err != nil {
			return NewError(
				hierr.Errorf(
					err,
					`unable to resolve absolute path to file: %q`,
					file,
				),

				`Check, that file exists and you have proper permissions `+
					`to access it.`,
			)
		}

		if !filepath.HasPrefix(name, base) {
			return NewError(
				errors.New(
					`you are trying to push file outside project directory`,
				),

				`Check file path and path to configuration file and try again.`,
			)
		}

		name, err = filepath.Rel(base, name)
		if err != nil {
			return NewError(
				hierr.Errorf(
					err,
					`unable to resolve relative path to file: %q`,
					file,
				),

				`Check, that file exists and you have proper permissions `+
					`to access it.`,
			)
		}

		if !useURI {
			uri = name
		}

		fileConfig, err := config.GetFileConfig(file)
		if err != nil {
			return NewError(
				hierr.Errorf(
					err,
					`unable to retrieve file specific configuration`,
				),

				``,
			)
		}

		contents, err := os.ReadFile(file)
		if err != nil {
			return NewError(
				hierr.Errorf(
					err,
					`unable to read file contents "%s"`,
					file,
				),

				`Check that file exists and readable by current user.`,
			)
		}

		request := smartling.FileUploadRequest{
			File:               contents,
			Authorize:          authorize,
			LocalesToAuthorize: locales,
		}

		request.FileURI = branch + uri

		if fileConfig.Push.Type == "" {
			if fileType == "" {
				request.FileType = smartling.GetFileTypeByExtension(
					filepath.Ext(file),
				)

				if request.FileType == smartling.FileTypeUnknown {
					return NewError(
						fmt.Errorf(
							"unable to deduce file type from extension: %q",
							filepath.Ext(file),
						),

						`You need to specify file type via --type option.`,
					)
				}
			} else {
				request.FileType = smartling.FileType(fileType)
			}
		} else {
			request.FileType = smartling.FileType(fileConfig.Push.Type)
		}

		request.Smartling.Directives = fileConfig.Push.Directives

		for _, directive := range directives {
			spec := strings.SplitN(directive, "=", 2)
			if len(spec) != 2 {
				return NewError(
					fmt.Errorf(
						"invalid directive specification: %q",
						directive,
					),

					`Should be in the form of <name>=<value>.`,
				)
			}

			if request.Smartling.Directives == nil {
				request.Smartling.Directives = map[string]string{}
			}

			request.Smartling.Directives[spec[0]] = spec[1]
		}

		response, err := client.UploadFile(project, request)

		if err != nil {
			if returnError(err) {
				return NewError(
					err,
					fmt.Sprintf(`unable to upload file "%s"`, file),
					`Check, that you have enough permissions to upload file to`+
						` the specified project`,
				)
			}
			_, _ = fmt.Fprintln(os.Stderr, "Unable to upload file "+file)
			failedFiles = append(failedFiles, file)
		} else {
			status := "new"
			if response.Overwritten {
				status = "overwritten"
			}

			fmt.Printf(
				"%s (%s) %s [%d strings %d words]\n",
				uri,
				request.FileType,
				status,
				response.StringCount,
				response.WordCount,
			)
		}
	}

	if len(failedFiles) != 0 {
		result = NewError(fmt.Errorf("failed to upload %d files", len(failedFiles)), "failed to upload files "+strings.Join(failedFiles, ", "))
	}

	return result
}

func returnError(err error) bool {
	if errors.Is(err, smartling.NotAuthorizedError{}) {
		return true
	}

	for {
		smartlingApiError, isSmartlingApiError := err.(smartling.APIError)
		if isSmartlingApiError {
			reasons := map[string]struct{}{
				"AUTHENTICATION_ERROR":   {},
				"AUTHORIZATION_ERROR":    {},
				"MAINTENANCE_MODE_ERROR": {},
			}

			_, stopExecution := reasons[smartlingApiError.Code]
			return stopExecution
		}
		if err = errors.Unwrap(err); err == nil {
			return false
		}
	}
}
