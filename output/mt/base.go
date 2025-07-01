package mt

import (
	"fmt"

	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type Model struct {
	table  table.Model
	output string
	err    error
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case UpdateRow:
		RenderUpdates(&m.table, msg.Coords, msg.Updates)
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
	default:
		rlog.Debugf("unexpected message of type %T: %#v\n", msg, msg)
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	s := m.table.View()

	if m.err != nil {
		return RenderError(m.err)
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
