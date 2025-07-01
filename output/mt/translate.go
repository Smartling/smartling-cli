package mt

import (
	"path/filepath"
	"strings"

	"github.com/Smartling/smartling-cli/services/helpers/pointer"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	"github.com/Smartling/smartling-cli/services/mt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	DefaultTranslateTemplate = "{{name .File}}_{{.Locale}}{{ext .File}}"
	done                     = "✓"
)

func RenderUpdates(t *table.Model, cellCoords CellCoords, val mt.TranslateUpdates) {
	rows := t.Rows()
	if val.ID < 0 || val.ID >= uint32(len(rows)) {
		rlog.Debugf("row out of range: %d > %d", val.ID, len(rows))
		return
	}

	updatedRow := make([]string, len(rows[val.ID]))
	copy(updatedRow, rows[val.ID])

	if cellCoords.LocaleCol != nil && val.Locale != nil {
		updatedRow[*cellCoords.LocaleCol] = *val.Locale
	}
	if cellCoords.UploadCol != nil && val.Upload != nil && *val.Upload {
		updatedRow[*cellCoords.UploadCol] = done
	}
	if cellCoords.TranslateCol != nil && val.Translate != nil {
		updatedRow[*cellCoords.TranslateCol] = *val.Translate
	}
	if cellCoords.DownloadCol != nil && val.Download != nil && *val.Download {
		updatedRow[*cellCoords.DownloadCol] = done
	}

	updatedRows := make([]table.Row, len(rows))
	copy(updatedRows, rows)
	updatedRows[val.ID] = updatedRow

	t.SetRows(updatedRows)
}

type CellCoords struct {
	LocaleCol    *uint8
	UploadCol    *uint8
	TranslateCol *uint8
	DownloadCol  *uint8
}

type UpdateRow struct {
	Coords  CellCoords
	Updates mt.TranslateUpdates
}

// RenderFiles renders files
func RenderFiles(files []string, outputFormat, outputTemplate string) (*tea.Program, CellCoords, error) {
	columns := []table.Column{
		{Title: "File", Width: 10},
		{Title: "Locale", Width: 10},
		{Title: "Name", Width: 10},
		{Title: "Ext", Width: 10},
		{Title: "Directory", Width: 10},
		{Title: "Upload", Width: 10},
		{Title: "Translate", Width: 10},
		{Title: "Download", Width: 10},
	}
	cellCoords := CellCoords{
		LocaleCol:    pointer.NewP(uint8(1)),
		UploadCol:    pointer.NewP(uint8(5)),
		TranslateCol: pointer.NewP(uint8(6)),
		DownloadCol:  pointer.NewP(uint8(7)),
	}
	rows := toTableRows(files)
	t := table.New(
		table.WithColumns(columns),
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
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := Model{table: t}
	program := tea.NewProgram(m)

	return program, cellCoords, nil
}

func toTableRows(files []string) []table.Row {
	res := make([]table.Row, len(files))
	for i, v := range files {
		res[i] = toTableRow(v)
	}
	return res
}

func toTableRow(file string) table.Row {
	filename := filepath.Base(file)
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)
	dir := filepath.Dir(file)
	return table.Row{
		filename,
		"",
		name,
		ext,
		dir,
		"",
		"",
		"",
	}
}
