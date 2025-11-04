package mt

import (
	"path/filepath"
	"strings"

	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	"github.com/Smartling/smartling-cli/services/mt"

	"github.com/charmbracelet/bubbles/table"
)

const (
	// DefaultDetectTemplate is the default template used for rendering detected files.
	DefaultDetectTemplate = "{{.File}}\t{{.Language}}" + "\n"
	// DefaultShortDetectTemplate is the default template used for rendering detected files in short format.
	DefaultShortDetectTemplate = "{{.Language}}" + "\n"
)

// DetectCellCoords represents the column positions (if present) for each detect-related action.
type DetectCellCoords struct {
	LanguageCol *uint8
	UploadCol   *uint8
	DetectCol   *uint8
}

// DetectUpdateRow defines a row update operation
type DetectUpdateRow struct {
	RowByHeader RowByHeaderName
	Updates     mt.DetectUpdates
}

// RenderDetectUpdates applies detect updates to the given table model row
func RenderDetectUpdates(t *table.Model, rowByHeader RowByHeaderName, val mt.DetectUpdates) {
	rows := t.Rows()
	if val.ID >= uint32(len(rows)) {
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

func toDetectTableRows(files []string, targetLocalesQnt uint8) []table.Row {
	res := make([]table.Row, len(files)*int(targetLocalesQnt))
	for i, v := range files {
		for j := uint8(0); j < targetLocalesQnt; j++ {
			res[2*i+int(j)] = toDetectTableRow(v)
		}
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

// DetectDataProvider defines data provider for detect flow
type DetectDataProvider struct{}

// Headers returns headers
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

// RowByHeaderName returns a mapping from header names by their column indices
func (t DetectDataProvider) RowByHeaderName() RowByHeaderName {
	return RowByHeaderName{
		"upload":   4,
		"detect":   5,
		"language": 6,
	}
}

// ToTableRows converts slice with files to slice with table rows
func (t DetectDataProvider) ToTableRows(files []string, targetLocalesQnt uint8) []table.Row {
	return toDetectTableRows(files, targetLocalesQnt)
}
