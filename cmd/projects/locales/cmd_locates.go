package locales

import (
	"github.com/spf13/cobra"
)

var (
	short  bool
	format string
)

func NewLocatesCmd() *cobra.Command {
	locatesCmd := &cobra.Command{
		Use:   "locates",
		Short: "Display list of target locales.",
		Long:  `Display list of target locales.`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	locatesCmd.Flags().BoolVarP(&short, "short", "s", false, "Display only target locale IDs.")
	locatesCmd.Flags().StringVar(&format, "format", "", `Use specified format for listing locales.
                           [format: $PROJECTS_LOCALES_FORMAT]`)

	return locatesCmd
}
