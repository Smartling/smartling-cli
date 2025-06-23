package mt

import (
	"context"

	globfiles "github.com/Smartling/smartling-cli/services/helpers/glob_files"

	api "github.com/Smartling/api-sdk-go/api/mt"
)

// Service defines behavior for interacting with Smartling MT.
type Service interface {
	RunDetect(ctx context.Context, p DetectParams, listAllFilesFn globfiles.ListFilesFn) ([]DetectOutput, error)
	RunTranslate(ctx context.Context, p TranslateParams) (TranslateOutput, error)
}

// NewService creates a new implementation of the Service
func NewService(downloader api.Downloader, fileTranslator api.FileTranslator,
	uploader api.Uploader, translationControl api.TranslationControl) Service {
	return service{
		downloader:         downloader,
		fileTranslator:     fileTranslator,
		uploader:           uploader,
		translationControl: translationControl,
	}
}

type service struct {
	downloader         api.Downloader
	fileTranslator     api.FileTranslator
	uploader           api.Uploader
	translationControl api.TranslationControl
}
