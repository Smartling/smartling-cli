package main

import (
	"fmt"
	"github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/config"
	table2 "github.com/Smartling/smartling-cli/services/helpers/table"
	"os"

	smartling "github.com/Smartling/api-sdk-go"
	"github.com/reconquest/hierr-go"
)

func doProjectsInfo(client *smartling.Client, config config.Config) error {
	details, err := client.GetProjectDetails(config.ProjectID)
	if err != nil {
		if _, ok := err.(smartling.NotFoundError); ok {
			return clierror.ProjectNotFoundError{}
		}

		return hierr.Errorf(
			err,
			`unable to get project "%s" details`,
			config.ProjectID,
		)
	}

	table := table2.NewTableWriter(os.Stdout)

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
		fmt.Fprintf(
			table,
			"%s\t%s\n",
			row...,
		)
	}

	err = table2.Render(table)
	if err != nil {
		return err
	}

	return nil
}
