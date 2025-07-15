package mt

import (
	"context"
	"time"

	api "github.com/Smartling/api-sdk-go/api/mt"
)

var (
	pollingInterval = time.Second
	pollingDuration = 5 * time.Minute
)

// Service defines behavior for interacting with Smartling MT.
type Service interface {
	RunDetect(ctx context.Context, p DetectParams, files []string, updates chan any) ([]DetectOutput, error)
	RunTranslate(ctx context.Context, p TranslateParams, files []string, updates chan any) ([]TranslateOutput, error)
	GetFiles(inputDirectory, fileOrPattern string) ([]string, error)
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
