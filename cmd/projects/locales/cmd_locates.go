package locales

import (
	"github.com/Smartling/smartling-cli/services/projects"

	"github.com/spf13/cobra"
)

var (
	short  bool
	source bool
	format string
)

func NewLocatesCmd(s *projects.Service) *cobra.Command {
	locatesCmd := &cobra.Command{
		Use:   "locates",
		Short: "Display list of target locales.",
		Long:  `Display list of target locales.`,
		Run: func(cmd *cobra.Command, args []string) {
			params := projects.LocalesParams{
				Format: format,
				Short:  short,
				Source: source,
			}
			err := s.RunLocales(params)
			if err != nil {
				// TODO log it
			}
		},
	}
	locatesCmd.Flags().BoolVarP(&short, "short", "s", false, "Display only target locale IDs.")
	locatesCmd.Flags().BoolVar(&source, "source", false, "Source.")
	locatesCmd.Flags().StringVar(&format, "format", "", `Use specified format for listing locales.
                           [format: $PROJECTS_LOCALES_FORMAT]`)

	return locatesCmd
}
