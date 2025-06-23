package mt

import (
	"context"

	api "github.com/Smartling/api-sdk-go/api/mt"
)

// Service defines behavior for interacting with Smartling MT.
type Service interface {
	RunDetect(ctx context.Context, p DetectParams) (DetectOutput, error)
	RunTranslate(ctx context.Context, p TranslateParams) (TranslateOutput, error)
}

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
