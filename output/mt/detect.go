package mt

import (
	"path/filepath"
	"strings"

	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	"github.com/Smartling/smartling-cli/services/mt"

	"github.com/charmbracelet/bubbles/table"
)

const DefaultDetectTemplate = `{{.File}}\t{{.Language}}\n`

type DetectCellCoords struct {
	LanguageCol *uint8
	UploadCol   *uint8
	DetectCol   *uint8
}

type DetectUpdateRow struct {
	RowByHeader RowByHeaderName
	Updates     mt.DetectUpdates
}

func RenderDetectUpdates(t *table.Model, rowByHeader RowByHeaderName, val mt.DetectUpdates) {
	rows := t.Rows()
	if val.ID < 0 || val.ID >= uint32(len(rows)) {
		rlog.Debugf("row out of range: %d > %d", val.ID, len(rows))
		return
	}

	t.SetCursor(int(val.ID))

	updatedRow := make([]string, len(rows[val.ID]))
	copy(updatedRow, rows[val.ID])

	if row, found := rowByHeader["language"]; found && val.Language != nil {
		updatedRow[row] = *val.Language
	}
	if row, found := rowByHeader["upload"]; found {
		updatedRow[row] = done
	}
	if row, found := rowByHeader["detect"]; found && val.Detect != nil {
		updatedRow[row] = *val.Detect
	}

	updatedRows := make([]table.Row, len(rows))
	copy(updatedRows, rows)
	updatedRows[val.ID] = updatedRow

	t.SetRows(updatedRows)
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

type DetectDataProvider struct {
	data []table.Row
}

func (t DetectDataProvider) Headers() []table.Column {
	return []table.Column{
		{Title: "File", Width: 10},
		{Title: "Name", Width: 10},
		{Title: "Ext", Width: 10},
		{Title: "Directory", Width: 10},
		{Title: "Upload", Width: 10},
		{Title: "Detect", Width: 10},
		{Title: "Language", Width: 10},
	}
}

func (t DetectDataProvider) RowByHeaderName() RowByHeaderName {
	return RowByHeaderName{
		"upload":   4,
		"detect":   5,
		"language": 6,
	}
}

func (t DetectDataProvider) GetRows() []table.Row {
	return t.data
}

func (t DetectDataProvider) SetRows(rows []table.Row) {
	t.data = rows
}

func (t DetectDataProvider) UpdateCell(i, j uint, val string) {
	t.data[i][j] = val
}

func (t DetectDataProvider) ToTableRows(files []string) []table.Row {
	return toDetectTableRows(files)
}
