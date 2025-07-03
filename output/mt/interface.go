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
	GetRows() []table.Row
	SetRows(rows []table.Row)
	UpdateCell(i, j uint, val string)
	ToTableRows(files []string) []table.Row
}

type Static struct{}
