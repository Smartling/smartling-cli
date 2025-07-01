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

const DefaultDetectTemplate = `{{.File}}\t{{.Language}}\n`

type DetectCellCoords struct {
	LanguageCol *uint8
	UploadCol   *uint8
	DetectCol   *uint8
}

type DetectUpdateRow struct {
	Coords  DetectCellCoords
	Updates mt.DetectUpdates
}

func RenderDetectUpdates(t *table.Model, cellCoords DetectCellCoords, val mt.DetectUpdates) {
	rows := t.Rows()
	if val.ID < 0 || val.ID >= uint32(len(rows)) {
		rlog.Debugf("row out of range: %d > %d", val.ID, len(rows))
		return
	}

	t.SetCursor(int(val.ID))

	updatedRow := make([]string, len(rows[val.ID]))
	copy(updatedRow, rows[val.ID])

	if cellCoords.LanguageCol != nil && val.Language != nil {
		updatedRow[*cellCoords.LanguageCol] = *val.Language
	}
	if cellCoords.UploadCol != nil && val.Upload != nil && *val.Upload {
		updatedRow[*cellCoords.UploadCol] = done
	}
	if cellCoords.DetectCol != nil && val.Detect != nil {
		updatedRow[*cellCoords.DetectCol] = *val.Detect
	}

	updatedRows := make([]table.Row, len(rows))
	copy(updatedRows, rows)
	updatedRows[val.ID] = updatedRow

	t.SetRows(updatedRows)
}

// RenderDetectFiles renders files
func RenderDetectFiles(files []string, outputFormat, outputTemplate string) (*tea.Program, DetectCellCoords, error) {
	columns := []table.Column{
		{Title: "File", Width: 10},
		{Title: "Name", Width: 10},
		{Title: "Ext", Width: 10},
		{Title: "Directory", Width: 10},
		{Title: "Upload", Width: 10},
		{Title: "Detect", Width: 10},
		{Title: "Language", Width: 10},
	}
	cellCoords := DetectCellCoords{
		UploadCol:   pointer.NewP(uint8(4)),
		DetectCol:   pointer.NewP(uint8(5)),
		LanguageCol: pointer.NewP(uint8(6)),
	}
	rows := toDetectTableRows(files)
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
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#c5c5c5")).
		Bold(false)
	t.SetStyles(s)

	m := Model{table: t}
	program := tea.NewProgram(m)

	return program, cellCoords, nil
}

func toDetectTableRows(files []string) []table.Row {
	res := make([]table.Row, len(files))
	for i, v := range files {
		res[i] = toDetectTableRow(v)
	}
	return res
}

func toDetectTableRow(file string) table.Row {
	filename := filepath.Base(file)
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)
	dir := filepath.Dir(file)
	return table.Row{
		filename,
		name,
		ext,
		dir,
		"",
		"",
		"",
	}
}
