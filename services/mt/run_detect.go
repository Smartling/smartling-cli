package mt

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/pointer"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	api "github.com/Smartling/api-sdk-go/api/mt"
	"github.com/reconquest/hierr-go"
)

// DetectParams is the parameters for the RunDetect method.
type DetectParams struct {
	InputDirectory string
	FileType       string
	FileOrPattern  string
	AccountUID     api.AccountUID
}

func (s service) RunDetect(ctx context.Context, files []string, p DetectParams, updates chan any) ([]DetectOutput, error) {
	var res []DetectOutput
	for fileID, file := range files {
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
		request := api.UploadFileRequest{
			File:     contents,
			FileType: fileType,
		}
		update := DetectUpdates{ID: uint32(fileID)}
		uploadFileResponse, err := s.uploader.UploadFile(p.AccountUID, filepath.Base(file), request)
		if err != nil {
			return nil, err
		}

		update.Upload = pointer.NewP(true)
		updates <- update

		detectFileLanguageResponse, err := s.translationControl.DetectFileLanguage(p.AccountUID, uploadFileResponse.FileUID)
		if err != nil {
			return nil, err
		}

		update.Detect = pointer.NewP("start")
		updates <- update

		//
		var processed bool
		for !processed {
			detectionProgressResponse, err := s.translationControl.DetectionProgress(p.AccountUID, uploadFileResponse.FileUID, detectFileLanguageResponse.LanguageDetectionUID)
			if err != nil {
				return nil, err
			}

			update.Detect = pointer.NewP(detectionProgressResponse.State)
			updates <- update

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
			var languageIDs []string
			for _, detectedSourceLanguages := range detectionProgressResponse.DetectedSourceLanguages {
				filename := filepath.Base(file)
				res = append(res, DetectOutput{
					File:     filename,
					Language: detectedSourceLanguages.LanguageID,
				})
				languageIDs = append(languageIDs, detectedSourceLanguages.LanguageID)
			}

			update.Language = pointer.NewP(strings.Join(languageIDs, ","))
			updates <- update

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

type DetectUpdates struct {
	ID       uint32
	Language *string
	Upload   *bool
	Detect   *string
}
