package files

import (
	"github.com/Smartling/smartling-cli/services/helpers/config"

	sdk "github.com/Smartling/api-sdk-go"
)

// Service defines behaviors to interact with Smartling files.
type Service interface {
	RunDelete(uri string) error
	RunImport(params ImportParams) error
	RunList(formatType string, short bool, uri string) error
	RunPull(params PullParams) error
	RunPush(params PushParams) error
	RunRename(oldURI, newURI string) error
	RunStatus(params StatusParams) error
}

// service provides methods to interact with Smartling files.
type service struct {
	Client     sdk.ClientInterface
	Config     config.Config
	FileConfig config.FileConfig
}

// NewService creates a new instance of the Service with the provided client, and configurations.
func NewService(client sdk.ClientInterface, config config.Config, fileConfig config.FileConfig) Service {
	return &service{
		Client:     client,
		Config:     config,
		FileConfig: fileConfig,
	}
}
