package projects

import (
	"fmt"
	"io"
	"os"

	"github.com/Smartling/smartling-cli/services/helpers/table"
	projectconfig "github.com/Smartling/smartling-cli/services/projects/config"
)

func RenderTable(config projectconfig.Extended) error {
	tableWriter := table.NewTableWriter(os.Stdout)
	for _, row := range buildInfoRows(config) {
		if _, err := fmt.Fprintf(tableWriter, "%s\t%s\n", row...); err != nil {
			return err
		}
	}
	return table.Render(tableWriter)
}

// RenderPlain writes a multi-line summary of the resolved configuration to w.
// Used by the --show-config persistent flag. Every line is prefixed with "> "
// so the banner is easy to grep out of mixed output if needed.
func RenderPlain(w io.Writer, config projectconfig.Extended) error {
	lines := []string{
		"Smartling CLI configuration:",
		fmt.Sprintf("  Config file: %s", config.ConfigFile),
		fmt.Sprintf("  User:     %s", config.UserID),
		fmt.Sprintf("  Account:  %s", config.AccountUID),
		fmt.Sprintf("  Project:  %s", config.ProjectID),
		fmt.Sprintf("  Project Name:  %s", config.Name),
		fmt.Sprintf("  Locale:  %s", config.Locale),
		fmt.Sprintf("  Status:  %s", config.Status),
		fmt.Sprintf("  Sources:  %s", config.Sources),
	}
	for _, line := range lines {
		if _, err := fmt.Fprintf(w, "> %s\n", line); err != nil {
			return err
		}
	}
	return nil
}

func buildInfoRows(config projectconfig.Extended) [][]any {
	return [][]any{
		{"ID", config.ProjectID},
		{"ACCOUNT", config.AccountUID},
		{"NAME", config.Name},
		{"LOCALE", config.Locale},
		{"STATUS", config.Status},
		{"USER", config.UserID},
		{"CONFIG", config.ConfigFile},
		{"SOURCES", config.Sources},
	}
}
