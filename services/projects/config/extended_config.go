package projectconfig

import (
	"context"

	sdkerror "github.com/Smartling/api-sdk-go/helpers/sm_error"
	"github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/config"

	sdk "github.com/Smartling/api-sdk-go"
	"github.com/reconquest/hierr-go"
)

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

func toExtended(project sdk.ProjectDetails, cfg config.Config) Extended {
	return Extended{
		ProjectID:  project.ProjectID,
		AccountUID: project.AccountUID,
		Name:       project.ProjectName,
		Locale:     project.SourceLocaleID + ": " + project.SourceLocaleDescription,
		Status:     getStatus(project),
		UserID:     cfg.UserID,
		ConfigFile: cfg.Path,
		Sources:    cfg.Sources.String(),
	}
}

func getStatus(details sdk.ProjectDetails) string {
	if details.Archived {
		return "archived"
	}
	return "active"
}
