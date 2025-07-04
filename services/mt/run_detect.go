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
		contents, err := getContent(p.InputDirectory, file)
		if err != nil {
			return nil, err
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

func getContent(inputDirectory string, file string) ([]byte, error) {
	name, err := filepath.Abs(file)
	if err != nil {
		return nil, clierror.UIError{
			Err:         err,
			Operation:   "filepath.Abs",
			Description: "Check, that file exists and you have proper permissions to access it.",
			Fields:      map[string]string{"file": file},
		}
	}

	if relPath, err := filepath.Rel(inputDirectory, name); err != nil || strings.HasPrefix(relPath, "..") {
		return nil, clierror.UIError{
			Err:         errors.New(`you are trying to push file outside project directory`),
			Operation:   "filepath.Rel",
			Description: "Check file path and path to configuration file and try again.",
			Fields:      map[string]string{"name": name},
		}
	}

	name, err = filepath.Rel(inputDirectory, name)
	if err != nil {
		return nil, clierror.UIError{
			Err:       err,
			Operation: "filepath.Rel",
			Description: `Unable to resolve relative path to file.
Check, that file exists and you have proper permissions to access it.`,
			Fields: map[string]string{"file": file},
		}
	}

	contents, err := os.ReadFile(file)
	if err != nil {
		return nil, clierror.UIError{
			Err:       err,
			Operation: "os.ReadFile",
			Description: `Unable to read file contents.
Check that file exists and readable by current user.`,
			Fields: map[string]string{"file": file},
		}
	}
	return contents, nil
}

// DetectOutput represents the result of a language detection process for a file
type DetectOutput struct {
	File       string
	Language   string
	Confidence string
}

// DetectUpdates defines updates
type DetectUpdates struct {
	ID       uint32
	Language *string
	Upload   *bool
	Detect   *string
}
