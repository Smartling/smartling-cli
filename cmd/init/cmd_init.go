package init

import (
	"github.com/Smartling/smartling-cli/services/init"

	"github.com/spf13/cobra"
)

var (
	dryRun bool
)

func NewInitCmd(s init.Service, verbose uint8) *cobra.Command {
	initCmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Prepares project to work with Smartling",
		Long: `Prepares project to work with Smartling,
essentially, assisting user in creating
configuration file.`,
		Run: func(cmd *cobra.Command, args []string) {
			err := s.RunInit(dryRun, verbose)
			if err != nil {
				// TODO log it
			}
		},
	}
	initCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Do not actually write file, just output it on stdout.")

	return initCmd
}
