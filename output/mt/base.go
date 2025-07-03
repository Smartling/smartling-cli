package mt

import (
	"fmt"
	"os"

	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Headers        []table.Column
	Data           []table.Row
	RowByHeader    RowByHeaderName
	OutputFormat   string
	OutputTemplate string
}

func RenderAndExitIfErr(err error) {
	if err == nil {
		return
	}
	fmt.Println(RenderError(err))
	os.Exit(1)
}

func RenderError(err error) string {
	uiErr, isUIError := err.(clierror.UIError)
	if !isUIError {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F5F")).
			MarginTop(1)
		return errorStyle.Render(fmt.Sprintf(" !ERR Error: %s", err.Error()))
	}

	header := fmt.Sprintf(" !ERR [%s]: %s", uiErr.Operation, uiErr.Err.Error())
	var fields string
	for k, v := range uiErr.Fields {
		fields += fmt.Sprintf("  • %s: %s\n", k, v)
	}

	text := fmt.Sprintf("%s\n\n%s", header, fields)

	if uiErr.Description != "" {
		text += fmt.Sprintf("\nDescription:\n%s\n", uiErr.Description)
	}

	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#FF5F5F")).
		Padding(1, 2).
		Render(text)
}
