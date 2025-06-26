package projects

import (
	"github.com/Smartling/smartling-cli/services/helpers/config"

	sdk "github.com/Smartling/api-sdk-go"
)

// Service defines behavior for interacting with Smartling projects.
type Service interface {
	RunInfo() error
	RunList(short bool) error
	RunLocales(params LocalesParams) error
}

// service provides methods to interact with Smartling projects.
type service struct {
	Client sdk.ClientInterface
	Config config.Config
}

// NewService creates a new instance of the Service with the provided client and configuration.
func NewService(client sdk.ClientInterface, config config.Config) Service {
	return &service{
		Client: client,
		Config: config,
	}
}
