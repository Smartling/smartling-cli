package files

import (
	"github.com/Smartling/smartling-cli/services/helpers/config"

	sdk "github.com/Smartling/api-sdk-go"
)

type Service struct {
	ClientI sdk.ClientInterface
	Client  *sdk.Client
	Config  config.Config
}

func NewService(clientI sdk.ClientInterface, client *sdk.Client, config config.Config) *Service {
	return &Service{ClientI: clientI, Client: client, Config: config}
}
