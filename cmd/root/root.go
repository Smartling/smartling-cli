package root

import (
	"github.com/Smartling/smartling-cli/cmd/root/files"
	"github.com/Smartling/smartling-cli/cmd/root/init"
	"github.com/spf13/cobra"
)

var (
	verbose bool
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "smartling-cli",
		Short:   "Manage translation files using Smartling CLI.",
		Version: "1.7",
		Long: `Manage translation files using Smartling CLI.
                Complete documentation is available at https://www.smartling.com`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose logging")

	rootCmd.AddCommand(init.NewInitCmd())
	rootCmd.AddCommand(files.NewFilesCmd())

	return rootCmd
}
