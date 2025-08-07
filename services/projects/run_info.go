package projects

import (
	"fmt"
	"os"

	"github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/table"

	sdkerror "github.com/Smartling/api-sdk-go/helpers/sm_error"
	"github.com/reconquest/hierr-go"
)

// RunInfo retrieves and output project details.
// Returns an error if any
func (s service) RunInfo() error {
	details, err := s.Client.GetProjectDetails(s.Config.ProjectID)
	if err != nil {
		if _, ok := err.(sdkerror.NotFoundError); ok {
			return clierror.ProjectNotFoundError{}
		}

		return hierr.Errorf(
			err,
			`unable to get project "%s" details`,
			s.Config.ProjectID,
		)
	}

	tableWriter := table.NewTableWriter(os.Stdout)

	status := "active"

	if details.Archived {
		status = "archived"
	}

	info := [][]interface{}{
		{"ID", details.ProjectID},
		{"ACCOUNT", details.AccountUID},
		{"NAME", details.ProjectName},
		{
			"LOCALE",
			details.SourceLocaleID + ": " + details.SourceLocaleDescription,
		},
		{"STATUS", status},
	}

	for _, row := range info {
		if _, err := fmt.Fprintf(tableWriter, "%s\t%s\n", row...); err != nil {
			return err
		}
	}

	err = table.Render(tableWriter)
	if err != nil {
		return err
	}

	return nil
}
