package mt

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	globfiles "github.com/Smartling/smartling-cli/services/helpers/glob_files"

	api "github.com/Smartling/api-sdk-go/api/mt"
	sdkfile "github.com/Smartling/api-sdk-go/helpers/sm_file"
	"github.com/reconquest/hierr-go"
)

// TranslateParams is the parameters for the RunTranslate method.
type TranslateParams struct {
	SourceLocale     string
	DetectLanguage   bool
	TargetLocales    []string
	OutputDirectory  string
	Directives       map[string]string
	Progress         bool
	OverrideFileType string
	FileOrPattern    string
	ProjectID        string
	AccountUID       api.AccountUID
}

func (s service) RunTranslate(ctx context.Context, p TranslateParams) ([]TranslateOutput, error) {
	var res []TranslateOutput

	base, pattern := globfiles.GetDirectoryFromPattern(p.FileOrPattern)
	files, err := globfiles.LocallyFunc(
		p.OutputDirectory,
		base,
		pattern,
	)
	if err != nil {
		return nil, clierror.NewError(
			hierr.Errorf(
				err,
				`unable to find matching files to upload`,
			),

			`Check, that specified pattern is valid and refer to help for`+
				` more information about glob patterns.`,
		)
	}

	if len(files) == 0 {
		return nil, clierror.NewError(
			fmt.Errorf(`no files found by specified patterns`),

			`Check command line pattern if any and configuration file for`+
				` more patterns to search for.`,
		)
	}

	base = filepath.Dir(base)

	for _, file := range files {
		name, err := filepath.Abs(file)
		if err != nil {
			return nil, clierror.NewError(
				hierr.Errorf(
					err,
					`unable to resolve absolute path to file: %q`,
					file,
				),

				`Check, that file exists and you have proper permissions `+
					`to access it.`,
			)
		}

		if relPath, err := filepath.Rel(base, name); err != nil || strings.HasPrefix(relPath, "..") {
			return nil, clierror.NewError(
				errors.New(
					`you are trying to push file outside project directory`,
				),
				`Check file path and path to configuration file and try again.`,
			)
		}

		name, err = filepath.Rel(base, name)
		if err != nil {
			return nil, clierror.NewError(
				hierr.Errorf(
					err,
					`unable to resolve relative path to file: %q`,
					file,
				),

				`Check, that file exists and you have proper permissions `+
					`to access it.`,
			)
		}

		contents, err := os.ReadFile(file)
		if err != nil {
			return nil, clierror.NewError(
				hierr.Errorf(
					err,
					`unable to read file contents "%s"`,
					file,
				),

				`Check that file exists and readable by current user.`,
			)
		}
		request := sdkfile.FileUploadRequest{
			File:               contents,
			LocalesToAuthorize: []string{p.SourceLocale},
		}
		uploadFileResponse, err := s.uploader.UploadFile(p.AccountUID, p.ProjectID, request)
		if err != nil {
			return nil, err
		}
		translatorStartResponse, err := s.fileTranslator.Start(p.AccountUID, uploadFileResponse.FileUID)
		if err != nil {
			return nil, err
		}

		var processed bool
		for !processed {
			progressResponse, err := s.fileTranslator.Progress(p.AccountUID, uploadFileResponse.FileUID, translatorStartResponse.MtUID)
			if err != nil {
				return nil, err
			}
			switch progressResponse.State {
			case api.QueuedTranslatedState, api.ProcessingTranslatedState:
				continue
			case api.FailedTranslatedState:
				processed = true
			case api.CanceledTranslatedState:
				processed = true
			case api.CompletedTranslatedState:
				processed = true
			}
			if processed && progressResponse.State == api.CompletedTranslatedState {
				for _, localeProcessStatus := range progressResponse.LocaleProcessStatuses {
					res = append(res, TranslateOutput{
						File:      file,
						Locale:    localeProcessStatus.LocaleID,
						Name:      "",
						Ext:       "",
						Directory: "",
					})
				}
			}
		}

		return res, nil
	}
	return res, nil
}

// TranslateOutput is translate output
type TranslateOutput struct {
	File      string
	Locale    string
	Name      string
	Ext       string
	Directory string
}
