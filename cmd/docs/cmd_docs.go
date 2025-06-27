package docs

import (
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// NewDocsCmd creates a new command to generate docs.
func NewDocsCmd() *cobra.Command {
	docsCmd := &cobra.Command{
		Use:   "docs",
		Short: "Generate markdown docs for CLI commands",
		Run: func(cmd *cobra.Command, args []string) {
			err := doc.GenMarkdownTree(cmd.Root(), "./docs")
			if err != nil {
				rlog.Errorf("failed to generate docs: %v", err)
				os.Exit(1)
			}
			rlog.Infof("markdown docs generated in ./docs/")
		},
	}

	return docsCmd
}
