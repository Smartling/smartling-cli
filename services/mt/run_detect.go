package mt

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/pointer"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	api "github.com/Smartling/api-sdk-go/api/mt"
	sdktype "github.com/Smartling/api-sdk-go/helpers/file"
)

// DetectParams is the parameters for the RunDetect method.
type DetectParams struct {
	InputDirectory string
	FileType       string
	FileOrPattern  string
	AccountUID     api.AccountUID
}

func (s service) RunDetect(ctx context.Context, p DetectParams, files []string, updates chan any) ([]DetectOutput, error) {
	var res []DetectOutput
	for fileID, file := range files {
		rlog.Debugf("Running detect for file %s", file)
		contents, err := getContent(p.InputDirectory, file)
		if err != nil {
			return nil, err
		}

		fileType, found := sdktype.TypeByExt[filepath.Ext(file)]
		if !found {
			rlog.Debugf("unknown file type: %s", file)
		}
		request := api.UploadFileRequest{
			File:     contents,
			FileType: fileType,
		}
		update := DetectUpdates{ID: uint32(fileID)}
		rlog.Debugf("start upload")
		uploadFileResponse, err := s.uploader.UploadFile(p.AccountUID, filepath.Base(file), request)
		if err != nil {
			return nil, err
		}
		rlog.Debugf("finish upload")

		update.Upload = pointer.NewP(true)
		updates <- update

		rlog.Debugf("detect language")
		detectFileLanguageResponse, err := s.translationControl.DetectFileLanguage(p.AccountUID, uploadFileResponse.FileUID)
		if err != nil {
			return nil, err
		}

		update.Detect = pointer.NewP("start")
		updates <- update

		started := time.Now()
		var processed bool
		for !processed {
			if time.Since(started) > pollingDuration {
				return nil, errors.New("timeout exceeded for polling detection progress of LanguageDetectionUID: " + detectFileLanguageResponse.LanguageDetectionUID)
			}
			rlog.Debugf("check detection progress")
			detectionProgressResponse, err := s.translationControl.DetectionProgress(p.AccountUID, uploadFileResponse.FileUID, detectFileLanguageResponse.LanguageDetectionUID)
			if err != nil {
				return nil, err
			}

			update.Detect = pointer.NewP(detectionProgressResponse.State)
			updates <- update

			rlog.Debugf("progress state: %s", detectionProgressResponse.State)
			switch strings.ToUpper(detectionProgressResponse.State) {
			case api.QueuedTranslatedState, api.ProcessingTranslatedState:
				time.Sleep(pollingInterval)
				continue
			case api.FailedTranslatedState, api.CanceledTranslatedState, api.CompletedTranslatedState:
				processed = true
			default:
				processed = true
			}
			if detectionProgressResponse.State != api.CompletedTranslatedState {
				break
			}

			if len(detectionProgressResponse.DetectedSourceLanguages) > 0 {
				update.Language = pointer.NewP(detectionProgressResponse.DetectedSourceLanguages[0].LanguageID)
			}

			updates <- update
		}

		res = append(res, DetectOutput{
			File:       string(uploadFileResponse.FileUID),
			Language:   detectFileLanguageResponse.Code,
			Confidence: "",
		})
	}

	return res, nil
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

	if strings.HasPrefix(name, "..") {
		return nil, clierror.UIError{
			Err:         errors.New("file name should not start with '..'"),
			Operation:   "filepath.Rel",
			Description: "Check file and directory and try again.",
			Fields: map[string]string{
				"name":           name,
				"inputDirectory": inputDirectory,
			},
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
