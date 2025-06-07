package files

import (
	"github.com/Smartling/smartling-cli/services/helpers/config"

	sdk "github.com/Smartling/api-sdk-go"
)

type Service struct {
	Client     sdk.ClientInterface
	Config     config.Config
	FileConfig config.FileConfig
}

func NewService(client sdk.ClientInterface, config config.Config, fileConfig config.FileConfig) *Service {
	return &Service{
		Client:     client,
		Config:     config,
		FileConfig: fileConfig,
	}
}
