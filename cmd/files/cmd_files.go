package files

import (
	"github.com/Smartling/smartling-cli/cmd/files/delete"
	importcmd "github.com/Smartling/smartling-cli/cmd/files/import"
	"github.com/Smartling/smartling-cli/cmd/files/list"
	"github.com/Smartling/smartling-cli/cmd/files/pull"
	"github.com/Smartling/smartling-cli/cmd/files/push"
	"github.com/Smartling/smartling-cli/cmd/files/rename"
	"github.com/Smartling/smartling-cli/cmd/files/status"
	"github.com/Smartling/smartling-cli/services/files"

	"github.com/spf13/cobra"
)

func NewFilesCmd(s *files.Service) *cobra.Command {
	filesCmd := &cobra.Command{
		Use:     "files",
		Aliases: []string{"f"},
		Short:   "Used to access various files sub-commands.",
		Long:    `Used to access various files sub-commands.`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	filesCmd.AddCommand(delete.NewDeleteCmd(s))
	filesCmd.AddCommand(importcmd.NewImportCmd(s))
	filesCmd.AddCommand(list.NewListCmd(s))
	filesCmd.AddCommand(pull.NewPullCmd(s))
	filesCmd.AddCommand(push.NewPushCmd(s))
	filesCmd.AddCommand(rename.NewRenameCmd(s))
	filesCmd.AddCommand(status.NewStatusCmd(s))

	return filesCmd
}
