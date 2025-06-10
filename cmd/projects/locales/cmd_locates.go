package locales

import (
	projectscmd "github.com/Smartling/smartling-cli/cmd/projects"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	"github.com/Smartling/smartling-cli/services/projects"

	"github.com/spf13/cobra"
)

var (
	short  bool
	source bool
	format string
)

func NewLocatesCmd() *cobra.Command {
	locatesCmd := &cobra.Command{
		Use:   "locates",
		Short: "Display list of target locales.",
		Long:  `Display list of target locales.`,
		Run: func(cmd *cobra.Command, args []string) {
			s, err := projectscmd.GetService()
			if err != nil {
				rlog.Errorf("failed to get project service: %s", err)
				return
			}

			params := projects.LocalesParams{
				Format: format,
				Short:  short,
				Source: source,
			}
			err = s.RunLocales(params)
			if err != nil {
				rlog.Errorf("failed to run locates: %s", err)
				return
			}
		},
	}
	locatesCmd.Flags().BoolVarP(&short, "short", "s", false, "Display only target locale IDs.")
	locatesCmd.Flags().BoolVar(&source, "source", false, "Source.")
	locatesCmd.Flags().StringVar(&format, "format", "", `Use specified format for listing locales.
                           [format: $PROJECTS_LOCALES_FORMAT]`)

	return locatesCmd
}
