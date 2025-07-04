package initialize

import (
	"github.com/Smartling/smartling-cli/services/helpers/config"

	sdk "github.com/Smartling/api-sdk-go"
)

// Service defines behavior for initializing the Smartling CLI.
type Service interface {
	RunInit(dryRun bool) error
}

// service provides methods to init Smartling CLI.
type service struct {
	Client sdk.APIClient
	Config config.Config
}

// NewService creates a new instance of the Service with the provided client and configuration.
func NewService(client sdk.APIClient, config config.Config) Service {
	return &service{Client: client, Config: config}
}
