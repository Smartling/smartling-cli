package projects

import (
	"context"

	projectconfig "github.com/Smartling/smartling-cli/services/projects/config"
)

// RunInfo retrieves and outputs project details, including the resolved
// local configuration. Returns an error if any.
func (s service) RunInfo(ctx context.Context) (projectconfig.Extended, error) {
	return projectconfig.FetchExtendedConfig(ctx, s.Config, s.Client.GetProjectDetails)
}
