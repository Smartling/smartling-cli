package locales

import (
	projectscmd "github.com/Smartling/smartling-cli/cmd/projects"
	"github.com/Smartling/smartling-cli/services/helpers/format"
	"github.com/Smartling/smartling-cli/services/helpers/help"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	"github.com/Smartling/smartling-cli/services/projects"

	"github.com/spf13/cobra"
)

var (
	short      bool
	source     bool
	formatType string
)

// NewLocalesCmd creates a new command to list locales.
func NewLocalesCmd(initializer projectscmd.SrvInitializer) *cobra.Command {
	localesCmd := &cobra.Command{
		Use:   "locales",
		Short: "Display list of target locales.",
		Long: `smartling-cli projects locales — list target locales.

Lists target locales from specified project.

To list only locale IDs --short option can be used.
` + help.FormatOption + `
Following variables are available:

  > .LocaleID — target locale ID to translate into;
  > .Description — human-readable locale description;
  > .Enabled — true/false specifying is locale active or not;


Available options:
  -p --project <project>
    Specify project to use.

  -s --short
    List only locale IDs.

  --format
    Use specific output format instead of default.
` + help.AuthenticationOptions,
		Run: func(_ *cobra.Command, _ []string) {
			s, err := initializer.InitProjectsSrv()
			if err != nil {
				rlog.Errorf("failed to get project service: %s", err)
				return
			}

			params := projects.LocalesParams{
				Format: formatType,
				Short:  short,
				Source: source,
			}
			err = s.RunLocales(params)
			if err != nil {
				rlog.Errorf("failed to run locales: %s", err)
				return
			}
		},
	}
	localesCmd.Flags().BoolVarP(&short, "short", "s", false, "Display only target locale IDs.")
	localesCmd.Flags().BoolVar(&source, "source", false, "Source.")
	localesCmd.Flags().StringVar(&formatType, "format", "", `Use specified format for listing locales.
                           Format: `+format.DefaultProjectsLocalesFormat)

	return localesCmd
}
