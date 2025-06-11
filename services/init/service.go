package initialize

import (
	"github.com/Smartling/smartling-cli/services/helpers/config"

	sdk "github.com/Smartling/api-sdk-go"
)

// Service provides methods to init Smartling CLI.
type Service struct {
	Client sdk.ClientInterface
	Config config.Config
}

// NewService creates a new instance of the Service with the provided client and configuration.
func NewService(client sdk.ClientInterface, config config.Config) *Service {
	return &Service{Client: client, Config: config}
}
