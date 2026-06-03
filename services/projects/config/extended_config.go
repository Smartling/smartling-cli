package projectconfig

import (
	"context"

	sdkerror "github.com/Smartling/api-sdk-go/helpers/sm_error"
	"github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/config"

	sdk "github.com/Smartling/api-sdk-go"
	"github.com/reconquest/hierr-go"
)

// Extended is local config joined with API-fetched project facts.
type Extended struct {
	ProjectID  string
	Name       string
	AccountUID string
	Locale     string
	Status     string
	UserID     string
	ConfigFile string
	Sources    string
}

// InjectConfig injects fields from config.
func (e *Extended) InjectConfig(cfg config.Config) {
	e.AccountUID = cfg.AccountID
	e.UserID = cfg.UserID
	e.ProjectID = cfg.ProjectID
	e.ConfigFile = cfg.Path
	e.Sources = cfg.Sources.String()
}

// InjectProject injects fields from project details.
func (e *Extended) InjectProject(project sdk.ProjectDetails) {
	e.ProjectID = project.ProjectID
	e.AccountUID = project.AccountUID
	e.Name = project.ProjectName
	e.Locale = project.SourceLocaleID + ": " + project.SourceLocaleDescription
	e.Status = getStatus(project)
}

// FetchExtendedConfig fetches project details and merges with local config.
func FetchExtendedConfig(ctx context.Context, config config.Config,
	projectFetcher func(ctx context.Context, projectID string) (*sdk.ProjectDetails, error),
) (Extended, error) {
	details, err := projectFetcher(ctx, config.ProjectID)
	if err != nil {
		if _, ok := err.(sdkerror.NotFoundError); ok {
			return Extended{}, clierror.ProjectNotFoundError{}
		}

		return Extended{}, hierr.Errorf(
			err,
			`unable to get project "%s" details`,
			config.ProjectID,
		)
	}

	infoOutput := toExtended(*details, config)
	return infoOutput, nil
}

// toExtended flattens project details + local config into Extended.
func toExtended(project sdk.ProjectDetails, cfg config.Config) Extended {
	var res Extended
	res.InjectConfig(cfg)
	res.InjectProject(project)
	return res
}

// getStatus returns "archived" or "active".
func getStatus(details sdk.ProjectDetails) string {
	if details.Archived {
		return "archived"
	}
	return "active"
}
