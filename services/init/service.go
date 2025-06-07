package init

import (
	"github.com/Smartling/smartling-cli/services/helpers/client"
	"github.com/Smartling/smartling-cli/services/helpers/config"
)

type Service struct {
	Config          config.Config
	CliClientConfig client.Config
}

func NewService(config config.Config, cliClientConfig client.Config) *Service {
	return &Service{Config: config, CliClientConfig: cliClientConfig}
}
