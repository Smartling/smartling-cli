package projects

import (
	"context"

	"github.com/Smartling/smartling-cli/services/helpers/config"
	projectconfig "github.com/Smartling/smartling-cli/services/projects/config"

	sdk "github.com/Smartling/api-sdk-go"
)

// Service defines behavior for interacting with Smartling projects.
type Service interface {
	RunInfo(ctx context.Context) (projectconfig.Extended, error)
	RunList(ctx context.Context, short bool) error
	RunLocales(ctx context.Context, params LocalesParams) error
}

// service provides methods to interact with Smartling projects.
type service struct {
	Client sdk.APIClient
	Config config.Config
}

// NewService creates a new instance of the Service with the provided client and configuration.
func NewService(client sdk.APIClient, config config.Config) Service {
	return &service{
		Client: client,
		Config: config,
	}
}
