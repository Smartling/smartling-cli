package initialize

import (
	"context"

	"github.com/Smartling/smartling-cli/services/helpers/client"
	"github.com/Smartling/smartling-cli/services/helpers/config"
)

// Service defines behavior for initializing the Smartling CLI.
type Service interface {
	RunInit(ctx context.Context, dryRun bool) error
}

// service provides methods to init Smartling CLI.
type service struct {
	ClientConfig client.Config
	Config       config.Config
	Verbose      uint8
}

// NewService creates a new instance of the Service with the provided client settings and configuration.
func NewService(clientConfig client.Config, config config.Config, verbose uint8) Service {
	return &service{ClientConfig: clientConfig, Config: config, Verbose: verbose}
}
