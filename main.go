package main

import (
	"github.com/Smartling/smartling-cli/cmd"
	"github.com/Smartling/smartling-cli/cmd/build"
	"github.com/Smartling/smartling-cli/cmd/docs"
	"github.com/Smartling/smartling-cli/cmd/files"
	deletecmd "github.com/Smartling/smartling-cli/cmd/files/delete"
	importcmd "github.com/Smartling/smartling-cli/cmd/files/import"
	"github.com/Smartling/smartling-cli/cmd/files/list"
	"github.com/Smartling/smartling-cli/cmd/files/pull"
	"github.com/Smartling/smartling-cli/cmd/files/push"
	"github.com/Smartling/smartling-cli/cmd/files/rename"
	"github.com/Smartling/smartling-cli/cmd/files/status"
	"github.com/Smartling/smartling-cli/cmd/glossaries"
	glcreate "github.com/Smartling/smartling-cli/cmd/glossaries/create"
	glexport "github.com/Smartling/smartling-cli/cmd/glossaries/export"
	glimport "github.com/Smartling/smartling-cli/cmd/glossaries/import"
	gllist "github.com/Smartling/smartling-cli/cmd/glossaries/list"
	initialize "github.com/Smartling/smartling-cli/cmd/init"
	"github.com/Smartling/smartling-cli/cmd/jobs"
	jobfiles "github.com/Smartling/smartling-cli/cmd/jobs/files"
	jobfileadd "github.com/Smartling/smartling-cli/cmd/jobs/files/add"
	jobfilelist "github.com/Smartling/smartling-cli/cmd/jobs/files/list"
	jobfileremove "github.com/Smartling/smartling-cli/cmd/jobs/files/remove"
	joblist "github.com/Smartling/smartling-cli/cmd/jobs/list"
	joblocales "github.com/Smartling/smartling-cli/cmd/jobs/locales"
	joblocaleadd "github.com/Smartling/smartling-cli/cmd/jobs/locales/add"
	joblocaleremove "github.com/Smartling/smartling-cli/cmd/jobs/locales/remove"
	"github.com/Smartling/smartling-cli/cmd/jobs/progress"
	jobstrings "github.com/Smartling/smartling-cli/cmd/jobs/strings"
	jobstringadd "github.com/Smartling/smartling-cli/cmd/jobs/strings/add"
	jobstringlist "github.com/Smartling/smartling-cli/cmd/jobs/strings/list"
	jobstringremove "github.com/Smartling/smartling-cli/cmd/jobs/strings/remove"
	jobview "github.com/Smartling/smartling-cli/cmd/jobs/view"
	"github.com/Smartling/smartling-cli/cmd/mt"
	"github.com/Smartling/smartling-cli/cmd/mt/detect"
	"github.com/Smartling/smartling-cli/cmd/mt/translate"
	"github.com/Smartling/smartling-cli/cmd/projects"
	"github.com/Smartling/smartling-cli/cmd/projects/info"
	listprojects "github.com/Smartling/smartling-cli/cmd/projects/list"
	"github.com/Smartling/smartling-cli/cmd/projects/locales"
	output "github.com/Smartling/smartling-cli/output/mt"
)

func main() {
	cmd.ConfigureLogger()

	rootCmd := cmd.NewRootCmd()

	docsCmd := docs.NewDocsCmd()
	rootCmd.AddCommand(docsCmd)

	buildCmd := build.NewBuildCmd()
	rootCmd.AddCommand(buildCmd)

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
	projectsCmd.AddCommand(locales.NewLocalesCmd(projectsSrvInitializer))

	mtCmd := mt.NewMTCmd()
	rootCmd.AddCommand(mtCmd)
	mtInitializer := mt.NewSrvInitializer()
	mtCmd.AddCommand(detect.NewDetectCmd(mtInitializer))
	mtCmd.AddCommand(translate.NewTranslateCmd(mtInitializer))

	jobsCmd := jobs.NewJobsCmd()
	rootCmd.AddCommand(jobsCmd)
	jobInitializer := jobs.NewSrvInitializer()
	jobsCmd.AddCommand(progress.NewProgressCmd(jobInitializer))
	jobsCmd.AddCommand(joblist.NewListCmd(jobInitializer))
	jobsCmd.AddCommand(jobview.NewViewCmd(jobInitializer))
	jobFiles := jobfiles.NewJobFilesCmd()
	jobFilesInitializer := jobfiles.NewSrvInitializer()
	jobFiles.AddCommand(jobfilelist.NewListCmd(jobFilesInitializer))
	jobFiles.AddCommand(jobfileadd.NewJobFilesAddCmd(jobFilesInitializer))
	jobFiles.AddCommand(jobfileremove.NewJobFilesRemoveCmd(jobFilesInitializer))
	jobsCmd.AddCommand(jobFiles)
	jobLocales := joblocales.NewJobLocalesCmd()
	jobLocalesInitializer := joblocales.NewSrvInitializer()
	jobLocales.AddCommand(joblocaleadd.NewJobLocalesAddCmd(jobLocalesInitializer))
	jobLocales.AddCommand(joblocaleremove.NewJobLocalesRemoveCmd(jobLocalesInitializer))
	jobsCmd.AddCommand(jobLocales)
	jobStrings := jobstrings.NewJobStringsCmd()
	jobStringsInitializer := jobstrings.NewSrvInitializer()
	jobStrings.AddCommand(jobstringadd.NewJobStringsAddCmd(jobStringsInitializer))
	jobStrings.AddCommand(jobstringremove.NewJobStringsRemoveCmd(jobStringsInitializer))
	jobStrings.AddCommand(jobstringlist.NewJobStringsListCmd(jobStringsInitializer))
	jobsCmd.AddCommand(jobStrings)

	glossariesCmd := glossaries.NewGlossariesCmd()
	rootCmd.AddCommand(glossariesCmd)
	glossarySrvInitializer := glossaries.NewSrvInitializer()
	glossaryImport := glimport.NewImportCmd(glossarySrvInitializer)
	glossaryExport := glexport.NewExportCmd(glossarySrvInitializer)
	glossaryCreate := glcreate.NewCreateCmd(glossarySrvInitializer)
	glossaryList := gllist.NewListCmd(glossarySrvInitializer)
	glossariesCmd.AddCommand(glossaryImport)
	glossariesCmd.AddCommand(glossaryExport)
	glossariesCmd.AddCommand(glossaryCreate)
	glossariesCmd.AddCommand(glossaryList)

	if err := rootCmd.Execute(); err != nil {
		output.RenderAndExitIfErr(err)
	}
}
