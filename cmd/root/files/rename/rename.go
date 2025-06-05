package rename

import (
	"github.com/spf13/cobra"
)

var (
	old string
	new string
)

func NewRenameCmd() *cobra.Command {
	renameCmd := &cobra.Command{
		Use:   "rename <old> <new>",
		Short: "Renames given file by old URI into new URI.",
		Long:  `Renames given file by old URI into new URI.`,
		Run: func(cmd *cobra.Command, args []string) {
			old = args[0]
			new = args[1]
		},
	}

	return renameCmd
}
