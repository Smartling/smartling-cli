package mt

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	api "github.com/Smartling/api-sdk-go/api/mt"
	sdkfile "github.com/Smartling/api-sdk-go/helpers/sm_file"
	"github.com/reconquest/hierr-go"
)

// DetectParams is the parameters for the RunDetect method.
type DetectParams struct {
	InputDirectory string
	FileType       string
	FileOrPattern  string
	ProjectID      string
	AccountUID     api.AccountUID
	URI            string
}

func (s service) RunDetect(ctx context.Context, files []string, p DetectParams) ([]DetectOutput, error) {
	var res []DetectOutput
	for _, file := range files {
		name, err := filepath.Abs(file)
		if err != nil {
			return nil, clierror.NewError(
				hierr.Errorf(
					err,
					`unable to resolve absolute path to file: %q`,
					file,
				),

				`Check, that file exists and you have proper permissions to access it.`,
			)
		}

		if relPath, err := filepath.Rel(p.InputDirectory, name); err != nil || strings.HasPrefix(relPath, "..") {
			return nil, clierror.NewError(
				errors.New(
					`you are trying to push file outside project directory`,
				),
				`Check file path and path to configuration file and try again.`,
			)
		}

		name, err = filepath.Rel(p.InputDirectory, name)
		if err != nil {
			return nil, clierror.NewError(
				hierr.Errorf(
					err,
					`unable to resolve relative path to file: %q`,
					file,
				),

				`Check, that file exists and you have proper permissions to access it.`,
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
		fileType, found := api.FileTypeByExt[filepath.Ext(file)]
		if !found {
			rlog.Debugf("unknown file type: %s", file)
		}
		request := sdkfile.FileUploadRequest{
			File: contents,

			FileType: sdkfile.FileType(fileType.String()),
		}
		uploadFileResponse, err := s.uploader.UploadFile(p.AccountUID, filepath.Base(file), request)
		if err != nil {
			return nil, err
		}
		detectFileLanguageResponse, err := s.translationControl.DetectFileLanguage(p.AccountUID, uploadFileResponse.FileUID)
		if err != nil {
			return nil, err
		}

		//
		var processed bool
		for !processed {
			detectionProgressResponse, err := s.translationControl.DetectionProgress(p.AccountUID, uploadFileResponse.FileUID, detectFileLanguageResponse.LanguageDetectionUID)
			if err != nil {
				return nil, err
			}

			switch strings.ToUpper(detectionProgressResponse.State) {
			case api.QueuedTranslatedState, api.ProcessingTranslatedState:
				continue
			case api.FailedTranslatedState, api.CanceledTranslatedState, api.CompletedTranslatedState:
				processed = true
			default:
				processed = true
			}
			if detectionProgressResponse.State != api.CompletedTranslatedState {
				break
			}
			for _, detectedSourceLanguages := range detectionProgressResponse.DetectedSourceLanguages {
				filename := filepath.Base(file)
				res = append(res, DetectOutput{
					File:     filename,
					Language: detectedSourceLanguages.LanguageID,
				})
			}

			/*update.Locale = pointer.NewP(strings.Join(localeIDs, ","))
			updates <- update*/

		}
		//

		res = append(res, DetectOutput{
			File:       string(uploadFileResponse.FileUID),
			Language:   detectFileLanguageResponse.Code,
			Confidence: "",
		})
	}

	return res, nil
}

type DetectOutput struct {
	File       string
	Language   string
	Confidence string
}
