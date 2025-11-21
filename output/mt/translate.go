package mt

import (
	"path/filepath"
	"strings"

	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	"github.com/Smartling/smartling-cli/services/mt"

	"github.com/charmbracelet/bubbles/table"
)

const (
	// DefaultTranslateTemplate is the default template used for rendering translated files.
	DefaultTranslateTemplate = "{{.File}}\t{{.Locale}}\t{{.TranslatedFile}}" + "\n"
)

// TranslateCellCoords represents the column positions (if present) for each translation-related action.
type TranslateCellCoords struct {
	LocaleCol    *uint8
	UploadCol    *uint8
	TranslateCol *uint8
	DownloadCol  *uint8
}

// TranslateUpdateRow defines a row update operation
type TranslateUpdateRow struct {
	RowByHeader RowByHeaderName
	Updates     mt.TranslateUpdates
}

// RowByHeaderName defines row position by name
type RowByHeaderName map[string]uint8

// RenderTranslateUpdates applies translation updates to the given table model row
func RenderTranslateUpdates(t *table.Model, rowByHeader RowByHeaderName, val mt.TranslateUpdates) {
	rows := t.Rows()
	if val.ID >= uint32(len(rows)) {
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
	if row, found := rowByHeader["translated_file"]; found && val.TranslatedFile != nil {
		updatedRow[row] = *val.TranslatedFile
	}
	if row, found := rowByHeader["download"]; found && val.Download != nil {
		updatedRow[row] = done
	}

	updatedRows := make([]table.Row, len(rows))
	copy(updatedRows, rows)
	updatedRows[val.ID] = updatedRow

	t.SetRows(updatedRows)
}

func toTranslateTableRows(files []string, targetLocalesQnt uint8) []table.Row {
	res := make([]table.Row, len(files)*int(targetLocalesQnt))
	for i, v := range files {
		for j := uint8(0); j < targetLocalesQnt; j++ {
			res[int(targetLocalesQnt)*i+int(j)] = toTranslateTableRow(v)
		}
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
		"",
	}
}

// TranslateDataProvider defines data provider for translate flow
type TranslateDataProvider struct{}

// Headers returns headers
func (t TranslateDataProvider) Headers() []table.Column {
	return []table.Column{
		{Title: "File", Width: 10},
		{Title: "Locale", Width: 10},
		{Title: "Name", Width: 10},
		{Title: "Ext", Width: 10},
		{Title: "Directory", Width: 10},
		{Title: "Upload", Width: 10},
		{Title: "Translate", Width: 10},
		{Title: "TranslatedFile", Width: 10},
		{Title: "Download", Width: 10},
	}
}

// RowByHeaderName returns a mapping from header names by their column indices
func (t TranslateDataProvider) RowByHeaderName() RowByHeaderName {
	return RowByHeaderName{
		"locale":          1,
		"upload":          5,
		"translate":       6,
		"translated_file": 7,
		"download":        8,
	}
}

// ToTableRows converts slice with files to slice with table rows
func (t TranslateDataProvider) ToTableRows(files []string, targetLocalesQnt uint8) []table.Row {
	return toTranslateTableRows(files, targetLocalesQnt)
}
