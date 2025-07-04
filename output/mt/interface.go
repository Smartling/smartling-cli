package mt

import "github.com/charmbracelet/bubbles/table"

type Renderer interface {
	Init(dataProvider TableDataProvider, files []string, outputFormat, outputTemplate string)
	Run() error
	Update(updates chan any)
	End()
}

type TableDataProvider interface {
	Headers() []table.Column
	RowByHeaderName() RowByHeaderName
	ToTableRows(files []string) []table.Row
}
