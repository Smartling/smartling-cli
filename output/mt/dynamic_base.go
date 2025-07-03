package mt

import (
	"fmt"

	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const done = "✓"

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type DynamicTableModel struct {
	table table.Model
	err   error
}

func (m DynamicTableModel) Init() tea.Cmd { return nil }

func (m DynamicTableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case TranslateUpdateRow:
		RenderTranslateUpdates(&m.table, msg.RowByHeader, msg.Updates)
	case DetectUpdateRow:
		RenderDetectUpdates(&m.table, msg.RowByHeader, msg.Updates)
	case clierror.UIError:
		m.err = msg
		return m, tea.Quit
	case error:
		m.err = msg
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m DynamicTableModel) View() string {
	s := m.table.View()

	if m.err != nil {
		s += "\n\n" + RenderError(m.err)
	}

	s += "\n\nPress 'q' to quit."
	return s + "\n"
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
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#FF5F5F")).
		Padding(1, 2).
		Render(fmt.Sprintf("%s\n\n%s", header, fields))
}
