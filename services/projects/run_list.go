package projects

import (
	"fmt"
	"os"

	"github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/table"

	sdk "github.com/Smartling/api-sdk-go"
	"github.com/reconquest/hierr-go"
)

func (s Service) RunList(short bool) error {
	projects, err := s.Client.ListProjects(s.Config.AccountID, sdk.ProjectsListRequest{})
	if err != nil {
		return clierror.NewError(
			hierr.Errorf(err, "unable to list projects"),
			"",
		)
	}

	tableWriter := table.NewTableWriter(os.Stdout)

	for _, project := range projects.Items {
		if short {
			fmt.Fprintln(tableWriter, project.ProjectID)
		} else {
			fmt.Fprintf(
				tableWriter,
				"%s\t%s\t%s\n",
				project.ProjectID,
				project.ProjectName,
				project.SourceLocaleID,
			)
		}
	}

	err = table.Render(tableWriter)
	if err != nil {
		return err
	}

	return nil
}
