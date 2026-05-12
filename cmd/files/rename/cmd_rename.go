package rename

import (
	"os"

	filescmd "github.com/Smartling/smartling-cli/cmd/files"
	"github.com/Smartling/smartling-cli/services/helpers/help"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/spf13/cobra"
)

var (
	old string
	new string
)

// NewRenameCmd creates a new command to rename files.
func NewRenameCmd(initializer filescmd.SrvInitializer) *cobra.Command {
	renameCmd := &cobra.Command{
		Use:   "rename <old> <new>",
		Short: "Renames given file by old URI into new URI.",
		Long: `smartling-cli files rename — rename specified file.

Renames specified file URI into new file URI.

Available options:
  -p --project <project>
    Specify project to use.
` + help.AuthenticationOptions,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			if len(args) > 0 {
				old = args[0]
			}
			if len(args) > 1 {
				new = args[1]
			}

			s, err := initializer.InitFilesSrv(ctx)
			if err != nil {
				rlog.Errorf("failed to get files service: %s", err)
				os.Exit(1)
			}

			err = s.RunRename(ctx, old, new)
			if err != nil {
				rlog.Errorf("failed to run rename: %s", err)
				os.Exit(1)
			}
		},
	}

	return renameCmd
}
