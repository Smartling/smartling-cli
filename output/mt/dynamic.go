package mt

import (
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/mt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Dynamic struct {
	model        Model
	program      *tea.Program
	dataProvider TableDataProvider
}

func (d *Dynamic) Init(dataProvider TableDataProvider, files []string, outputFormat, outputTemplate string) {
	d.model.OutputFormat = outputFormat
	d.model.OutputTemplate = outputTemplate
	d.model.Headers = dataProvider.Headers()
	d.model.RowByHeader = dataProvider.RowByHeaderName()

	rows := dataProvider.ToTableRows(files)
	//dataProvider.SetRows(rows)

	d.model.Data = rows

	t := table.New(
		table.WithColumns(d.model.Headers),
		table.WithRows(rows),
		table.WithFocused(true),
		//table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderTop(true).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#c5c5c5")).
		Bold(false)
	t.SetStyles(s)

	m := DynamicTableModel{table: t}
	d.program = tea.NewProgram(m)

}

func (d *Dynamic) Run() error {
	_, err := d.program.Run()
	return err
}

func (d *Dynamic) Update(updates chan any) error {
	for update := range updates {
		switch update := update.(type) {
		case mt.TranslateUpdates:
			updateRow := TranslateUpdateRow{
				RowByHeader: d.model.RowByHeader,
				Updates:     update,
			}
			d.program.Send(updateRow)
		case mt.DetectUpdates:
			updateRow := DetectUpdateRow{
				RowByHeader: d.model.RowByHeader,
				Updates:     update,
			}
			d.program.Send(updateRow)
		case clierror.UIError:
			d.program.Send(update)
		case error:
			d.program.Send(update)
		}
	}
	return nil
}

func (d *Dynamic) End() {
	d.program.Quit()
}
