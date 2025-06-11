package files

import (
	"github.com/Smartling/smartling-cli/services/helpers/config"

	sdk "github.com/Smartling/api-sdk-go"
)

// Service provides methods to interact with Smartling files.
type Service struct {
	Client     sdk.ClientInterface
	Config     config.Config
	FileConfig config.FileConfig
}

// NewService creates a new instance of the Service with the provided client, and configurations.
func NewService(client sdk.ClientInterface, config config.Config, fileConfig config.FileConfig) *Service {
	return &Service{
		Client:     client,
		Config:     config,
		FileConfig: fileConfig,
	}
}
