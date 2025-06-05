package files

import (
	"github.com/Smartling/smartling-cli/cmd/files/delete"
	importcmd "github.com/Smartling/smartling-cli/cmd/files/import"
	"github.com/Smartling/smartling-cli/cmd/files/list"
	"github.com/Smartling/smartling-cli/cmd/files/pull"
	"github.com/Smartling/smartling-cli/cmd/files/push"
	"github.com/Smartling/smartling-cli/cmd/files/rename"
	"github.com/Smartling/smartling-cli/cmd/files/status"

	"github.com/spf13/cobra"
)

func NewFilesCmd() *cobra.Command {
	filesCmd := &cobra.Command{
		Use:     "files",
		Aliases: []string{"f"},
		Short:   "Used to access various files sub-commands.",
		Long:    `Used to access various files sub-commands.`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	filesCmd.AddCommand(delete.NewDeleteCmd())
	filesCmd.AddCommand(importcmd.NewImportCmd())
	filesCmd.AddCommand(list.NewListCmd())
	filesCmd.AddCommand(pull.NewPullCmd())
	filesCmd.AddCommand(push.NewPushCmd())
	filesCmd.AddCommand(rename.NewRenameCmd())
	filesCmd.AddCommand(status.NewStatusCmd())

	return filesCmd
}
