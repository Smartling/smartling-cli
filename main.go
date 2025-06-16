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
	"github.com/Smartling/smartling-cli/cmd/projects/locates"
)

func main() {
	cmd.ConfigureLogger()
	rootCmd, err := cmd.NewRootCmd()
	if err != nil {
		panic(err)
	}

	initSrvInitializer := initialize.NewSrvInitializer()
	initCmd := initialize.NewInitCmd(initSrvInitializer)
	rootCmd.AddCommand(initCmd)

	filesCmd := files.NewFilesCmd()
	rootCmd.AddCommand(filesCmd)
	filesSrvInitializer := files.NewSrvInitializer()
	filesCmd.AddCommand(deletecmd.NewDeleteCmd(filesSrvInitializer))
	filesCmd.AddCommand(importcmd.NewImportCmd(filesSrvInitializer))
	filesCmd.AddCommand(list.NewListCmd(filesSrvInitializer))
	filesCmd.AddCommand(pull.NewPullCmd(filesSrvInitializer))
	filesCmd.AddCommand(push.NewPushCmd(filesSrvInitializer))
	filesCmd.AddCommand(rename.NewRenameCmd(filesSrvInitializer))
	filesCmd.AddCommand(status.NewStatusCmd(filesSrvInitializer))

	projectsCmd := projects.NewProjectsCmd()
	rootCmd.AddCommand(projectsCmd)
	projectsSrvInitializer := projects.NewSrvInitializer()
	projectsCmd.AddCommand(listprojects.NewListCmd(projectsSrvInitializer))
	projectsCmd.AddCommand(info.NewInfoCmd(projectsSrvInitializer))
	projectsCmd.AddCommand(locates.NewLocatesCmd(projectsSrvInitializer))

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
