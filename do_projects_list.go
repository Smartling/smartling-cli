package main

import (
	"fmt"
	"github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/config"
	table2 "github.com/Smartling/smartling-cli/services/helpers/table"
	"os"

	"github.com/Smartling/api-sdk-go"
	"github.com/reconquest/hierr-go"
)

func doProjectsList(
	client *smartling.Client,
	config config.Config,
	args map[string]interface{},
) error {
	var (
		short = args["--short"].(bool)
	)

	projects, err := client.ListProjects(
		config.AccountID,
		smartling.ProjectsListRequest{},
	)
	if err != nil {
		return clierror.NewError(
			hierr.Errorf(err, "unable to list projects"),
			"",
		)
	}

	table := table2.NewTableWriter(os.Stdout)

	for _, project := range projects.Items {
		if short {
			fmt.Fprintln(table, project.ProjectID)
		} else {
			fmt.Fprintf(
				table,
				"%s\t%s\t%s\n",
				project.ProjectID,
				project.ProjectName,
				project.SourceLocaleID,
			)
		}
	}

	err = table2.Render(table)
	if err != nil {
		return err
	}

	return nil
}
