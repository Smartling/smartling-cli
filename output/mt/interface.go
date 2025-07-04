package mt

import "github.com/charmbracelet/bubbles/table"

// Renderer defines behaviour for output rendering.
type Renderer interface {
	Init(dataProvider TableDataProvider, files []string, outputFormat, outputTemplate string)
	Run() error
	Update(updates chan any) error
	End()
}

// TableDataProvider defines behaviour for providing tabular data
type TableDataProvider interface {
	Headers() []table.Column
	RowByHeaderName() RowByHeaderName
	ToTableRows(files []string) []table.Row
}
