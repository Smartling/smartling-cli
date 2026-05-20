package build

import (
	"fmt"

	"github.com/Smartling/smartling-cli/cmd/helpers/build"

	"github.com/spf13/cobra"
)

// NewBuildCmd creates a new build command.
func NewBuildCmd() *cobra.Command {
	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "Print the build information",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(build.String())
		},
	}
	return buildCmd
}
