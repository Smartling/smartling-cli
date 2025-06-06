package files

import (
	"github.com/Smartling/smartling-cli/services/helpers/config"

	sdk "github.com/Smartling/api-sdk-go"
)

type Service struct {
	ClientI    sdk.ClientInterface
	Client     *sdk.Client
	Config     config.Config
	FileConfig config.FileConfig
}

func NewService(clientI sdk.ClientInterface, client *sdk.Client, config config.Config, fileConfig config.FileConfig) *Service {
	return &Service{
		ClientI:    clientI,
		Client:     client,
		Config:     config,
		FileConfig: fileConfig,
	}
}
