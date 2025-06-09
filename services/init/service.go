package initialize

import (
	"github.com/Smartling/smartling-cli/services/helpers/client"
	"github.com/Smartling/smartling-cli/services/helpers/config"

	sdk "github.com/Smartling/api-sdk-go"
)

type Service struct {
	Client sdk.ClientInterface
	Config config.Config
}

func NewService(client sdk.ClientInterface, config config.Config, cliClientConfig client.Config) *Service {
	return &Service{Client: client, Config: config}
}
