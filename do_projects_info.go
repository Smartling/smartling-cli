package main

import (
	"fmt"
	"os"

	smartling "github.com/Smartling/api-sdk-go"
	"github.com/reconquest/hierr-go"
)

func doProjectsInfo(client *smartling.Client, config Config) error {
	details, err := client.GetProjectDetails(config.ProjectID)
	if err != nil {
		if _, ok := err.(smartling.NotFoundError); ok {
			return ProjectNotFoundError{}
		}

		return hierr.Errorf(
			err,
			`unable to get project "%s" details`,
			config.ProjectID,
		)
	}

	table := NewTableWriter(os.Stdout)

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

	err = RenderTable(table)
	if err != nil {
		return err
	}

	return nil
}
