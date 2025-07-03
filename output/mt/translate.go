package mt

import (
	"path/filepath"
	"strings"

	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	"github.com/Smartling/smartling-cli/services/mt"

	"github.com/charmbracelet/bubbles/table"
)

const (
	DefaultTranslateTemplate = "{{name .File}}_{{.Locale}}{{ext .File}}"
)

type TranslateCellCoords struct {
	LocaleCol    *uint8
	UploadCol    *uint8
	TranslateCol *uint8
	DownloadCol  *uint8
}
type RowByHeaderName map[string]uint8

type TranslateUpdateRow struct {
	RowByHeader RowByHeaderName
	Updates     mt.TranslateUpdates
}

func RenderTranslateUpdates(t *table.Model, rowByHeader RowByHeaderName, val mt.TranslateUpdates) {
	rows := t.Rows()
	if val.ID < 0 || val.ID >= uint32(len(rows)) {
		rlog.Debugf("row out of range: %d > %d", val.ID, len(rows))
		return
	}

	t.SetCursor(int(val.ID))

	updatedRow := make([]string, len(rows[val.ID]))
	copy(updatedRow, rows[val.ID])

	if row, found := rowByHeader["locale"]; found && val.Locale != nil {
		updatedRow[row] = *val.Locale
	}
	if row, found := rowByHeader["upload"]; found {
		updatedRow[row] = done
	}
	if row, found := rowByHeader["translate"]; found && val.Translate != nil {
		updatedRow[row] = *val.Translate
	}
	if row, found := rowByHeader["download"]; found && val.Download != nil {
		updatedRow[row] = done
	}

	updatedRows := make([]table.Row, len(rows))
	copy(updatedRows, rows)
	updatedRows[val.ID] = updatedRow

	t.SetRows(updatedRows)
}

func toTranslateTableRows(files []string) []table.Row {
	res := make([]table.Row, len(files))
	for i, v := range files {
		res[i] = toTranslateTableRow(v)
	}
	return res
}

func toTranslateTableRow(file string) table.Row {
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

type TranslateDataProvider struct {
	data []table.Row
}

func (t TranslateDataProvider) Headers() []table.Column {
	return []table.Column{
		{Title: "File", Width: 10},
		{Title: "Locale", Width: 10},
		{Title: "Name", Width: 10},
		{Title: "Ext", Width: 10},
		{Title: "Directory", Width: 10},
		{Title: "Upload", Width: 10},
		{Title: "Translate", Width: 10},
		{Title: "Download", Width: 10},
	}
}

func (t TranslateDataProvider) RowByHeaderName() RowByHeaderName {
	return RowByHeaderName{
		"locale":    1,
		"upload":    5,
		"translate": 6,
		"download":  7,
	}
}

func (t TranslateDataProvider) GetRows() []table.Row {
	return t.data
}

func (t TranslateDataProvider) SetRows(rows []table.Row) {
	t.data = rows
}

func (t TranslateDataProvider) UpdateCell(i, j uint, val string) {
	t.data[i][j] = val
}

func (t TranslateDataProvider) ToTableRows(files []string) []table.Row {
	return toTranslateTableRows(files)
}
