package main

import (
	"github.com/Smartling/smartling-cli/cmd"
	"github.com/Smartling/smartling-cli/cmd/files"
	deletecmd "github.com/Smartling/smartling-cli/cmd/files/delete"
	importcmd "github.com/Smartling/smartling-cli/cmd/files/import"
	"github.com/Smartling/smartling-cli/cmd/files/list"
	"github.com/Smartling/smartling-cli/cmd/files/pull"
	"github.com/Smartling/smartling-cli/cmd/files/push"
	"github.com/Smartling/smartling-cli/cmd/files/rename"
	"github.com/Smartling/smartling-cli/cmd/files/status"
	initialize "github.com/Smartling/smartling-cli/cmd/init"
	"github.com/Smartling/smartling-cli/cmd/projects"
	"github.com/Smartling/smartling-cli/cmd/projects/info"
	listprojects "github.com/Smartling/smartling-cli/cmd/projects/list"
	"github.com/Smartling/smartling-cli/cmd/projects/locales"
)

func main() {
	cmd.ConfigureLogger()
	rootCmd, err := cmd.NewRootCmd()
	if err != nil {
		panic(err)
	}

	initCmd := initialize.NewInitCmd()
	rootCmd.AddCommand(initCmd)

	filesCmd := files.NewFilesCmd()
	rootCmd.AddCommand(filesCmd)
	filesCmd.AddCommand(deletecmd.NewDeleteCmd())
	filesCmd.AddCommand(importcmd.NewImportCmd())
	filesCmd.AddCommand(list.NewListCmd())
	filesCmd.AddCommand(pull.NewPullCmd())
	filesCmd.AddCommand(push.NewPushCmd())
	filesCmd.AddCommand(rename.NewRenameCmd())
	filesCmd.AddCommand(status.NewStatusCmd())

	projectsCmd := projects.NewProjectsCmd()
	rootCmd.AddCommand(projectsCmd)
	projectsCmd.AddCommand(listprojects.NewListCmd())
	projectsCmd.AddCommand(info.NewInfoCmd())
	projectsCmd.AddCommand(locales.NewLocatesCmd())

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
