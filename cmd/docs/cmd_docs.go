package docs

import (
	"fmt"

	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// NewDocsCmd creates a new command to generate docs.
func NewDocsCmd() *cobra.Command {
	docsCmd := &cobra.Command{
		Use:   "docs",
		Short: "Generate markdown docs for CLI commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := doc.GenMarkdownTree(cmd.Root(), "./docs")
			if err != nil {
				return fmt.Errorf("failed to generate docs: %v", err)
			}
			rlog.Infof("markdown docs generated in ./docs/")
			return nil
		},
	}

	return docsCmd
}
