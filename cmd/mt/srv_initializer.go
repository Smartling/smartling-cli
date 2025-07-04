package mt

import (
	rootcmd "github.com/Smartling/smartling-cli/cmd"
	srv "github.com/Smartling/smartling-cli/services/mt"

	api "github.com/Smartling/api-sdk-go/api/mt"
)

// SrvInitializer defines files service initializer
type SrvInitializer interface {
	InitMTSrv() (srv.Service, error)
}

// NewSrvInitializer returns new SrvInitializer implementation
func NewSrvInitializer() SrvInitializer {
	return srvInitializer{}
}

type srvInitializer struct{}

// InitMTSrv initializes `mt` service with the client and configuration.
func (i srvInitializer) InitMTSrv() (srv.Service, error) {
	client, err := rootcmd.Client()
	if err != nil {
		return nil, err
	}
	downloader := api.NewDownloader(client.Client)
	fileTranslator := api.NewFileTranslator(client.Client)
	uploader := api.NewUploader(client.Client)
	translationControl := api.NewTranslationControl(client.Client)
	mtSrv := srv.NewService(downloader, fileTranslator, uploader, translationControl)
	return mtSrv, nil
}
